# raracandy 🍬
A CLI tool to safely edit **Pokémon Gen 1** save files — starting with **Pokémon Yellow**.

> **Goal:** dump a `.sav` from an original cartridge, edit it (e.g., add Rare Candies), recalculate checksums, and write it back — **without touching the ROM**.

---

## Why this exists
I have an original Pokémon Yellow cartridge and I want a **reproducible, auditable, command-line workflow** to:
- Backup my save to a computer
- Inject items (Rare Candy, etc.) into the save file
- Keep the save valid by updating checksums
- Re-write the modified save back to the cartridge using external hardware dumpers

This is a personal project that’s public on GitHub: fun, niche, and portfolio-friendly.

---

## What raracandy edits (and what it does NOT)
### ✅ In scope
- `.sav` files (Game Boy Pokémon Gen 1)
- Items in the **Bag** ✓ Implemented
- Player **Money** ✓ Implemented
- Save validation + checksum recalculation ✓ Implemented

### ❌ Out of scope / non-goals
- ROM distribution or ROM modification
- Real-time “cheat engine” / live memory editing
- GameShark code generation
- Anything piracy-related

raracandy only operates on save files you personally dumped.

---

## Supported games
### v1
- Pokémon Yellow (Gen 1)

### Planned
- Pokémon Red / Blue
- More Gen 1 save features (PC items, party editor, flags/events)

---

## Real hardware workflow (how you’d actually use this)
1. **Dump** your save from the original cartridge to `yellow.sav`
   - Using a dumper such as **GBxCart RW** or **GB Operator**
2. Run raracandy to create a modified save:
   - `yellow_mod.sav`
3. **Write** the modified save back to the cartridge with the same dumper
4. Play normally on real hardware

raracandy does not include dumping/writing functionality; it edits the `.sav`.

---

## Installation

### Build from source

```bash
git clone https://github.com/abravonunez/raracandy.git
cd raracandy
go build -o raracandy ./cmd/raracandy
```

## Usage

### Verify save integrity (Recommended first step)

```bash
raracandy yellow verify pokemon_yellow.sav
```

Output example:
```
Save File: pokemon_yellow.sav
Size: 32 KB

Detected Version: Pokémon Yellow (North America)

Checksum:
  Stored:     0xA7
  Calculated: 0xA7
  Status:     ✓ Valid

Bag Structure:
  Status:     ✓ Valid

Money Format:
  Status:     ✓ Valid BCD encoding

SHA256: 3a4f2c1b8e9d...

Overall Status: ✓ VALID - Safe to modify
```

### Inspect a save file

```bash
raracandy yellow inspect pokemon_yellow.sav
```

Output example:
```
Save File: pokemon_yellow.sav
Size: 32 KB

Checksum:
  Stored:     0xA9
  Calculated: 0xA9
  Status:     ✓ Valid

Money: ¥12,345

Bag (5/20 items):
  - Potion x10
  - Rare Candy x3
  - Master Ball x1
  - Ultra Ball x20
  - Revive x5
```

### Add or modify items

```bash
# Add 99 Rare Candies
raracandy yellow add-item pokemon_yellow.sav \
  --item rare_candy \
  --qty 99 \
  --out pokemon_yellow_modified.sav

# Preview changes without writing (dry-run)
raracandy yellow add-item pokemon_yellow.sav \
  --item master_ball \
  --qty 50 \
  --out pokemon_yellow_modified.sav \
  --dry-run
```

Supported item names: `rare_candy`, `master_ball`, `ultra_ball`, `great_ball`, `poke_ball`, `potion`, `super_potion`, `hyper_potion`, `max_potion`, `full_restore`, `revive`, `max_revive`, and more.

### Set money

```bash
# Set money to 999,999 (max)
raracandy yellow set-money pokemon_yellow.sav \
  --amount 999999 \
  --out pokemon_yellow_rich.sav

# Preview changes
raracandy yellow set-money pokemon_yellow.sav \
  --amount 500000 \
  --out pokemon_yellow_rich.sav \
  --dry-run
```

## Project Structure

