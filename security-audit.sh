#!/bin/bash

# Security Audit Script for olhourbano2
# This script performs comprehensive security checks to ensure
# no sensitive information is exposed or leaked

# Remove set -e to prevent early exit on grep commands that find no matches
# set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Emojis for better visual feedback
PASS="✅"
FAIL="❌"
WARNING="⚠️"
INFO="ℹ️"

# Global variables
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
WARNINGS=0

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "PASS")
            echo -e "${GREEN}${PASS} PASS${NC}: $message"
            ((PASSED_TESTS++))
            ;;
        "FAIL")
            echo -e "${RED}${FAIL} FAIL${NC}: $message"
            ((FAILED_TESTS++))
            ;;
        "WARN")
            echo -e "${YELLOW}${WARNING} WARN${NC}: $message"
            ((WARNINGS++))
            ;;
        "INFO")
            echo -e "${BLUE}${INFO} INFO${NC}: $message"
            ;;
    esac
    ((TOTAL_TESTS++))
}

# Function to check if Docker is running
check_docker() {
    if ! docker compose ps >/dev/null 2>&1; then
        print_status "WARN" "Docker containers are not running. Some tests will be skipped."
        return 1
    fi
    return 0
}

# Header
echo "🔍 Security Audit for olhourbano2"
echo "================================="
echo "Timestamp: $(date)"
echo ""

# Test 1: Git Repository Security
echo "🔧 1. Git Repository Security"
echo "------------------------------"

# Check if git repo exists
if [ ! -d ".git" ]; then
    print_status "WARN" "Not a git repository - skipping git history checks"
else
    # Check git history for secrets (exclude common security-related terms)
    SUSPICIOUS_COMMITS=$(git log --all --full-history --oneline | grep -i "password\|secret\|key" | grep -v -i "use.*secret\|add.*secret\|implement.*secret\|docker.*secret\|secret.*management\|secret.*config\|api.*key.*config\|session.*key.*config" || true)
    if [ -n "$SUSPICIOUS_COMMITS" ]; then
        print_status "FAIL" "Potential secrets found in git commit messages"
        echo "$SUSPICIOUS_COMMITS" | head -3
    else
        print_status "PASS" "No suspicious secrets found in git commit history"
    fi

    # Check for actual secret values in git-tracked files (look for assignments)
    SECRET_ASSIGNMENTS=$(git ls-files | xargs grep -E "(password|secret|key)\s*=\s*['\"][^'\"]{8,}" 2>/dev/null || true)
    if [ -n "$SECRET_ASSIGNMENTS" ]; then
        print_status "FAIL" "Potential secret assignments found in tracked files"
        echo "$SECRET_ASSIGNMENTS" | head -3
    else
        print_status "PASS" "No secret assignments found in git-tracked files"
    fi
fi

# Check .gitignore for secrets directory
if grep -q "secrets/" .gitignore 2>/dev/null; then
    print_status "PASS" "secrets/ directory is properly gitignored"
else
    print_status "FAIL" "secrets/ directory not found in .gitignore"
fi

echo ""

# Test 2: Source Code Security
echo "🔧 2. Source Code Security"
echo "----------------------------"

# Check for hardcoded secret values in Go files (look for actual assignments)
HARDCODED_SECRETS=$(grep -r -E "(password|secret|key)\s*[:=]\s*['\"][^'\"]{8,}" --include="*.go" . | grep -v "PASSWORD_FILE\|SECRET_FILE\|KEY_FILE" || true)
if [ -n "$HARDCODED_SECRETS" ]; then
    print_status "FAIL" "Potential hardcoded secret values found in Go files"
    echo "$HARDCODED_SECRETS" | head -5
else
    print_status "PASS" "No hardcoded secret values found in source code"
fi

# Check String() method excludes sensitive fields
if grep -A 10 "func.*String()" config/config.go 2>/dev/null | grep -q "DBPassword\|SMTPPassword\|SessionKey\|APIKey"; then
    print_status "FAIL" "String() method may expose sensitive fields"
else
    print_status "PASS" "String() method properly excludes sensitive fields"
fi

echo ""

# Test 3: Configuration File Security
echo "🔧 3. Configuration File Security"
echo "----------------------------------"

# Check .env files for real secrets
ENV_FILES=$(find . -name "*.env*" 2>/dev/null || true)
if [ -n "$ENV_FILES" ]; then
    for file in $ENV_FILES; do
        if [[ "$file" == *".example"* ]]; then
            # Check if example file contains real-looking secrets
            if grep -E "password.*[a-zA-Z0-9]{8,}|key.*[a-zA-Z0-9]{20,}" "$file" >/dev/null 2>&1; then
                print_status "WARN" "Example file $file may contain real secrets"
            else
                print_status "PASS" "Example file $file contains only placeholders"
            fi
        else
            # Real env file - should not be committed
            if [ -f "$file" ] && git ls-files --error-unmatch "$file" >/dev/null 2>&1; then
                print_status "FAIL" "Environment file $file is tracked by git"
            else
                print_status "PASS" "Environment file $file is not tracked by git"
            fi
        fi
    done
else
    print_status "INFO" "No .env files found"
fi

echo ""

# Test 4: Docker Container Security (only if containers are running)
echo "🔧 4. Docker Container Security"
echo "--------------------------------"

