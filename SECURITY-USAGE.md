# Security Audit Usage Guide

## Quick Start

### 1. Run Security Audit
```bash
# Make script executable (first time only)
chmod +x security-audit.sh

# Run complete security audit
./security-audit.sh
```

### 2. Interpret Results

The script will show results with colored indicators:
- ✅ **PASS**: Security test passed
- ❌ **FAIL**: Critical security issue found
- ⚠️ **WARN**: Potential concern, review recommended
- ℹ️ **INFO**: Informational message

### 3. Example Output
```
🔍 Security Audit for olhourbano2
=================================

🔧 1. Git Repository Security
✅ PASS: No secrets found in git commit history
✅ PASS: secrets/ directory is properly gitignored

🔧 2. Source Code Security  
✅ PASS: No hardcoded secrets found in source code
✅ PASS: String() method properly excludes sensitive fields

📊 Security Audit Summary
=========================
Total Tests: 15
Passed: 13
Failed: 0
Warnings: 2

😊 GOOD: No critical failures, but 2 warning(s) found.
```

## What Each Test Checks

| Test Category | What It Checks | Why It Matters |
|---------------|----------------|----------------|
| **Git Security** | Secrets in commit history, .gitignore setup | Prevents accidental secret exposure in version control |
| **Source Code** | Hardcoded passwords, safe string methods | Ensures no secrets in source code |
| **Config Files** | .env files, example files | Prevents config file secret leaks |
| **Docker Security** | Environment variables, secret permissions | Ensures container security |
| **Endpoints** | Debug/config endpoint exposure | Prevents runtime secret exposure |
| **Logs** | Sensitive data in application logs | Prevents log-based secret leaks |
| **File System** | File/directory permissions | Ensures proper access controls |

## Common Issues & Fixes

### ❌ "Secrets found in git commit messages"
**Cause**: Commit messages contain words like "password", "secret", or "key"
**Fix**: This may be a false positive if describing security improvements
**Action**: Review the commit messages - if they're just descriptions, this is safe

### ❌ "Environment file tracked by git"
**Cause**: `.env` file is committed to git
**Fix**: 
```bash
git rm --cached .env
echo ".env" >> .gitignore
```

### ❌ "Raw secrets in environment variables"
**Cause**: Secrets stored directly in env vars instead of files
**Fix**: Use `*_FILE` environment variables pointing to secret files

### ⚠️ "Docker containers not running"
**Cause**: Application not started
**Fix**: 
```bash
docker compose up -d
./security-audit.sh
```

## Automated Security Checks

### CI/CD Integration
Add to your CI pipeline:
```yaml
# .github/workflows/security.yml
- name: Security Audit
  run: |
    chmod +x security-audit.sh
    ./security-audit.sh
```

### Pre-commit Hook
```bash
# Add to .git/hooks/pre-commit
#!/bin/bash
./security-audit.sh
if [ $? -ne 0 ]; then
    echo "Security audit failed. Commit aborted."
    exit 1
fi
```

### Scheduled Checks
```bash
# Add to crontab for weekly checks
0 9 * * 1 cd /path/to/olhourbano2 && ./security-audit.sh > security-report.log 2>&1
```

## Security Best Practices

### ✅ Do's
- Use Docker secrets for sensitive data
- Store secrets in separate files with 600 permissions
- Use `*_FILE` environment variables
- Keep secrets out of version control
- Regular security audits
- Review code changes for secret exposure

### ❌ Don'ts
- Never commit secrets to git
- Don't use raw passwords in environment variables
- Don't expose debug endpoints in production
- Don't log sensitive information
- Don't share secret files via insecure channels

## Emergency Response

### If Secrets Are Exposed in Git
1. **Immediately rotate** all exposed credentials
2. **Remove secrets** from git history:
   ```bash
   git filter-branch --force --index-filter \
   'git rm --cached --ignore-unmatch secrets/exposed-file.txt' \
   --prune-empty --tag-name-filter cat -- --all
   ```
3. **Force push** to remote repository
4. **Notify team** of credential rotation

### If Secrets Are in Application Logs
1. **Stop application** immediately
2. **Rotate exposed credentials**
3. **Clear/secure log files**
4. **Fix logging code** to exclude secrets
5. **Restart with fixed code**

## Support

For security questions or issues:
1. Run `./security-audit.sh` first
2. Review `SECURITY-AUDIT.md` for detailed information
3. Check this usage guide for common solutions
4. If critical: immediately rotate affected credentials

---

**Remember**: Security is an ongoing process, not a one-time check. Run audits regularly! 