#!/usr/bin/env bash
# mock-anode-sys.sh
#
# Creates a mocked ANODE sysfs/dev tree under /tmp/sys (or --root).
# Adds /dev/i2c-* mock endpoints under that root for tests.
#
# Optional: install i2c-tools wrappers into $ROOT/bin and prepend to PATH in tests:
#   export PATH="/tmp/sys/bin:$PATH"
#
set -euo pipefail

ROOT="/tmp/sys"
MAKE_WRAPPERS=0
DEV_MODE="fifo"   # fifo|file  (fifo is nicer for emulation, file is simplest)
VERBOSE=1

usage() {
  cat <<EOF
Usage: $0 [--root <path>] [--clean] [--wrappers] [--dev-mode fifo|file]

Options:
  --root <path>        Root directory for mock tree (default: /tmp/sys)
  --clean              Remove root directory and exit
  --wrappers           Install i2cget/i2cset/i2cdetect wrapper scripts in \$ROOT/bin
  --dev-mode fifo|file Create /dev/i2c-* as FIFO or plain file (default: fifo)
EOF
}

log() { [[ "$VERBOSE" -eq 1 ]] && echo "$@"; }
mkdirp() { mkdir -p "$1"; }
mkfile() {
  local path="$1" content="${2:-0}"
  mkdir -p "$(dirname "$path")"
  printf "%s\n" "$content" > "$path"
}
mklink() {
  local target="$1" link="$2"
  mkdir -p "$(dirname "$link")"
  ln -snf "$target" "$link"
}

mk_gpio_class() {
  local n="$1"
  local base="$ROOT/class/gpio/gpio${n}"
  mkdirp "$base"
  mkfile "$base/active_low" "0"
  mkfile "$base/direction" "in"
  mkfile "$base/edge" "none"
  mkfile "$base/polairy" "0"  # keep legacy spelling
  mkfile "$base/value" "0"
}

mk_platform_gpios() {
  local devname="$1"; shift
  local base="$ROOT/devices/platform/${devname}"
  mkdirp "$base"
  # placeholders
  mkfile "$base/subsystem" ""
  mkfile "$base/driver" ""
  mkfile "$base/driver_override" ""
  mkfile "$base/modalias" ""
  mkfile "$base/uevent" ""
  mkfile "$base/of_node" ""
  # attrs
  for f in "$@"; do
    mkfile "$base/$f" "0"
  done
}

mk_i2c_bus() {
  local bus="$1"
  mkdirp "$ROOT/bus/i2c/devices/i2c-${bus}"
}

mk_i2c_dev() {
  local bus="$1" addr_hex="$2"
  local addr_dec=$((addr_hex))
  local addr4; addr4="$(printf "%04x" "$addr_dec")"
  local devdir="$ROOT/bus/i2c/devices/i2c-${bus}/${bus}-${addr4}"
  mkdirp "$devdir"
  mkfile "$devdir/name" ""
  mkfile "$devdir/uevent" ""
  echo "$devdir"
}

mk_eeprom_file() {
  local devdir="$1"
  mkdirp "$devdir"
  # 256B inventory is enough for most tests; bump if needed.
  dd if=/dev/zero of="$devdir/eeprom" bs=1 count=256 status=none
}

mk_hwmon() {
  local base="$ROOT/class/hwmon/hwmon0"
  mkdirp "$base"
  mkfile "$base/4" "0"
  mkfile "$base/5" "0"
  mkfile "$base/name" "anode-mock"
  mkfile "$base/temp1_input" "42000"
  mkfile "$base/temp2_input" "43000"
  mkfile "$base/temp3_input" "44000"
}

mk_leds() {
  for i in 0 1 2 3; do
    mkdirp "$ROOT/class/led/led${i}"
    mkfile "$ROOT/class/led/led${i}/red" "0"
  done
}

mk_dev_i2c_endpoints() {
  mkdirp "$ROOT/dev"
  for b in 0 1 2; do
    local p="$ROOT/dev/i2c-$b"
    rm -f "$p"
    if [[ "$DEV_MODE" == "fifo" ]]; then
      mkfifo "$p"
    else
      : > "$p"
    fi
  done
  log "Created mock /dev endpoints: $ROOT/dev/i2c-{0,1,2} ($DEV_MODE)"
}

