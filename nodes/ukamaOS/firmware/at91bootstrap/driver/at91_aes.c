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
#include "hardware.h"
#include "pmc.h"
#include "arch/at91_aes.h"
#include "aes.h"
#include "debug.h"
#include "board.h"
#include "string.h"

#define swab32(x) (			\
	(((x) & 0x000000ffUL) << 24) |	\
	(((x) & 0x0000ff00UL) <<  8) |	\
	(((x) & 0x00ff0000UL) >>  8) |	\
	(((x) & 0xff000000UL) >> 24))

static inline unsigned int aes_readl(unsigned int reg)
{
	return readl(AT91C_BASE_AES + reg);
}

static inline void aes_writel(unsigned int reg, unsigned int value)
{
	writel(value, AT91C_BASE_AES + reg);
}

static inline unsigned short aes_readw(unsigned int reg)
{
	return readw(AT91C_BASE_AES + reg);
}

static inline void aes_writew(unsigned int reg, unsigned short value)
{
	writew(value, AT91C_BASE_AES + reg);
}

static inline unsigned char aes_readb(unsigned int reg)
{
	return readb(AT91C_BASE_AES + reg);
}

static inline void aes_writeb(unsigned int reg, unsigned char value)
{
	writeb(value, AT91C_BASE_AES + reg);
}


void at91_aes_init(void)
{
	/* Enable peripheral clock */
	pmc_enable_periph_clock(AT91C_ID_AES, PMC_PERIPH_CLK_DIVIDER_NA);

	/* Reset AES */
	aes_writel(AES_CR, AES_CR_SWRST);
}

void at91_aes_cleanup(void)
{
	/* Reset AES */
	aes_writel(AES_CR, AES_CR_SWRST);

	/* Disable peripheral clock */
	pmc_disable_periph_clock(AT91C_ID_AES);
}

static inline void at91_aes_set_iv(const unsigned int *iv)
{
	unsigned int i, reg = AES_IVR0;

	for (i = 0; i < AT91_AES_IV_SIZE_WORD; ++i, reg += 4) {
		unsigned int iv_be = swab32(iv[i]);
		aes_writel(reg, iv_be);
	}
}

static inline int at91_aes_set_opmode(at91_aes_operation_t operation,
				      at91_aes_mode_t mode,
				      at91_aes_key_size_t key_size,
				      unsigned int *data_width,
				      unsigned int *chunk_size)
{
	unsigned int mr = AES_MR_CKEY_PASSWD | AES_MR_SMOD_AUTO_START;

	switch (operation) {
	case AT91_AES_OP_DECRYPT:
		mr |= AES_MR_CIPHER_DECRYPT;
		break;

	case AT91_AES_OP_ENCRYPT:
		mr |= AES_MR_CIPHER_ENCRYPT;
		break;

	case AT91_AES_OP_MAC:
		mr |= AES_MR_CIPHER_ENCRYPT | AES_MR_LOD;
		break;

	default:
		return -1;
	}

	switch (mode) {
	case AT91_AES_MODE_ECB:
		mr |= AES_MR_OPMOD_ECB;
		*chunk_size = 4;
		*data_width = 4;
		break;

	case AT91_AES_MODE_CBC:
		mr |= AES_MR_OPMOD_CBC;
		*chunk_size = 4;
		*data_width = 4;
		break;

	case AT91_AES_MODE_OFB:
		mr |= AES_MR_OPMOD_OFB;
		*chunk_size = 4;
		*data_width = 4;
		break;

	case AT91_AES_MODE_CFB_128:
		mr |= AES_MR_OPMOD_CFB | AES_MR_CFBS_128BIT;
		*chunk_size = 4;
		*data_width = 4;
		break;

	case AT91_AES_MODE_CFB_64:
		mr |= AES_MR_OPMOD_CFB | AES_MR_CFBS_64BIT;
		*chunk_size = 2;
		*data_width = 4;
		break;

	case AT91_AES_MODE_CFB_32:
		mr |= AES_MR_OPMOD_CFB | AES_MR_CFBS_32BIT;
		*chunk_size = 1;
		*data_width = 4;
		break;

	case AT91_AES_MODE_CFB_16:
		mr |= AES_MR_OPMOD_CFB | AES_MR_CFBS_16BIT;
		*chunk_size = 1;
		*data_width = 2;
		break;

	case AT91_AES_MODE_CFB_8:
		mr |= AES_MR_OPMOD_CFB | AES_MR_CFBS_8BIT;
		*chunk_size = 1;
		*data_width = 1;
		break;

	case AT91_AES_MODE_CTR:
		mr |= AES_MR_OPMOD_CTR;
		*chunk_size = 4;
		*data_width = 4;
		break;

	case AT91_AES_MODE_GCM:
		mr |= AES_MR_OPMOD_GCM;
		*chunk_size = 4;
		*data_width = 4;
		break;

	default:
		return -1;
	}

	switch (key_size) {
	case AT91_AES_KEY_SIZE_128:
		mr |= AES_MR_KEYSIZE_AES128;
		break;

	case AT91_AES_KEY_SIZE_192:
		mr |= AES_MR_KEYSIZE_AES192;
		break;

	case AT91_AES_KEY_SIZE_256:
		mr |= AES_MR_KEYSIZE_AES256;
		break;

	default:
		return -1;
	}

	aes_writel(AES_MR, mr);
	return 0;
}

