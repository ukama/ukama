# Build Container image for emulating device nodes 

Remember to adjust path according to your directory structure.

## Build Steps

### Build UkamaOS.

   Follow steps from [build Ukama OS](https://github.com/ukama/ukamaOS/tree/lwm2m_E2E#readme)

### Create device sysfs:

** For HNODE **

```
./lib/ubsp/utils/prepare_env.sh -u cnode-lte
```

** For ANODE **

```
./lib/ubsp/utils/prepare_env.sh -u anode
```

### Create Manufacturing Schema 
Manufacturing json schema for EEPROM with serail Id's.Only last four digit are provided for now.

** For HNODE **

```
./lib/ubsp/build/schema -n ComV1 -u UK-0001-HNODE-SA03-1102 -m UK-1001-COM-1102 -f mfgdata/schema/com.json -n LTE -m UK-2001-LTE-1102 -f mfgdata/schema/lte.json -n MASK -m UK-3001-MSK-1102 -f mfgdata/schema/mask.json
```

** For ANODE **

```
./lib/ubsp/build/schema -u UK-5001-ANODE-SA03-1102 -m UK-5001-RFC-1102 -n "RF CTRL BOARD" -f mfgdata/schema/rfctrl.json -m UK-4001-RFA-1102 -n "RF BOARD" -f mfgdata/schema/rffe.json 
```


### Create EEPROM Database

**For HNODE board with LTE and Mask board**

```
./lib/ubsp/build/mfgutil -n ComV1 -m UK-1001-COM-1102 -s mfgdata/schema/com.json -n LTE -m UK-2001-LTE-1102 -s mfgdata/schema/lte.json -n MASK -m UK-3001-MSK-1102 -s mfgdata/schema/mask.json
```

** For ANODE **

```
./lib/ubsp/build/mfgutil -n "RF CTRL BOARD" -m UK-5001-RFC-1102 -s mfgdata/schema/rfctrl.json -n "RF BOARD" -m UK-4001-RFA-1102 -s mfgdata/schema/rffe.json
```

### Copy the SYS FS of the device to container

```
cp -rf /tmp/sys ./container/sys
```