install_i2c_wrappers() {
  mkdirp "$ROOT/bin"

  # i2cdetect wrapper: prints busses that "exist" in the mock
  cat > "$ROOT/bin/i2cdetect" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
ROOT="${UKAMA_SYSROOT:-/tmp/sys}"

if [[ "${1:-}" == "-l" ]]; then
  for b in 0 1 2; do
    echo "i2c-$b\tunknown\tI2C adapter\tMock"
  done
  exit 0
fi

echo "mock i2cdetect: only '-l' is supported here" >&2
exit 2
SH
  chmod +x "$ROOT/bin/i2cdetect"

  # i2cget wrapper: reads from a deterministic per-device "reg" file if present
  cat > "$ROOT/bin/i2cget" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
ROOT="${UKAMA_SYSROOT:-/tmp/sys}"

# minimal parsing for: i2cget -y <bus> <addr> <reg> w
# We don't emulate ADC conversion; we just return a placeholder value.
if [[ "${1:-}" == "-y" ]]; then shift; fi
BUS="${1:-}"; ADDR="${2:-}"; REG="${3:-0x00}"; MODE="${4:-}"

if [[ -z "$BUS" || -z "$ADDR" ]]; then
  echo "usage: i2cget -y <bus> <addr> <reg> [mode]" >&2
  exit 2
fi

# normalize addr and reg to lowercase hex without leading zeros issues
ADDR_DEC=$((ADDR))
ADDR4=$(printf "%04x" "$ADDR_DEC")
DEV="$ROOT/bus/i2c/devices/i2c-$BUS/$BUS-$ADDR4"

# If a mock register file exists, return it. Else return 0x0000.
REG_DEC=$((REG))
REG2=$(printf "%02x" "$REG_DEC")
REGF="$DEV/reg_$REG2"

if [[ -f "$REGF" ]]; then
  cat "$REGF"
else
  # "w" mode expects a 16-bit word printed as 0xNNNN in i2c-tools
  echo "0x0000"
fi
SH
  chmod +x "$ROOT/bin/i2cget"

  # i2cset wrapper: stores written values into reg_<xx> files in the mock tree
  cat > "$ROOT/bin/i2cset" <<'SH'
#!/usr/bin/env bash
set -euo pipefail
ROOT="${UKAMA_SYSROOT:-/tmp/sys}"

# supports: i2cset -y <bus> <addr> <reg> <val1> [val2 ...] [i]
if [[ "${1:-}" == "-y" ]]; then shift; fi
BUS="${1:-}"; ADDR="${2:-}"; REG="${3:-}"

if [[ -z "$BUS" || -z "$ADDR" || -z "$REG" ]]; then
  echo "usage: i2cset -y <bus> <addr> <reg> <val...> [mode]" >&2
  exit 2
fi
shift 3

ADDR_DEC=$((ADDR))
ADDR4=$(printf "%04x" "$ADDR_DEC")
DEV="$ROOT/bus/i2c/devices/i2c-$BUS/$BUS-$ADDR4"
mkdir -p "$DEV"

REG_DEC=$((REG))
REG2=$(printf "%02x" "$REG_DEC")
REGF="$DEV/reg_$REG2"

# collect values until last arg if it's "i" or similar mode
vals=()
for a in "$@"; do
  if [[ "$a" == "i" || "$a" == "b" || "$a" == "w" || "$a" == "s" ]]; then
    break
  fi
  vals+=("$a")
done

# Store as a single line; keep i2c-tools-like "0xNNNN" if caller passes hex.
printf "%s " "${vals[@]}" | sed 's/[[:space:]]*$//' > "$REGF"
SH
  chmod +x "$ROOT/bin/i2cset"

  log "Installed wrappers in: $ROOT/bin"
  log "Use in tests:"
  log "  export UKAMA_SYSROOT=\"$ROOT\""
  log "  export PATH=\"$ROOT/bin:\$PATH\""
}