```bash
raracandy/
  cmd/raracandy/              # CLI entrypoint
    main.go                   # Main entry point
    yellow.go                 # Yellow subcommand
    inspect.go                # Inspect command
    add_item.go               # Add item command
    set_money.go              # Set money command
  internal/
    gen1/
      save/                   # Save file parser/serializer
      items/                  # Item management & mappings
      money/                  # Money management (BCD encoding)
    backup/                   # Backup system
  testdata/
    fixtures/                 # Test save files
  README.md
  LICENSE
```

## SRE-Focused Features

This tool is designed with reliability and safety in mind:

**Core Safety:**
- **Never overwrites input files** - Always requires explicit `--out` flag
- **Automatic backups** - Creates `.bak` files before any modification
- **Checksum validation** - Verifies save integrity before and after editing
- **Dry-run mode** - Preview changes before applying them (`--dry-run`)

**New in v1.1 (Security Update):**
- **Pre-modification integrity check** - Refuses to modify corrupted saves
- **Game version detection** - Prevents using wrong offsets
- **SHA256 hashing** - Cryptographic verification of backups
- **Interactive confirmation** - Asks permission before modifying (bypass with `--force`)
- **Post-write verification** - Ensures modifications succeeded
- **Comprehensive verify command** - Check save health anytime

**Development:**
- **Clear error messages** - Descriptive errors for debugging
- **Type-safe operations** - Go's type system prevents common mistakes
- **Extensive testing** - Unit tests + integration tests

**Confidence Level: 99%+** when following recommended workflow

## CLI Principles

- Never overwrite input files by default
- Always require an explicit output path via `--out`
- Create backups automatically (e.g. `input.sav.bak`)
- Support `--dry-run` to preview changes without writing
- Produce clear, script-friendly output

---

## v1 Scope (Minimal but Complete)

### Must-have

- Read `.sav` files  
  - Expected size for Gen 1: **32 KB**
- Basic validation heuristics:
  - File size checks
  - Sanity checks for expected save blocks
- Locate **Bag items** and modify them:
  - Insert item if missing (when there is available space)
  - Update quantity if item is already present
- Recalculate save **checksum(s)** correctly
- Write an `output.sav` that the game accepts on real hardware

---

### Nice-to-have (Still Small)

- `inspect` command prints current Bag items in a simple, table-like output ✓
- `--dry-run` shows a diff-style summary of the changes that would be applied ✓

---

## Technical Details

### Save File Format

- **Size**: 32 KB (0x8000 bytes)
- **Structure**: 4 banks of 8 KB each
- **Main data**: Bank 1 (0x2000-0x3FFF)
- **Checksum**: 1 byte at offset 0x3523
  - Calculated by summing bytes 0x2598-0x3522 and applying bitwise NOT

### Memory Offsets (Pokemon Yellow NA)

| Data | Save Offset | RAM Address | Size |
|------|-------------|-------------|------|
| Bag Count | 0x25C9 | D31D | 1 byte |
| Bag Items | 0x25CA | D31E | 40 bytes (20 items × 2 bytes) |
| Money | 0x25F3 | D347 | 3 bytes (BCD encoded) |
| Checksum | 0x3523 | - | 1 byte |

Conversion formula: `Save Offset = RAM Address - 0xAD54`

### References

- [Bulbapedia - Save data structure (Generation I)](https://bulbapedia.bulbagarden.net/wiki/Save_data_structure_(Generation_I))
- [Data Crystal - Pokemon Yellow RAM map](https://datacrystal.tcrf.net/wiki/Pokémon_Yellow/RAM_map)
- [Gen1Py - Python save editor](https://github.com/micah-raney/Gen1Py)
- [Rhydon - C# save editor](https://github.com/SciresM/Rhydon)

---

## Future Enhancements

- Support for Pokemon Red/Blue
- PC items editing
- Party Pokemon editing (levels, moves, stats)
- Event flags/badges modification
- Compare two save files (diff mode)

---

## License

MIT License - see [LICENSE](LICENSE) for details.

## Disclaimer

This tool is for educational and personal use only. It operates exclusively on save files you personally own. No ROMs are distributed or modified. Pokémon is a trademark of Nintendo/Game Freak/The Pokémon Company.








# rarecandy
# rarecandy
