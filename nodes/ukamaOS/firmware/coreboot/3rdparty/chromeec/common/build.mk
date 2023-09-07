# -*- makefile -*-
# Copyright 2014 The Chromium OS Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
#
# Common files build
#

# Note that this variable includes the trailing "/"
_common_dir:=$(dir $(lastword $(MAKEFILE_LIST)))

common-y=util.o
common-y+=version.o printf.o queue.o queue_policies.o

common-$(CONFIG_ACCELGYRO_BMA255)+=math_util.o
common-$(CONFIG_ACCELGYRO_BMI160)+=math_util.o
common-$(CONFIG_ACCELGYRO_LSM6DS0)+=math_util.o
common-$(CONFIG_ACCELGYRO_LSM6DSM)+=math_util.o
common-$(CONFIG_ACCELGYRO_LSM6DSO)+=math_util.o
common-$(CONFIG_ACCEL_FIFO)+=motion_sense_fifo.o
common-$(CONFIG_ACCEL_LIS2DW12)+=math_util.o
common-$(CONFIG_ACCEL_LIS2DH)+=math_util.o
common-$(CONFIG_ACCEL_KXCJ9)+=math_util.o
common-$(CONFIG_ACCEL_KX022)+=math_util.o
ifneq ($(CORE),cortex-m)
common-$(CONFIG_AES)+=aes.o
endif
common-$(CONFIG_AES_GCM)+=aes-gcm.o
common-$(CONFIG_CMD_ADC)+=adc.o
common-$(HAS_TASK_ALS)+=als.o
common-$(CONFIG_AP_HANG_DETECT)+=ap_hang_detect.o
common-$(CONFIG_AUDIO_CODEC)+=audio_codec.o
common-$(CONFIG_AUDIO_CODEC_DMIC)+=audio_codec_dmic.o
common-$(CONFIG_AUDIO_CODEC_I2S_RX)+=audio_codec_i2s_rx.o
common-$(CONFIG_AUDIO_CODEC_WOV)+=audio_codec_wov.o
common-$(CONFIG_BACKLIGHT_LID)+=backlight_lid.o
common-$(CONFIG_BASE32)+=base32.o
common-$(CONFIG_DETACHABLE_BASE)+=base_state.o
common-$(CONFIG_BATTERY)+=battery.o
common-$(CONFIG_BATTERY_FUEL_GAUGE)+=battery_fuel_gauge.o
common-$(CONFIG_BLUETOOTH_LE)+=bluetooth_le.o
common-$(CONFIG_BLUETOOTH_LE_STACK)+=btle_hci_controller.o btle_ll.o
common-$(CONFIG_CAPSENSE)+=capsense.o
common-$(CONFIG_CASE_CLOSED_DEBUG_V1)+=ccd_config.o
common-$(CONFIG_CEC)+=cec.o
common-$(CONFIG_CROS_BOARD_INFO)+=cbi.o
common-$(CONFIG_CHARGE_MANAGER)+=charge_manager.o
common-$(CONFIG_CHARGE_RAMP_HW)+=charge_ramp.o
common-$(CONFIG_CHARGE_RAMP_SW)+=charge_ramp.o charge_ramp_sw.o
common-$(CONFIG_CMD_CHARGEN) += chargen.o
common-$(CONFIG_CHARGER)+=charger.o charge_state_v2.o
common-$(CONFIG_CHARGER_PROFILE_OVERRIDE_COMMON)+=charger_profile_override.o
common-$(CONFIG_CMD_I2CWEDGE)+=i2c_wedge.o
common-$(CONFIG_COMMON_GPIO)+=gpio.o gpio_commands.o
common-$(CONFIG_IO_EXPANDER)+=ioexpander.o
common-$(CONFIG_COMMON_PANIC_OUTPUT)+=panic_output.o
common-$(CONFIG_COMMON_RUNTIME)+=hooks.o main.o system.o peripheral.o
common-$(CONFIG_COMMON_TIMER)+=timer.o
common-$(CONFIG_CRC8)+= crc8.o
common-$(CONFIG_CURVE25519)+=curve25519.o
ifneq ($(CORE),cortex-m0)
common-$(CONFIG_CURVE25519)+=curve25519-generic.o
endif
common-$(CONFIG_DEDICATED_RECOVERY_BUTTON)+=button.o
common-$(CONFIG_DEVICE_EVENT)+=device_event.o
common-$(CONFIG_DEVICE_STATE)+=device_state.o
common-$(CONFIG_DPTF)+=dptf.o
common-$(CONFIG_EC_EC_COMM_MASTER)+=ec_ec_comm_master.o
common-$(CONFIG_EC_EC_COMM_SLAVE)+=ec_ec_comm_slave.o
common-$(CONFIG_HOSTCMD_ESPI)+=espi.o
common-$(CONFIG_EXTENSION_COMMAND)+=extension.o
common-$(CONFIG_EXTPOWER_GPIO)+=extpower_gpio.o
common-$(CONFIG_FANS)+=fan.o pwm.o
common-$(CONFIG_FACTORY_MODE)+=factory_mode.o
common-$(CONFIG_FLASH)+=flash.o
common-$(CONFIG_FLASH_LOG)+=flash_log.o flash_log_vc.o
common-$(CONFIG_FLASH_NVMEM)+=nvmem.o
common-$(CONFIG_FLASH_NVMEM)+=new_nvmem.o
common-$(CONFIG_FLASH_NVMEM_VARS)+=nvmem_vars.o
common-$(CONFIG_FMAP)+=fmap.o
common-$(CONFIG_GESTURE_SW_DETECTION)+=gesture.o
common-$(CONFIG_HOSTCMD_EVENTS)+=host_event_commands.o
common-$(CONFIG_HOSTCMD_GET_UPTIME_INFO)+=uptime.o
common-$(CONFIG_HOSTCMD_PD)+=host_command_master.o
common-$(CONFIG_HOSTCMD_RTC)+=rtc.o
common-$(CONFIG_I2C_DEBUG)+=i2c_trace.o
common-$(CONFIG_I2C_MASTER)+=i2c_master.o
common-$(CONFIG_I2C_SLAVE)+=i2c_slave.o
common-$(CONFIG_I2C_VIRTUAL_BATTERY)+=virtual_battery.o
common-$(CONFIG_INDUCTIVE_CHARGING)+=inductive_charging.o
common-$(CONFIG_KEYBOARD_PROTOCOL_8042)+=keyboard_8042.o \
	keyboard_8042_sharedlib.o
