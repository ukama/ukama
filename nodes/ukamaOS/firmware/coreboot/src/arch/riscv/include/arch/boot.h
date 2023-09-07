/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Jonathan Neuschäfer
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

#ifndef ARCH_RISCV_INCLUDE_ARCH_BOOT_H
#define ARCH_RISCV_INCLUDE_ARCH_BOOT_H

#define RISCV_PAYLOAD_MODE_U 0
#define RISCV_PAYLOAD_MODE_S 1
#define RISCV_PAYLOAD_MODE_M 3

struct prog;
void run_payload(struct prog *prog, void *fdt, int payload_mode);
void run_payload_opensbi(struct prog *prog, void *fdt, struct prog *opensbi, int payload_mode);

void run_opensbi(const int hart_id,
		 const void *opensbi,
		 const void *fdt,
		 const void *payload,
		 const int payload_mode);
#endif
