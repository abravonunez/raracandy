# Workflow SEGURO para tu Save Original de Pokemon Yellow

## ⚠️ LEE ESTO ANTES DE MODIFICAR TU SAVE ORIGINAL

**IMPORTANTE:** Este workflow asume que tienes un GBxCart RW u otro dumper de cartuchos.

---

## Pre-requisitos

- ✅ GBxCart RW (o GB Operator, Joey Jr, etc.)
- ✅ Software del dumper instalado
- ✅ Cartucho original de Pokemon Yellow (NA/US)
- ✅ raracandy compilado y testeado en emulador primero

---

## Paso 1: Backup MÚLTIPLE de tu Save Original

```bash
# 1. Usa GBxCart para extraer el save del cartucho
# Resultado: pokemon_yellow_original.sav

# 2. Crea MÚLTIPLES backups con timestamp
cp pokemon_yellow_original.sav pokemon_yellow_original.sav.backup_$(date +%Y%m%d_%H%M%S)

# 3. Crea un backup en otra ubicación (disco externo, cloud, etc.)
cp pokemon_yellow_original.sav ~/Dropbox/pokemon_backups/pokemon_yellow_$(date +%Y%m%d).sav

# 4. Verifica que los backups se crearon
ls -lh pokemon_yellow_original.sav*
```

**REGLA DE ORO:** Nunca modifiques el archivo original directamente.

---

## Paso 2: Verificar el Save Original

```bash
# Inspeccionar el save antes de modificar
./raracandy yellow inspect pokemon_yellow_original.sav

# Output esperado:
# - Size: 32 KB ✓
# - Checksum: Status: ✓ Valid
# - Money: (tu dinero actual)
# - Bag: (tus items actuales)
```

**Si el checksum NO es válido:** Tu save original podría estar corrupto. DETENTE y consulta.

---

## Paso 3: Modificación (Dry-run primero)

```bash
# SIEMPRE haz dry-run primero
./raracandy yellow add-item pokemon_yellow_original.sav \
  --item rare_candy \
  --qty 50 \
  --out pokemon_yellow_mod.sav \
  --dry-run

# Lee el output cuidadosamente:
# - Verifica los cambios propuestos
# - Confirma que es lo que quieres

# Si todo se ve bien, ejecuta SIN dry-run
./raracandy yellow add-item pokemon_yellow_original.sav \
  --item rare_candy \
  --qty 50 \
  --out pokemon_yellow_mod.sav
```

---

## Paso 4: Verificación del Save Modificado

```bash
# Inspeccionar el save modificado
./raracandy yellow inspect pokemon_yellow_mod.sav

# VERIFICA ESTOS PUNTOS CRÍTICOS:
# 1. Size: 32 KB (debe ser exacto)
# 2. Checksum: ✓ Valid (MUY IMPORTANTE!)
# 3. Money: (debe ser correcto)
# 4. Bag: (debe incluir tus nuevos items)
```

**Si el checksum NO es válido:** NO inyectes este save. Hay un problema.

---

## Paso 5: Test en Emulador ANTES de Hardware Real

```bash
# CRÍTICO: Prueba el save modificado en mGBA primero
cp pokemon_yellow_mod.sav ~/roms/gameboy/pokemon_yellow.sav

mgba ~/roms/gameboy/pokemon_yellow.gbc

# En el emulador:
# 1. ¿El juego carga sin errores? ✓
# 2. ¿El menú START funciona? ✓
# 3. ¿Los items están en tu bag? ✓
# 4. ¿Puedes guardar el juego? ✓
# 5. ¿El save persiste después de cerrar/reabrir? ✓
```

**Si ALGO falla en el emulador:** NO lo inyectes al cartucho real.

---

## Paso 6: Inyección al Cartucho (ÚLTIMO PASO)

```bash
# Solo si TODO funcionó en el emulador:

# 1. Abre el software de GBxCart
# 2. Selecciona "Write Save to Cartridge"
# 3. Selecciona: pokemon_yellow_mod.sav
# 4. Confirma la escritura
# 5. Espera a que termine

# El software debería verificar la escritura
```

---

## Paso 7: Verificación en Game Boy Real

