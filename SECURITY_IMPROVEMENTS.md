# Mejoras de Seguridad para raracandy

## Nivel Actual: ~95% confianza
## Objetivo: >99% confianza

---

## 1. Detección Automática de Versión del Juego ⭐ CRÍTICO

### Problema
Los offsets son diferentes entre versiones (NA/JP/EU). Usar offsets incorrectos → corrupción.

### Solución
```go
// Detectar la versión del juego leyendo "magic bytes" o patrones conocidos
type GameVersion int

const (
    VersionUnknown GameVersion = iota
    VersionYellowNA  // North America
    VersionYellowJP  // Japan
    VersionYellowEU  // Europe
)

func DetectGameVersion(s *Save) GameVersion {
    // Leer player name offset, checksums, patrones específicos
    // Retornar la versión detectada
}
```

**Impacto:** Evita el 80% de errores de usuario

---

## 2. Verificación de Integridad Pre-Modificación ⭐ CRÍTICO

### Problema
Si el save original está corrupto, modificarlo lo empeorará.

### Solución
```go
type IntegrityReport struct {
    IsValid        bool
    Errors         []string
    Warnings       []string
    GameVersion    GameVersion
    PlayerName     string
    SaveCount      int
    BagItemsValid  bool
    MoneyValid     bool
}

func (s *Save) CheckIntegrity() IntegrityReport {
    report := IntegrityReport{}

    // 1. Checksum
    if !s.ValidateChecksum() {
        report.Errors = append(report.Errors, "Invalid checksum")
        report.IsValid = false
    }

    // 2. Bag count
    bagCount := s.GetByte(OffsetBagCount)
    if bagCount > MaxBagItems {
        report.Errors = append(report.Errors, "Bag count exceeds max")
    }

    // 3. Money range
    money := money.GetMoney(s)
    if money > MaxMoney {
        report.Errors = append(report.Errors, "Money exceeds max")
    }

    // 4. Item IDs válidos
    items := items.GetBagItems(s)
    for _, item := range items {
        if !items.IsValidItemID(item.ID) {
            report.Warnings = append(report.Warnings,
                fmt.Sprintf("Unknown item ID: 0x%02X", item.ID))
        }
    }

    return report
}
```

---

## 3. Sistema de SHA256 Hashing ⭐ IMPORTANTE

### Problema
No hay forma de verificar que el archivo no se modificó accidentalmente.

### Solución
```go
import "crypto/sha256"

func (s *Save) GetHash() string {
    hash := sha256.Sum256(s.data)
    return fmt.Sprintf("%x", hash)
}

// Guardar hash del original
func CreateBackupWithHash(path string) (string, error) {
    s, _ := Load(path)
    hash := s.GetHash()

    // Guardar en backup.sav.sha256
    hashFile := path + ".bak.sha256"
    os.WriteFile(hashFile, []byte(hash), 0644)

    return hash, CreateBackup(path)
}
```

**Uso:**
```bash
./raracandy yellow inspect pokemon.sav
# Output: SHA256: 3a4f2c1b8e9d... (guarda esto)

# Después de modificar, verifica que el original no cambió
./raracandy verify-original pokemon.sav --expected-hash 3a4f2c1b8e9d...
```

---

## 4. Modo "Verify-Only" (Sin Escritura) ⭐ IMPORTANTE

### Problema
Los usuarios pueden equivocarse y sobrescribir sin querer.

### Solución
```go
var verifyOnlyMode bool

// Nuevo comando
var verifyCmd = &cobra.Command{
    Use:   "verify <save-file>",
    Short: "Verify save file integrity without modification",
    RunE:  runVerify,
}

func runVerify(cmd *cobra.Command, args []string) error {
    s, err := save.Load(args[0])
    report := s.CheckIntegrity()

    // Output detallado
    fmt.Printf("Game Version: %s\n", report.GameVersion)
    fmt.Printf("Checksum: %s\n", report.ChecksumStatus)
    fmt.Printf("Integrity: %s\n", report.IsValid)

    if !report.IsValid {
        return fmt.Errorf("save file has integrity issues")
    }

    return nil
}
```

---

## 5. Comparación Before/After con Diff ⭐ ÚTIL

