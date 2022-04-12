/*
 * Video Decoder (VDEC) - System peripherals registers.
 *
 * Copyright (C) 2009  Hantro Products Oy.
 *
 * Based on SAMA5D4 datasheet.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */

#ifndef AT91_VDEC_H
#define AT91_VDEC_H

#define VDEC_IDR      0x00       /* ID Register (read-only) */
#define   VDEC_IDR_BUILD_VER              0xf /* Build Version is 0x02. */
#define   VDEC_IDR_MINOR_VER      (0xff << 4) /* Minor Version is 0x88. */
#define   VDEC_IDR_MAJOR_VER      (0xf << 12) /* Major Version is 0x08. */
#define   VDEC_IDR_PROD_ID     (0xffff << 16) /* Product ID is 0x6731. */

#define VDEC_DIR      0x04       /* Decoder Interrupt Register */
#define   VDEC_DIR_DE                       1 /* 1: Enable decoder; 0: Disable decoder. */
#define   VDEC_DIR_ID                    0x10 /* 1: Disable interrupts for decoder; 0: Enable interrupts. */
#define   VDEC_DIR_ABORT                 0x20
#define   VDEC_DIR_ISET                 0x100 /* Decoder Interrupt Set. 0: Clears the Decoder Interrupt. */

#define VDEC_PPIR     0xF0       /* Post Processor Interrupt Register */
#define   VDEC_PPIR_PPE                     1 /* 1: Enable post-processor; 0: Disable post-processor */
#define   VDEC_PPIR_ID                   0x10 /* 1: Disable interrupts for post-processor; 0: Enable interrupts. */
#define   VDEC_PPIR_ISET                0x100 /* Post-processor Interrupt Set. 0: Clears the post-processor Interrupt. */

#endif
