#!/bin/bash
# Script de testing completo para raracandy

set -e

echo "=== Test Workflow para raracandy ==="
echo ""

# Colores para output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 1. Inspeccionar save original
echo -e "${BLUE}1. Inspeccionando save original...${NC}"
./raracandy yellow inspect testdata/fixtures/test.sav
echo ""

# 2. Test: Agregar Rare Candies
echo -e "${BLUE}2. Agregando 99 Rare Candies...${NC}"
./raracandy yellow add-item testdata/fixtures/test.sav \
  --item rare_candy \
  --qty 99 \
  --out testdata/test_candies.sav
echo ""

# 3. Verificar cambios
echo -e "${BLUE}3. Verificando que se agregaron las Rare Candies...${NC}"
./raracandy yellow inspect testdata/test_candies.sav | grep -A 5 "Bag"
echo ""

# 4. Test: Agregar más items al save modificado
echo -e "${BLUE}4. Agregando Master Balls al save ya modificado...${NC}"
./raracandy yellow add-item testdata/test_candies.sav \
  --item master_ball \
  --qty 50 \
  --out testdata/test_multi.sav
echo ""

# 5. Test: Modificar dinero
echo -e "${BLUE}5. Estableciendo dinero a 999,999...${NC}"
./raracandy yellow set-money testdata/test_multi.sav \
  --amount 999999 \
  --out testdata/test_complete.sav
echo ""

# 6. Inspección final
echo -e "${BLUE}6. Save final con todos los cambios:${NC}"
./raracandy yellow inspect testdata/test_complete.sav
echo ""

# 7. Test dry-run
echo -e "${BLUE}7. Test de dry-run (no debe crear archivo):${NC}"
./raracandy yellow add-item testdata/fixtures/test.sav \
  --item ultra_ball \
  --qty 99 \
  --out testdata/should_not_exist.sav \
  --dry-run
echo ""

if [ -f testdata/should_not_exist.sav ]; then
  echo "❌ ERROR: Dry-run creó un archivo!"
  exit 1
else
  echo -e "${GREEN}✓ Dry-run funcionó correctamente (no creó archivo)${NC}"
fi
echo ""

# 8. Verificar backups
echo -e "${BLUE}8. Verificando que se crearon backups:${NC}"
ls -lh testdata/*.bak 2>/dev/null | wc -l | xargs echo "Archivos .bak creados:"
echo ""

echo -e "${GREEN}=== ✓ Todos los tests pasaron ===${NC}"
