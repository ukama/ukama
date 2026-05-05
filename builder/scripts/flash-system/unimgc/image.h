#if !defined(UNIMGC_H_IMAGE)
#define UNIMGC_H_IMAGE

#include <stdint.h>


struct pascal_str {
    unsigned char length;
    char data[0xFF];
};


#define IMGC_HEADER_SIZE 0x1000

struct imgc_header {
    /* creator software metadata */
    struct {
        /* "HDD Raw Copy Tool" */
        struct pascal_str name;
        /* "1.10" */
        struct pascal_str version;
    } software;
    /* original volume metadata */
    struct {
        /* "SMI USB DISK" */
        struct pascal_str model;
        /* "0CB0" */
        struct pascal_str revision;
        /* "624811604181" */
        struct pascal_str serial;
    } volume;
    /* image properties */
    struct {
        uint64_t sector_count; /* ? */
        uint64_t sector_size;  /* ? */
        uint64_t unk1;
        uint64_t unk2;
        uint8_t  unk3;
    } image;
};

enum imgc_block_type {
    IMGC_BLOCK_COMPRESSED = 1,
    IMGC_BLOCK_ZERO
};

#define IMGC_BLOCK_HEADER_SIZE  8

struct imgc_block_header {
    enum imgc_block_type type;
    uint32_t size;
};


char *pascal_to_cstr(struct pascal_str *p);
int pascal_from_cstr(struct pascal_str *p, const char *s);

int imgc_parse(const uint8_t *buf, size_t len, struct imgc_header *hdr);
int imgc_parse_block(const uint8_t *buf, size_t len, struct imgc_block_header *hdr);
size_t imgc_decompress_block(const uint8_t *buf, size_t len, uint8_t *out, size_t outlen);

#endif /* !defined(UNIMGC_H_IMAGE) */