common-$(CONFIG_KEYBOARD_PROTOCOL_MKBP)+=keyboard_mkbp.o
common-$(CONFIG_KEYBOARD_TEST)+=keyboard_test.o
common-$(CONFIG_LED_COMMON)+=led_common.o
common-$(CONFIG_LED_POLICY_STD)+=led_policy_std.o
common-$(CONFIG_LED_PWM)+=led_pwm.o
common-$(CONFIG_LED_ONOFF_STATES)+=led_onoff_states.o
common-$(CONFIG_LID_ANGLE)+=motion_lid.o math_util.o
common-$(CONFIG_LID_ANGLE_UPDATE)+=lid_angle.o
common-$(CONFIG_LID_SWITCH)+=lid_switch.o
common-$(CONFIG_HOSTCMD_X86)+=acpi.o port80.o ec_features.o
common-$(CONFIG_MAG_CALIBRATE)+= mag_cal.o math_util.o vec3.o mat33.o mat44.o
common-$(CONFIG_MKBP_EVENT)+=mkbp_event.o
common-$(CONFIG_ONEWIRE)+=onewire.o
common-$(CONFIG_PECI_COMMON)+=peci.o
common-$(CONFIG_PHYSICAL_PRESENCE)+=physical_presence.o
common-$(CONFIG_PINWEAVER)+=pinweaver.o
common-$(CONFIG_POWER_BUTTON)+=power_button.o
common-$(CONFIG_POWER_BUTTON_X86)+=power_button_x86.o
common-$(CONFIG_PSTORE)+=pstore_commands.o
common-$(CONFIG_PWM)+=pwm.o
common-$(CONFIG_PWM_KBLIGHT)+=pwm_kblight.o
common-$(CONFIG_KEYBOARD_BACKLIGHT)+=keyboard_backlight.o
common-$(CONFIG_RMA_AUTH)+=rma_auth.o
common-$(CONFIG_RSA)+=rsa.o
common-$(CONFIG_ROLLBACK)+=rollback.o
common-$(CONFIG_RWSIG)+=rwsig.o vboot/common.o
common-$(CONFIG_RWSIG_TYPE_RWSIG)+=vboot/vb21_lib.o
common-$(CONFIG_MATH_UTIL)+=math_util.o
common-$(CONFIG_SHA1)+= sha1.o
common-$(CONFIG_SHA256)+=sha256.o
common-$(CONFIG_SOFTWARE_CLZ)+=clz.o
common-$(CONFIG_SOFTWARE_CTZ)+=ctz.o
common-$(CONFIG_CMD_SPI_XFER)+=spi_commands.o
common-$(CONFIG_SPI_FLASH)+=spi_flash.o spi_flash_reg.o
common-$(CONFIG_SPI_FLASH_REGS)+=spi_flash_reg.o
common-$(CONFIG_SPI_NOR)+=spi_nor.o
common-$(CONFIG_SWITCH)+=switch.o
common-$(CONFIG_SW_CRC)+=crc.o
common-$(CONFIG_TABLET_MODE)+=tablet_mode.o
common-$(CONFIG_TEMP_SENSOR)+=temp_sensor.o
common-$(CONFIG_THROTTLE_AP)+=thermal.o throttle_ap.o
common-$(CONFIG_THROTTLE_AP_ON_BAT_DISCHG_CURRENT)+=throttle_ap.o
common-$(CONFIG_THROTTLE_AP_ON_BAT_VOLTAGE)+=throttle_ap.o
common-$(CONFIG_TPM_I2CS)+=i2cs_tpm.o
common-$(CONFIG_U2F)+=u2f.o
common-$(CONFIG_USB_CHARGER)+=usb_charger.o
common-$(CONFIG_USB_CONSOLE_STREAM)+=usb_console_stream.o
common-$(CONFIG_USB_I2C)+=usb_i2c.o
common-$(CONFIG_USB_PORT_POWER_DUMB)+=usb_port_power_dumb.o
common-$(CONFIG_USB_PORT_POWER_SMART)+=usb_port_power_smart.o
common-$(CONFIG_USB_POWER_DELIVERY)+=usb_common.o
ifeq ($(CONFIG_USB_SM_FRAMEWORK),)
common-$(CONFIG_USB_POWER_DELIVERY)+=usb_pd_protocol.o usb_pd_policy.o
endif
common-$(CONFIG_USB_PD_LOGGING)+=event_log.o pd_log.o
common-$(CONFIG_USB_PD_TCPC)+=usb_pd_tcpc.o
common-$(CONFIG_USB_UPDATE)+=usb_update.o update_fw.o
common-$(CONFIG_USBC_PPC)+=usbc_ppc.o
common-$(CONFIG_VBOOT_EFS)+=vboot/vboot.o
common-$(CONFIG_VBOOT_HASH)+=sha256.o vboot_hash.o
common-$(CONFIG_VOLUME_BUTTONS)+=button.o
common-$(CONFIG_VSTORE)+=vstore.o
common-$(CONFIG_WEBUSB_URL)+=webusb_desc.o
common-$(CONFIG_WIRELESS)+=wireless.o
common-$(HAS_TASK_CHIPSET)+=chipset.o
common-$(HAS_TASK_CONSOLE)+=console.o console_output.o uart_buffering.o
common-$(CONFIG_CMD_MEM)+=memory_commands.o
common-$(HAS_TASK_HOSTCMD)+=host_command.o ec_features.o
common-$(HAS_TASK_PDCMD)+=host_command_pd.o
common-$(HAS_TASK_KEYSCAN)+=keyboard_scan.o
common-$(HAS_TASK_LIGHTBAR)+=lb_common.o lightbar.o
common-$(HAS_TASK_MOTIONSENSE)+=motion_sense.o
common-$(HAS_TASK_TPM)+=tpm_registers.o

