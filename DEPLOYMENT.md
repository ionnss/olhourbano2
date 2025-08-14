# VPS Deployment Guide

This guide explains how to securely deploy your Olho Urbano application to a VPS using git.

## 🚨 IMPORTANT: Secrets Are Not in Git

Your secrets are **intentionally excluded** from git (see `.gitignore`):
```
secrets/          # ← All secret files are ignored
.env              # ← Environment file is ignored
```

This means you need to **manually handle secrets** on your VPS.

## 📋 Deployment Checklist

### 1. **Clone Repository to VPS**
```bash
# On your VPS
git clone https://github.com/yourusername/olhourbano2.git
cd olhourbano2
```

### 2. **Create Secrets Directory**
```bash
# Create the secrets directory
mkdir -p secrets

# Set proper permissions
chmod 700 secrets
```

### 3. **Add Secret Files**
You need to manually create these files on your VPS:

```bash
# Database password
echo "your-production-db-password" > secrets/db_password.txt

# Postgres password (for Docker)
echo "your-production-postgres-password" > secrets/postgres_password.txt

# Session key (generate a strong one)
echo "your-strong-session-key-for-production" > secrets/session_key.txt

# SMTP password (Gmail app password)
echo "your-gmail-app-password" > secrets/smtp_password.txt

# CPF API key
echo "fa018a9cd28f31e28..." > secrets/cpfhub_api_key.txt

# Google Maps API key
echo "AIzaSyA9IfFHzj9hV..." > secrets/google_maps_api_key.txt
```

### 4. **Set Secure Permissions**
```bash
# Restrict access to secrets
chmod 600 secrets/*.txt

# Verify permissions
ls -la secrets/
# Should show: -rw------- (owner read/write only)
```

### 5. **Create Production Environment File**
```bash
# Create .env for production
cat > .env << 'EOF'
# Database Configuration
DB_HOST=db
DB_PORT=5432
DB_USER=olhourbano
DB_NAME=olhourbanovault

# Postgres Service Configuration
POSTGRES_USER=olhourbano
POSTGRES_DB=olhourbanovault
POSTGRES_PASSWORD=your-production-postgres-password

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=olhourbano.contato@gmail.com

# App Configuration
COOKIE_DOMAIN=.olhourbano.com.br
APP_VERSION=2.0.0

# API URLs
CPFHUB_API_URL=https://api.cpfhub.io/api/cpf
GOOGLE_MAPS_API_URL=https://maps.googleapis.com/maps/api
EOF
```

### 6. **Deploy with Docker Compose**
```bash
# Build and start services
docker-compose up -d --build

# Check logs
docker-compose logs -f backend
```

## 🔐 Security Best Practices for VPS

### File Permissions
```bash
# Project directory
chmod 755 /path/to/olhourbano2

# Secrets directory
chmod 700 secrets/

# Secret files
chmod 600 secrets/*.txt

# Environment file
chmod 600 .env
```

### User Management
```bash
# Create dedicated user for the application
sudo useradd -m -s /bin/bash olhourbano
sudo usermod -aG docker olhourbano

# Deploy as this user
sudo su - olhourbano
```

### Firewall Configuration
```bash
# Basic firewall setup
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP
sudo ufw allow 443   # HTTPS
sudo ufw enable
```

## 🔄 Updates and Maintenance

### Updating Code
```bash
# Pull latest changes
git pull origin main

# Rebuild and restart
docker-compose down
docker-compose up -d --build
```

### Secret Rotation
```bash
# Update a secret (example: API key)
echo "new-api-key" > secrets/cpfhub_api_key.txt
chmod 600 secrets/cpfhub_api_key.txt

# Restart to reload secrets
docker-compose restart backend
```

### Backup Secrets
```bash
# Create encrypted backup of secrets
tar -czf secrets-backup-$(date +%Y%m%d).tar.gz secrets/

# Store backup securely (encrypted storage)
gpg -c secrets-backup-$(date +%Y%m%d).tar.gz
rm secrets-backup-$(date +%Y%m%d).tar.gz
```

