/*
 * This file is part of the coreboot project.
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

#ifndef __EM64T_SAVE_STATE_H__
#define __EM64T_SAVE_STATE_H__

#include <types.h>
#include <cpu/x86/smm.h>

/* Intel Core 2 (EM64T) SMM State-Save Area
 * starts @ 0x7c00
 */
#define SMM_EM64T_ARCH_OFFSET 0x7c00
#define SMM_EM64T_SAVE_STATE_OFFSET \
	SMM_SAVE_STATE_BEGIN(SMM_EM64T_ARCH_OFFSET)
typedef struct {
	u8	reserved0[256];
	u8	reserved1[208];

	u32	gdtr_upper_base;
	u32	ldtr_upper_base;
	u32	idtr_upper_base;

	u8	reserved2[4];

	u64	io_rdi;
	u64	io_rip;
	u64	io_rcx;
	u64	io_rsi;
	u64	cr4;

	u8	reserved3[68];

	u64	gdtr_base;
	u64	idtr_base;
	u64	ldtr_base;

	u8	reserved4[84];

	u32	smm_revision;
	u32	smbase;

	u16	io_restart;
	u16	autohalt_restart;

	u8	reserved5[24];

	u64	r15;
	u64	r14;
	u64	r13;
	u64	r12;
	u64	r11;
	u64	r10;
	u64	r9;
	u64	r8;

	u64	rax;
	u64	rcx;
	u64	rdx;
	u64	rbx;

	u64	rsp;
	u64	rbp;
	u64	rsi;
	u64	rdi;


	u64	io_mem_addr;
	u32	io_misc_info;

	u32	es_sel;
	u32	cs_sel;
	u32	ss_sel;
	u32	ds_sel;
	u32	fs_sel;
	u32	gs_sel;

	u32	ldtr_sel;
	u32	tr_sel;

	u64	dr7;
	u64	dr6;
	u64	rip;
	u64	efer;
	u64	rflags;

	u64	cr3;
	u64	cr0;
} __packed em64t_smm_state_save_area_t;

#endif
