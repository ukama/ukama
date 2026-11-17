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

echo "Enable PA1 Bias settings"

echo "Writing PA1 init"
i2cset -y 1 0x0c 0x7F 0xFF 0xFF i 

echo "Writing PA1 carrier data."
i2cset -y 1 0x0c 0x59 0x78 0x15 i 

echo "Writing PA1 peak"
i2cset -y 1 0x0c 0x58 0x4C 0xDC i 

echo "Enable PA2 Bias settings"

echo "Writing PA2 init"
i2cset -y 2 0x0c 0x7F 0xFF 0xFF i

echo "Writing PA2 carrier data."
i2cset -y 2 0x0c 0x59 0x78 0x15 i

echo "Writing PA2 peak"
i2cset -y 2 0x0c 0x58 0x4C 0xDC i

exit 0
