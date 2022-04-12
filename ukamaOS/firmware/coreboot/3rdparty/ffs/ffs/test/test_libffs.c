/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ffs/test/test_libffs.c $                                      */
/*                                                                        */
/* OpenPOWER FFS Project                                                  */
/*                                                                        */
/* Contributors Listed Below - COPYRIGHT 2014,2015                        */
/* [+] International Business Machines Corp.                              */
/*                                                                        */
/*                                                                        */
/* Licensed under the Apache License, Version 2.0 (the "License");        */
/* you may not use this file except in compliance with the License.       */
/* You may obtain a copy of the License at                                */
/*                                                                        */
/*     http://www.apache.org/licenses/LICENSE-2.0                         */
/*                                                                        */
/* Unless required by applicable law or agreed to in writing, software    */
/* distributed under the License is distributed on an "AS IS" BASIS,      */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or        */
/* implied. See the License for the specific language governing           */
/* permissions and limitations under the License.                         */
/*                                                                        */
/* IBM_PROLOG_END_TAG                                                     */

/*
 *    File: test_libffs.c
 *  Author: Shekar Babu S <shekbabu@in.ibm.com>
 *   Descr: unit test tool for api's in libffs.so
 *    Date: 06/26/2012
 */
#include <sys/types.h>
#include <sys/stat.h>

#include <fcntl.h>
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include <unistd.h>
#include <getopt.h>
#include <errno.h>
#include <ctype.h>

#include <clib/exception.h>
#include <clib/attribute.h>
#include <clib/min.h>
#include <sys/xattr.h>

#include <clib/bb_trace.h>

#include "test_libffs.h"

FILE*
log_open(void) {

        FILE *logfile = NULL;
        logfile = fopen("test_libffs.log", "w");
        if (logfile == NULL) {
                perror("logfile");
                exit(EXIT_FAILURE);
        }
        setvbuf(logfile, NULL, _IOLBF, 0);
        return logfile;
}

void
create_regular_file(ffs_ops_t *ffs_ops) {
	int rc = 0;
	size_t size = ffs_ops->device_size;
	mode_t mode = (S_IRUSR | S_IWUSR) | (S_IRGRP | S_IRGRP) | S_IROTH;

	fprintf(ffs_ops->log, "%s: Creating regular file\n", __func__);

	int fd = open(ffs_ops->nor_image, O_RDWR | O_CREAT | O_TRUNC, mode);
	if (fd == -1) {
		fprintf(ffs_ops->log, "%s: Error creating regular file '%s'",
			 __func__, ffs_ops->nor_image);
		rc = FFS_ERROR;
		return;
	}

	if (ftruncate(fd, size) != 0) {
		fprintf(ffs_ops->log, "%s: Error truncating '%s'",
                         __func__, ffs_ops->nor_image);
                rc = FFS_ERROR;
		return;
        }

	uint32_t page_size = sysconf(_SC_PAGESIZE);
	char buf[page_size];
	memset(buf, 0xFF, page_size);

	while (0 < size) {
		ssize_t rc = write(fd, buf, min(sizeof(buf), size));
		if (rc == -1) {
			fprintf(ffs_ops->log, "%s: Error writing to '%s'",
				 __func__, ffs_ops->nor_image);
			rc = FFS_ERROR;
			return;
		} else if(rc == 0) {
			break;
		}
		size -= rc;
	}

	if(fd == 0) {
		close(fd);
	}
}

int
create_partition(ffs_ops_t *ffs_ops) {

	uint32_t rc = 0;
	ffs_t *ffs = NULL;

	create_regular_file(ffs_ops);

	ffs = __ffs_create(ffs_ops->nor_image, ffs_ops->part_off,
			ffs_ops->blk_sz,
			ffs_ops->device_size / ffs_ops->blk_sz);

	if(ffs != NULL) {
		__ffs_close(ffs);
		fprintf(ffs_ops->log, "%s: Creating pnor image '%s' success\n",
			__func__, ffs_ops->nor_image);
	} else {
		fprintf(ffs_ops->log, "%s: Creating pnor image '%s' failed\n",
			__func__, ffs_ops->nor_image);
		rc = FFS_ERROR;
	}

	return rc;
}

