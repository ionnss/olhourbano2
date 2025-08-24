# Olho Urbano - Migration and Backup Guide

## Overview

This guide provides step-by-step instructions for migrating the Olho Urbano application between VPS servers, including database backup/restore and file transfer procedures.

## Prerequisites

- Access to both source and destination VPS servers
- Docker and Docker Compose installed on both servers
- SSH access to both servers
- Sufficient disk space for backups

## Migration Process

### Step 1: Prepare Source VPS (Old Server)

#### 1.1 SSH into the source VPS
```bash
ssh root@OLD_VPS_IP
```

#### 1.2 Navigate to project directory
```bash
cd olhourbano2
```

#### 1.3 Stop the application
```bash
docker compose down
```

#### 1.4 Start only the database
```bash
docker compose up -d db
```

#### 1.5 Wait for database to be ready
```bash
sleep 15
```

#### 1.6 Create database backup
```bash
docker exec olhourbano2-db-1 pg_dump -U [DB_USER] [DB_NAME] > backup.sql
```

#### 1.7 Verify backup was created
```bash
ls -la backup.sql
```

### Step 2: Transfer Files to Local Machine

#### 2.1 Download database backup
```bash
# From your local machine
scp root@OLD_VPS_IP:/root/olhourbano2/backup.sql ./
```

#### 2.2 Download uploads folder
```bash
# From your local machine
scp -r root@OLD_VPS_IP:/root/olhourbano2/uploads/ ./
```

### Step 3: Transfer Files to Destination VPS

#### 3.1 Upload database backup
```bash
# From your local machine
scp backup.sql root@NEW_VPS_IP:/root/olhourbano2/
```

#### 3.2 Upload uploads folder
```bash
# From your local machine
scp -r uploads/ root@NEW_VPS_IP:/root/olhourbano2/
```

### Step 4: Restore on Destination VPS

#### 4.1 SSH into destination VPS
```bash
ssh root@NEW_VPS_IP
```

#### 4.2 Navigate to project directory
```bash
cd olhourbano2
```

#### 4.3 Stop the application
```bash
docker compose down
```

#### 4.4 Start only the database
```bash
docker compose up -d db
```

#### 4.5 Wait for database to be ready
```bash
sleep 15
```

#### 4.6 Restore database
```bash
docker exec -i olhourbano2-db-1 psql -U [DB_USER] [DB_NAME] < backup.sql
```

**Note**: You may see some errors about duplicate constraints or existing tables. This is normal if the database already has some structure. The important thing is that the data is restored.

#### 4.7 Verify data restoration
```bash
docker exec olhourbano2-db-1 psql -U [DB_USER] [DB_NAME] -c "SELECT COUNT(*) FROM reports;"
```

#### 4.8 Start the full application
```bash
docker compose up -d
```

#### 4.9 Verify application status
```bash
docker compose ps
```

## Alternative Migration Methods

### Method 1: Direct VPS-to-VPS Migration

If both VPS can communicate directly:

```bash
# On destination VPS
docker exec olhourbano2-db-1 pg_dump -U [DB_USER] -h OLD_VPS_IP [DB_NAME] | docker exec -i olhourbano2-db-1 psql -U [DB_USER] [DB_NAME]
```

**Note**: This method requires PostgreSQL to be configured to accept external connections on the old VPS.

### Method 2: Using rsync for File Transfer

For large uploads folders:

```bash
# From local machine
rsync -avz -e ssh root@OLD_VPS_IP:/root/olhourbano2/uploads/ ./uploads/
rsync -avz -e ssh ./uploads/ root@NEW_VPS_IP:/root/olhourbano2/uploads/
```

## Backup Procedures

### Automated Backup Script

Create a backup script on your VPS:

```bash
#!/bin/bash
# backup_olhourbano.sh

BACKUP_DIR="/root/backups"
DATE=$(date +%Y%m%d_%H%M%S)
PROJECT_DIR="/root/olhourbano2"

# Create backup directory
mkdir -p $BACKUP_DIR

# Database backup
cd $PROJECT_DIR
docker exec olhourbano2-db-1 pg_dump -U [DB_USER] [DB_NAME] | gzip > $BACKUP_DIR/backup_$DATE.sql.gz

# Uploads backup
tar -czf $BACKUP_DIR/uploads_$DATE.tar.gz -C $PROJECT_DIR uploads/

# Keep only last 7 backups
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete
find $BACKUP_DIR -name "uploads_*.tar.gz" -mtime +7 -delete

echo "Backup completed: backup_$DATE.sql.gz, uploads_$DATE.tar.gz"
```