```bash
# 1. Inserta el cartucho en tu Game Boy
# 2. Enciende
# 3. ¿Sale "The file data is destroyed!"?
#    - NO → ✓ Éxito
#    - SÍ → Problema con el checksum

# 4. Si cargó bien:
#    - Verifica tus items (START > BAG)
#    - Verifica tu dinero
#    - Guarda el juego (in-game)
#    - Apaga y enciende de nuevo
#    - Verifica que el save persiste
```

---

## Plan de Recuperación si Algo Sale Mal

### Si el Game Boy dice "The file data is destroyed!"

```bash
# 1. NO ENTRES EN PÁNICO
# 2. El cartucho NO está dañado físicamente
# 3. Solo necesitas restaurar el backup

# Pasos de recuperación:
# 1. Usa GBxCart para escribir el backup original
cp pokemon_yellow_original.sav.backup_[TIMESTAMP] pokemon_yellow_recovery.sav

# 2. Inyecta pokemon_yellow_recovery.sav al cartucho
# 3. Tu save original debería volver

# 4. Investiga qué salió mal antes de intentar de nuevo
```

### Si perdiste el backup

**Esto NO debería pasar si seguiste el Paso 1 correctamente.**

Pero si pasa:
- El cartucho físico NO está dañado
- La batería interna del cartucho podría tener tu save original
- Contacta a comunidades de Pokemon speedrunning/modding para ayuda

---

## Checklist Pre-Inyección

Antes de escribir al cartucho, verifica:

- [ ] Tienes al menos 2 backups del save original
- [ ] El save modificado tiene checksum válido
- [ ] El save modificado funciona en emulador (mGBA)
- [ ] El tamaño es exactamente 32 KB (32768 bytes)
- [ ] Tu cartucho es Pokemon Yellow NA/US (no JP/EU)
- [ ] La batería del cartucho está bien (puede guardar)
- [ ] Has leído toda esta guía

**Si marcaste TODO:** Adelante, debería ser seguro.

**Si falta ALGO:** Detente y revisa.

---

## Consejos Adicionales

### Versión del Juego

```bash
# raracandy está probado para:
# - Pokemon Yellow (North America / US)

# Versiones NO probadas (pueden tener offsets diferentes):
# - Pokemon Yellow (Japan)
# - Pokemon Yellow (Europe)
# - Pokemon Red/Blue (cualquier región)

# Verifica tu versión en el cartucho
```

### Modificaciones Conservadoras

Para tu primer uso, **empieza pequeño**:

```bash
# En vez de:
--qty 99  # Máximo, obvio

# Prueba:
--qty 10  # Más sutil, menos obvio si falla

# En vez de:
--amount 999999  # Máximo dinero

# Prueba:
--amount 50000   # Cantidad razonable
```

### Qué NO hacer

- ❌ NO uses raracandy en un save de competición oficial
- ❌ NO modifiques directamente el archivo original
- ❌ NO saltes el paso del emulador
- ❌ NO inyectes un save con checksum inválido
- ❌ NO uses en versiones JP/EU sin verificar offsets
- ❌ NO modifiques datos fuera del bag/money (por ahora)

---

## Nivel de Confianza

**Basado en testing:**
- ✅ Checksums: Probados y funcionan
- ✅ Items: Offset verificado contra fuentes técnicas
- ✅ Money: BCD encoding probado
- ✅ Emulador: Funciona en mGBA

**Riesgos residuales:**
- ⚠️ Versiones no-NA podrían tener offsets diferentes
- ⚠️ Saves corruptos originalmente no se pueden "arreglar"
- ⚠️ Errores en el hardware dumper (GBxCart) son posibles

**Confianza general:** 95% si sigues este workflow.

---

## ¿Preguntas antes de empezar?

Antes de inyectar a tu cartucho original:

1. ¿Has probado en emulador primero?
2. ¿Tienes backups múltiples?
3. ¿El checksum es válido?
4. ¿Tu cartucho es Pokemon Yellow NA/US?

**Si respondiste SÍ a todo:** Adelante con confianza.
**Si respondiste NO a algo:** Revisa ese paso primero.
