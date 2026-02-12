#!/usr/bin/env bash
# mock-sysfs-anode.sh
#
# Mock sysfs tree for Amplifier Node (ANode) under /tmp/sys (or --root).
# Models:
#  - Controller: I2C bus 0  (0x48 TMP10x) + controller-gpios (platform)
#  - FEM1:       I2C bus 1  (0x49 LM75A, 0x48 ADS1015, 0x0C AD5667, 0x50 EEPROM) + fema1-gpios (platform)
#  - FEM2:       I2C bus 2  (0x49 LM75A, 0x48 ADS1015, 0x0C AD5667, 0x50 EEPROM) + fema2-gpios (platform)
#
# Also creates inventory symlinks required by your tests:
#   /tmp/sys/inventory_db        -> .../i2c-0/0-0051/eeprom
#   /tmp/sys/anode_inventory_db  -> .../i2c-0/0-0051/eeprom
#
# NOTE:
# - This script only creates filesystem objects (dirs/files/symlinks).
# - It does NOT emulate actual i2cget/i2cset behavior. Your unit tests that
#   call i2c-tools must be stubbed/mocked separately, or use a fake i2c layer.
#
set -euo pipefail

ROOT="/tmp/sys"

usage() {
    cat <<EOF
Usage: $0 [--root <path>] [--clean] [--init]

Options:
  --root <path>   Sysfs root directory (default: /tmp/sys)
  --clean         Remove root directory and exit
  --init          Create the mocked sysfs tree (default action)
EOF
}

clean_root() {
    rm -rf "$ROOT"
    echo "Cleaned: $ROOT"
}

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
# Create a "sysfs-like" gpio class entry
mk_gpio_class() {
    local n="$1"
    local base="$ROOT/class/gpio/gpio${n}"
    mkdirp "$base"
    mkfile "$base/active_low" "0"
    mkfile "$base/direction" "in"
    mkfile "$base/edge" "none"
    # keep your existing typo/compat: "polairy"
    mkfile "$base/polairy" "0"
    mkfile "$base/value" "0"
}

# Create a platform "gpios driver" directory and leaf files (simple 0/1 files)
mk_platform_gpios() {
    local devname="$1"; shift
    local base="$ROOT/devices/platform/${devname}"
    mkdirp "$base"

    # sysfs-ish placeholders (some code checks these exist)
    mkfile "$base/subsystem" ""
    mkfile "$base/driver" ""
    mkfile "$base/driver_override" ""
    mkfile "$base/modalias" ""
    mkfile "$base/uevent" ""
    mkfile "$base/of_node" ""

    # actual gpio attributes your system reads/writes
    for f in "$@"; do
        mkfile "$base/$f" "0"
    done
}

# Create i2c adapter and device directories under /bus/i2c/devices
# Format: /bus/i2c/devices/i2c-<bus>/<bus>-<addr4>/
mk_i2c_bus() {
    local bus="$1"
    mkdirp "$ROOT/bus/i2c/devices/i2c-${bus}"
}

mk_i2c_dev() {
    local bus="$1" addr_hex="$2" devdir
    # addr_hex like 0x48, 0x0C, 0x50 etc
    local addr_dec=$((addr_hex))
    local addr4
    addr4="$(printf "%04x" "$addr_dec")"   # e.g. 0048, 000c, 0050

    devdir="$ROOT/bus/i2c/devices/i2c-${bus}/${bus}-${addr4}"
    mkdirp "$devdir"

    # common placeholders
    mkfile "$devdir/name" ""
    mkfile "$devdir/uevent" ""

    # For EEPROM-type devices, create an "eeprom" file (binary-ish ok; we just need it to exist)
    # We'll create for any device you pass as EEPROM. Call separately.
    echo "$devdir"
}

mk_eeprom_file() {
    local devdir="$1"
    # Make a deterministic file; adjust size if your code expects more.
    # 256 bytes is often enough for inventory in tests.
    mkdirp "$devdir"
    dd if=/dev/zero of="$devdir/eeprom" bs=1 count=256 status=none
}

# hwmon: keep backward-compat with your current mock (hwmon0 contains files "4" and "5")
# and also add more standard hwmon files (name/temp1_input) for future-proofing.
mk_hwmon() {
    local base="$ROOT/class/hwmon/hwmon0"
    mkdirp "$base"

    # legacy files your current tree shows
    mkfile "$base/4" "0"
    mkfile "$base/5" "0"

    # more standard-ish
    mkfile "$base/name" "anode-mock"
    mkfile "$base/temp1_input" "42000"  # milli-C
    mkfile "$base/temp2_input" "43000"
    mkfile "$base/temp3_input" "44000"
}

