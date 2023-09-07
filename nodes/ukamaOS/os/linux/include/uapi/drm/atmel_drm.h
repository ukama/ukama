/*
 * Copyright (C) 2014 Atmel
 *
 * Author: Mohamed Jamsheeth <mohamedjamsheeth.hajanajubudeen@atmel.com>
 *
 * This program is free software; you can redistribute it and/or modify it
 * under the terms of the GNU General Public License version 2 as published by
 * the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
 * more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef _UAPI_ATMEL_DRM_H_
#define _UAPI_ATMEL_DRM_H_

#include <drm/drm.h>

#if defined(__cplusplus)
extern "C" {
#endif

struct drm_gfx2d_submit {
	__u32 flags;
	__u32 size;
	__u32 buf;
};

struct drm_gfx2d_gem_addr {
	__u32 name;
	__u32 paddr;
	__u32 size;
};

#define DRM_ATMEL_GEM_GET		0x00
#define DRM_GFX2D_SUBMIT		0x01
#define DRM_GFX2D_FLUSH			0x02
#define DRM_GFX2D_GEM_ADDR		0x03

#define DRM_IOCTL_ATMEL_GEM_GET		DRM_IOWR(DRM_COMMAND_BASE +	\
						 DRM_ATMEL_GEM_GET, struct drm_mode_map_dumb)
#define DRM_IOCTL_GFX2D_SUBMIT		DRM_IOWR(DRM_COMMAND_BASE +	\
						 DRM_GFX2D_SUBMIT, struct drm_gfx2d_submit)
#define DRM_IOCTL_GFX2D_FLUSH		DRM_IO(DRM_COMMAND_BASE + DRM_GFX2D_FLUSH)
#define DRM_IOCTL_GFX2D_GEM_ADDR	DRM_IOWR(DRM_COMMAND_BASE +	\
						 DRM_GFX2D_GEM_ADDR, struct drm_gfx2d_gem_addr)

#if defined(__cplusplus)
}
#endif

#endif /* _UAPI_ATMEL_DRM_H_ */