if check_docker; then
    # Check environment variables don't contain raw secrets
    RAW_SECRETS=$(docker compose exec -T backend env | grep -E "(PASSWORD|SECRET|KEY)" | grep -v "_FILE=" || true)
    if [ -n "$RAW_SECRETS" ]; then
        print_status "FAIL" "Raw secrets found in environment variables"
        echo "$RAW_SECRETS"
    else
        print_status "PASS" "Only file paths found in environment variables"
    fi

    # Check secret file permissions
    SECRET_PERMS=$(docker compose exec -T backend ls -la /run/secrets/ 2>/dev/null | grep "^-rw-------" | wc -l || echo "0")
    TOTAL_SECRETS=$(docker compose exec -T backend ls -la /run/secrets/ 2>/dev/null | grep "^-" | wc -l || echo "0")
    
    if [ "$SECRET_PERMS" -eq "$TOTAL_SECRETS" ] && [ "$TOTAL_SECRETS" -gt 0 ]; then
        print_status "PASS" "All secret files have correct permissions (600)"
    elif [ "$TOTAL_SECRETS" -eq 0 ]; then
        print_status "WARN" "No secret files found in container"
    else
        print_status "FAIL" "Some secret files have incorrect permissions"
    fi

    # Check for secrets in process environment
    PROC_ENV_SECRETS=$(docker compose exec -T backend cat /proc/1/environ 2>/dev/null | tr '\0' '\n' | grep -E "(PASSWORD|SECRET|KEY)" | grep -v "_FILE=" | wc -l || echo "0")
    if [ "$PROC_ENV_SECRETS" -eq 0 ]; then
        print_status "PASS" "No raw secrets in process environment"
    else
        print_status "FAIL" "$PROC_ENV_SECRETS raw secrets found in process environment"
    fi
else
    print_status "INFO" "Skipping Docker container tests (containers not running)"
fi

echo ""

# Test 5: Application Endpoint Security
echo "🔧 5. Application Endpoint Security"
echo "------------------------------------"

if check_docker; then
    # Test for exposed debug/config endpoints
    ENDPOINTS=("/config" "/debug" "/health" "/env" "/status" "/info")
    
    for endpoint in "${ENDPOINTS[@]}"; do
        RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8081$endpoint" 2>/dev/null || echo "000")
        if [ "$RESPONSE" = "200" ]; then
            # Check if response contains sensitive data
            RESPONSE_BODY=$(curl -s "http://localhost:8081$endpoint" 2>/dev/null || echo "")
            if echo "$RESPONSE_BODY" | grep -qi "password\|secret\|key"; then
                print_status "FAIL" "Endpoint $endpoint exposes sensitive information"
            else
                print_status "WARN" "Endpoint $endpoint is accessible but appears safe"
            fi
        elif [ "$RESPONSE" = "404" ]; then
            print_status "PASS" "Endpoint $endpoint properly returns 404"
        elif [ "$RESPONSE" = "000" ]; then
            print_status "INFO" "Could not connect to application (may not be running)"
            break
        else
            print_status "INFO" "Endpoint $endpoint returns HTTP $RESPONSE"
        fi
    done
else
    print_status "INFO" "Skipping endpoint tests (containers not running)"
fi

echo ""

# Test 6: Application Log Security
echo "🔧 6. Application Log Security"
echo "-------------------------------"

if check_docker; then
    # Check application logs for sensitive data
    LOG_SECRETS=$(docker compose logs backend 2>&1 | grep -E "(password|secret|key)" -i | grep -v "PASSWORD_FILE\|SECRET_FILE\|KEY_FILE" || true)
    if [ -n "$LOG_SECRETS" ]; then
        print_status "FAIL" "Sensitive information found in application logs"
        echo "$LOG_SECRETS" | head -3
    else
        print_status "PASS" "No sensitive information found in application logs"
    fi
else
    print_status "INFO" "Skipping log analysis (containers not running)"
fi

echo ""

# Test 7: File System Security
echo "🔧 7. File System Security"
echo "---------------------------"

# Check secrets directory permissions
if [ -d "secrets" ]; then
    SECRETS_PERM=$(stat -f "%p" secrets 2>/dev/null || stat -c "%a" secrets 2>/dev/null || echo "unknown")
    if [[ "$SECRETS_PERM" == *"700"* ]] || [[ "$SECRETS_PERM" == "700" ]]; then
        print_status "PASS" "secrets/ directory has correct permissions"
    else
        print_status "WARN" "secrets/ directory permissions: $SECRETS_PERM (recommended: 700)"
    fi

    # Check individual secret files
    SECRET_FILES_COUNT=$(find secrets -type f -name "*.txt" | wc -l)
    SECURE_FILES_COUNT=$(find secrets -type f -name "*.txt" -perm 600 | wc -l)
    
    if [ "$SECRET_FILES_COUNT" -eq "$SECURE_FILES_COUNT" ] && [ "$SECRET_FILES_COUNT" -gt 0 ]; then
        print_status "PASS" "All secret files have correct permissions (600)"
    elif [ "$SECRET_FILES_COUNT" -eq 0 ]; then
        print_status "INFO" "No secret files found"
    else
        print_status "WARN" "Some secret files may have incorrect permissions"
    fi
else
    print_status "WARN" "secrets/ directory not found"
fi

echo ""

# Summary
echo "📊 Security Audit Summary"
echo "========================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
echo -e "${RED}Failed: $FAILED_TESTS${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo ""

# Overall assessment
if [ "$FAILED_TESTS" -eq 0 ]; then
    if [ "$WARNINGS" -eq 0 ]; then
        echo -e "${GREEN}🎉 EXCELLENT: All security tests passed!${NC}"
        exit 0
    else
        echo -e "${YELLOW}😊 GOOD: No critical failures, but $WARNINGS warning(s) found.${NC}"
        exit 0
    fi
else
    echo -e "${RED}⚠️  ATTENTION NEEDED: $FAILED_TESTS critical security issue(s) found!${NC}"
    echo ""
    echo "Please review the failed tests and fix the issues before deploying to production."
    exit 1
fi 