ifneq ($(HAVE_PRIVATE_AUDIO_CODEC_WOV_LIBS),y)
common-$(CONFIG_AUDIO_CODEC_WOV)+=hotword_dsp_api.o
endif

ifneq ($(CONFIG_COMMON_RUNTIME),)
common-$(CONFIG_MALLOC)+=shmalloc.o
common-$(call not_cfg,$(CONFIG_MALLOC))+=shared_mem.o
endif

ifeq ($(CTS_MODULE),)
common-$(TEST_BUILD)+=test_util.o
else
common-y+=test_util.o
endif

ifneq ($(CONFIG_RSA_OPTIMIZED),)
$(out)/RW/common/rsa.o: CFLAGS+=-O3
$(out)/RO/common/rsa.o: CFLAGS+=-O3
endif

# AES-GCM code needs C99, else we'd have to move many variables declarations
# around.
$(out)/RW/common/aes-gcm.o: CFLAGS+=-std=c99 -Wno-declaration-after-statement
$(out)/RO/common/aes-gcm.o: CFLAGS+=-std=c99 -Wno-declaration-after-statement

ifneq ($(CONFIG_BOOTBLOCK),)
build-util-bin += gen_emmc_transfer_data

# Bootblock is only packed in RO image.
$(out)/util/gen_emmc_transfer_data: BUILD_LDFLAGS += -DSECTION_IS_RO=$(EMPTY)
$(out)/bootblock_data.h: $(out)/util/gen_emmc_transfer_data $(out)/.bootblock
	$(call quiet,emmc_bootblock,BTBLK  )

