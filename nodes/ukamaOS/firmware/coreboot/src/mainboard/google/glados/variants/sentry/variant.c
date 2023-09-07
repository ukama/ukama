/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2014 Google Inc.
 * Copyright (C) 2015 Intel Corporation.
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

#include <stdint.h>
#include <string.h>
#include <baseboard/variant.h>
#include <fsp/soc_binding.h>

#define K4E6E304EE_MEM_ID	0x3

void variant_memory_init_params(FSPM_UPD *mupd, const int spd_index)
{
	FSP_M_CONFIG *mem_cfg;
	mem_cfg = &mupd->FspmConfig;

	/* DQ byte map */
	const u8 dq_map[2][12] = {
		  { 0x0F, 0xF0, 0x00, 0xF0, 0x0F, 0xF0,
		    0x0F, 0x00, 0xFF, 0x00, 0xFF, 0x00 },
		  { 0x0F, 0xF0, 0x00, 0xF0, 0x0F, 0xF0,
		    0x0F, 0x00, 0xFF, 0x00, 0xFF, 0x00 } };
	/* DQS CPU<>DRAM map */
	const u8 dqs_map[2][8] = {
		{ 0, 1, 3, 2, 6, 5, 4, 7 },
		{ 2, 3, 0, 1, 6, 7, 4, 5 } };

	/* Rcomp resistor */
	const u16 RcompResistor[3] = { 200, 81, 162 };

	/* Rcomp target */
	const u16 RcompTarget[5] = { 100, 40, 40, 23, 40 };

	/* Strengthen the Rcomp Target Ctrl for K4E6E304EE_MEM_ID */
	const u16 StrengthendRcompTarget[5] = { 100, 40, 40, 21, 40 };

	/* Default Rcomp Target assignment */
	const u16 *targeted_rcomp = RcompTarget;

	/* Override Rcomp Target assignment for specific SKU(s) */
	if (spd_index == K4E6E304EE_MEM_ID)
		targeted_rcomp = StrengthendRcompTarget;

	memcpy(mem_cfg->DqByteMapCh0, dq_map[0], sizeof(mem_cfg->DqByteMapCh0));
	memcpy(mem_cfg->DqByteMapCh1, dq_map[1], sizeof(mem_cfg->DqByteMapCh1));
	memcpy(mem_cfg->DqsMapCpu2DramCh0, dqs_map[0], sizeof(mem_cfg->DqsMapCpu2DramCh0));
	memcpy(mem_cfg->DqsMapCpu2DramCh1, dqs_map[1], sizeof(mem_cfg->DqsMapCpu2DramCh1));
	memcpy(mem_cfg->RcompResistor, RcompResistor, sizeof(mem_cfg->RcompResistor));
	memcpy(mem_cfg->RcompTarget, RcompTarget, sizeof(mem_cfg->RcompTarget));
}
