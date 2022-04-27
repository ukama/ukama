/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/test/list.c $                                            */
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
#include <string.h>

#include <clib/slab.h>
#include <clib/list.h>
#include <clib/list_iter.h>

struct data {
    list_node_t node;
    int i;
    float f;
};
typedef struct data data_t;

int main (void) {
//    slab_t s = INIT_SLAB;
//    slab_init(&s, "my_slab", sizeof(data_t), 0);

    list_t l = INIT_LIST;
    list_init(&l);

    int i;
    for (i=0; i<10; i++) {
//      data_t * d = (data_t *)slab_alloc(&s);
	data_t * d = (data_t *)malloc(sizeof(*d));

        d->i = i;
        d->f = (float)i;

        list_add_tail(&l, &d->node);
    }

    list_iter_t it;
    list_iter_init(&it, &l, LI_FLAG_FWD);

    data_t * d;
    list_for_each(&it, d, node) {
        printf("i: %d f: %f\n", d->i, d->f);
    }

    while (list_empty(&l) == false) {
	data_t * d = container_of(list_remove_tail(&l), data_t, node);
        printf("i: %d f: %f\n", d->i, d->f);
    }

//    slab_dump(&s, stdout);
//    slab_delete(&s);

    return 0;
}