int
add_partition_entry(ffs_ops_t *ffs_ops) {

	uint32_t rc = 0;
	ffs_t *ffs = NULL;
	ffs = __ffs_open(ffs_ops->nor_image, ffs_ops->part_off);
	if(ffs != NULL) {
		__ffs_entry_add(ffs, ffs_ops->part_entry, ffs_ops->entry_off,
				ffs_ops->entry_sz, ffs_ops->type, 0x00);
		__ffs_close(ffs);
		fprintf(ffs_ops->log, "%s: Adding partition entry '%s' "
			"success\n", __func__, ffs_ops->part_entry);
	} else {
		fprintf(ffs_ops->log, "%s: Adding partition entry '%s' "
			"failed\n", __func__, ffs_ops->part_entry);
	}
	return rc;
}

uint32_t
read_partition(ffs_ops_t *ffs_ops) {

	uint32_t rc = 0;
	ffs_t *ffs = NULL;

	FILE *fp = fopen(ffs_ops->o_file, "w+");
	if (fp == NULL) {
		fprintf(ffs_ops->log, "%s: Error opening file '%s'", __func__,
			ffs_ops->o_file);
		rc = FFS_ERROR;
		goto error;
	}

	ffs = __ffs_open(ffs_ops->nor_image, ffs_ops->part_off);
	if(ffs == NULL) {
		fprintf(ffs_ops->log, "%s: Error, opening nor "
			"image '%s'\n", __func__, ffs_ops->nor_image);
		rc = FFS_ERROR;
		goto error;
	}
	fprintf(ffs_ops->log, "%s: Successfully opened "
		"nor image\n", __func__);
	ffs_entry_t entry;

	if (__ffs_entry_find(ffs, ffs_ops->part_entry, &entry) == false) {
		fprintf(ffs_ops->log, "%s: Error, '%s' not found\n", __func__,
			ffs_ops->part_entry);
		rc = FFS_ERROR;
		goto error;
	}

	if(false) {
error:
		if(fp != NULL)
			fclose(fp);
		if(ffs != NULL)
			__ffs_close(ffs);
		return rc;
	}

	fprintf(ffs_ops->log, "%s: Finding entry '%s' success\n", __func__,
		ffs_ops->part_entry);

	uint32_t block_size = ffs->hdr->block_size;
	char block[block_size];
	memset(block, 0, block_size);

	if (setvbuf(fp, NULL, _IOFBF, block_size) != 0) {
		fprintf(ffs_ops->log, "%s: Error, setvbuf failed "
			"%s (errno=%d)", __func__, strerror(errno), errno);
		rc = FFS_ERROR;
		goto out;
	}

	for (uint32_t i=0; i<entry.size; i++) {
		size_t rc = __ffs_entry_read(ffs, ffs_ops->part_entry, block,
						i * block_size, block_size);
		if(rc != block_size) {
			fprintf(ffs_ops->log, "%s: Error, ffs_entry_read"
				" '%s'\n", __func__, ffs_ops->part_entry);
	                rc = FFS_ERROR;
			goto out;
		}

		rc = fwrite(block, rc, 1, fp);
	        if (rc != 1) {
			fprintf(ffs_ops->log, "%s: Error, fwrite "
				"%s (errno=%d)\n", __func__,
				strerror(ferror(fp)), ferror(fp));
			rc = FFS_ERROR;
			goto out;
		}
		else if (rc == 0)
	            break;
	}

	uint32_t bytes_read = entry.size * block_size;

	fprintf(ffs_ops->log, "%s: Read %d bytes from partition '%s' Success\n",
		 __func__, bytes_read, ffs_ops->part_entry);

	if (fclose(fp) == EOF) {
	        fprintf(ffs_ops->log, "%s: Error, flose '%s' "
			"=> %s (errno=%d)", __func__, ffs_ops->o_file,
			strerror(errno), errno);
		fp = NULL;
		rc = FFS_ERROR;
		goto out;
	}

	__ffs_close(ffs);

	fprintf(ffs_ops->log, "%s: Writing %d bytes from partition '%s' to "
		"file '%s' Success\n", __func__, bytes_read,
		ffs_ops->part_entry, ffs_ops->o_file);

	if(false) {
out:
		if(fp != NULL)
			fclose(fp);
		if(ffs != NULL)
			__ffs_close(ffs);
	}
	return rc;
}

