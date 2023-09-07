// SPDX-License-Identifier: GPL-2.0+
/*
 * Copyright (C) 2017 Microchip
 *		      Wenyou Yang <wenyou.yang@microchip.com>
 */

#include <common.h>
#include <dm.h>
#include <environment.h>
#include <i2c_eeprom.h>
#include <netdev.h>

#define BOOTDEV_ADD 0x00
#define BOOTPART_ADD 0x01
#define PREVBOOTDEV_ADD 0x02
#define PREVBOOTPART_ADD 0x03
#define EXPECTEDBOOTDEV_ADD 0x04 
#define EXPECTEDBOOTPART_ADD 0x05
#define UPGRADEFLAG_ADD 0x06
#define BOOTDEVFALLBACK_ADD 0x07
#define BOOTPARTFALLBACK_ADD 0x08
#define MAX_BOOT_CFG_IDX          0x09

#define MIN_CFG_VAL_IDX 0
#define MAX_CFG_VAL_IDX 1
uint8_t validBootCfg[MAX_BOOT_CFG_IDX][2] = { 
					{0, 1},
					{1, 2},
					{0, 1},
					{1, 2},
					{0, 1},
					{1, 2},
					{0, 1},
					{0, 1},
					{0, 1}
				   };

					
int validate_bootcfg(int bootcfg, int val)
{	
	int ret = 1;
	if( (val>=validBootCfg[bootcfg][MIN_CFG_VAL_IDX]) && (val<=validBootCfg[bootcfg][MAX_CFG_VAL_IDX]) ) {
		ret = 0;
	}
	return ret;
}

int validate_bootpart(int part)
{	
	int ret = 1;
	if( (part>=1) && (part<=2) ) {
		ret = 0;
	}
	return ret;
}

