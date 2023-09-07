# Adding new devices to a device tree

## Introduction

ACPI exposes a platform-independent interface for operating systems to perform
power management and other platform-level functions.  Some operating systems
also use ACPI to enumerate devices that are not immediately discoverable, such
as those behind I2C or SPI busses (in contrast to PCI).  This document discusses
the way that coreboot uses the concept of a "device tree" to generate ACPI
tables for usage by the operating system.

## Devicetree and overridetree (if applicable)

For mainboards that are organized around a "reference board" or "baseboard"
model (see ``src/mainboard/google/octopus`` or ``hatch`` for examples), there is
typically a devicetree.cb file that all boards share, and any differences for a
specific board ("variant") are captured in the overridetree.cb file.  Any
settings changed in the overridetree take precedence over those in the main
devicetree.  Note, not all mainboards will have the devicetree/overridetree
distinction, and may only have a devicetree.cb file.  Or you can always just
write the ASL (ACPI Source Language) code yourself.

## Device drivers

Let's take a look at an example entry from
``src/mainboard/google/hatch/variant/hatch/overridetree.cb``:

```
device pci 15.0 on
	chip drivers/i2c/generic
		register "hid" = ""ELAN0000""
		register "desc" = ""ELAN Touchpad""
		register "irq" = "ACPI_IRQ_WAKE_EDGE_LOW(GPP_A21_IRQ)"
		register "wake" = "GPE0_DW0_21"
		device i2c 15 on end
	end
end # I2C #0
```

When this entry is processed during ramstage, it will create a device in the
ACPI SSDT table (all devices in devicetrees end up in the SSDT table).  The ACPI
generation routines in coreboot actually generate the raw bytecode that
represents the device's structure, but looking at ASL code is easier to
understand; see below for what the disassembled bytecode looks like:

```
Scope (\_SB.PCI0.I2C0)
{
    Device (D015)
    {
        Name (_HID, "ELAN0000")  // _HID: Hardware ID
        Name (_UID, Zero)  // _UID: Unique ID
        Name (_DDN, "ELAN Touchpad")  // _DDN: DOS Device Name
        Method (_STA, 0, NotSerialized)  // _STA: Status
        {
            Return (0x0F)
        }
        Name (_CRS, ResourceTemplate ()  // _CRS: Current Resource Settings
        {
            I2cSerialBusV2 (0x0015, ControllerInitiated, 400000,
                AddressingMode7Bit, "\\_SB.PCI0.I2C0",
                0x00, ResourceConsumer, , Exclusive, )
            Interrupt (ResourceConsumer, Edge, ActiveLow, ExclusiveAndWake, ,, )
            {
                0x0000002D,
            }
        })
        Name (_S0W, 0x04)  // _S0W: S0 Device Wake State
        Name (_PRW, Package (0x02)  // _PRW: Power Resources for Wake
        {
            0x15, // GPE #21
            0x03  // Sleep state S3
        })
    }
}
```

You can see it generates _HID, _UID, _DDN, _STA, _CRS, _S0W, and _PRW
names/methods in the Device's scope.

## Utilizing a device driver

The device driver must be enabled for your build.  There will be a CONFIG option
in the Kconfig file in the directory that the driver is in (e.g.,
``src/drivers/i2c/generic`` contains a Kconfig file; the option here is named
CONFIG_DRIVERS_I2C_GENERIC).  The config option will need to be added to your
mainboard's Kconfig file (e.g., ``src/mainboard/google/hatch/Kconfig``) in order
to be compiled into your build.

## Diving into the above example:

Let's take a look at how the devicetree language corresponds to the generated
ASL.

First, note this:

```
    chip drivers/i2c/generic
```

This means that the device driver we're using has a corresponding structure,
located at ``src/drivers/i2c/generic/chip.h``, named **struct
drivers_i2c_generic_config** and it contains many properties you can specify to
be included in the ACPI table.

### hid

```
    register "hid" = ""ELAN0000""
```