init_anode_tree() {
  mkdirp "$ROOT/bus/i2c/devices"
  mkdirp "$ROOT/class/gpio" "$ROOT/class/hwmon" "$ROOT/class/led"
  mkdirp "$ROOT/devices/platform"

  # i2c busses
  mk_i2c_bus 0
  mk_i2c_bus 1
  mk_i2c_bus 2

  # Controller (bus 0): TMP10x @ 0x48
  local d0_tmp; d0_tmp="$(mk_i2c_dev 0 0x48)"
  mkfile "$d0_tmp/name" "tmp10x"

  # Inventory EEPROM expected @ bus0:0x51
  local d0_inv; d0_inv="$(mk_i2c_dev 0 0x51)"
  mkfile "$d0_inv/name" "inventory-eeprom"
  mk_eeprom_file "$d0_inv"

  # FEM1 (bus 1): 0x49 LM75A, 0x48 ADS1015, 0x0C AD5667, 0x50 EEPROM
  local d1_lm75 d1_ads d1_dac d1_eep
  d1_lm75="$(mk_i2c_dev 1 0x49)"; mkfile "$d1_lm75/name" "lm75a"
  d1_ads="$(mk_i2c_dev 1 0x48)";  mkfile "$d1_ads/name" "ads1015"
  d1_dac="$(mk_i2c_dev 1 0x0C)";  mkfile "$d1_dac/name" "ad5667"
  d1_eep="$(mk_i2c_dev 1 0x50)";  mkfile "$d1_eep/name" "eeprom"; mk_eeprom_file "$d1_eep"

  # FEM2 (bus 2): same
  local d2_lm75 d2_ads d2_dac d2_eep
  d2_lm75="$(mk_i2c_dev 2 0x49)"; mkfile "$d2_lm75/name" "lm75a"
  d2_ads="$(mk_i2c_dev 2 0x48)";  mkfile "$d2_ads/name" "ads1015"
  d2_dac="$(mk_i2c_dev 2 0x0C)";  mkfile "$d2_dac/name" "ad5667"
  d2_eep="$(mk_i2c_dev 2 0x50)";  mkfile "$d2_eep/name" "eeprom"; mk_eeprom_file "$d2_eep"

  # inventory symlinks required by your code/tests
  mklink "$ROOT/bus/i2c/devices/i2c-0/0-0051/eeprom" "$ROOT/inventory_db"
  mklink "$ROOT/bus/i2c/devices/i2c-0/0-0051/eeprom" "$ROOT/anode_inventory_db"

  # platform gpios
  mk_platform_gpios "controller-gpios" \
    "psu_pgood" "ctrlr_therm_alert" "ctrlr_eeprom_wp"

  mk_platform_gpios "fema1-gpios" \
    "pa_disable" "pg_reg_5v" "rx_rf_enable" "eeprom_wp_enable" \
    "pa_vds_enable" "rf_pal_enable" "tx_rf_enable"

  mk_platform_gpios "fema2-gpios" \
    "pa_disable" "pg_reg_5v" "rx_rf_enable" "eeprom_wp_enable" \
    "pa_vds_enable" "rf_pal_enable" "tx_rf_enable"

  # legacy gpio class entries (keep your current ones)
  for n in 34 35 38 40 61 63; do mk_gpio_class "$n"; done

  mk_hwmon
  mk_leds
  mk_dev_i2c_endpoints

  if [[ "$MAKE_WRAPPERS" -eq 1 ]]; then
    install_i2c_wrappers
  fi

  log "ANode mock created at: $ROOT"
}

ACTION="init"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --root) ROOT="$2"; shift 2;;
    --clean) ACTION="clean"; shift;;
    --wrappers) MAKE_WRAPPERS=1; shift;;
    --dev-mode) DEV_MODE="$2"; shift 2;;
    -h|--help) usage; exit 0;;
    *) echo "Unknown arg: $1"; usage; exit 1;;
  esac
done

case "$ACTION" in
  clean) rm -rf "$ROOT"; echo "Cleaned: $ROOT";;
  init) rm -rf "$ROOT"; init_anode_tree;;
esac
