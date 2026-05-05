#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>
#include <string.h>
#include <stdint.h>
#include <inttypes.h>
#include <unistd.h>
#include <errno.h>

#include "endian.h"
#include "image.h"

#define UNIMGC_VERSION "0.1"


enum unimgc_error {
    UNIMGC_ERROR_NONE = 0,
    UNIMGC_ERROR_IO = 1,
    UNIMGC_ERROR_ALLOC = 2,
    UNIMGC_ERROR_CORRUPTED_FILE = 3,
    UNIMGC_ERROR_CORRUPTED_BLOCK = 4,
    UNIMGC_ERROR_CORRUPTED_DATA = 5
};

static void fatal(enum unimgc_error code, const char *msg, ...)
{
    fprintf(stderr, "unimgc: fatal: ");
    va_list ap;
    va_start(ap, msg);
    vfprintf(stderr, msg, ap);
    va_end(ap);
    exit(code);
}


static double si_ify(uint64_t n, char *u)
{
    static const char prefixes[] = {'k', 'M', 'G', 'T', 'P', 'E'};
    for (int i = sizeof(prefixes); i; i--) {
        uint64_t div = 1ULL << (10 * i);
        if (n > div) {
            *u = prefixes[i - 1];
            return n / (float)div;
        }
    }
    *u = ' ';
    return (double)n;
}

static void dump_header(struct imgc_header *hdr)
{
    fprintf(stderr, "volume metadata:\n");
    fprintf(stderr, "  model: %s\n", pascal_to_cstr(&hdr->volume.model));
    fprintf(stderr, "  revision: %s\n", pascal_to_cstr(&hdr->volume.revision));
    fprintf(stderr, "  serial number: %s\n", pascal_to_cstr(&hdr->volume.serial));
    fprintf(stderr, "software metadata:\n");
    fprintf(stderr, "  name: %s\n", pascal_to_cstr(&hdr->software.name));
    fprintf(stderr, "  version: %s\n", pascal_to_cstr(&hdr->software.version));

    char unit;
    double size = si_ify(hdr->image.sector_count * hdr->image.sector_size, &unit);
    fprintf(stderr, "image metadata:\n");
    fprintf(stderr, "  size: %.02f %ciB (%" PRIu64 " sectors * %" PRIu64 " bytes)\n",
        size, unit, hdr->image.sector_count, hdr->image.sector_size);
    fprintf(stderr, "  unk1: %016" PRIx64 "\n", hdr->image.unk1);
    fprintf(stderr, "  unk2: %016" PRIx64 "\n", hdr->image.unk2);
    fprintf(stderr, "  unk3: %02" PRIx8 "\n",  hdr->image.unk3);
}

static void dump_block_header(struct imgc_block_header *bhdr, size_t pos)
{
    fprintf(stderr, "block @ 0x%zx\n", pos);
    fprintf(stderr, "  type: %s\n", bhdr->type == IMGC_BLOCK_COMPRESSED ? "compressed" : "zeroed");
    fprintf(stderr, "  size: %u\n", bhdr->size);
}



static struct {
    int only_info;
    int verbose;
} options;