uint32_t
write_partition(ffs_ops_t * ffs_ops) {

	struct stat st;
	uint32_t rc = 0;
	ffs_t * ffs = NULL;


	if (stat(ffs_ops->i_file, &st) != 0) {
	        fprintf(ffs_ops->log, "%s: '%s' => %s (errno=%d)", __func__,
			ffs_ops->i_file, strerror(errno), errno);
		rc = FFS_ERROR;
		goto error;
	}

	ffs = __ffs_open(ffs_ops->nor_image, ffs_ops->part_off);
	if(ffs == NULL) {
		fprintf(ffs_ops->log, "%s: Error, opening nor "
			"image '%s'\n", __func__, ffs_ops->nor_image);
		rc = FFS_ERROR;
		goto error;
	}

	fprintf(ffs_ops->log, "%s: Successfully opened nor image '%s' for "
		"writing\n", __func__, ffs_ops->nor_image);

	ffs_entry_t entry;
	if (__ffs_entry_find(ffs, ffs_ops->part_entry, &entry) == false) {
		fprintf(ffs_ops->log, "%s: Error, '%s' not found\n", __func__,
			ffs_ops->part_entry);
		rc = FFS_ERROR;
		goto error;
	}

	size_t entry_size = entry.size * ffs->hdr->block_size;


	fprintf(ffs_ops->log, "%s: Found entry '%s' with size=%d\n", __func__,
		ffs_ops->part_entry, entry_size);

	if (entry_size < st.st_size) {
		fprintf(ffs_ops->log, "%s: '%s' of size '%lld' too big for "
			"partition '%s' of size '%d'\n", __func__,
			ffs_ops->i_file, st.st_size, ffs_ops->nor_image,
			entry_size);
		rc = FFS_ERROR;
		goto error;
	}

	FILE * fp = fopen(ffs_ops->i_file, "r+");
	if (fp == NULL) {
		fprintf(ffs_ops->log, "%s: Error opening file '%s'", __func__,
			ffs_ops->i_file);
		rc = FFS_ERROR;
		goto error;
	}

	if(false) {
error:
		if(ffs != NULL)
			__ffs_close(ffs);
		return rc;
	}

	uint32_t block_size = ffs->hdr->block_size;
	char block[block_size];

	if (setvbuf(fp, NULL, _IOFBF, block_size) != 0) {
		fprintf(ffs_ops->log, "%s: Error, setvbuf failed "
			"%s (errno=%d)", __func__, strerror(errno), errno);
		rc = FFS_ERROR;
		goto out;
	}

	fprintf(ffs_ops->log, "%s: Writing data file into partition\n",
		__func__);

	for (uint32_t j=0; j<entry.size; j++) {
		clearerr(fp);

	        size_t bytes_read = fread(block, 1, block_size, fp);
		__ffs_entry_write(ffs, ffs_ops->part_entry, block,
					j * block_size,	bytes_read);
		if (bytes_read == 0) {
			int err = ferror(fp);
			if (err) {
				fprintf(ffs_ops->log, "%s: Error, setvbuf "
					"failed %s (errno=%d)", __func__,
					strerror(errno), errno);
				rc = FFS_ERROR;
				goto out;
			}
			else {
				break;
			}
	        }
	}

	if (fclose(fp) == EOF) {
	        fprintf(ffs_ops->log, "%s: Error, flose '%s' "
			"=> %s (errno=%d)", __func__, ffs_ops->o_file,
			strerror(errno), errno);
		fp = NULL;
		rc = FFS_ERROR;
		goto out;
	}
	fprintf(ffs_ops->log, "%s: Writing to partition '%s' from data file "
		"'%s' Success\n", __func__, ffs_ops->part_entry,
		ffs_ops->i_file);

	if(false) {
out:
		__ffs_close(ffs);
	}
	return rc;
}

