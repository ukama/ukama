/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/cunit/ecc.c $                                            */
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

#include <stdio.h>
#include <stdlib.h>

#include <clib/type.h>

#include <clib/ecc.h>
#include <clib/misc.h>

#include <CUnit/Basic.h>

#define COUNT	10000
#define SEED    41

static int init_ecc(void) {
    return 0;
}

static int clean_ecc(void) {
    return 0;
}

static void ecc_1(void) {
    int size[] = {7, 8, 16, 1024+16, 4096+32, 8*1024*4+64, 0};

    for (int * s = size; *s != 0; s++) {
        unsigned char in[*s];
        unsigned char out[*s + (*s / 8)];

        ssize_t in_sz = sizeof in;
        ssize_t out_sz = sizeof out;

        memset(in, 0, in_sz);
        memset(out, 0, out_sz);

        FILE * f = fopen("/dev/urandom", "r");
        CU_ASSERT_FATAL(f != NULL);
        CU_ASSERT_FATAL(fread(in, in_sz, 1, f) == 1);
        fclose(f);

        ssize_t rc = sfc_ecc_inject(out, out_sz, in, in_sz);
        if ((*s % 8) != 0) {
            CU_ASSERT(rc == -1);
        } else {
            CU_ASSERT(rc == in_sz + (in_sz / 8));
        }

        unsigned char cmp[*s];
        ssize_t cmp_sz = sizeof cmp;
        memset(cmp, 0, cmp_sz);

        rc = sfc_ecc_remove(cmp, cmp_sz, out, out_sz);
        if ((out_sz % 9) != 0) {
            CU_ASSERT(rc == -1);
        } else {
            CU_ASSERT(rc == in_sz)
            CU_ASSERT_FATAL(memcmp(in, cmp, in_sz) == 0);
#ifdef DEBUG
            dump_memory(stdout, 0, in, in_sz);
            dump_memory(stdout, 0, cmp, cmp_sz);
#endif
        }
    }
}

void ecc_test(void) {
    CU_pSuite suite = CU_add_suite("ecc", init_ecc, clean_ecc);
    if (NULL == suite)
	return;

    if (CU_add_test(suite, "test of --> ecc_1", ecc_1) == NULL) return;
}