### Cron Job for Automated Backups

Add to crontab:

```bash
# Edit crontab
crontab -e

# Add this line for daily backups at 2 AM
0 2 * * * /root/backup_olhourbano.sh
```

## Verification Checklist

After migration, verify:

- [ ] Database contains expected number of reports
- [ ] All images in uploads folder are accessible
- [ ] Application starts without errors
- [ ] All containers are running
- [ ] Website is accessible
- [ ] User registrations are preserved
- [ ] Uploaded files display correctly

## Troubleshooting

### Common Issues

#### 1. Permission Denied Errors
```bash
# Ensure proper file permissions
chmod 755 /root/olhourbano2/uploads/
chown -R root:root /root/olhourbano2/uploads/
```

#### 2. Database Connection Errors
```bash
# Check if database is running
docker compose ps

# Check database logs
docker compose logs db
```

#### 3. Insufficient Disk Space
```bash
# Check available space
df -h

# Clean up old backups
find /root/backups -name "*.gz" -mtime +7 -delete
```

#### 4. Migration Errors
If you see errors about duplicate constraints:
- This is normal if the database already has structure
- The important data (reports, votes, comments) should still be restored
- Verify with: `SELECT COUNT(*) FROM reports;`

### Emergency Recovery

If the migration fails:

1. **Stop all containers**:
   ```bash
   docker compose down
   ```

2. **Remove existing database**:
   ```bash
   docker volume rm olhourbano2_db_data
   ```

3. **Start fresh database**:
   ```bash
   docker compose up -d db
   ```

4. **Restore from backup**:
   ```bash
   docker exec -i olhourbano2-db-1 psql -U [DB_USER] [DB_NAME] < backup.sql
   ```

## Security Considerations

### Before Migration
- Ensure both VPS have proper firewall rules
- Use SSH keys instead of passwords when possible
- Verify backup integrity before deleting old data

### After Migration
- Update DNS records if applicable
- Update any hardcoded IP addresses in configuration
- Test all functionality thoroughly
- Monitor logs for any issues

## Post-Migration Tasks

1. **Update DNS** (if applicable):
   - Point domain to new VPS IP
   - Update any hardcoded references

2. **Test Functionality**:
   - Submit new reports
   - Upload files
   - Test voting system
   - Verify email functionality

3. **Monitor Performance**:
   - Check application logs
   - Monitor resource usage
   - Verify database performance

4. **Clean Up**:
   - Remove old VPS (after confirming everything works)
   - Archive old backups
   - Update documentation

## Example Migration Commands

Here's a complete example for migrating from VPS A to VPS B:

```bash
# On VPS A (source)
ssh root@VPS_A_IP
cd olhourbano2
docker compose down
docker compose up -d db
sleep 15
docker exec olhourbano2-db-1 pg_dump -U [DB_USER] [DB_NAME] > backup.sql

# On local machine
scp root@VPS_A_IP:/root/olhourbano2/backup.sql ./
scp -r root@VPS_A_IP:/root/olhourbano2/uploads/ ./
scp backup.sql root@VPS_B_IP:/root/olhourbano2/
scp -r uploads/ root@VPS_B_IP:/root/olhourbano2/

# On VPS B (destination)
ssh root@VPS_B_IP
cd olhourbano2
docker compose down
docker compose up -d db
sleep 15
docker exec -i olhourbano2-db-1 psql -U [DB_USER] [DB_NAME] < backup.sql
docker compose up -d
docker exec olhourbano2-db-1 psql -U [DB_USER] [DB_NAME] -c "SELECT COUNT(*) FROM reports;"
```

## Support

For issues with migration:
- Check Docker logs: `docker compose logs`
- Verify database connectivity: `docker exec olhourbano2-db-1 psql -U [DB_USER] [DB_NAME] -c "SELECT 1;"`
- Review this guide for troubleshooting steps

---

**Last Updated**: August 24, 2025  
**Version**: 1.0  
**Tested On**: Ubuntu 24.04 LTS with Docker Compose
