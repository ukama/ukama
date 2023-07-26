/* Formatted output to UTF-8 strings.
   Copyright (C) 1999, 2002, 2006-2007, 2009-2022 Free Software Foundation,
   Inc.

   This file is free software.
   It is dual-licensed under "the GNU LGPLv3+ or the GNU GPLv2+".
   You can redistribute it and/or modify it under either
     - the terms of the GNU Lesser General Public License as published
       by the Free Software Foundation, either version 3, or (at your
       option) any later version, or
     - the terms of the GNU General Public License as published by the
       Free Software Foundation; either version 2, or (at your option)
       any later version, or
     - the same dual license "the GNU LGPLv3+ or the GNU GPLv2+".

   This file is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
   Lesser General Public License and the GNU General Public License
   for more details.

   You should have received a copy of the GNU Lesser General Public
   License and of the GNU General Public License along with this
   program.  If not, see <https://www.gnu.org/licenses/>.  */

/* Tell glibc's <stdio.h> to provide a prototype for snprintf().
   This must come before <config.h> because <config.h> may include
   <features.h>, and once <features.h> has been included, it's too late.  */
#ifndef _GNU_SOURCE
# define _GNU_SOURCE    1
#endif

#include <config.h>
#include <alloca.h>

/* Specification.  */
#include "unistdio.h"

#include "uniconv.h"
#include "unistr.h"
#include "unistdio/u-printf-parse.h"

#define VASNPRINTF u8_vasnprintf
#define FCHAR_T char
#define DIRECTIVE char_directive
#define DIRECTIVES char_directives
#define FCHAR_T_ONLY_ASCII 1
#define PRINTF_PARSE ulc_printf_parse
#define DCHAR_T uint8_t
#define DCHAR_T_IS_UINT8_T 1
#define DCHAR_CPY u8_cpy
#define DCHAR_SET u8_set
#define DCHAR_MBSNLEN u8_mbsnlen
#define DCHAR_IS_UINT8_T 1
#define U16_TO_DCHAR u16_to_u8
#define U32_TO_DCHAR u32_to_u8
#define DCHAR_CONV_FROM_ENCODING u8_conv_from_encoding
#define TCHAR_T char
#include "vasnprintf.c"