# We only want to repack the bootblock if: $(BOOTBLOCK) variable value has
# changed, or the file pointed at by $(BOOTBLOCK) has changed. We do this
# by recording the latest $(BOOTBLOCK) file information in .bootblock
# TODO: Need a better makefile tricks to do this.

bootblock_ls := $(shell ls -l "$(BOOTBLOCK)" 2>&1)
old_bootblock_ls := $(shell cat $(out)/.bootblock 2>/dev/null)

$(out)/.bootblock: $(BOOTBLOCK)
	@echo "$(bootblock_ls)" > $@

ifneq ($(bootblock_ls),$(old_bootblock_ls))
.PHONY: $(out)/.bootblock
endif
endif # CONFIG_BOOTBLOCK

ifneq ($(CONFIG_TOUCHPAD_HASH_FW),)
$(out)/RO/common/update_fw.o: $(out)/touchpad_fw_hash.h
$(out)/RW/common/update_fw.o: $(out)/touchpad_fw_hash.h

$(out)/touchpad_fw_hash.h: $(out)/util/gen_touchpad_hash $(out)/.touchpad_fw
	$(call quiet,tp_hash,TPHASH )

# We only want to recompute the hash if: $(TOUCHPAD_FW) variable value has
# changed, or the file pointed at by $(TOUCHPAD_FW) has changed. We do this
# by recording the latest $(TOUCHPAD_FW) file information in .touchpad_fw.

touchpad_fw_ls := $(shell ls -l "$(TOUCHPAD_FW)" 2>&1)
old_touchpad_fw_ls := $(shell cat $(out)/.touchpad_fw 2>/dev/null)

$(out)/.touchpad_fw: $(TOUCHPAD_FW)
	@echo "$(touchpad_fw_ls)" > $@

ifneq ($(touchpad_fw_ls),$(old_touchpad_fw_ls))
.PHONY: $(out)/.touchpad_fw
endif
endif

ifeq ($(TEST_BUILD),)

ifeq ($(CONFIG_RMA_AUTH_USE_P256),)
BLOB_FILE = rma_key_blob.x25519.test
else
BLOB_FILE = rma_key_blob.p256.test
endif

$(out)/RW/common/rma_auth.o: $(out)/rma_key_from_blob.h

$(out)/rma_key_from_blob.h: board/$(BOARD)/$(BLOB_FILE) util/bin2h.sh
	$(Q)util/bin2h.sh RMA_KEY_BLOB $< $@

endif

ifeq ($(CONFIG_LIBCRYPTOC),y)
CRYPTOCLIB := $(realpath ../../third_party/cryptoc)
ifneq ($(BOARD),host)
CPPFLAGS += -I$(abspath ./builtin)
endif
CPPFLAGS += -I$(CRYPTOCLIB)/include
CRYPTOC_LDFLAGS := -L$(out)/cryptoc -lcryptoc

# Force the external build each time, so it can look for changed sources.
.PHONY: $(out)/cryptoc/libcryptoc.a
$(out)/cryptoc/libcryptoc.a:
	$(MAKE) obj=$(realpath $(out))/cryptoc SUPPORT_UNALIGNED=1 \
		CONFIG_UPTO_SHA512=$(CONFIG_UPTO_SHA512) -C $(CRYPTOCLIB)

# Link RO and RW against cryptoc.
$(out)/RO/ec.RO.elf $(out)/RO/ec.RO_B.elf: LDFLAGS_EXTRA += $(CRYPTOC_LDFLAGS)
$(out)/RO/ec.RO.elf $(out)/RO/ec.RO_B.elf: $(out)/cryptoc/libcryptoc.a
$(out)/RW/ec.RW.elf $(out)/RW/ec.RW_B.elf: LDFLAGS_EXTRA += $(CRYPTOC_LDFLAGS)
$(out)/RW/ec.RW.elf $(out)/RW/ec.RW_B.elf: $(out)/cryptoc/libcryptoc.a
# Host test executables (including fuzz tests).
$(out)/$(PROJECT).exe: LDFLAGS_EXTRA += $(CRYPTOC_LDFLAGS)
$(out)/$(PROJECT).exe: $(out)/cryptoc/libcryptoc.a
endif

include $(_common_dir)fpsensor/build.mk
include $(_common_dir)usbc/build.mk

include $(_common_dir)mock/build.mk
common-y+=$(foreach m,$(mock-y),mock/$(m))
