#
# Copyright (c) 2018-2019, Arm Limited. All rights reserved.
#
# SPDX-License-Identifier: BSD-3-Clause
#

include plat/arm/css/sgi/sgi-common.mk

RDE1EDGE_BASE		=	plat/arm/board/rde1edge

PLAT_INCLUDES		+=	-I${RDE1EDGE_BASE}/include/

SGI_CPU_SOURCES		:=	lib/cpus/aarch64/neoverse_e1.S

BL1_SOURCES		+=	${SGI_CPU_SOURCES}			\
				${RDE1EDGE_BASE}/rde1edge_err.c

BL2_SOURCES		+=	${RDE1EDGE_BASE}/rde1edge_plat.c	\
				${RDE1EDGE_BASE}/rde1edge_security.c	\
				${RDE1EDGE_BASE}/rde1edge_err.c		\
				drivers/arm/tzc/tzc_dmc620.c		\
				lib/utils/mem_region.c			\
				plat/arm/common/arm_nor_psci_mem_protect.c

BL31_SOURCES		+=	${SGI_CPU_SOURCES}			\
				${RDE1EDGE_BASE}/rde1edge_plat.c	\
				drivers/cfi/v2m/v2m_flash.c		\
				lib/utils/mem_region.c			\
				plat/arm/common/arm_nor_psci_mem_protect.c

# Add the FDT_SOURCES and options for Dynamic Config
FDT_SOURCES		+=	${RDE1EDGE_BASE}/fdts/${PLAT}_tb_fw_config.dts
TB_FW_CONFIG		:=	${BUILD_PLAT}/fdts/${PLAT}_tb_fw_config.dtb

# Add the TB_FW_CONFIG to FIP and specify the same to certtool
$(eval $(call TOOL_ADD_PAYLOAD,${TB_FW_CONFIG},--tb-fw-config))

FDT_SOURCES		+=	${RDE1EDGE_BASE}/fdts/${PLAT}_nt_fw_config.dts
NT_FW_CONFIG		:=	${BUILD_PLAT}/fdts/${PLAT}_nt_fw_config.dtb

# Add the NT_FW_CONFIG to FIP and specify the same to certtool
$(eval $(call TOOL_ADD_PAYLOAD,${NT_FW_CONFIG},--nt-fw-config))

override CTX_INCLUDE_AARCH32_REGS	:= 0
