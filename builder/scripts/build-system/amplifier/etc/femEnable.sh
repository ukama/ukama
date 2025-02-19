#!/bin/sh
set -x
cd /sys/devices/platform
echo "Enable PA1 VDS"
echo 1 > ./fema1-gpios/pa_vds_enable

echo "Enable PA1"
echo 1 > ./fema1-gpios/rx_rf_enable 
echo 1 > ./fema1-gpios/tx_rf_enable

echo "Enable PA2 VDS"
echo 1 > ./fema2-gpios/pa_vds_enable

echo "Enable PA2"
echo 1 > ./fema2-gpios/rx_rf_enable
echo 1 > ./fema2-gpios/tx_rf_enable

echo "Remove control for 0x0C"

cd /sys/bus/i2c/devices/i2c-1/1-000c/driver
echo 1-000c > unbind 

cd /sys/bus/i2c/devices/i2c-2/2-000c/driver
echo 2-000c > unbind

echo "Enable PA1 Bias settings"

echo "Writing PA1 init"
i2cset -y 1 0x0c 0x7F 0xFF 0xFF i 

echo "Writing PA1 carrier data."
i2cset -y 1 0x0c 0x59 0x83 0x15 i 

echo "Writing PA1 peak"
i2cset -y 1 0x0c 0x58 0x3F 0xDC i 

echo "Enable PA2 Bias settings"

echo "Writing PA2 init"
i2cset -y 2 0x0c 0x7F 0xFF 0xFF i

echo "Writing PA2 carrier data."
i2cset -y 2 0x0c 0x59 0x83 0x15 i

echo "Writing PA2 peak"
i2cset -y 2 0x0c 0x58 0x3F 0xDC i

exit 0
