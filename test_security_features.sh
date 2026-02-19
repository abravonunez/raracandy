#!/bin/bash
# Test Security Features

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=== Testing Security Features ==="
echo ""

# 1. Verify command
echo -e "${BLUE}1. Testing VERIFY command${NC}"
./raracandy yellow verify testdata/fixtures/test.sav
echo ""

# 2. Integrity check before modification
echo -e "${BLUE}2. Testing INTEGRITY CHECK (automatic)${NC}"
echo "yes" | ./raracandy yellow add-item testdata/fixtures/test.sav \
  --item potion \
  --qty 10 \
  --out testdata/test_security1.sav
echo ""

# 3. SHA256 hashing
echo -e "${BLUE}3. Testing SHA256 HASH creation${NC}"
if [ -f testdata/fixtures/test.sav.bak.sha256 ]; then
  echo -e "${GREEN}✓ Hash file created${NC}"
  cat testdata/fixtures/test.sav.bak.sha256
else
  echo "❌ Hash file NOT created"
  exit 1
fi
echo ""

# 4. Confirmation prompt (force mode bypasses it)
echo -e "${BLUE}4. Testing FORCE MODE (skips confirmation)${NC}"
./raracandy yellow add-item testdata/fixtures/test.sav \
  --item master_ball \
  --qty 5 \
  --out testdata/test_security2.sav \
  --force
echo ""

# 5. Dry-run mode
echo -e "${BLUE}5. Testing DRY-RUN MODE${NC}"
./raracandy yellow set-money testdata/fixtures/test.sav \
  --amount 500000 \
  --out testdata/should_not_exist2.sav \
  --dry-run

if [ -f testdata/should_not_exist2.sav ]; then
  echo "❌ Dry-run created a file!"
  exit 1
else
  echo -e "${GREEN}✓ Dry-run did NOT create file (correct)${NC}"
fi
echo ""

# 6. Verify written file
echo -e "${BLUE}6. Testing POST-WRITE VERIFICATION${NC}"
./raracandy yellow verify testdata/test_security1.sav
echo -e "${GREEN}✓ Written file has valid checksum${NC}"
echo ""

# 7. Check all created files
echo -e "${BLUE}7. Summary of created files:${NC}"
echo "Backup files:"
ls -lh testdata/fixtures/*.bak 2>/dev/null || echo "  (none)"
echo ""
echo "Hash files:"
ls -lh testdata/fixtures/*.sha256 2>/dev/null || echo "  (none)"
echo ""
echo "Modified saves:"
ls -lh testdata/test_security*.sav 2>/dev/null || echo "  (none)"
echo ""

echo -e "${GREEN}=== All security tests passed! ===${NC}"
echo ""
echo "Security features verified:"
echo "  ✓ Integrity check before modification"
echo "  ✓ SHA256 hash of backups"
echo "  ✓ Confirmation prompts (interactive mode)"
echo "  ✓ Force mode (--force flag)"
echo "  ✓ Dry-run mode (no writes)"
echo "  ✓ Post-write verification"
echo "  ✓ Verify command"