mk_leds() {
    for i in 0 1 2 3; do
        mkdirp "$ROOT/class/led/led${i}"
        # your current tree has ledX/red; keep that
        mkfile "$ROOT/class/led/led${i}/red" "0"
    done
}

init_anode_tree() {
    # base dirs
    mkdirp "$ROOT/bus/i2c/devices"
    mkdirp "$ROOT/class/gpio"
    mkdirp "$ROOT/class/hwmon"
    mkdirp "$ROOT/class/led"
    mkdirp "$ROOT/devices/platform"

    # i2c busses
    mk_i2c_bus 0
    mk_i2c_bus 1
    mk_i2c_bus 2

    # Controller devices (bus 0)
    # 0x48 TMP10x (controller temp sensor)
    local d0_tmp
    d0_tmp="$(mk_i2c_dev 0 0x48)"
    mkfile "$d0_tmp/name" "tmp10x"

    # Inventory EEPROM currently expected at 0x51 on bus 0 (per your existing symlink)
    local d0_inv
    d0_inv="$(mk_i2c_dev 0 0x51)"
    mkfile "$d0_inv/name" "inventory-eeprom"
    mk_eeprom_file "$d0_inv"

    # FEM1 devices (bus 1)
    local d1_lm75 d1_ads d1_dac d1_eep
    d1_lm75="$(mk_i2c_dev 1 0x49)"; mkfile "$d1_lm75/name" "lm75a"
    d1_ads="$(mk_i2c_dev 1 0x48)";  mkfile "$d1_ads/name" "ads1015"
    d1_dac="$(mk_i2c_dev 1 0x0C)";  mkfile "$d1_dac/name" "ad5667"
    d1_eep="$(mk_i2c_dev 1 0x50)";  mkfile "$d1_eep/name" "eeprom"; mk_eeprom_file "$d1_eep"

    # FEM2 devices (bus 2)
    local d2_lm75 d2_ads d2_dac d2_eep
    d2_lm75="$(mk_i2c_dev 2 0x49)"; mkfile "$d2_lm75/name" "lm75a"
    d2_ads="$(mk_i2c_dev 2 0x48)";  mkfile "$d2_ads/name" "ads1015"
    d2_dac="$(mk_i2c_dev 2 0x0C)";  mkfile "$d2_dac/name" "ad5667"
    d2_eep="$(mk_i2c_dev 2 0x50)";  mkfile "$d2_eep/name" "eeprom"; mk_eeprom_file "$d2_eep"

    # inventory symlinks (required)
    mklink "$ROOT/bus/i2c/devices/i2c-0/0-0051/eeprom" "$ROOT/inventory_db"
    mklink "$ROOT/bus/i2c/devices/i2c-0/0-0051/eeprom" "$ROOT/anode_inventory_db"

    # platform gpios (based on your real /sys/devices/platform listing)
    mk_platform_gpios "controller-gpios" \
                      "psu_pgood" "ctrlr_therm_alert" "ctrlr_eeprom_wp"

    mk_platform_gpios "fema1-gpios" \
                      "pa_disable" "pg_reg_5v" "rx_rf_enable" "eeprom_wp_enable" \
                      "pa_vds_enable" "rf_pal_enable" "tx_rf_enable"

    mk_platform_gpios "fema2-gpios" \
                      "pa_disable" "pg_reg_5v" "rx_rf_enable" "eeprom_wp_enable" \
                      "pa_vds_enable" "rf_pal_enable" "tx_rf_enable"

    # class gpio entries (keep your existing ones so older tests don’t break)
    # If you know the exact mapping from platform gpios -> gpio numbers, replace these.
    for n in 34 35 38 40 61 63; do
        mk_gpio_class "$n"
    done

    mk_hwmon
    mk_leds

    echo "ANode mock sysfs created at: $ROOT"
    echo "Key symlinks:"
    echo "  $ROOT/inventory_db -> $(readlink -f "$ROOT/inventory_db" || true)"
    echo "  $ROOT/anode_inventory_db -> $(readlink -f "$ROOT/anode_inventory_db" || true)"
}

# --- arg parsing ---
ACTION="init"
while [[ $# -gt 0 ]]; do
  case "$1" in
      --root) ROOT="$2"; shift 2;;
      --clean) ACTION="clean"; shift;;
      --init) ACTION="init"; shift;;
      -h|--help) usage; exit 0;;
      *) echo "Unknown arg: $1"; usage; exit 1;;
  esac
done

case "$ACTION" in
    clean) clean_root;;
    init)
        # If it already exists, wipe it so tests are deterministic
        rm -rf "$ROOT"
        init_anode_tree
        ;;
esac

exit 0
