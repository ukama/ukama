/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/cunit/splay.c $                                          */
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

#include <clib/libclib.h>

#include <clib/slab.h>
#include <clib/tree.h>
#include <clib/tree_iter.h>

#include <CUnit/Basic.h>

#define COUNT	30000
#define SEED    22

slab_t slab;

typedef struct {
    tree_node_t node;
    int i;
    float f;
} data_t;

static int init_splay(void) {
    slab_init(&slab, "my_slab", sizeof(data_t), 4096);
    return 0;
}

static int clean_splay(void) {
    slab_delete(&slab);
    return 0;
}

static void __insert(tree_t * t, int i) {
    data_t * d = (data_t *)slab_alloc(&slab);

    d->i = i;
    d->f = (float)i;

    tree_node_init(&d->node, (const void *)(d->i));

    if (splay_insert(t, &d->node) < 0) {
        tree_dump(t, stdout);

        fprintf(stdout, "key: %d root->key: %d\n",
                i, (int)tree_root(t)->key);

        err_t * err = err_get();
        fprintf(stderr, "%s(%d): %.*s\n",
                err_file(err), err_line(err), err_size(err),
		(const char *)err_data(err));
    }
}

static data_t * __remove(tree_t * t, int i) {
    tree_node_t * n = tree_find(t, (const void *)i);
    if (n == NULL) tree_dump(t, stdout);
    CU_ASSERT_PTR_NOT_NULL_FATAL(n);

    splay_remove(t, n);
    CU_ASSERT_PTR_NULL(tree_node_parent(n));
    CU_ASSERT_PTR_NULL(tree_node_left(n));
    CU_ASSERT_PTR_NULL(tree_node_right(n));
    CU_ASSERT_PTR_NOT_NULL(tree_node_key(n));

    data_t * d = container_of(n, data_t, node);
    CU_ASSERT_PTR_NOT_NULL_FATAL(n);

    if (0 <= i)
        CU_ASSERT(d->i == i);

    return d;
}

static int compare(const void * v1, const void * v2) {
    const int i1 = (const int)v1, i2 = (const int)v2;
    return i1 - i2;
}

static void splay_1(void) {
    tree_t t;
    tree_init(&t, compare);

    CU_ASSERT(tree_min(&t) == NULL);
    CU_ASSERT(tree_max(&t) == NULL);
    CU_ASSERT(t.compare != NULL);

    CU_ASSERT(tree_root(&t) == NULL);
    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);
}

static void splay_2(void) {
    tree_t t;
    tree_init(&t, compare);

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    for (int i=1; i<=COUNT; i++)
        __insert(&t, i);

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_root(&t) != NULL);
    CU_ASSERT(tree_size(&t) == COUNT);

    for (int i=1; i<=COUNT; i++) {
        CU_ASSERT(i == (int)tree_min(&t)->key);
        CU_ASSERT(COUNT == (int)tree_max(&t)->key);
        __remove(&t, (int)tree_min(&t)->key);
        CU_ASSERT(tree_size(&t) + i == COUNT);
    }

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    for (int i=1; i<=COUNT; i++)
        __insert(&t, i);

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_root(&t) != NULL);
    CU_ASSERT(tree_size(&t) == COUNT);

    for (int i=1; i<=COUNT; i++) {
        CU_ASSERT(1 == (int)tree_min(&t)->key);
        CU_ASSERT(COUNT - i + 1 == (int)tree_max(&t)->key);
        __remove(&t, (int)tree_max(&t)->key);
        CU_ASSERT(tree_size(&t) + i == COUNT);
    }

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);
}

static void splay_3(void) {
    tree_t t;
    tree_init(&t, compare);

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    for (int i=1; i<=COUNT; i++)
        __insert(&t, COUNT - i + 1);

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_root(&t) != NULL);
    CU_ASSERT(tree_size(&t) == COUNT);

    for (int i=1; i<=COUNT; i++) {
        CU_ASSERT(1 == (int)tree_min(&t)->key);
        CU_ASSERT(COUNT - i + 1 == (int)tree_max(&t)->key);
        __remove(&t, (int)tree_max(&t)->key);
        CU_ASSERT(tree_size(&t) + i == COUNT);
    }

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    for (int i=1; i<=COUNT; i++)
        __insert(&t, COUNT - i + 1);

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_root(&t) != NULL);
    CU_ASSERT(tree_size(&t) == COUNT);

    for (int i=1; i<=COUNT; i++) {
        CU_ASSERT(i == (int)tree_min(&t)->key);
        CU_ASSERT(COUNT == (int)tree_max(&t)->key);
        __remove(&t, (int)tree_min(&t)->key);
        CU_ASSERT(tree_size(&t) + i == COUNT);
    }

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);
}

