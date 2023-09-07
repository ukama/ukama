/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 Google LLC
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <boardid.h>
#include "commands.h"

uint32_t board_id(void)
{
	MAYBE_STATIC_NONZERO uint32_t id = BOARD_ID_INIT;

	if (id == BOARD_ID_INIT) {
		uint8_t ec_id;
		if (wilco_ec_get_board_id(&ec_id) <= 0)
			id = BOARD_ID_UNKNOWN;
		else
			id = ec_id;
	}

	return id;
}
