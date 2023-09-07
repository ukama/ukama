/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2011,2012 Advanced Micro Devices, Inc.
 * Copyright (C) 2013 Sage Electronic Engineering, LLC
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#ifndef __CALLOUTS_AMD_AGESA_H__
#define __CALLOUTS_AMD_AGESA_H__

#include <amdblocks/agesawrapper.h>
#include <stdint.h>

#define BIOS_HEAP_SIZE			0x30000
#define BSP_STACK_BASE_ADDR		0x30000

typedef struct _BIOS_HEAP_MANAGER {
	uint32_t StartOfAllocatedNodes;
	uint32_t StartOfFreedNodes;
} BIOS_HEAP_MANAGER;

typedef struct _BIOS_BUFFER_NODE {
	uint32_t BufferHandle;
	uint32_t BufferSize;
	uint32_t NextNodeOffset;
} BIOS_BUFFER_NODE;

AGESA_STATUS agesa_GetTempHeapBase(uint32_t Func, uintptr_t Data,
							void *ConfigPtr);
AGESA_STATUS agesa_HeapRebase(uint32_t Func, uintptr_t Data, void *ConfigPtr);

AGESA_STATUS agesa_AllocateBuffer(uint32_t Func, uintptr_t Data,
							void *ConfigPtr);
AGESA_STATUS agesa_DeallocateBuffer(uint32_t Func, uintptr_t Data,
							void *ConfigPtr);
AGESA_STATUS agesa_LocateBuffer(uint32_t Func, uintptr_t Data, void *ConfigPtr);

AGESA_STATUS agesa_NoopUnsupported(uint32_t Func, uintptr_t Data,
							void *ConfigPtr);
AGESA_STATUS agesa_NoopSuccess(uint32_t Func, uintptr_t Data, void *ConfigPtr);
AGESA_STATUS agesa_EmptyIdsInitData(uint32_t Func, uintptr_t Data,
							void *ConfigPtr);
AGESA_STATUS agesa_Reset(uint32_t Func, uintptr_t Data, void *ConfigPtr);
AGESA_STATUS agesa_RunFuncOnAp(uint32_t Func, uintptr_t Data, void *ConfigPtr);
AGESA_STATUS agesa_GfxGetVbiosImage(uint32_t Func, uintptr_t FchData,
							void *ConfigPrt);

AGESA_STATUS agesa_ReadSpd(uint32_t Func, uintptr_t Data, void *ConfigPtr);
AGESA_STATUS agesa_RunFcnOnAllAps(uint32_t Func, uintptr_t Data,
							void *ConfigPtr);
AGESA_STATUS agesa_PcieSlotResetControl(uint32_t Func, uintptr_t Data,
							void *ConfigPtr);
AGESA_STATUS agesa_WaitForAllApsFinished(uint32_t Func, uintptr_t Data,
							void *ConfigPtr);
AGESA_STATUS agesa_IdleAnAp(uint32_t Func, uintptr_t Data, void *ConfigPtr);

AGESA_STATUS GetBiosCallout(uint32_t Func, uintptr_t Data, void *ConfigPtr);

AGESA_STATUS agesa_fch_initreset(uint32_t Func, uintptr_t FchData,
							void *ConfigPtr);
AGESA_STATUS agesa_fch_initenv(uint32_t Func, uintptr_t FchData,
							void *ConfigPtr);
AGESA_STATUS agesa_HaltThisAp(uint32_t Func, uintptr_t Data, void *ConfigPtr);

void platform_FchParams_reset(FCH_RESET_DATA_BLOCK *FchParams_reset);
void platform_FchParams_env(FCH_DATA_BLOCK *FchParams_env);
AGESA_STATUS platform_PcieSlotResetControl(uint32_t Func, uintptr_t Data,
	void *ConfigPtr);
typedef struct {
	uint32_t CalloutName;
	CALLOUT_ENTRY CalloutPtr;
} BIOS_CALLOUT_STRUCT;

extern const BIOS_CALLOUT_STRUCT BiosCallouts[];
extern const int BiosCalloutsLen;

#endif /* __CALLOUTS_AMD_AGESA_H__ */