uint32_t
list_partition(ffs_ops_t * ffs_ops) {

	uint32_t rc = 0;
	ffs_t * ffs = NULL;

	ffs = __ffs_open(ffs_ops->nor_image, ffs_ops->part_off);
	if(ffs == NULL) {
		fprintf(ffs_ops->log, "%s: Error, opening nor "
			"image '%s'\n", __func__, ffs_ops->nor_image);
		fprintf(ffs_ops->log, "%s: Listing partition entries in '%s'"
			" image failed\n", __func__, ffs_ops->nor_image);
		rc = FFS_ERROR;
		return rc;
	}
	//!< List all entries in the partition table
	__ffs_list_entries(ffs, ".*", true, stdout);
	fprintf(ffs_ops->log, "%s: Listing partition entries in '%s' image "
				"success\n", __func__, ffs_ops->nor_image);

	__ffs_close(ffs);

	return rc;
}

uint32_t
hexdump_entry(ffs_ops_t * ffs_ops) {

        uint32_t rc = 0;
        ffs_t * ffs = NULL;

        ffs = __ffs_open(ffs_ops->nor_image, ffs_ops->part_off);
        if(ffs == NULL) {
                fprintf(ffs_ops->log, "%s: Error, opening nor "
                        "image '%s'\n", __func__, ffs_ops->nor_image);
		fprintf(ffs_ops->log, "%s: Hexdump of partition entries in '%s'"
			" image failed\n", __func__, ffs_ops->nor_image);
                rc = FFS_ERROR;
                return rc;
        }
        //!< Hexdump the entry in the partition table
	__ffs_entry_hexdump(ffs, ffs_ops->part_entry, stdout);
	fprintf(ffs_ops->log, "%s: Hexdump of partition entries in '%s' image "
				"success\n", __func__, ffs_ops->nor_image);
        __ffs_close(ffs);

        return rc;
}

uint32_t
delete_entry(ffs_ops_t * ffs_ops) {

        uint32_t rc = 0;
        ffs_t * ffs = NULL;

        ffs = __ffs_open(ffs_ops->nor_image, ffs_ops->part_off);
        if(ffs == NULL) {
                fprintf(ffs_ops->log, "%s: Error, opening nor "
                        "image '%s'\n", __func__, ffs_ops->nor_image);
		fprintf(ffs_ops->log, "%s: Delete entry '%s' failed\n",
			__func__, ffs_ops->part_entry);
                rc = FFS_ERROR;
                return rc;
        }
        //!< Delete the entry from the partition table
        __ffs_entry_delete(ffs, ffs_ops->part_entry);
	fprintf(ffs_ops->log, "%s: Delete entry '%s' success\n", __func__,
			ffs_ops->part_entry);

        __ffs_close(ffs);

        return rc;
}

