/* punycode.h - prototypes for internal punycode functions
   Copyright (C) 2011-2021 Simon Josefsson

   Libidn2 is free software: you can redistribute it and/or modify it
   under the terms of either:

     * the GNU Lesser General Public License as published by the Free
       Software Foundation; either version 3 of the License, or (at
       your option) any later version.

   or

     * the GNU General Public License as published by the Free
       Software Foundation; either version 2 of the License, or (at
       your option) any later version.

   or both in parallel, as here.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received copies of the GNU General Public License and
   the GNU Lesser General Public License along with this program.  If
   not, see <http://www.gnu.org/licenses/>.
*/

#ifndef LIBIDN2_PUNYCODE_H
#define LIBIDN2_PUNYCODE_H

#include <stddef.h>
#include <stdint.h>

extern int
_idn2_punycode_encode_internal (size_t input_length,
				const uint32_t input[],
				size_t *output_length, char output[]);

extern int
_idn2_punycode_decode_internal (size_t input_length,
				const char input[],
				size_t *output_length, uint32_t output[]);

#endif /* LIBIDN2_PUNYCODE_H */
