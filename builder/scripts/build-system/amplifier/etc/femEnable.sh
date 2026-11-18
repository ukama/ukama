#!/bin/sh
set -x
cd /sys/devices/platform
echo "Enable PA1 VDS"
echo 1 > ./fema1-gpios/pa_vds_enable

echo "Enable PA1"
echo 1 > ./fema1-gpios/rx_rf_enable 
echo 1 > ./fema1-gpios/tx_rf_enable
echo 1 > ./fema1-gpios/rf_pal_enable

echo "Enable PA2 VDS"
echo 1 > ./fema2-gpios/pa_vds_enable

echo "Enable PA2"
echo 1 > ./fema2-gpios/rx_rf_enable
echo 1 > ./fema2-gpios/tx_rf_enable
echo 1 > ./fema2-gpios/rf_pal_enable

echo "Remove kernel control for the DAC, temperature sensors and ADC"

for dev in 1-000c 1-0018 1-0049 1-0048 2-000c 2-0018 2-0049 2-0048; do
    if [ -e "/sys/bus/i2c/devices/${dev}/driver" ]; then
        echo "${dev}" > "/sys/bus/i2c/devices/${dev}/driver/unbind"
    fi
done

band=$(sed -n 's/.*"FEM_BAND": *"\([^"]*\)".*/\1/p' /ukama/manifest.json 2>/dev/null | head -n 1)

case "$band" in
B41)
    carrier="0x73 0x33"
    peak="0x4C 0xCD"
    ;;
B1)
    carrier="0x80 0x00"
    peak="0x51 0xEC"
    ;;
*)
    echo "No PA bias values for band '${band}', the PA bias is not set"
    exit 0
    ;;
esac

for fe in 1 2; do
    echo "Writing PA${fe} init"
    i2cset -y ${fe} 0x0c 0x7F 0xFF 0xFF i

    echo "Writing PA${fe} ${band} carrier ${carrier}"
    i2cset -y ${fe} 0x0c 0x59 ${carrier} i

    echo "Writing PA${fe} ${band} peak ${peak}"
    i2cset -y ${fe} 0x0c 0x58 ${peak} i
done

exit 0
