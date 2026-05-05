#if !defined(UNIMGC_H_LZO)
#define UNIMGC_H_LZO

#include <stdint.h>
#include <sys/types.h>

size_t lzo_decompress(const uint8_t *buf, size_t len, uint8_t *out, size_t outlen);

#endif /* !defined(UNIMGC_H_LZO) */