static void unimgc_data(struct imgc_header *hdr, FILE *in, FILE *out)
{
    uint8_t *bbuf = NULL, *dbuf = NULL;
    size_t bsize = 0, dsize = 0;
    char zeros[4096] = { 0 };

    uint64_t total_size = hdr->image.sector_count * hdr->image.sector_size;
    for (;;) {
        if (options.verbose >= 1)
            fprintf(stderr, "\r%04.2f%% (%" PRIu64 " / %" PRIu64 " bytes)...",
                (100.0 * ftello(out)) / total_size, ftello(out), total_size);

        /* read block header buf */
        uint8_t bhbuf[IMGC_BLOCK_HEADER_SIZE];
        size_t nread = fread(bhbuf, 1, sizeof(bhbuf), in);
        if (nread != sizeof(bhbuf)) {
            if (feof(in)) {
                if (nread == 0)
                    break;
                else
                    fatal(UNIMGC_ERROR_CORRUPTED_BLOCK, "truncated input file (block header)\n");
            } else {
                fatal(UNIMGC_ERROR_IO, "could not read input file: %s\n", strerror(errno));
            }
        }

        /* parse block header */
        struct imgc_block_header bhdr;
        if (imgc_parse_block(bhbuf, sizeof(bhbuf), &bhdr) < 0)
            fatal(UNIMGC_ERROR_CORRUPTED_BLOCK, "invalid IMGC block header\n");
        if (options.verbose >= 2)
            dump_block_header(&bhdr, ftello(in) - IMGC_BLOCK_HEADER_SIZE);

        /* read block data */
        size_t sz = bhdr.size - IMGC_BLOCK_HEADER_SIZE;
        if (sz > bsize) {
            if (!(bbuf = realloc(bbuf, sz)))
                fatal(UNIMGC_ERROR_ALLOC, "could not allocate %zu bytes\n", sz);
            bsize = sz;
        }
        if (fread(bbuf, 1, sz, in) != sz)
            fatal(UNIMGC_ERROR_CORRUPTED_DATA, "truncated input file (data)\n");
        
        /* process block */
        switch (bhdr.type) {
        case IMGC_BLOCK_COMPRESSED: {
            /* uncompress pseudo-LZO block */
            size_t dsz = imgc_decompress_block(bbuf, sz, NULL, 0);
            if (options.verbose >= 2)
                fprintf(stderr, "   dec: %zd\n", dsz);

            if (dsz > dsize) {
                if (!(dbuf = realloc(dbuf, dsz)))
                    fatal(UNIMGC_ERROR_ALLOC, "could not allocate %zd bytes\n", dsz);
                dsize = dsz;
            }
            if (!imgc_decompress_block(bbuf, sz, dbuf, dsz))
                fatal(UNIMGC_ERROR_CORRUPTED_DATA, "corrupted IMGC compressed block\n");
            if (fwrite(dbuf, 1, dsz, out) != dsz)
                fatal(UNIMGC_ERROR_IO, "could not write %zd bytes to target file: %s\n", dsz, strerror(errno));
            break;
        }
        case IMGC_BLOCK_ZERO: {
            size_t dsz = le64toh(*(uint64_t *)bbuf);
            if (options.verbose >= 2)
                fprintf(stderr, "  zero: %zu\n", dsz);

            while (dsz > sizeof(zeros))
                dsz -= fwrite(zeros, 1, sizeof(zeros), out);
            if (fwrite(zeros, 1, dsz, out) != dsz)
                fatal(UNIMGC_ERROR_IO, "could not write %zu bytes to target file: %s\n", dsz, strerror(errno));
            break;
        }
        default:
            fatal(UNIMGC_ERROR_CORRUPTED_BLOCK, "unknown IMGC block type: %d\n", bhdr.type);
        }
    }

    if (options.verbose >= 1)
        fputc('\n', stderr);
}

static void unimgc_header(struct imgc_header *hdr, FILE *in)
{
    uint8_t hbuf[IMGC_HEADER_SIZE];
    fread(hbuf, 1, sizeof(hbuf), in);

    if (imgc_parse(hbuf, sizeof(hbuf), hdr) < 0)
        fatal(UNIMGC_ERROR_CORRUPTED_FILE, "invalid IMGC header\n");
    
    if (options.only_info || options.verbose >= 1)
        dump_header(hdr);
}


static FILE *open_or(const char *name, const char *mode, FILE *alt)
{
    if (!name || !*name || !strcmp(name, "-"))
        return alt;
    return fopen(name, mode);
}

static void usage(const char *prog)
{
    printf("usage: %s [-ihvV] [IN] [OUT]\n", prog);
    puts("    -h \t show usage");
    puts("    -V \t show version");
    puts("    -v \t increase verbosity (level 1: progress, level 2: header dumps)");
    puts("    -i \t only show image information, don't decompress");
    puts("    IN \t input file; defaults to standard input");
    puts("   OUT \t output file; defaults to standard output");
}

int main(int argc, char **argv)
{
    int opt;

    while ((opt = getopt(argc, argv, "ihvV")) != -1) {
        switch(opt) {
        case 'i':
            options.only_info = 1;
            break;
        case 'v':
            options.verbose++;
            break;
        case 'V':
            printf("unimgc v%s\n", UNIMGC_VERSION);
            puts("copyright (c) 2019 shiz; released under the WTFPL");
            exit(0);
            break;
        case 'h':
            usage(argv[0]);
            exit(0);
            break;
        default:
            usage(argv[0]);
            exit(255);
            break;
        }
    }
    argc -= optind;

    /* open files */
    FILE *in = open_or(argc > 0 ? argv[optind] : NULL, "rb", stdin);
    if (!in)
        fatal(UNIMGC_ERROR_IO, "could not open input %s: %s\n", argv[optind], strerror(errno));
    FILE *out = open_or(argc > 1 ? argv[optind + 1] : NULL, "wb", stdout);
    if (!out)
        fatal(UNIMGC_ERROR_IO, "could not open output %s: %s\n", argv[optind + 1], strerror(errno));

    /* do the boogie */
    struct imgc_header hdr;
    unimgc_header(&hdr, in);
    if (!options.only_info)
        unimgc_data(&hdr, in, out);
}
