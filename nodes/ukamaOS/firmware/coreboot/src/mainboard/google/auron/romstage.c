/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007-2010 coresystems GmbH
 * Copyright (C) 2014 Google Inc.
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

#include <console/console.h>
#include <ec/google/chromeec/ec.h>
#include <soc/pei_data.h>
#include <soc/pei_wrapper.h>
#include <soc/romstage.h>
#include <variant/spd.h>
#include "variant.h"

__weak void variant_romstage_entry(struct romstage_params *rp)
{
}

void mainboard_pre_raminit(struct romstage_params *rp)
{
	/* Fill out PEI DATA */
	mainboard_fill_pei_data(&rp->pei_data);
	mainboard_fill_spd_data(&rp->pei_data);

}

void mainboard_post_raminit(struct romstage_params *rp)
{
	/* Do variant-specific init */
	variant_romstage_entry(rp);
}
