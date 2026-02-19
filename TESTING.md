# Guía de Testing para raracandy

## Testing Automático

Ejecuta el script de testing completo:

```bash
./test_workflow.sh
```

Este script verifica:
- ✓ Inspección de saves
- ✓ Agregar items (simple y múltiple)
- ✓ Modificar dinero
- ✓ Dry-run mode
- ✓ Backups automáticos
- ✓ Validación de checksums

## Testing con Emulador (Recomendado)

### Opción 1: mGBA (Recomendado)

1. **Instalar mGBA**:
   ```bash
   # macOS
   brew install mgba

   # O descarga desde: https://mgba.io/downloads.html
   ```

2. **Obtener una ROM de Pokemon Yellow** (debes poseer el juego original)

3. **Generar un save inicial**:
   ```bash
   # Abre mGBA y juega un poco para crear un save
   mgba pokemon_yellow.gbc
   # Guarda el juego (in-game save)
   # Cierra el emulador
   ```

4. **Modificar el save con raracandy**:
   ```bash
   # Encuentra el save (usualmente pokemon_yellow.sav en el mismo directorio)
   ./raracandy yellow inspect pokemon_yellow.sav

   # Agregar items
   ./raracandy yellow add-item pokemon_yellow.sav \
     --item rare_candy \
     --qty 99 \
     --out pokemon_yellow_modified.sav

   # Reemplazar el save original (BACKUP PRIMERO!)
   cp pokemon_yellow.sav pokemon_yellow.sav.original
   cp pokemon_yellow_modified.sav pokemon_yellow.sav
   ```

5. **Cargar en el emulador**:
   ```bash
   mgba pokemon_yellow.gbc
   ```

6. **Verificar en el juego**:
   - Presiona START
   - Ve a BAG
   - Verifica que tienes 99 Rare Candies ✓

### Opción 2: BGB (Windows/Wine)

Similar a mGBA pero específico para Game Boy.

### Opción 3: Emuladores móviles

Muchos emuladores de Game Boy en iOS/Android soportan importar saves.

## Testing con Hardware Real (GBxCart)

Si tienes un GBxCart RW:

### 1. Dump del save original

```bash
# Usa el software de GBxCart para extraer el save
# Resultado: pokemon_yellow_original.sav (32 KB)
```

### 2. Modificar con raracandy

```bash
./raracandy yellow inspect pokemon_yellow_original.sav

./raracandy yellow add-item pokemon_yellow_original.sav \
  --item rare_candy \
  --qty 50 \
  --out pokemon_yellow_mod.sav

./raracandy yellow set-money pokemon_yellow_mod.sav \
  --amount 999999 \
  --out pokemon_yellow_final.sav
```

### 3. Verificar antes de inyectar

```bash
./raracandy yellow inspect pokemon_yellow_final.sav

# Output esperado:
# Checksum: ✓ Valid (MUY IMPORTANTE!)
# Money: ¥999,999
# Bag: Rare Candy x50
```

### 4. Inyectar al cartucho

```bash
# Usa el software de GBxCart para escribir pokemon_yellow_final.sav
# al cartucho original
```

### 5. Probar en Game Boy real

- Enciende el Game Boy
- Carga el save
- Si el checksum es correcto, el juego NO dirá "The file data is destroyed!"
- Verifica tu bag e items

## Tests Unitarios

```bash
# Ejecutar todos los tests de Go
go test ./...

# Tests específicos
go test ./internal/gen1/save/
go test ./internal/gen1/money/
go test ./internal/gen1/items/

# Con coverage
go test -cover ./...
```

## Verificaciones Importantes

### Antes de usar en hardware real:

1. **Checksum válido**:
   ```bash
   ./raracandy yellow inspect modified.sav | grep "Status"
   # Debe mostrar: Status: ✓ Valid
   ```

2. **Tamaño correcto**:
   ```bash
   ls -lh modified.sav
   # Debe ser exactamente 32 KB (32768 bytes)
   ```

3. **Backup del original**:
   ```bash
   # SIEMPRE haz backup antes de inyectar a hardware real
   cp original.sav original.sav.backup.$(date +%Y%m%d_%H%M%S)
   ```

## Items Soportados para Testing

```bash
# Items más útiles:
- rare_candy (Rare Candy)
- master_ball (Master Ball)
- max_potion (Max Potion)
- max_revive (Max Revive)
- full_restore (Full Restore)

# Lista completa:
./raracandy yellow add-item --help
# O ver: internal/gen1/items/itemdb.go
```

## Troubleshooting

### "The file data is destroyed!" en el juego

- El checksum es incorrecto
- Verifica: `./raracandy yellow inspect save.sav`
- Solución: Usa el backup (.bak) y vuelve a intentar

### Save no carga en emulador

- Verifica el tamaño: debe ser exactamente 32 KB
- Verifica que el emulador busca el save en el directorio correcto

### Items no aparecen en el juego

- Verifica los offsets (pueden ser diferentes para versiones JP/EU)
- Esta herramienta está probada para Pokemon Yellow NA (North America)

## Workflow Completo de Testing

```bash
# 1. Build
go build -o raracandy ./cmd/raracandy

# 2. Tests automáticos
./test_workflow.sh

# 3. Tests unitarios
go test ./...

# 4. Test con emulador (mGBA)
# - Crear save inicial en mGBA
# - Modificar con raracandy
# - Verificar en mGBA

# 5. (Opcional) Test en hardware real
# - Dump save con GBxCart
# - Modificar con raracandy
# - Inyectar con GBxCart
# - Probar en Game Boy
```

## ¿Qué NO testear?

- ❌ No uses saves de otras versiones (Red, Blue, Crystal)
- ❌ No modifiques ROMs (solo saves)
- ❌ No uses en competiciones oficiales de Pokemon
- ❌ No intentes valores fuera de rango (>99 items, >999999 dinero)