static inline int at91_aes_set_key(at91_aes_key_size_t key_size,
				   const unsigned int *key)
{
	unsigned int i, num_words, reg = AES_KEYWR0;

	switch (key_size) {
	case AT91_AES_KEY_SIZE_128:
		num_words = 4;
		break;

	case AT91_AES_KEY_SIZE_192:
		num_words = 6;
		break;

	case AT91_AES_KEY_SIZE_256:
		num_words = 8;
		break;

	default:
		return -1;
	}

	for (i = 0; i < num_words; ++i, reg += 4) {
		unsigned int key_be = swab32(key[i]);
		aes_writel(reg, key_be);
	}

	return 0;
}


static void at91_aes_compute_pio_long(unsigned int chunk_size,
				      unsigned int num_blocks,
				      unsigned int is_mac,
				      const unsigned int *input,
				      unsigned int *output)
{
	unsigned int n, i, reg;

	for (n = 0; n < num_blocks; ++n) {
		reg = AES_IDATAR0;
		for (i = 0; i < chunk_size; ++i, reg += 4)
			aes_writel(reg, *input++);

		while (!(aes_readl(AES_ISR) & AES_INT_DATRDY));
		if (is_mac)
			continue;

		reg = AES_ODATAR0;
		for (i = 0; i < chunk_size; ++i, reg += 4)
			*output++ = aes_readl(reg);
	}
}

static void at91_aes_compute_pio_word(unsigned int num_blocks,
				      unsigned int is_mac,
				      const unsigned short *input,
				      unsigned short *output)
{
	unsigned int n;

	for (n = 0; n < num_blocks; ++n) {
		aes_writew(AES_IDATAR0, *input++);

		while (!(aes_readl(AES_ISR) & AES_INT_DATRDY));
		if (is_mac)
			continue;

		*output++ = aes_readw(AES_ODATAR0);
	}
}

static void at91_aes_compute_pio_byte(unsigned int num_blocks,
				      unsigned int is_mac,
				      const unsigned char *input,
				      unsigned char *output)
{
	unsigned int n;

	for (n = 0; n < num_blocks; ++n) {
		aes_writeb(AES_IDATAR0, *input++);

		while (!(aes_readl(AES_ISR) & AES_INT_DATRDY));
		if (is_mac)
			continue;

		*output++ = aes_readb(AES_ODATAR0);
	}
}

static void at91_aes_compute_pio(unsigned int data_width,
				 unsigned int chunk_size,
				 unsigned int num_blocks,
				 unsigned int is_mac,
				 const void *input,
				 void *output)
{
	switch (data_width) {
	case 4:
		at91_aes_compute_pio_long(chunk_size,
					  num_blocks,
					  is_mac,
					  (const unsigned int *)input,
					  (unsigned int *)output);
		break;

	case 2:
		at91_aes_compute_pio_word(num_blocks,
					  is_mac,
					  (const unsigned short *)input,
					  (unsigned short *)output);
		break;

	case 1:
		at91_aes_compute_pio_byte(num_blocks,
					  is_mac,
					  (const unsigned char *)input,
					  (unsigned char *)output);
		break;

	default:
		break;
	}
}

static inline unsigned int at91_aes_length2blocks(unsigned int data_length,
						  unsigned int block_size)
{
	unsigned int mask = 0;
	unsigned int shift = 0;

	switch (block_size) {
	default:
	case 1:
		break;

	case 2:
		mask = 0x1;
		shift = 1;
		break;

	case 4:
		mask = 0x3;
		shift = 2;
		break;

	case 8:
		mask = 0x7;
		shift = 3;
		break;

	case 16:
		mask = 0xf;
		shift = 4;
		break;
	}

	return (data_length >> shift) + ((data_length & mask) ? 1 : 0);
}