This corresponds to **const char *hid** in the struct.  In the ACPI ASL, it
translates to:

```
    Name (_HID, "ELAN0000") // _HID: Hardware ID
```

under the device.  **This property is used to match the device to its driver
during enumeration in the OS.**

### desc

```
    register "desc" = ""ELAN Touchpad""
```

corresponds to **const char *desc** and in ASL:

```
    Name (_DDN, "ELAN Touchpad") // _DDN: DOS Device Name
```

### irq

It also adds the interrupt,

```
    Interrupt (ResourceConsumer, Edge, ActiveLow, ExclusiveAndWake, ,, )
    {
        0x0000002D,
    }
```

which comes from:

```
    register "irq" = "ACPI_IRQ_WAKE_EDGE_LOW(GPP_A21_IRQ)"
```

The GPIO pin IRQ settings control the "Edge", "ActiveLow", and
"ExclusiveAndWake" settings seen above (edge means it is an edge-triggered
interrupt as opposed to level-triggered; active low means the interrupt is
triggered on a falling edge).

Note that the ACPI_IRQ_WAKE_EDGE_LOW macro informs the platform that the GPIO
will be routed through SCI (ACPI's System Control Interrupt) for use as a wake
source.  Also note that the IRQ names are SoC-specific, and you will need to
find the names in your SoC's header file.  The ACPI_* macros are defined in
``src/arch/x86/include/arch/acpi_device.h``.

Using a GPIO as an IRQ requires that it is configured in coreboot correctly.
This is often done in a mainboard-specific file named ``gpio.c``.

### wake

The last register is:

```
    register "wake" = "GPE0_DW0_21"
```

which indicates that the method of waking the system using the touchpad will be
through a GPE, #21 associated with DW0, which is set up in devicetree.cb from
this example.  The "21" indicates GPP_X21, where GPP_X is mapped onto DW0
elsewhere in the devicetree.

The last bit of the definition of that device includes:

```
    device i2c 15 on end
```

which means it's an I2C device, with 7-bit address 0x15, and the device is "on",
meaning it will be exposed in the ACPI table.  The PCI device that the
controller is located in determines which I2C bus the device is expected to be
found on.  In this example, this is I2C bus 0.  This also determines the ACPI
"Scope" that the device names and methods will live under, in this case
"\_SB.PCI0.I2C0".

## Other auto-generated names

(see [ACPI specification
6.3](https://uefi.org/sites/default/files/resources/ACPI_6_3_final_Jan30.pdf)
for more details on ACPI methods)

### _S0W (S0 Device Wake State)
_S0W indicates the deepest S0 sleep state this device can wake itself from,
which in this case is 4, representing _D3cold_.

### _PRW (Power Resources for Wake)
_PRW indicates the power resources and events required for wake.  There are no
dependent power resources, but the GPE (GPE0_DW0_21) is mentioned here (0x15),
as well as the deepest sleep state supporting waking the system (3), which is
S3.

### _STA (Status)
The _STA method is generated automatically, and its values, 0xF, indicates the
following:

    Bit [0] – Set if the device is present.
    Bit [1] – Set if the device is enabled and decoding its resources.
    Bit [2] – Set if the device should be shown in the UI.
    Bit [3] – Set if the device is functioning properly (cleared if device failed its diagnostics).

### _CRS (Current resource settings)
The _CRS method is generated automatically, as the driver knows it is an I2C
controller, and so specifies how to configure the controller for proper
operation with the touchpad.

```
Name (_CRS, ResourceTemplate ()  // _CRS: Current Resource Settings
{
    I2cSerialBusV2 (0x0015, ControllerInitiated, 400000,
                    AddressingMode7Bit, "\\_SB.PCI0.I2C0",
                    0x00, ResourceConsumer, , Exclusive, )
```

## Notes

 - **All fields that are left unspecified in the devicetree are initialized to
   zero.**
 - **All devices in devicetrees end up in the SSDT table, and are generated in
   coreboot's ramstage**
