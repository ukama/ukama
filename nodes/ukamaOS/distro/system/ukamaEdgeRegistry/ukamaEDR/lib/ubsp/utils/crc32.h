#include <stdio.h>
#include <stddef.h>
#include <stdint.h>
#ifndef DEF_LIBCRC_CHECKSUM_H
#define DEF_LIBCRC_CHECKSUM_H

#define		CRC_POLY_32			0xEDB88320ul
#define		CRC_START_32		0xFFFFFFFFul

uint32_t crc_32(const unsigned char *input_str, size_t num_bytes);
uint32_t update_crc_32(uint32_t crc, unsigned char c);

#endif  // DEF_LIBCRC_CHECKSUM_H
