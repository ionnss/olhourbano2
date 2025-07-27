# Security Documentation

This document outlines the security measures implemented in the Olho Urbano project.

## ✅ Security Measures Implemented

### 1. Secret Management
- **Secrets stored in files**, not environment variables
- **Docker secrets** used in production environments
- **Local secret files** for development (`.secrets/` directory)
- **Never committed to version control** (`.gitignore` protection)

### 2. File Permissions
- Secret files restricted to owner-only access (`600` permissions)
- Prevents other users from reading sensitive data

### 3. Configuration Security
- **Separation of concerns**: URLs in env vars, secrets in files
- **No secrets in environment variables** - only file paths
- **Secure error handling** - no file paths or secret details in error messages
- **Safe config representation** - `String()` method excludes secrets

### 4. Network Security
- **HTTPS enforced** in production (Caddy configuration)
- **Blocked suspicious paths** (`.env`, `.git`, etc.)
- **Health checks** for service monitoring

### 5. Database Security
- **Password stored in secret file**, not hardcoded
- **SSL mode configurable** (currently disabled for local dev)
- **Connection pooling** and proper connection handling

## 🛡️ Security Best Practices

### Secret Files
```bash
# Correct permissions
chmod 600 secrets/*.txt

# Verify no secrets in git
git status --porcelain | grep -E "secrets/.*\.txt"
```

### Environment Variables
```bash
# Only file paths, never actual secrets
CPFHUB_API_KEY_FILE=/run/secrets/cpfhub_api_key
GOOGLE_MAPS_API_KEY_FILE=/run/secrets/google_maps_api_key
```

### Code Usage
```go
// ✅ CORRECT: Load config once
cfg, err := config.Load()
apiKey := cfg.CPFHubAPIKey

// ❌ WRONG: Don't log config directly
fmt.Printf("Config: %+v", cfg) // This would expose secrets

// ✅ CORRECT: Use safe string representation
fmt.Printf("Config: %s", cfg.String()) // No secrets exposed
```

## 🔍 Security Checklist

### Development Environment
- [ ] Secret files have `600` permissions
- [ ] `.env` file doesn't contain secrets
- [ ] `secrets/` directory in `.gitignore`
- [ ] No hardcoded passwords in code

### Production Environment
- [ ] Docker secrets properly mounted
- [ ] Environment variables contain only file paths
- [ ] HTTPS enabled and enforced
- [ ] Suspicious paths blocked by reverse proxy

### Code Review
- [ ] No secrets in log statements
- [ ] Error messages don't expose file paths
- [ ] Database connections use secret files
- [ ] API keys loaded from secret files

## 🚨 Security Vulnerabilities to Avoid

### ❌ Common Mistakes
```go
// DON'T: Log config with secrets
log.Printf("Config: %+v", config)

// DON'T: Put secrets in environment variables
os.Setenv("API_KEY", "secret-value")

// DON'T: Hardcode secrets
apiKey := "REDACTED..."

// DON'T: Expose file paths in errors
return fmt.Errorf("failed to read %s", secretFile)
```

### ✅ Secure Alternatives
```go
// DO: Use safe config representation
log.Printf("Config loaded: %s", config.String())

// DO: Read secrets from files
apiKey, err := readSecretFile(keyFile)

// DO: Use configuration system
cfg, err := config.Load()
apiKey := cfg.CPFHubAPIKey

// DO: Generic error messages
return fmt.Errorf("failed to load API key")
```

## 📊 Security Status: SECURE ✅

Your application implements industry-standard security practices:
- **Secrets properly isolated** from code and logs
- **Docker secrets integration** for production
- **Secure file permissions** preventing unauthorized access
- **No sensitive data exposure** in error messages or logs
- **HTTPS enforced** with automatic certificate management

## 🔄 Regular Security Maintenance

1. **Rotate secrets** periodically (API keys, passwords)
2. **Review file permissions** after system updates
3. **Monitor access logs** for suspicious activity
4. **Update dependencies** regularly for security patches
5. **Audit environment variables** to ensure no secrets leak 