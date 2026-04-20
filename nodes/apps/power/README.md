On dev:

export POWER_MOCK=1
export POWER_BOARD=dev
export POWER_SAMPLE_MS=1000
./power-monitor -l DEBUG


On real hardware:

export POWER_MOCK=0
export POWER_BOARD=tower
export POWER_SAMPLE_MS=1000

export POWER_LM75_DEV=/dev/i2c-1
export POWER_LM75_ADDR=0x48

export POWER_LM25066_DEV=/dev/i2c-1
export POWER_LM25066_ADDR=0x40
export POWER_LM25066_CL_HIGH=0
export POWER_LM25066_RS_MOHM=1

export POWER_ADS1015_DEV=/dev/i2c-1
export POWER_ADS1015_ADDR=0x49
export POWER_ADS1015_CHMAP=vin=0,vpa=1,aux=2

./power-monitor -l DEBUG