### Problema
Usuario no sabe exactamente qué cambió.

### Solución
```bash
./raracandy yellow diff original.sav modified.sav

# Output:
# Differences found:
#
# Offset 0x25CA (Bag Items):
#   - Rare Candy: 3 → 99 (+96)
#   + Master Ball: 0 → 50 (new)
#
# Offset 0x25F3 (Money):
#   12,345 → 999,999 (+987,654)
#
# Offset 0x3523 (Checksum):
#   0xA7 → 0xF2
#
# Total bytes changed: 8
# Integrity: ✓ Valid
```

---

## 6. Rollback Automático ⭐ CRÍTICO

### Problema
Si falla la escritura, el archivo podría quedar corrupto.

### Solución
```go
func (s *Save) WriteWithRollback(path string) error {
    // 1. Crear backup automático
    backupPath := path + ".rollback.tmp"
    if err := CreateBackup(path); err != nil {
        return err
    }

    // 2. Escribir a archivo temporal primero
    tmpPath := path + ".tmp"
    if err := s.Write(tmpPath); err != nil {
        return err
    }

    // 3. Verificar el archivo temporal
    verify, err := Load(tmpPath)
    if err != nil || !verify.ValidateChecksum() {
        os.Remove(tmpPath)
        return fmt.Errorf("verification failed, rollback")
    }

    // 4. Solo entonces reemplazar el original
    if err := os.Rename(tmpPath, path); err != nil {
        return err
    }

    return nil
}
```

---

## 7. Límites de Seguridad Configurables ⭐ ÚTIL

### Problema
Modificaciones extremas (999 items) son obviamente errores.

### Solución
```go
type SafetyLimits struct {
    MaxItemQuantity  byte   // Default: 99
    MaxMoney         uint32 // Default: 999999
    AllowedItemIDs   []byte // Whitelist de items "seguros"
    WarnOnSuspicious bool   // Avisar si se modifican >50% de items
}

var defaultLimits = SafetyLimits{
    MaxItemQuantity: 99,
    MaxMoney:        999999,
    WarnOnSuspicious: true,
    AllowedItemIDs: []byte{
        items.IDRareCandy,
        items.IDMasterBall,
        // etc
    },
}

func ValidateModification(s *Save, limits SafetyLimits) error {
    items := items.GetBagItems(s)

    for _, item := range items {
        if item.Quantity > limits.MaxItemQuantity {
            return fmt.Errorf("item quantity %d exceeds limit %d",
                item.Quantity, limits.MaxItemQuantity)
        }
    }

    return nil
}
```

---

## 8. Sistema de Confirmación Interactivo ⭐ IMPORTANTE

### Problema
Usuarios ejecutan comandos sin leer el output.

### Solución
```go
func ConfirmDangerousOperation(message string) bool {
    fmt.Printf("\n⚠️  WARNING: %s\n", message)
    fmt.Print("Type 'yes' to continue: ")

    var response string
    fmt.Scanln(&response)

    return strings.ToLower(response) == "yes"
}

// Uso
if !dryRun && !forceFlag {
    if !ConfirmDangerousOperation("You are about to modify your save file") {
        return fmt.Errorf("operation cancelled by user")
    }
}
```

**Flag para scripts:**
```bash
# Modo interactivo (default)
./raracandy yellow add-item save.sav --item rare_candy --qty 99 --out modified.sav
# Pregunta confirmación ⚠️

# Modo forzado (para scripts)
./raracandy yellow add-item save.sav --item rare_candy --qty 99 --out modified.sav --force
# No pregunta
```

---

## 9. Log de Auditoría ⭐ ÚTIL

### Problema
No hay registro de qué se modificó y cuándo.

### Solución
```go
type AuditLog struct {
    Timestamp    time.Time
    InputFile    string
    OutputFile   string
    Command      string
    Changes      []string
    OldChecksum  byte
    NewChecksum  byte
    OldHash      string
    NewHash      string
}

func LogOperation(log AuditLog) {
    logFile := "raracandy_audit.log"

    entry := fmt.Sprintf(
        "[%s] %s → %s | Checksum: 0x%02X → 0x%02X | Changes: %v\n",
        log.Timestamp.Format(time.RFC3339),
        log.InputFile,
        log.OutputFile,
        log.OldChecksum,
        log.NewChecksum,
        log.Changes,
    )

    f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()
    f.WriteString(entry)
}
```

