#include <string.h>
#include <stdint.h>
#include <stddef.h>
#include "endian.h"

#define min(a, b) ((a) > (b) ? (b) : (a))
/* extract n bits starting from position p from value v */
#define bits(v, p, n) (((v) >> (p)) & ((1 << (n)) - 1))
/* extract bit from position p from value v */
#define bit(v, p) bits((v), (p), 1)


/* references:
 * - https://www.infradead.org/~mchehab/kernel_docs/unsorted/lzo.html
 * - https://github.com/synopse/mORMot/blob/master/SynLZO.pas
 */

static inline unsigned lzo_parse_length(const uint8_t **buf, unsigned bits)
{
    unsigned mask = (1 << bits) - 1;

    const uint8_t *p = *buf;
    unsigned len = *p++ & mask;
    if (!len) {
        for (; !*p; p++)
            len += 0xFF;
        len += *p++ + mask;
    }

    *buf = p;
    return len;
}

static inline void lzo_copy(uint8_t **out, const uint8_t **in, size_t n)
{
    memcpy(*out, *in, n);
    *out += n;
    *in += n;
}

static inline void lzo_copy_distance(uint8_t **out, ptrdiff_t dist, size_t n)
{
    /* interestingly, memmove() does *not* work here, at least on macOS */
    uint8_t *p = *out;
    const uint8_t *in = p - dist;
    while (n--)
        *p++ = *in++;
    *out = p;
}

size_t lzo_decompress(const uint8_t *buf, size_t len, uint8_t *out, size_t outlen)
{
    const uint8_t *p = buf, *end = buf + len, *oend = out + outlen;
    uint8_t *op = out;
    uint8_t state = 0;

    while (p < end) {
        uint8_t instr = *p;
        /* first command is special */
        if (p == buf) {
            if (instr > 17) {
                lzo_copy(&op, &p, instr - 17);
                state = 4;
                continue;
            }
        }

        uint32_t length = 0;
        uint16_t follow = 0;
        ptrdiff_t distance = 0;
        /* general notation for following formats:
         * L: length bits
         * D: distance bits
         * S: state bits
         */
        if (instr >= 64) {
            /* format: L L L D D D S S;
             * follow: D D D D D D D D */
            p++;
            length = bits(instr, 5, 3) + 1;
            follow = *p++;
            distance = (follow << 3) + bits(instr, 2, 3) + 1;
            state = bits(instr, 0, 2);
        } else if (instr >= 32) {
            /* format: 0 0 1 L L L L L;
             * follow: D D D D D D D D D D D D D D S S */
            length = lzo_parse_length(&p, 5) + 2;
            follow = le16toh(*(const uint16_t *)p);
            distance = (follow >> 2) + 1;
            state = bits(follow, 0, 2);
            p += 2;
        } else if (instr >= 16) {
            /* format: 0 0 0 1 D L L L;
             * follow: D D D D D D D D D D D D D D S S */
            length = lzo_parse_length(&p, 3) + 2;
            follow = le16toh(*(const uint16_t *)p);
            distance = (bit(instr, 3) << 14) + (follow >> 2) + 0x4000;
            state = bits(follow, 0, 2);
            p += 2;
        } else if (!state) {
            /* back up and parse length properly */
            length = lzo_parse_length(&p, 4) + 3;
            lzo_copy(&op, &p, length);
            state = 4;
            continue;
        } else {
            /* this shouldn't happen */
            return 0;
        }

        length = min(length, oend - op);
        if (length)
            lzo_copy_distance(&op, distance, length);
        if (state > 0 && state < 4)
            lzo_copy(&op, &p, state);
    }
    return op - out;
}
