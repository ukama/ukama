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

#include <stdint.h>
#include <cpu/amd/microcode.h>

struct id_mapping {
	uint32_t orig_id;
	uint16_t new_id;
};

static u16 get_equivalent_processor_rev_id(u32 orig_id)
{
	static const struct id_mapping id_mapping_table[] = {
		/* Family 16h */

	/* TODO This equivalent processor revisions ID needs verification */
		{ 0x730f01, 0x7301 },

		/* Array terminator */
		{ 0xffffff, 0x0000 },
	};

	u32 new_id;
	int i;

	new_id = 0;

	for (i = 0; id_mapping_table[i].orig_id != 0xffffff; i++) {
		if (id_mapping_table[i].orig_id == orig_id) {
			new_id = id_mapping_table[i].new_id;
			break;
		}
	}

	return new_id;
}

void update_microcode(u32 cpu_deviceid)
{
	u32 equivalent_processor_rev_id =
				get_equivalent_processor_rev_id(cpu_deviceid);
	amd_update_microcode_from_cbfs(equivalent_processor_rev_id);
}
