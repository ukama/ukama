#include <string.h>
#include <stdint.h>

#include "endian.h"
#include "image.h"
#include "lzo.h"


/* return a non-owned reference to the C string */
char *pascal_to_cstr(struct pascal_str *p)
{
    return p->data;
}

/* initialize pascal string */
int pascal_from_data(struct pascal_str *p, const char *d, size_t len)
{
    if (len + 1 > sizeof(p->data)) {
        return 0;
    }
    p->length = len;
    memcpy(p->data, d, len);
    p->data[len] = 0;
    return 1;
}

/* initialize owned pascal string from C string contents */
int pascal_from_cstr(struct pascal_str *p, const char *s)
{
    return pascal_from_data(p, s, strlen(s));
}


int imgc_parse(const uint8_t *buf, size_t len, struct imgc_header *hdr)
{
    if (len < 0x521)
        return -1;
    memcpy(&hdr->software.name,    &buf[0x000], 0x100);
    memcpy(&hdr->software.version, &buf[0x100], 0x100);
    memcpy(&hdr->volume.model,     &buf[0x200], 0x100);
    memcpy(&hdr->volume.revision,  &buf[0x300], 0x100);
    memcpy(&hdr->volume.serial,    &buf[0x400], 0x100);
    hdr->image.sector_count = le64toh(*(uint64_t *)(buf + 0x500));
    hdr->image.sector_size  = le64toh(*(uint64_t *)(buf + 0x508));
    hdr->image.unk1         = le64toh(*(uint64_t *)(buf + 0x510));
    hdr->image.unk2         = le64toh(*(uint64_t *)(buf + 0x518));
    hdr->image.unk3         = buf[0x520];
    return 0;
}

int imgc_parse_block(const uint8_t *buf, size_t len, struct imgc_block_header *hdr)
{
    /* block type:
     * "omg!": zero'd block
     * "lol!": LZO-decompress following contents
     */
    if (len < 8)
        return -1;

    if (!strncmp((const char *)buf, "omg!", 4))
        hdr->type = IMGC_BLOCK_ZERO;
    else if (!strncmp((const char *)buf, "lol!", 4))
        hdr->type = IMGC_BLOCK_COMPRESSED;
    else
        return -2;

    hdr->size = le32toh(*(const uint32_t *)(buf + 4));
    return 0;
}

size_t imgc_decompress_block(const uint8_t *buf, size_t len, uint8_t *out, size_t outlen)
{
    if (len < 2)
        /* input too small */
        return 0;

    uint32_t size = le16toh(*(const uint16_t *)buf);
    buf += 2;
    if (size & 0x8000) {
        if (len < 4)
            /* input too small */
            return 0;
        size &= 0x7FFF;
        size |= le16toh(*(const uint16_t *)buf) << 15;
        buf += 2;
    }

    if (!out)
        return size;

    if (size > outlen)
        /* output too small */
        return 0;

    return lzo_decompress(buf, len - 2 - 2 * (size > 0x7FFF), out, outlen);
}
