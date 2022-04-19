#include <stdio.h>
#include <stddef.h>
#include <stdint.h>
#ifndef DEF_LIBCRC_CHECKSUM_H
#define DEF_LIBCRC_CHECKSUM_H

#define		CRC_POLY_32			0xEDB88320ul
#define		CRC_START_32		0xFFFFFFFFul

/**
 * @fn      uint32_t crc_32(const unsigned char*, size_t)
 * @brief   calculate crc32 for the numbytes bytes of inpuStr
 *
 * @param   inputStr
 * @param   numBytes
 * @return  unsigned int crc32.
 */
uint32_t crc_32(const unsigned char *input_str, size_t num_bytes);

/**
 * @fn      uint32_t update_crc_32(uint32_t, unsigned char)
 * @brief   calculate new crc value based on previous value of crc.
 *
 * @param   crc
 * @param   c
 * @return  new crc32 value
 */
uint32_t update_crc_32(uint32_t crc, unsigned char c);

#endif  // DEF_LIBCRC_CHECKSUM_H