uint32_t
modify_entry_get(ffs_ops_t * ffs_ops) {

        uint32_t rc = 0;
        uint32_t value = 0;
        ffs_t * ffs = NULL;

        ffs = __ffs_open(ffs_ops->nor_image, ffs_ops->part_off);
        if(ffs == NULL) {
                fprintf(ffs_ops->log, "%s: Error, opening nor "
                        "image '%s'\n", __func__, ffs_ops->nor_image);
		fprintf(ffs_ops->log, "%s: Get user word at '%d' failed\n",
			__func__, ffs_ops->user);
                rc = FFS_ERROR;
                return rc;
        }
	//!< Get the user word at index
	__ffs_entry_user_get(ffs, ffs_ops->part_entry, ffs_ops->user, &value);
	fprintf(ffs_ops->log, "%s: Get user word at '%d' --> %x success\n",
			__func__, ffs_ops->user, value);
	//!< Required to check what is put
	fprintf(stdout, "UW: %d-->%d\n", ffs_ops->user, value);

        __ffs_close(ffs);

        return rc;
}

uint32_t
modify_entry_put(ffs_ops_t * ffs_ops) {

        uint32_t rc = 0;
        ffs_t * ffs = NULL;

        ffs = __ffs_open(ffs_ops->nor_image, ffs_ops->part_off);
        if(ffs == NULL) {
                fprintf(ffs_ops->log, "%s: Error, opening nor "
                        "image '%s'\n", __func__, ffs_ops->nor_image);
		fprintf(ffs_ops->log, "%s: Put user word at '%d' failed\n",
			__func__, ffs_ops->user);
                rc = FFS_ERROR;
                return rc;
        }
        //!< Put the user word at index
        __ffs_entry_user_put(ffs, ffs_ops->part_entry, ffs_ops->user,
                                ffs_ops->value);
	fprintf(ffs_ops->log, "%s: Put user word at '%d' --> %d success\n",
			__func__, ffs_ops->user, ffs_ops->value);
	//!< Required to check what is get
	fprintf(stdout, "UW: %d-->%d\n", ffs_ops->user, ffs_ops->value);

        __ffs_close(ffs);

        return rc;
}

void
usage(void) {
	printf("This program is a unit test tool, its callers responsibility "
		"to pass the correct parameters.\nNote: No usage errors are "
		"displayed, any mistake in params may result in unexpected "
		"results\n");
}