int at91_set_rootpartition(int offset)
{
	const int BOOTPART_LEN = 2;
	const int MMCIFACE_LEN = 4;
	const int BOOTDEV_LEN = 2;
	const int RAMARGS_LEN = 128;
	
	uint8_t bootpart=2;
	uint8_t bootdev=0;
	uint8_t expbootpart=2;
	uint8_t expbootdev=0;
	uint8_t upgradeFlag=0;
	uint8_t partFallbackFlag=0;
	uint8_t devFallbackFlag=0;

	char sbootpart[BOOTPART_LEN];
	char smmciface[MMCIFACE_LEN];
	char sbootdev[BOOTDEV_LEN];
	char sramargs[RAMARGS_LEN];

	const char *BOOTDEV_NAME = "bootdev";
	const char *BOOTPART_NAME = "bootpart";
	const char *MMCIFACE_NAME = "mmciface";
	const char *RAMARGS_NAME = "ramargs";
	
	struct udevice *dev;
	int ret;
	printf("Loading boot config eeprom.\n");
	//if (env_get(ROOT_PART_NAME))
	//	return 0;

	ret = uclass_first_device_err(UCLASS_I2C_EEPROM, &dev);
	if ( ret )
		return ret;
	
	printf("Reading boot device and root partition config from boot config eeprom.\n");
	
	ret = i2c_eeprom_read(dev, BOOTDEV_ADD, &bootdev, sizeof(uint8_t));	
	if ( ret ) {
		return ret;
	}

	ret = validate_bootcfg(BOOTDEV_ADD, bootdev);
	if ( ret ) {
		printf("No valid boot device (%d) configuration found. Setting default %d.\n",bootdev, CONFIG_SYS_BOOT_DEVICE);
		bootdev = 0;
		if( CONFIG_SYS_BOOT_DEVICE == 0 ) {
			bootdev = 0;
		} else {
			bootdev = 1;
		}
	}
	
		
	ret = i2c_eeprom_read(dev, EXPECTEDBOOTDEV_ADD, &expbootdev, sizeof(uint8_t));	
	if ( ret ) {
		return ret;
	}

	ret = validate_bootcfg(EXPECTEDBOOTDEV_ADD, expbootdev);
	if ( ret ) {
		printf("No valid expected boot device (%d) configuration found. Setting default 0.\n",expbootdev);
		expbootdev = 0;
	}
	
	if( bootdev != expbootdev ) {
		printf("Fallback on device has happened due to some failure.\n");
		devFallbackFlag = 1;
		ret = i2c_eeprom_write(dev, BOOTDEVFALLBACK_ADD, &devFallbackFlag, sizeof(uint8_t));
		if ( ret ) {
			printf("EEPROM write failure(%d).\n", ret); /* TODO: No eeprom wirte implemented yet.*/
			//return ret;
		}
	}
	
	ret = i2c_eeprom_read(dev, BOOTPART_ADD, &bootpart, sizeof(uint8_t));
	if ( ret ) {
		return ret;
	}

	ret = validate_bootcfg(BOOTPART_ADD, bootpart);
	if ( ret ) {
		printf("No valid boot partition configuration(%d) found. Setting default partition 2.\n",bootpart);
		bootpart = 2; 
	}

	ret = i2c_eeprom_read(dev, EXPECTEDBOOTPART_ADD, &expbootpart, sizeof(uint8_t));
	if ( ret ) {
		return ret;
	}

	ret = validate_bootcfg(EXPECTEDBOOTPART_ADD, expbootpart);
	if ( ret ) {
		printf("No valid previous boot partition configuration(%d) found. Setting default previous boot device 2.\n",expbootpart);
		expbootpart = 2; 
	}
	
	ret = i2c_eeprom_read(dev, UPGRADEFLAG_ADD, &upgradeFlag, sizeof(uint8_t));
        if (ret) {
                return ret;
        }
	
	ret = validate_bootcfg(UPGRADEFLAG_ADD,upgradeFlag);
	if( ret ) {
		printf("No valid upgradeFlag configuration(%d) found. Setting default upgradeFlag to 0.\n",upgradeFlag);
		upgradeFlag = 0;
	}

	if( bootpart != expbootpart ) {
		if( upgradeFlag != 1 ) {
			printf("Fallback on partition has happened due to some failure.\n");
			partFallbackFlag = 1;
			ret = i2c_eeprom_write(dev, BOOTPARTFALLBACK_ADD, &partFallbackFlag, sizeof(uint8_t));
			if ( ret ) {
				printf("EEPROM write failure(%d).\n", ret);
				//return ret;
			}
		}
	}

	ret = i2c_eeprom_read(dev, BOOTPART_ADD, &bootpart, sizeof(uint8_t));
	if ( ret ) {
		return ret;
	}

	ret = validate_bootcfg(BOOTPART_ADD, bootpart);
	if ( ret ) {
		printf("No valid boot partition configuration(%d) found. Setting default partition 2.\n",bootpart);
		bootpart = 2; 
	}
	
	printf("Boot device is %d and boot partition is %d.\n",bootdev, bootpart);
	sprintf(sbootdev,"%d",bootdev);
	sprintf(sbootpart,"%d",bootpart);
	sprintf(smmciface,"%d:%d",bootdev,bootpart);
	sprintf(sramargs,"setenv bootargs console=ttyS0,115200 earlyprintk root=/dev/mmcblk%dp%d rw rootwait",bootdev,bootpart);

			
	printf("setenv %s=%s.\n",BOOTDEV_NAME, sbootdev);
	env_set(BOOTDEV_NAME, sbootdev);
	
	printf("setenv %s=%s \n",BOOTPART_NAME, sbootpart);
	env_set(BOOTPART_NAME, sbootpart);
	
	
	printf("setenv %s=%s \n",MMCIFACE_NAME, smmciface);
	env_set(MMCIFACE_NAME, smmciface);

	printf("setenv ramargs as %s.\n", sramargs);
	env_set(RAMARGS_NAME, sramargs);
	
	return 0;
}