static int at91_aes_process(const at91_aes_params_t *params)
{
	unsigned int data_width, chunk_size;
	unsigned int block_size, num_blocks;
	unsigned int is_mac = (params->operation == AT91_AES_OP_MAC);

	/* Reset AES */
	aes_writel(AES_CR, AES_CR_SWRST);

	if (at91_aes_set_opmode(params->operation, params->mode,
				params->key_size, &data_width, &chunk_size))
		return -1;

	if (at91_aes_set_key(params->key_size, params->key))
		return -1;

	if (params->mode != AT91_AES_MODE_ECB)
		at91_aes_set_iv(params->iv);


	aes_writel(AES_IER, AES_INT_DATRDY);

	block_size = data_width * chunk_size;
	num_blocks = at91_aes_length2blocks(params->data_length, block_size);
	at91_aes_compute_pio(data_width, chunk_size, num_blocks,
			     is_mac, params->input, params->output);

	if (is_mac) {
		unsigned int reg, i;

		reg = AES_ODATAR0;
		for (i = 0; i < AT91_AES_BLOCK_SIZE_WORD; ++i, reg += 4)
			((unsigned int *)params->output)[i] = aes_readl(reg);
	}

	return 0;
}

int at91_aes_cbc(unsigned int data_length,
		 const void *input,
		 void *output,
		 int encrypt,
		 at91_aes_key_size_t key_size,
		 const unsigned int *key,
		 const unsigned int *iv)
{
	at91_aes_params_t params;

	if (!data_length || !input || !output || !key || !iv)
		return -1;

	memset(&params, 0, sizeof(params));
	params.operation = (encrypt) ? AT91_AES_OP_ENCRYPT : AT91_AES_OP_DECRYPT;
	params.mode = AT91_AES_MODE_CBC;
	params.data_length = data_length;
	params.input = input;
	params.output = output;
	params.key_size = key_size;
	params.key = key;
	params.iv = iv;

	return at91_aes_process(&params);
}

int at91_aes_cmac(unsigned int data_length,
		  const void *data,
		  unsigned int *cmac,
		  at91_aes_key_size_t key_size,
		  const unsigned int *key)
{
	static const unsigned int null_block[AT91_AES_BLOCK_SIZE_WORD];
	static const unsigned int null_iv[AT91_AES_IV_SIZE_WORD];
	unsigned int last_input[AT91_AES_BLOCK_SIZE_WORD];
	unsigned int subkey[AT91_AES_BLOCK_SIZE_WORD];
	const unsigned int *input = (const unsigned int *)data;
	unsigned int num_blocks, offset;
	at91_aes_params_t params;
	unsigned char carry;
	int i; /* MUST be signed for the subkey loop */

	if (!data_length || !data || !cmac || !key)
		return -1;

	/* Set common parameters once for all */
	memset(&params, 0, sizeof(params));
	params.key_size = key_size;
	params.key = key;

	/* Generate the subkey */
	params.operation = AT91_AES_OP_ENCRYPT;
	params.mode = AT91_AES_MODE_ECB;
	params.data_length = AT91_AES_BLOCK_SIZE_BYTE;
	params.input = null_block;
	params.output = subkey;
	if (at91_aes_process(&params))
		return -1;

	carry = 0;
	for (i = AT91_AES_BLOCK_SIZE_BYTE-1; i >= 0; --i) {
		unsigned char tmp, next_carry;

		tmp = ((unsigned char *)subkey)[i];
		next_carry = ((tmp & 0x80) != 0);
		((unsigned char *)subkey)[i] = (tmp << 1) | carry;
		carry = next_carry;
	}
	carry = (0 - carry) & 0x87;
	((unsigned char *)subkey)[AT91_AES_BLOCK_SIZE_BYTE-1] ^= carry;

	/* Process the n-1 first blocks */
	num_blocks = at91_aes_length2blocks(data_length,
					    AT91_AES_BLOCK_SIZE_BYTE);
	if (num_blocks > 1) {
		params.operation = AT91_AES_OP_MAC;
		params.mode = AT91_AES_MODE_CBC;
		params.iv = null_iv;
		params.data_length = data_length - AT91_AES_BLOCK_SIZE_BYTE;
		params.input = data;
		params.output = cmac;
		if (at91_aes_process(&params))
			return -1;
	} else {
		memset(cmac, 0, AT91_AES_BLOCK_SIZE_BYTE);
	}

	/* Process the last block */
	offset = (num_blocks-1) * AT91_AES_BLOCK_SIZE_WORD;
	for (i = 0; i < AT91_AES_BLOCK_SIZE_WORD; ++i)
		last_input[i] = input[offset + i] ^ cmac[i] ^ subkey[i];

	params.operation = AT91_AES_OP_ENCRYPT;
	params.mode = AT91_AES_MODE_ECB;
	params.data_length = AT91_AES_BLOCK_SIZE_BYTE;
	params.input = last_input;
	params.output = cmac;
	return at91_aes_process(&params);
}
