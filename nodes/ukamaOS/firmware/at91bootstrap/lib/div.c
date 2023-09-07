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

int division(unsigned int dividend,
		unsigned int divisor,
		unsigned int *quotient,
		unsigned int *remainder)
{
	unsigned int shift;
	unsigned int divisor_shift;
	unsigned int factor = 0;
	unsigned char end_flag = 0;

	if (!divisor)
		return 0xffffffff;

	if (dividend < divisor) {
		*quotient = 0;
		*remainder = dividend;
		return 0;
	}

	while (dividend >= divisor) {
		for (shift = 0, divisor_shift = divisor;
			dividend >= divisor_shift;
			divisor_shift <<= 1, shift++) {
			if (dividend - divisor_shift < divisor_shift) {
				factor += 1 << shift;
				dividend -= divisor_shift;
				end_flag = 1;
				break;
			}
		}

		if (end_flag)
			continue;

		factor += 1 << (shift - 1);
		dividend -= divisor_shift >> 1;
	}

	if (quotient)
		*quotient = factor;

	if (remainder)
		*remainder = dividend;

	return 0;
}

unsigned int div(unsigned int dividend, unsigned int divisor)
{
	unsigned int quotient = 0;
	unsigned int remainder = 0;
	int ret;

	ret = division(dividend, divisor, &quotient, &remainder);
	if (ret)
		return 0xffffffff;

	return quotient;
}

unsigned int mod(unsigned int dividend, unsigned int divisor)
{
	unsigned int quotient = 0;
	unsigned int remainder = 0;
	int ret;

	ret = division(dividend, divisor, &quotient, &remainder);
	if (ret)
		return 0xffffffff;

	return remainder;
}
