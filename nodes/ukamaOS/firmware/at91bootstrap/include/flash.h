/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2013, Atmel Corporation
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Atmel's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY ATMEL "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL ATMEL BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
#ifndef __NORFLASH_H__
#define __NORFLASH_H__

#define AT91_NORFLASH_BASE	AT91C_EBI_CS0

/* SMC Chip select 0 timings */
#define AT91C_FLASH_NWE_SETUP           (4 << 0)
#define AT91C_FLASH_NCS_WR_SETUP        (2 << 8)
#define AT91C_FLASH_NRD_SETUP           (4 << 16)
#define AT91C_FLASH_NCS_RD_SETUP        (2 << 24)

#define AT91C_FLASH_NWE_PULSE           (8 << 0)
#define AT91C_FLASH_NCS_WR_PULSE        (10 << 8)
#define AT91C_FLASH_NRD_PULSE           (8 << 16)
#define AT91C_FLASH_NCS_RD_PULSE        (10 << 24)

#define AT91C_FLASH_NWE_CYCLE           (16 << 0)
#define AT91C_FLASH_NRD_CYCLE           (16 << 16)

int load_norflash(struct image_info *image);

#endif	/* #ifndef __NORFLASH_H__ */