**Output:**
```
raracandy_audit.log:
[2025-02-08T22:30:00Z] pokemon.sav → pokemon_mod.sav | Checksum: 0xA7 → 0xF2 | Changes: [+Rare Candy x99, Money: 12345 → 999999]
```

---

## 10. Test Suite Expandido ⭐ CRÍTICO

### Problema
No hay tests de corrupción intencional.

### Solución
```go
func TestCorruptedSaveDetection(t *testing.T) {
    // Crear save válido
    s := CreateTestSave()

    // Corromper el checksum
    s.SetByte(OffsetChecksum, 0xFF)

    // Debe detectar la corrupción
    if s.ValidateChecksum() {
        t.Error("Failed to detect corrupted checksum")
    }
}

func TestInvalidBagCount(t *testing.T) {
    s := CreateTestSave()

    // Bag count inválido (>20)
    s.SetByte(OffsetBagCount, 99)

    report := s.CheckIntegrity()
    if report.IsValid {
        t.Error("Failed to detect invalid bag count")
    }
}

func TestOffsetBoundaries(t *testing.T) {
    s := CreateTestSave()

    // Intentar escribir fuera de límites
    err := s.SetByte(0x8000, 0xFF)
    if err == nil {
        t.Error("Allowed out-of-bounds write")
    }
}
```

---

## Implementación por Prioridad

### Fase 1: CRÍTICO (Implementar YA)
1. ✅ Verificación de integridad pre-modificación
2. ✅ Rollback automático
3. ✅ Detección de versión del juego
4. ✅ Tests de corrupción

### Fase 2: IMPORTANTE (Próxima semana)
5. ✅ SHA256 hashing
6. ✅ Modo verify-only
7. ✅ Confirmación interactiva

### Fase 3: ÚTIL (Cuando sea necesario)
8. ✅ Diff before/after
9. ✅ Audit log
10. ✅ Safety limits configurables

---

## Mejoras de UX para Seguridad

```bash
# Comando más seguro por defecto
./raracandy yellow add-item pokemon.sav \
  --item rare_candy \
  --qty 99 \
  --out pokemon_mod.sav

# Output mejorado:
# ⚙️  Loading save...
# ✓ Save loaded: pokemon.sav (32 KB)
# ✓ Checksum valid: 0xA7
# ✓ Integrity check passed
#
# 🔍 Detected: Pokemon Yellow (North America)
# 📊 Current state:
#    - Money: ¥12,345
#    - Bag: 5/20 items
#
# 📝 Proposed changes:
#    - Rare Candy: 3 → 99 (+96)
#    - Checksum: 0xA7 → 0xF2 (will recalculate)
#
# ⚠️  WARNING: You are about to modify your save file
# Type 'yes' to continue: yes
#
# 💾 Creating backup: pokemon.sav.bak
# ✓ Backup created
# ✓ Backup hash: 3a4f2c1b8e9d...
#
# ✍️  Writing modified save...
# ✓ Save written: pokemon_mod.sav
# ✓ Checksum updated: 0xA7 → 0xF2
# ✓ Integrity verified
#
# ✅ Success! Your save is ready.
#
# 📋 Audit log: raracandy_audit.log
# 🔐 Backup hash saved: pokemon.sav.bak.sha256
```

---

## Nivel de Confianza Esperado

| Mejora | Nivel Antes | Nivel Después |
|--------|-------------|---------------|
| Base (actual) | 95% | 95% |
| + Detección versión | 95% | 97% |
| + Integrity check | 97% | 98% |
| + Rollback | 98% | 99% |
| + SHA256 hash | 99% | 99.5% |
| + Tests expandidos | 99.5% | 99.8% |

**Objetivo final: 99.8% confianza**

Los 0.2% restantes son:
- Errores de hardware (GBxCart falla)
- Batería del cartucho muerta
- ROM hacks desconocidos
- Acción del usuario (ignora todas las advertencias)