static void splay_4(void) {
    tree_t t;
    tree_init(&t, compare);

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    srandom(SEED);
    for (int i=1; i<=COUNT; i++)
        __insert(&t, random());

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_size(&t) == COUNT);

    for (int i=1; i<=COUNT; i++) {
         __remove(&t, (int)tree_min(&t)->key);
        CU_ASSERT(tree_size(&t) + i == COUNT);
    }

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    srandom(SEED);
    for (int i=1; i<=COUNT; i++)
        __insert(&t, random());

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_size(&t) == COUNT);

    for (int i=1; i<=COUNT; i++) {
         __remove(&t, (int)tree_max(&t)->key);
        CU_ASSERT(tree_size(&t) + i == COUNT);
    }

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    srandom(SEED);
    for (int i=1; i<=COUNT; i++)
        __insert(&t, random());

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_size(&t) == COUNT);

    for (int i=1; i<=COUNT; i++) {
         __remove(&t, (int)tree_root(&t)->key);
        CU_ASSERT(tree_size(&t) + i == COUNT);
    }

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    srandom(SEED);
    for (int i=1; i<=COUNT; i++)
        __insert(&t, random());

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_size(&t) == COUNT);

    srandom(SEED);
    for (int i=1; i<=COUNT; i++) {
         __remove(&t, random());
        CU_ASSERT(tree_size(&t) + i == COUNT);
    }

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);
}

static void splay_5(void) {
    tree_t t;
    tree_init(&t, compare);

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    srandom(SEED);
    for (int i=1; i<=COUNT; i++)
        __insert(&t, random());

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_size(&t) == COUNT);
    data_t * d;

    tree_iter_t it;

    int key = 0;
    tree_iter_init(&it, &t, TI_FLAG_FWD);
    tree_for_each(&it, d, node) {
        CU_ASSERT(key < d->i);
        key = d->i;
    }

    key = INT32_MAX;
    tree_iter_init(&it, &t, TI_FLAG_BWD);
    tree_for_each(&it, d, node) {
        CU_ASSERT(d->i < key);
        key = d->i;
    }
}

static void splay_6(void) {
    tree_t t;
    tree_init(&t, compare);

    data_t * d;

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    srandom(SEED);
    for (int i=1; i<=COUNT; i++)
        __insert(&t, random());

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_size(&t) == COUNT);

    for (int i=1; i<=COUNT; i++) {
        tree_iter_t it;
        tree_iter_init(&it, &t, TI_FLAG_FWD);

        int key = 0;
        tree_for_each(&it, d, node) {
            CU_ASSERT(key < d->i);
            key = d->i;
        }

        __remove(&t, (int)tree_min(&t)->key);

        if (0 < tree_size(&t)) {
            CU_ASSERT(tree_min(&t) != NULL);
        } else if (tree_size(&t) <= 0) {
            CU_ASSERT(tree_min(&t) == NULL);
        }

        CU_ASSERT(tree_size(&t) + i == COUNT);
    }
}

static void splay_7(void) {
    tree_t t;
    tree_init(&t, compare);

    data_t * d;

    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);

    srandom(SEED);
    for (int i=1; i<=COUNT; i++)
        __insert(&t, random());

    CU_ASSERT(tree_empty(&t) == false);
    CU_ASSERT(tree_size(&t) == COUNT);

    for (int i=1; i<=COUNT; i++) {
        tree_iter_t it;

        int key = INT32_MAX;
        tree_iter_init(&it, &t, TI_FLAG_BWD);
        tree_for_each(&it, d, node) {
            CU_ASSERT(d->i < key);
            key = d->i;
        }

        __remove(&t, (int)tree_max(&t)->key);

        if (0 < tree_size(&t)) {
            CU_ASSERT(tree_max(&t) != NULL);
        } else if ( tree_size(&t) <= 0) {
            CU_ASSERT(tree_max(&t) == NULL);
        }

        CU_ASSERT(tree_size(&t) + i == COUNT);
    }
}

void splay_test(void) {
    CU_pSuite suite = CU_add_suite("splay", init_splay, clean_splay);
    if (NULL == suite)
	return;

    if (CU_add_test(suite, "test of --> splay_1", splay_1) == NULL) return;
    if (CU_add_test(suite, "test of --> splay_2", splay_2) == NULL) return;
    if (CU_add_test(suite, "test of --> splay_3", splay_3) == NULL) return;
    if (CU_add_test(suite, "test of --> splay_4", splay_4) == NULL) return;
    if (CU_add_test(suite, "test of --> splay_5", splay_5) == NULL) return;
    if (CU_add_test(suite, "test of --> splay_6", splay_6) == NULL) return;
    if (CU_add_test(suite, "test of --> splay_7", splay_7) == NULL) return;
}