int
main(int argc, char * argv[]) {

	int32_t rc = 0;
	ffs_ops_t ffs_ops;

	memset(&ffs_ops, 0, sizeof(ffs_ops_t));
	ffs_ops.part_off   = PART_OFFSET;
	ffs_ops.log = log_open();

	while ((argc > 1) && (argv[1][0] == '-'))
	{
		switch (argv[1][1])
		{

//test_libffs -c pnor -O part_offset -s dev_size -b block size
//test_libffs -c pnor -O 4128768 -s 67108864 -b 65536
			case 'c':
				ffs_ops.nor_image  = argv[2]; //!< nor image
				ffs_ops.part_off = atoll(argv[4]);
				ffs_ops.device_size = atoi(argv[6]);
				ffs_ops.blk_sz = atoi(argv[8]);
				rc = create_partition(&ffs_ops);
				if(rc == FFS_ERROR) {
					goto out;
				}
				break;
//test_libffs -a pnor -O part_offset -n part_name -t part_type
//test_libffs -a sunray.pnor -O 4128768 -n boot0 -t logical
//test_libffs -a pnor -O part_off -n part_name -t type -s size -o entry_off
//test_libffs -a sunray.pnor -O 4128768 -n boot0/bootenv -t data -s 1048576 -o 0
			case 'a':
				ffs_ops.nor_image  = argv[2]; //!< nor image
				ffs_ops.part_off = atoll(argv[4]);
				ffs_ops.part_entry = argv[6];
				if (!strcasecmp(argv[8], "logical")) {
					ffs_ops.type = FFS_TYPE_LOGICAL;
					ffs_ops.entry_sz = 0;
					ffs_ops.entry_off = 0;
				}
				else if (!strcasecmp(argv[8], "data")) {
					ffs_ops.type = FFS_TYPE_DATA;
					ffs_ops.entry_sz = atol(argv[10]);
					ffs_ops.entry_off = atoll(argv[12]);
				}
				rc = add_partition_entry(&ffs_ops);
				if(rc == FFS_ERROR) {
					goto out;
				}
				break;
//test_libffs -r pnor -O part_off -n part_name -o out_file
//test_libffs -r sunray.pnor -O 4128768 -n boot0/bootenv -o out_file
			case 'r':
				ffs_ops.nor_image  = argv[2]; //!< nor image
				ffs_ops.part_off = atoll(argv[4]);
				ffs_ops.part_entry = argv[6]; //!< part entry
				ffs_ops.o_file     = argv[8]; //!< out put file
				ffs_ops.i_file     = NULL;
				rc = read_partition(&ffs_ops);
				if(rc == FFS_ERROR) {
					goto out;
				}
				break;
//test_libffs -w pnor -O part_off -n part_name -i in_file
//test_libffs -w sunray.pnor -O 4128768 -n boot0/bootenv -i in_file
			case 'w':
				ffs_ops.nor_image  = argv[2]; //!< nor image
				ffs_ops.part_off = atoll(argv[4]);
				ffs_ops.part_entry = argv[6]; //!< part entry
				ffs_ops.i_file     = argv[8]; //!< out put file
				ffs_ops.o_file     = NULL;
				rc = write_partition(&ffs_ops);
				if(rc == FFS_ERROR) {
					goto out;
				}
//test_libffs -l pnor -O part_off
//test_libffs -l sunray.pnor -O 4128768
			case 'l':
				ffs_ops.nor_image  = argv[2]; //!< nor image
				ffs_ops.part_off = atoll(argv[4]);
				rc = list_partition(&ffs_ops);
				if(rc == FFS_ERROR) {
					goto out;
				}
				break;
//test_libffs -h pnor -O part_off -n part_entry
//test_libffs -h sunray.pnor -O 4128768 -n boot0/bootenv
			case 'h':
				ffs_ops.nor_image  = argv[2]; //!< nor image
				ffs_ops.part_off = atoll(argv[4]);
				ffs_ops.part_entry = argv[6]; //!< part entry
				rc = hexdump_entry(&ffs_ops);
				if(rc == FFS_ERROR) {
					goto out;
				}
				break;
//test_libffs -d pnor -O part_off -n part_entry
//test_libffs -d sunray.pnor -O 4128768 -n boot0/bootenv
			case 'd':
				ffs_ops.nor_image  = argv[2]; //!< nor image
				ffs_ops.part_off = atoll(argv[4]);
				ffs_ops.part_entry = argv[6]; //!< part entry
				rc = delete_entry(&ffs_ops);
				if(rc == FFS_ERROR) {
					goto out;
				}
				break;
//test_libffs -m pnor -O part_off -n part_name -u index -g
//test_libffs -m sunray.pnor -O 4128768 -n boot0/bootenv -u 0 -g
//test_libffs -m pnor -O part_off -n part_name -u index -p -v some_value
//test_libffs -m sunray.pnor -O 4128768 -n boot0/bootenv -u 0 -p -v 1024
			case 'm':
				ffs_ops.nor_image  = argv[2]; //!< nor image
				ffs_ops.part_off = atoll(argv[4]);
				ffs_ops.part_entry = argv[6]; //!< part entry
				ffs_ops.user = atol(argv[8]);
				if(!strcmp(argv[9], "-g")) {
					rc = modify_entry_get(&ffs_ops);
					if(rc == FFS_ERROR) {
						goto out;
					}
				} else if(!strcmp(argv[9], "-p")) {
					ffs_ops.value = atol(argv[11]);
					rc = modify_entry_put(&ffs_ops);
					if(rc == FFS_ERROR) {
						goto out;
					}
				}
				break;
			default:
				usage();
				break;
		}
		break;
	}

out:
	if(ffs_ops.log != NULL) {
		fclose(ffs_ops.log);
	}
	return rc;
}