## 🚨 Security Warnings

### ❌ DON'T Do This:
```bash
# DON'T: Never commit secrets to git
git add secrets/
git commit -m "Add secrets"  # ← This exposes secrets!

# DON'T: Don't use weak permissions
chmod 777 secrets/  # ← Anyone can read secrets!

# DON'T: Don't put secrets in environment variables
export API_KEY="secret-value"  # ← Visible in process list!
```

### ✅ DO This Instead:
```bash
# DO: Only commit code, never secrets
git add . --exclude=secrets/
git commit -m "Update application"

# DO: Use restrictive permissions
chmod 600 secrets/*.txt

# DO: Use the configuration system
# Secrets are loaded from files by config.Load()
```

## 📊 Deployment Verification

### Check Secret Files
```bash
# Verify all secrets exist
ls -la secrets/
# Should show all .txt files with 600 permissions

# Verify secrets are readable by app
sudo -u olhourbano cat secrets/db_password.txt
```

### Check Application
```bash
# Check logs for configuration loading
docker-compose logs backend | grep -i "configuration"

# Test database connection
docker-compose exec backend ps aux
```

### Check Security
```bash
# Verify no secrets in git
git status --porcelain | grep secrets/
# Should return nothing

# Check file permissions
find secrets/ -type f ! -perm 600
# Should return nothing
```

## 🗄️ Database Migrations

Your application uses a custom migration system to manage database schema changes safely.

### Migration Commands

All migration commands must be run with the correct working directory (`-w /olhourbano2`):

```bash
# Check migration status
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:status

# Apply pending migrations
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate

# Validate migration files
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:validate

# Rollback to specific version (e.g., version 2)
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:rollback 2
```

### Migration Workflow for Deployments

**Before deploying new code:**
```bash
# 1. Check current migration status
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:status

# 2. Pull latest code
git pull

# 3. Rebuild containers (migrations run automatically during startup)
docker compose down
docker compose up --build -d

# 4. Verify migrations were applied
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:status
```

**Emergency rollback:**
```bash
# If you need to rollback database changes
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:rollback [target_version]

# Example: rollback to version 2
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:rollback 2
```

### Database Access

```bash
# Connect to database interactively
docker exec -it olhourbano2-db-1 psql -U olhourbano olhourbanovault

# Run single SQL commands
docker exec olhourbano2-db-1 psql -U olhourbano olhourbanovault -c "SELECT COUNT(*) FROM reports;"

# View migration history
docker exec olhourbano2-db-1 psql -U olhourbano olhourbanovault -c "SELECT * FROM schema_migrations ORDER BY version;"
```

## 🎯 Production Checklist

- [ ] Repository cloned to VPS
- [ ] All secret files created with correct values
- [ ] File permissions set to 600 for secrets
- [ ] .env file created for production
- [ ] Docker compose services running
- [ ] Application logs show successful configuration loading
- [ ] Database connection working
- [ ] **Database migrations applied successfully**
- [ ] **Migration status verified**
- [ ] HTTPS certificates obtained (Caddy handles this)
- [ ] Firewall configured
- [ ] Backups scheduled

## 🆘 Troubleshooting

### "Failed to load X secret" Error
```bash
# Check if file exists
ls -la secrets/

# Check permissions
ls -la secrets/*.txt

# Check file content (be careful not to log!)
wc -c secrets/*.txt  # Shows file sizes without content
```

### Docker Permission Issues
```bash
# Ensure docker group membership
sudo usermod -aG docker $USER
newgrp docker

# Check docker daemon
sudo systemctl status docker
```

### Network Issues
```bash
# Check if services are running
docker-compose ps

# Check port binding
sudo netstat -tlnp | grep :80
sudo netstat -tlnp | grep :443
``` 