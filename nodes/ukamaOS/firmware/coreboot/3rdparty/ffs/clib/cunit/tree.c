/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/cunit/tree.c $                                           */
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

#define COUNT 	30000
#define SEED    22

slab_t slab;

typedef struct {
    tree_node_t node;
    int i;
    float f;
} data_t;

static int init_tree(void) {
    slab_init(&slab, "my_slab", sizeof(data_t), 4096);
    return 0;
}

static int clean_tree(void) {
    slab_delete(&slab);
    return 0;
}

static void __insert(tree_t * t, int i) {
    data_t * d = (data_t *)slab_alloc(&slab);

    d->i = i;
    d->f = (float)i;

    tree_node_init(&d->node, (const void *)(d->i));

    if (tree_insert(t, &d->node) < 0) {
        err_t * err = err_get();
        fprintf(stderr, "%s(%d): UNEXPECTED: %.*s\n",
                err_file(err), err_line(err), err_size(err),
		(const char *)err_data(err));
    }
}

static data_t * __remove(tree_t * t, int i) {
    tree_node_t * n = tree_find(t, (const void *)i);
    if (n == NULL) tree_dump(t, stdout);
    CU_ASSERT_PTR_NOT_NULL_FATAL(n);

    tree_remove(t, n);
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

static void tree_1(void) {
    tree_t t;
    tree_init(&t, compare);

    CU_ASSERT(tree_min(&t) == NULL);
    CU_ASSERT(tree_max(&t) == NULL);
    CU_ASSERT(t.compare != NULL);

    CU_ASSERT(tree_root(&t) == NULL);
    CU_ASSERT(tree_empty(&t) == true);
    CU_ASSERT(tree_size(&t) == 0);
}

static void tree_2(void) {
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

static void tree_3(void) {
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

static void tree_4(void) {
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

static void tree_5(void) {
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

static void tree_6(void) {
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

static void tree_7(void) {
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


void tree_test(void) {
    CU_pSuite suite = CU_add_suite("tree", init_tree, clean_tree);
    if (NULL == suite)
	return;

    if (CU_add_test(suite, "test of --> tree_1", tree_1) == NULL)
        goto error;
    if (CU_add_test(suite, "test of --> tree_2", tree_2) == NULL)
        goto error;
    if (CU_add_test(suite, "test of --> tree_3", tree_3) == NULL)
        goto error;
    if (CU_add_test(suite, "test of --> tree_4", tree_4) == NULL)
        goto error;
    if (CU_add_test(suite, "test of --> tree_5", tree_5) == NULL)
        goto error;
    if (CU_add_test(suite, "test of --> tree_6", tree_6) == NULL)
        goto error;
    if (CU_add_test(suite, "test of --> tree_7", tree_7) == NULL)
        goto error;

    if (false) {
        err_t * err;
error:
        err = err_get();
        fprintf(stderr, "%s(%d): UNEXPECTED: %.*s\n",
                err_file(err), err_line(err), err_size(err),
		(const char *)err_data(err));
    }

}
