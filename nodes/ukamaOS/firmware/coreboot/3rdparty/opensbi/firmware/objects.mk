#
# SPDX-License-Identifier: BSD-2-Clause
#
# Copyright (c) 2019 Western Digital Corporation or its affiliates.
#
# Authors:
#   Anup Patel <anup.patel@wdc.com>
#

firmware-genflags-y =
firmware-cppflags-y +=
firmware-cflags-y +=
firmware-asflags-y +=
firmware-ldflags-y +=

ifdef FW_TEXT_START
firmware-genflags-y += -DFW_TEXT_START=$(FW_TEXT_START)
endif

firmware-bins-$(FW_DYNAMIC) += fw_dynamic.bin

firmware-bins-$(FW_JUMP) += fw_jump.bin
ifdef FW_JUMP_ADDR
firmware-genflags-$(FW_JUMP) += -DFW_JUMP_ADDR=$(FW_JUMP_ADDR)
endif
ifdef FW_JUMP_FDT_ADDR
firmware-genflags-$(FW_JUMP) += -DFW_JUMP_FDT_ADDR=$(FW_JUMP_FDT_ADDR)
endif

firmware-bins-$(FW_PAYLOAD) += fw_payload.bin
ifdef FW_PAYLOAD_PATH
FW_PAYLOAD_PATH_FINAL=$(FW_PAYLOAD_PATH)
else
FW_PAYLOAD_PATH_FINAL=$(platform_build_dir)/firmware/payloads/test.bin
endif
firmware-genflags-$(FW_PAYLOAD) += -DFW_PAYLOAD_PATH=\"$(FW_PAYLOAD_PATH_FINAL)\"
ifdef FW_PAYLOAD_OFFSET
firmware-genflags-$(FW_PAYLOAD) += -DFW_PAYLOAD_OFFSET=$(FW_PAYLOAD_OFFSET)
endif
ifdef FW_PAYLOAD_ALIGN
firmware-genflags-$(FW_PAYLOAD) += -DFW_PAYLOAD_ALIGN=$(FW_PAYLOAD_ALIGN)
endif

ifndef FW_PAYLOAD_FDT_PATH
ifdef FW_PAYLOAD_FDT
FW_PAYLOAD_FDT_PATH=$(platform_build_dir)/$(FW_PAYLOAD_FDT)
endif
endif
ifdef FW_PAYLOAD_FDT_PATH
firmware-genflags-$(FW_PAYLOAD) += -DFW_PAYLOAD_FDT_PATH=\"$(FW_PAYLOAD_FDT_PATH)\"
endif
ifdef FW_PAYLOAD_FDT_ADDR
firmware-genflags-$(FW_PAYLOAD) += -DFW_PAYLOAD_FDT_ADDR=$(FW_PAYLOAD_FDT_ADDR)
endif

ifdef FW_OPTIONS
firmware-genflags-y += -DFW_OPTIONS=$(FW_OPTIONS)
endif
