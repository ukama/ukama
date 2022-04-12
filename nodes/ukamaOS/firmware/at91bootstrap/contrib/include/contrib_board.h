/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2015, Atmel Corporation
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

#ifndef __CONTRIB_BOARD_H__
#define __CONTRIB_BOARD_H__

#ifdef CONFIG_VINCO
#include "vinco.h"
#endif

#ifdef CONFIG_AT91SAM9X5_ARIA
#include "at91sam9x5_aria.h"
#endif

#ifdef CONFIG_AT91SAM9X5_ARIETTA
#include "at91sam9x5_arietta.h"
#endif

#ifdef CONFIG_PICOSAM9G45
#include "pico_sam9g45.h"
#endif

#ifdef CONFIG_SAMA5D3_ACQUA
#include "sama5d3_acqua.h"
#endif

#ifdef CONFIG_SAMA5D3_STEPHAN
#include "sama5d3_stephan.h"
#endif

#ifdef CONFIG_SAMA5D2_ROADRUNNER
#include "sama5d2_roadrunner.h"
#endif

#ifdef CONFIG_CORE9G25
#include "core9g25.h"
#endif

#ifdef CONFIG_SAMA5D3_LINEA
#include "sama5d3_linea.h"
#endif

#endif
