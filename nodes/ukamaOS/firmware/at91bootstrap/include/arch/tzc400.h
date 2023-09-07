/* ----------------------------------------------------------------------------
 *         Microchip Technology AT91Bootstrap project
 * ----------------------------------------------------------------------------
 * Copyright (c) 2018, Microchip Technology Inc. and its subsidiaries
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Microchip's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY MICROCHIP "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL MICROCHIP BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

#ifndef TZC400_BASE
#error "Need to define TZC400_BASE"
#endif

#define TZC400_BUILDCONFIG				0x000

#define TZC_BUILD_CONFIG_REGIONS_MASK			(0x1F << 0)
#define TZC_BUILD_CONFIG_9REGIONS			(0x08 << 0)

#define TZC_BUILD_CONFIG_ADDRESS_WIDTH_MASK		(0x3F << 8)
#define TZC_BUILD_CONFIG_ADDRESS_WIDTH(x)		(((x) - 1) << 8)

#define TZC_BUILD_CONFIG_NO_OF_FILTERS_MASK		(0x03 << 24)
#define TZC_BUILD_CONFIG_NO_OF_FILTERS(x)		(((x) - 1) << 24)

#define TZC400_ACTION					0x004
#define TZC400_ACTION_REACTION_VALUE_MASK		(0x03 << 0)
#define TZC400_ACTION_REACTION_VALUE_INT_ERR		(0x03 << 0)

#define TZC400_GATE_KEEPER				0x008
#define TZC400_GATE_KEEPER_OPEN_ALL_FILTER		(0x0F << 0)

#define TZC400_SPECULATION_CTRL				0x00C
#define TZC400_SPECULATION_CTRL_READ_DIS		(0x01 << 0)
#define TZC400_SPECULATION_CTRL_WRITE_DIS		(0x01 << 1)

#define TZC400_INT_STATUS				0x010
#define TZC400_INT_CLEAR				0x014

#define TZC400_REGION_BASE_LOW(n)			(0x100 + 0x20 * (n))
#define TZC400_REGION_BASE_HIGH(n)			(0x104 + 0x20 * (n))
#define TZC400_REGION_TOP_LOW(n)			(0x108 + 0x20 * (n))
#define TZC400_REGION_TOP_HIGH(n)			(0x10C + 0x20 * (n))

#define TZC400_REGION_ATTRIBUTES(n)			(0x110 + 0x20 * (n))
#define TZC400_REGION_ATTRIBUTES_S_RD_EN		(0x01 << 30)
#define TZC400_REGION_ATTRIBUTES_S_WR_EN		(0x01 << 31)
#define TZC400_REGION_ATTRIBUTES_FILTER_EN_ALL		(0x0F << 0)

#define TZC400_REGION_ID_ACCESS(n)			(0x114 + 0x20 * (n))
#define TZC400_REGION_ID_ACCESS_NS_RD_EN		(0x01 << 0)
#define TZC400_REGION_ID_ACCESS_NS_WR_EN		(0x01 << 16)

#define write_tzc400(R, V)				writel(V, (R) + base)
#define read_tzc400(R)					readl((R) + base)
