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
#include "common.h"
#include "secure.h"
#include "aes.h"
#include "debug.h"
#include "string.h"
#include "autoconf.h"

static unsigned int cipher_key[8] = {
	CONFIG_AES_CIPHER_KEY_WORD0,
	CONFIG_AES_CIPHER_KEY_WORD1,
	CONFIG_AES_CIPHER_KEY_WORD2,
	CONFIG_AES_CIPHER_KEY_WORD3,
#if defined(CONFIG_AES_KEY_SIZE_192) || defined(CONFIG_AES_KEY_SIZE_256)
	CONFIG_AES_CIPHER_KEY_WORD4,
	CONFIG_AES_CIPHER_KEY_WORD5,
#endif
#if defined(CONFIG_AES_KEY_SIZE_256)
	CONFIG_AES_CIPHER_KEY_WORD6,
	CONFIG_AES_CIPHER_KEY_WORD7,
#endif
};

static unsigned int cmac_key[8] = {
	CONFIG_AES_CMAC_KEY_WORD0,
	CONFIG_AES_CMAC_KEY_WORD1,
	CONFIG_AES_CMAC_KEY_WORD2,
	CONFIG_AES_CMAC_KEY_WORD3,
#if defined(CONFIG_AES_KEY_SIZE_192) || defined(CONFIG_AES_KEY_SIZE_256)
	CONFIG_AES_CMAC_KEY_WORD4,
	CONFIG_AES_CMAC_KEY_WORD5,
#endif
#if defined(CONFIG_AES_KEY_SIZE_256)
	CONFIG_AES_CMAC_KEY_WORD6,
	CONFIG_AES_CMAC_KEY_WORD7,
#endif
};

static unsigned int iv[AT91_AES_IV_SIZE_WORD] = {
	CONFIG_AES_IV_WORD0,
	CONFIG_AES_IV_WORD1,
	CONFIG_AES_IV_WORD2,
	CONFIG_AES_IV_WORD3,
};

static int secure_decrypt(void *data, unsigned int data_length, int is_signed)
{
	at91_aes_key_size_t key_size;
	unsigned int computed_cmac[AT91_AES_BLOCK_SIZE_WORD];
	unsigned int fixed_length;
	const unsigned int *cmac;
	int rc = -1;

#if defined(CONFIG_AES_KEY_SIZE_128)
	key_size = AT91_AES_KEY_SIZE_128;
#elif defined(CONFIG_AES_KEY_SIZE_192)
	key_size = AT91_AES_KEY_SIZE_192;
#elif defined(CONFIG_AES_KEY_SIZE_256)
	key_size = AT91_AES_KEY_SIZE_256;
#else
#error "bad AES key size"
#endif

	/* Init periph */
	at91_aes_init();

	/* Check signature if required */
	if (is_signed) {
		/* Compute the CMAC */
		if (at91_aes_cmac(data_length, data, computed_cmac,
				  key_size, cmac_key))
			goto exit;

		/* Check the CMAC */
		fixed_length = at91_aes_roundup(data_length);
		cmac = (const unsigned int *)((char *)data + fixed_length);
		if (!consttime_memequal(cmac, computed_cmac, AT91_AES_BLOCK_SIZE_BYTE))
			goto exit;
	}

	/* Decrypt the whole file */
	if (at91_aes_cbc(data_length, data, data, 0,
			 key_size, cipher_key, iv))
		goto exit;

	rc = 0;
exit:
	/* Reset periph */
	at91_aes_cleanup();

	return rc;
}

static void wipe_keys()
{
	/* Reset keys */
	memset(cmac_key, 0, sizeof(cmac_key));
	memset(cipher_key, 0, sizeof(cipher_key));
	memset(iv, 0, sizeof(iv));
}

int secure_check(void *data)
{
	const at91_secure_header_t *header;
	void *file;
	int ret = -1;

	if (secure_decrypt(data, sizeof(*header), 0))
		goto secure_wipe_keys;

	header = (const at91_secure_header_t *)data;
	if (header->magic != AT91_SECURE_MAGIC)
		goto secure_wipe_keys;

	file = (unsigned char *)data + sizeof(*header);
	ret = secure_decrypt(file, header->file_size, 1);

secure_wipe_keys:
	wipe_keys();
	return ret;
}
