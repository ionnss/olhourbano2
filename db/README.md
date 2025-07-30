# Database & Migrations

This directory contains all database-related code for the Olho Urbano application.

## Directory Structure

```
db/
├── README.md              # This file
├── db.go                  # Database connection and configuration
├── migrations.go          # Custom migration system implementation
└── migrations/            # Migration files directory
    ├── 000001_create_reports_table.up.sql
    ├── 000001_create_reports_table.down.sql
    ├── 000002_create_votes_table.up.sql
    ├── 000002_create_votes_table.down.sql
    ├── 000003_add_reports_indexes.up.sql
    ├── 000003_add_reports_indexes.down.sql
    ├── 000004_add_votes_indexes.up.sql
    └── 000004_add_votes_indexes.down.sql
```

## Migration System

### Overview

The application uses a **custom migration system** built from scratch that provides:

- ✅ **Version tracking** via `schema_migrations` table
- ✅ **File integrity** with MD5 checksums
- ✅ **Transaction safety** for each migration
- ✅ **Rollback capability** using down migrations
- ✅ **Status reporting** to see what's applied
- ✅ **Validation** of migration files

### Migration File Naming Convention

Migration files follow this pattern:
```
{version}_{description}.{direction}.sql
```

- **Version**: 6-digit number (000001, 000002, etc.)
- **Description**: Snake_case description of the change
- **Direction**: `up` (apply) or `down` (rollback)

**Examples:**
- `000001_create_reports_table.up.sql` - Creates the reports table
- `000001_create_reports_table.down.sql` - Drops the reports table
- `000002_create_votes_table.up.sql` - Creates the votes table
- `000002_create_votes_table.down.sql` - Drops the votes table

### Migration Commands

All commands must be run from the correct Docker container with proper working directory:

```bash
# Check which migrations have been applied
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:status

# Apply all pending migrations
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate

# Validate all migration files
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:validate

# Rollback to a specific version
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:rollback 2
```

### Creating New Migrations

To create a new migration:

1. **Determine the next version number**:
   ```bash
   # Check current highest version
   docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:status
   ```

2. **Create both UP and DOWN files**:
   ```bash
   # Example: Adding a new table
   touch db/migrations/000005_create_categories_table.up.sql
   touch db/migrations/000005_create_categories_table.down.sql
   ```

3. **Write the UP migration** (applies the change):
   ```sql
   -- 000005_create_categories_table.up.sql
   CREATE TABLE categories (
       id SERIAL PRIMARY KEY,
       name VARCHAR(100) NOT NULL UNIQUE,
       description TEXT,
       created_at TIMESTAMP DEFAULT NOW()
   );

   CREATE INDEX idx_categories_name ON categories(name);
   ```

4. **Write the DOWN migration** (reverts the change):
   ```sql
   -- 000005_create_categories_table.down.sql
   DROP INDEX IF EXISTS idx_categories_name;
   DROP TABLE IF EXISTS categories;
   ```

5. **Test the migration**:
   ```bash
   # Validate the files
   docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:validate

   # Apply the migration
   docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate

   # Test rollback
   docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:rollback 4

   # Reapply
   docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate
   ```

### Best Practices

#### DO ✅
- **Always create both UP and DOWN migrations**
- **Test rollbacks** before deploying to production
- **Keep migrations atomic** - one logical change per migration
- **Use IF EXISTS/IF NOT EXISTS** for idempotent operations
- **Add comments** explaining complex changes
- **Backup database** before applying migrations in production
- **Run migrations in transactions** (handled automatically)

#### DON'T ❌
- **Mix schema and data changes** in the same migration
- **Reference application code** from migrations (migrations should be self-contained)
- **Use `CONCURRENTLY` in DOWN migrations** (not compatible with transactions)
- **Edit existing migration files** after they've been applied
- **Skip version numbers** (keep them sequential)

### Transaction Compatibility

The migration system runs each migration in its own transaction for safety. However, some PostgreSQL operations cannot run inside transactions:

**❌ Not allowed in transactions:**
- `CREATE INDEX CONCURRENTLY`
- `DROP INDEX CONCURRENTLY`
- `CREATE DATABASE`
- `DROP DATABASE`

