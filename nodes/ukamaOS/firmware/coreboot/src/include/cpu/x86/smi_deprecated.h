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

#ifndef __X86_SMI_DEPRECATED_H__
#define __X86_SMI_DEPRECATED_H__

#if CONFIG(PARALLEL_MP) || !CONFIG(HAVE_SMI_HANDLER)
/* Empty stubs for platforms without SMI handlers. */
static inline void smm_init(void) { }
static inline void smm_init_completion(void) { }
#else
void smm_init(void);
void smm_init_completion(void);
#endif

/* Entry from smmhandler.S. */
void smi_handler(u32 smm_revision);

#endif
