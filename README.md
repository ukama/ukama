# NodeD Service

NodeD service is to manage and support HW devices.

## Architecture Diagram
NodeD is divided into few logical blocks based on thier operations.
* Inventory

> Handles the manufacturing data related operations like creating invenrory database for modules, storing and retriving data from inventory database.<br>
Data stored in inventory data base is description of the hardware it's present on listingout sensors available of module,
it's manufacturing dates etc. Apart from Harware information it also stores configurations and certficates used by software to perform day to day operations.


* Ledger

> Handles sensor configurations such as updating behavior, reading status, enabling alerts, disabling alerts etc.

* Comm

> Provide a REST API based inetrface for other application to access and update harware status and certificates

![NodeD](docs/NodeD.jpg)

## Building
Preferred build method is to use UkamaOS build.

```
make
```

## Testing

#### Preparing setup

**Mocking SysFileSystem**

For testing purpose we can mock our sysfs under /tmp/sys directory using prepare_env.sh script

Example:

```
./utils/prepare_env.sh -u cnode-lte -u anode
```

**Generate Schema**

Dummy schema are provided under mfgdata/schema folder. Modification can be made to it on need basis.
If we just need to replicate these with updated serial numbers we could use ustility like genSchema.

Example:

```
./build/genSchema --n ComV1 --u UK-7001-HNODE-SA03-1102 --f mfgdata/schema/com.json --m UK-7001-COM-1102 --f mfgdata/schema/lte.json --f mfgdata/schema/mask.json
```
Could use this for more information

```
/build/genSchema --help
```

**Generate Inventory Database**

This utilty creates a inventory database for the modules you supplies as an argument to the utilty and place those under /tmp/sys directory

Example:

```
./build/genInventory --n COM --m UK-8001-COM-1102 --s mfgdata/schema/com.json -n LTE --m UK-8001-LTE-1102 --s mfgdata/schema/lte.json --n MASK -m UK-8001-MASK-1102 --s mfgdata/schema/mask.json
```

Again this could be used for more information

```
./build/genInventory --help
```

#### Run NodeD service

You can run noded service with deafult arguments but if wishes to change ineventory database or sensor related attributes
those can be provided in config files and supplied as argument to noded.

```
./build/noded
```