**✅ Safe for transactions:**
- `CREATE TABLE`
- `ALTER TABLE`
- `CREATE INDEX` (blocking)
- `DROP INDEX` (blocking)
- `INSERT`, `UPDATE`, `DELETE`

### Database Schema

The migration system creates and manages a `schema_migrations` table:

```sql
CREATE TABLE schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT NOW(),
    checksum VARCHAR(64)
);
```

This table tracks:
- **version**: Which migration was applied (e.g., "000001")
- **applied_at**: When the migration was applied
- **checksum**: MD5 hash of the migration file for integrity

### Current Database Schema

As of the latest migrations, the database contains:

#### Reports Table
```sql
CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    problem_type VARCHAR(100),
    hashed_cpf VARCHAR(64),
    email VARCHAR(100),
    location TEXT,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    description TEXT,
    photo_path TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    vote_count INTEGER DEFAULT 0,
    status VARCHAR(100) NOT NULL DEFAULT 'pending'
);
```

#### Votes Table
```sql
CREATE TABLE votes (
    id SERIAL PRIMARY KEY,
    report_id INTEGER REFERENCES reports(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    vote_hashed_cpf VARCHAR(64),
    UNIQUE(vote_hashed_cpf, report_id)
);
```

#### Indexes
- **Reports**: 8 indexes for optimal query performance
- **Votes**: 5 indexes for fast lookups and constraints
- See migration files for complete index definitions

### Troubleshooting

#### Migration Fails
```bash
# Check the error in logs
docker compose logs backend

# Verify file permissions and syntax
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:validate

# Check database connection
docker exec olhourbano2-db-1 psql -U olhourbano olhourbanovault -c "SELECT 1;"
```

#### Migration Status Issues
```bash
# Check what's in the tracking table
docker exec olhourbano2-db-1 psql -U olhourbano olhourbanovault -c "SELECT * FROM schema_migrations ORDER BY version;"

# Check if migration files exist
docker exec olhourbano2-backend-1 ls -la /olhourbano2/db/migrations/
```

#### Rollback Issues
```bash
# Ensure DOWN migration files exist and are valid
docker exec olhourbano2-backend-1 cat /olhourbano2/db/migrations/000004_add_votes_indexes.down.sql

# Check for transaction compatibility issues
# (CONCURRENTLY operations must be removed from DOWN migrations)
```

### Emergency Procedures

#### Manual Migration Table Reset
```sql
-- ⚠️ DANGEROUS: Only use in development or with proper backups
DELETE FROM schema_migrations WHERE version = '000004';
```

#### Force Migration Reapply
```bash
# 1. Remove from tracking table (database)
docker exec olhourbano2-db-1 psql -U olhourbano olhourbanovault -c "DELETE FROM schema_migrations WHERE version = '000004';"

# 2. Reapply the migration
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate
```

### Development vs Production

#### Development
- Migrations run automatically on container startup
- Feel free to experiment with rollbacks
- Can reset database entirely if needed

#### Production
- **Always backup** before running migrations
- **Test migrations** in staging environment first
- **Monitor logs** during migration application
- **Verify status** after deployment
- **Have rollback plan** ready

## Database Connection

The application connects to PostgreSQL using:
- **Host**: `db` (Docker service name)
- **Port**: `5432`
- **Database**: `olhourbanovault`
- **User**: `olhourbano`
- **Password**: Read from `/run/secrets/db_password`

Connection pooling and configuration are handled in `db.go`.

## Monitoring

### Check Migration Status
```bash
docker exec -w /olhourbano2 olhourbano2-backend-1 /usr/local/bin/app_olhourbano2 migrate:status
```

### View Migration History
```bash
docker exec olhourbano2-db-1 psql -U olhourbano olhourbanovault -c "
SELECT 
    version, 
    applied_at, 
    checksum
FROM schema_migrations 
ORDER BY version;
"
```

### Database Health Check
```bash
# Connection test
docker exec olhourbano2-db-1 psql -U olhourbano olhourbanovault -c "SELECT 1;"

# Table counts
docker exec olhourbano2-db-1 psql -U olhourbano olhourbanovault -c "
SELECT 
    'reports' as table_name, COUNT(*) as count FROM reports
UNION ALL
SELECT 
    'votes' as table_name, COUNT(*) as count FROM votes;
"
```

This migration system provides a robust, safe, and transparent way to manage database schema changes throughout the application lifecycle.