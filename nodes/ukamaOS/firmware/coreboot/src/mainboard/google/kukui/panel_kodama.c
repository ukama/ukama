/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 Bitland Tech Inc.
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

#include "panel.h"

static struct panel_description kodama_panels[] = {
	[1] = { .name = "AUO_B101UAN08_3", },
	[2] = { .name = "BOE_TV101WUM_N53", },
};

struct panel_description *get_panel_description(int panel_id)
{
	if (panel_id < 0 || panel_id >= ARRAY_SIZE(kodama_panels))
		return NULL;

	return get_panel_from_cbfs(&kodama_panels[panel_id]);
}
