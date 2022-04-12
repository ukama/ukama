/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/test/tree.c $                                            */
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

#include <clib/type.h>
#include <clib/slab.h>

#include <clib/tree.h>
#include <clib/tree_iter.h>

struct data {
    tree_node_t node;
    int i;
    float f;
};
typedef struct data data_t;

int main (void) {
    slab_t s = INIT_SLAB;
    slab_init(&s, "my_slab", sizeof(data_t), 4096);

    tree_t t = INIT_TREE;
    tree_init(&t, default_compare);

    int i;
    for (i=0; i<25; i++) {
        data_t * d = (data_t *)slab_alloc(&s);

        printf("insert i[%d] --> %p\n", i, d);
        d->i = i;
        d->f = (float)i;        /* key */

        tree_node_init(&d->node, (const void *)(d->i));
        tree_insert(&t, &d->node);
    }

#if 1
    i = 6;
    tree_node_t * n = tree_find(&t, (const void *)(i));
    tree_remove(&t, n);
    n = tree_find(&t, (const void *)i);

    i = 2;
    tree_find(&t, (const void *)i);
    i = 8;
    tree_find(&t, (const void *)i);
#endif

    tree_dump(&t, stdout);

    data_t * d;

    tree_iter_t it;
    tree_iter_init(&it, &t, TI_FLAG_FWD);

    tree_for_each(&it, d, node) {
        printf("depth first (FWD) i[%d] f[%f]\n", d->i, d->f);
    }

//    tree_dump(&t, stdout);
    slab_delete(&s);

    return 0;
}

