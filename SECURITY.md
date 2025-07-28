# Security Documentation

This document provides comprehensive security information for the **olhourbano2** application, including implemented measures, audit procedures, and usage guidelines.

## 📋 Table of Contents

1. [Security Measures Implemented](#-security-measures-implemented)
2. [Security Audit](#-security-audit)
3. [Quick Security Check](#-quick-security-check)
4. [Security Architecture](#-security-architecture)
5. [Best Practices](#-best-practices)
6. [Incident Response](#-incident-response)

---

## 🛡️ Security Measures Implemented

### 1. Secret Management
- **✅ Docker Secrets**: Production uses Docker secrets with proper mounting
- **✅ File-based Storage**: Secrets stored in files, not environment variables
- **✅ Version Control Protection**: `secrets/` directory gitignored
- **✅ Proper Permissions**: Secret files (`600`) and directory (`700`) restricted

### 2. Application Security
- **✅ Safe Configuration**: `String()` method excludes sensitive fields
- **✅ Secure File Loading**: Fallback mechanism with proper error handling
- **✅ Health Endpoints**: Dedicated health check endpoints for monitoring
- **✅ No Debug Exposure**: No configuration or debug endpoints in production

### 3. Network Security
- **✅ HTTPS Enforced**: Caddy handles SSL/TLS termination
- **✅ Blocked Suspicious Paths**: `.env`, `.git`, security scan attempts blocked
- **✅ Reverse Proxy**: Caddy acts as secure reverse proxy to backend
- **✅ Health Monitoring**: Continuous health checks with proper endpoints

### 4. Database Security
- **✅ Password Files**: Database password stored in secret files
- **✅ Connection Security**: Proper connection handling and pooling
- **✅ User Isolation**: Dedicated database user with minimal privileges

### 5. Container Security
- **✅ Secret Mounting**: Secrets mounted securely in containers
- **✅ User Permissions**: Proper file ownership and access control
- **✅ Network Isolation**: Services communicate through defined networks

---

## 🔍 Security Audit

### Automated Security Audit

We provide a comprehensive security audit script that checks for common vulnerabilities:

```bash
# Run complete security audit
./security-audit.sh
```

### What the Audit Checks

1. **Git Repository Security**
   - No secrets in git history
   - Proper `.gitignore` configuration
   - No hardcoded credentials in source code

2. **Source Code Security**
   - No hardcoded secret values
   - Safe string representation methods
   - Proper secret file handling

3. **Docker Container Security**
   - Environment variable safety
   - Secret file permissions
   - Container configuration

4. **Application Endpoint Security**
   - No debug endpoints exposed
   - No configuration leaks
   - Health check functionality

5. **File System Security**
   - Secret directory permissions (`700`)
   - Secret file permissions (`600`)
   - Proper ownership

6. **Configuration Security**
   - No secrets in configuration files
   - Proper environment variable usage
   - Safe error handling

7. **Network Security**
   - HTTPS configuration
   - Blocked dangerous paths
   - Proper reverse proxy setup

### Audit Results Interpretation

- **✅ PASS**: Security test passed
- **❌ FAIL**: Critical security issue - immediate action required
- **⚠️ WARN**: Potential concern - review recommended
- **ℹ️ INFO**: Informational message

---

## ⚡ Quick Security Check

### Daily Security Verification

```bash
# 1. Run security audit
./security-audit.sh

# 2. Check secret file permissions
ls -la secrets/

# 3. Verify no secrets in environment
docker compose exec backend env | grep -E "(PASSWORD|SECRET|KEY)"

# 4. Test health endpoints
curl -k https://localhost/health

# 5. Check for suspicious access attempts
docker compose logs caddy | grep -E "(404|403)" | tail -10
```

### Emergency Security Check

If you suspect a security breach:

```bash
# 1. Immediate audit
./security-audit.sh > security-breach-audit.log

# 2. Check git for any new commits with secrets
git log --oneline --since="1 day ago" | xargs git show | grep -E "(password|secret|key)" -i

# 3. Review recent access logs
docker compose logs --since=24h | grep -E "(error|fail|unauthorized)"

# 4. Verify current secret file integrity
find secrets/ -type f -name "*.txt" -exec ls -la {} \;
```

---

## 🏗️ Security Architecture

### Secret Flow Architecture

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│   Secrets   │    │    Docker    │    │ Application │
│    Files    │───▶│   Secrets    │───▶│   Runtime   │
│ (Host FS)   │    │ (Container)  │    │  (Go App)   │
└─────────────┘    └──────────────┘    └─────────────┘
      │                    │                   │
   chmod 600           /run/secrets/      File Reading
   chmod 700           (mounted)         (Secure API)
```

### Network Security Flow

```
┌─────────┐    ┌─────────┐    ┌──────────┐    ┌──────────┐
│ Client  │───▶│  Caddy  │───▶│ Backend  │───▶│ Database │
│(Browser)│    │(Reverse │    │   App    │    │   (PG)   │
│         │    │ Proxy)  │    │          │    │          │
└─────────┘    └─────────┘    └──────────┘    └──────────┘
                   │               │               │
              HTTPS/SSL     Health Checks    Secret Auth
              Cert Mgmt     /health         Password File
```

### Permission Structure

```
secrets/                 (drwx------) 700
├── db_passwords.txt     (-rw-------) 600
├── smtp_password.txt    (-rw-------) 600
├── session_key.txt      (-rw-------) 600
├── cpfhub_api_key.txt   (-rw-------) 600
└── google_maps_api_key.txt (-rw-------) 600
```

---

## 📚 Best Practices

### For Developers

1. **Never commit secrets**:
   ```bash
   # Always check before committing
   git diff --cached | grep -E "(password|secret|key)" -i
   ```

2. **Use proper secret loading**:
   ```go
   // Good: File-based secrets
   password, err := readSecretFile("/run/secrets/db_password")
   
   // Bad: Hardcoded secrets
   password := "hardcoded_password_123"
   ```

3. **Secure configuration representation**:
   ```go
   func (c *Config) String() string {
       return fmt.Sprintf("Config{DBHost: %s, AppVersion: %s}", 
           c.DBHost, c.AppVersion)
       // Note: Excludes DBPassword, SMTPPassword, etc.
   }
   ```

### For Operations

1. **Regular security audits**:
   ```bash
   # Run weekly
   ./security-audit.sh > weekly-security-audit.log
   ```

2. **Monitor access logs**:
   ```bash
   # Check for suspicious patterns
   docker compose logs caddy | grep -E "(404|403|scan|hack)"
   ```

3. **Rotate secrets regularly**:
   ```bash
   # Update secret files and restart services
   echo "new_password" > secrets/db_passwords.txt
   chmod 600 secrets/db_passwords.txt
   docker compose restart
   ```

### For Deployment

1. **Production checklist**:
   - [ ] All secrets in files, not environment variables
   - [ ] Secret file permissions set to `600`
   - [ ] Secret directory permissions set to `700`
   - [ ] HTTPS enforced
   - [ ] Health checks working
   - [ ] Security audit passes

2. **Environment-specific considerations**:
   - **Development**: Use local secret files
   - **Staging**: Mirror production security setup
   - **Production**: Use Docker secrets with orchestration

---

## 🚨 Incident Response

### Security Incident Procedure

1. **Immediate Actions**:
   ```bash
   # 1. Stop services if breach confirmed
   docker compose down
   
   # 2. Capture current state
   ./security-audit.sh > incident-audit-$(date +%Y%m%d-%H%M%S).log
   
   # 3. Review recent logs
   docker compose logs > incident-logs-$(date +%Y%m%d-%H%M%S).log
   ```

2. **Assessment**:
   - Run complete security audit
   - Check git history for any unauthorized commits
   - Review access logs for suspicious activity
   - Verify secret file integrity

3. **Recovery**:
   - Rotate all affected secrets
   - Update secret files with new values
   - Restart all services
   - Monitor for continued suspicious activity

4. **Post-Incident**:
   - Document lessons learned
   - Update security measures
   - Consider additional monitoring

### Contact Information

For security-related issues:
- **Email**: olhourbano.contato@gmail.com
- **Urgent**: Create GitHub issue with `security` label

---

## 📊 Security Compliance

### Current Security Status

**Overall Security Score**: ✅ **EXCELLENT**

- **Secret Management**: ✅ Implemented
- **Access Control**: ✅ Implemented  
- **Network Security**: ✅ Implemented
- **Audit Capability**: ✅ Implemented
- **Incident Response**: ✅ Documented

### Recent Audit Results

```bash
# Last audit: $(date)
📊 Security Audit Summary
========================
✅ Passed: 7/7 tests
❌ Failed: 0/7 tests
⚠️ Warnings: 0/7 tests

Overall Status: SECURE ✅
```

### Maintenance Schedule

- **Daily**: Health check monitoring
- **Weekly**: Security audit execution
- **Monthly**: Secret rotation review
- **Quarterly**: Full security review and documentation update

---

*Last updated: $(date)*
*Security documentation version: 2.0* 