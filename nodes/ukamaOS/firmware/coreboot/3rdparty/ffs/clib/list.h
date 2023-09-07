/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/list.h $                                                 */
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

/*!
 * @file list.h
 * @brief List container
 * @details
 * A List is a data structure (container) used for collecting a sequence of elements.
 * List allow for efficient insertion, removal and retreival of elements.
 *
 * @details For example,
 * @code
 * #include <clib/list.h>
 * #include <clib/list_iter.h>
 *
 * int main(const int argc, const char * argv[]) {
 *     typedef struct {
 *          list_node_t node;
 *          int i;
 *          float f;
 *     } data_t;
 *
 *     slab_t s;
 *     slab_init(&s, sizeof(data_t), 0);
 *
 *     list_t a;
 *     list_init(&a);
 *
 *     int i;
 *     for (i=0; i<10; i++) {
 *         data_t * d = (data_t *)slab_alloc(&s);
 *
 *         d->i = i;
 *         d->f = (float)i;
 *
 *         list_add_tail(&l, &d->node);
 *     }
 *
 *     data_t * d;
 *     list_for_each(&l, d, node) {
 *         printf("i: %d f: %f\n", d->i, d->f);
 *     }
 *
 *     list_dump(&l, stdout);
 *     slab_delete(&s);
 *
 *     return 0;
 * }
 * @endcode
 * @date 2010-2011
 */

#ifndef __LIST_H__
#define __LIST_H__

#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>
#include <assert.h>

#include "clib/builtin.h"
#include "type.h"

/* ==================================================================== */

#define INIT_LIST_NODE	{NULL,NULL}
#define INIT_LIST	{INIT_LIST_NODE}

typedef struct list_node list_node_t;	//!< Alias for the @em list_node class
typedef struct list list_t;	//!< Alias for the @em list class

/*!
 * @brief list node
 * @details Primitive types cannot be stored in the @em list container, instead the user must
 * embed a @em list_node object within the stored object.
 */
struct list_node {
	list_node_t *prev;	//!< Reference the previous list_node
	list_node_t *next;	//!< Reference the next list_node
};

/*!
 * @brief list container
 */
struct list {
	list_node_t node;
};

/* ==================================================================== */

/*!
 * @brief Return a pointer to the node's next node
 * @details For example,
 * @code
 * ...
 * list_node_t * node = list_head(&l);
 * list_node_t * next = list_node_next(node);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list node object @em self pointer
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline list_node_t *list_node_next(list_node_t * self)
/*! @cond */ __nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Return a pointer to the node's previous node
 * @details For example,
 * @code
 * ...
 * list_node_t * node = list_head(&l);
 * list_node_t * next = list_node_prev(node);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list node object @em self pointer
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline list_node_t *list_node_prev(list_node_t * self)
/*! @cond */ __nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Constructs an @em list container object
 * @details For example,
 * @code
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline void list_init(list_t * self)
/*! @cond */ __nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Inserts a node at the head of the @em list container
 * @details For example,
 * @code
 * ...
 * typedef struct {
 *    ...
 *    list_node_t node;
 *    ...
 * } data_t;
 * ...
 * list_t l;
 * list_init(&l);
 *
 * data_t * d = (data_t *) MALLOC(sizeof(*d));
 * list_add_head(&l, &d->node);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @param node [in] list_node object to insert
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline void list_add_head(list_t * self, list_node_t * node)
/*! @cond */ __nonnull((1, 2)) /*! @endcond */ ;

/*!
 * @brief Inserts a node at the tail of the @em list container
 * @details For example,
 * @code
 * ...
 * typedef struct {
 *    ...
 *    list_node_t node;
 *    ...
 * } data_t;
 * ...
 * list_t l;
 * list_init(&l);
 *
 * data_t * d = (data_t *) MALLOC(sizeof(*d));
 * list_add_tail(&l, &d->node);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @param node [in] list_node object to insert
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline void list_add_tail(list_t * self, list_node_t * node)
/*! @cond */ __nonnull((1, 2)) /*! @endcond */ ;

/*!
 * @brief Remove a node at the head of the @em list container
 * @details For example,
 * @code
 * ...
 * typedef struct {
 *    ...
 *    list_node_t node;
 *    ...
 * } data_t;
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * data_t * d = (data_t *)list_remove_head(&l);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @return non-NULL on success, NULL otherwise
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline list_node_t *list_remove_head(list_t * self)
/*! @cond */ __nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Remove a node at the tail of the @em list container
 * @details For example,
 * @code
 * ...
 * typedef struct {
 *    ...
 *    list_node_t node;
 *    ...
 * } data_t;
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * data_t * d = (data_t *)list_remove_tail(&l);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @return non-NULL on success, NULL otherwise
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline list_node_t *list_remove_tail(list_t * self)
/*! @cond */ __nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Remove the list_node referenced by @em node from the @em list container
 * @details For example,
 * @code
 * ...
 * typedef struct {
 *    ...
 *    list_node_t node;
 *    ...
 * } data_t;
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * data_t * d = (data_t *) MALLOC(sizeof(*d));
 * list_add_tail(&l, &d->node);
 * list_remove_node(&l, &d->node);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @return non-NULL on success, NULL otherwise
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline list_node_t *list_remove_node(list_t * self, list_node_t * node)
/*! @cond */ __nonnull((1, 2)) /*! @endcond */ ;

/*!
 * @brief Test whether a @em list container is empty
 * @details For example,
 * @code
 * ...
 * typedef struct {
 *    ...
 *    list_node_t node;
 *    ...
 * } data_t;
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * if (list_empty(&l)) {
 * ...
 * }
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @return true if empty, false otherwise
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline bool list_empty(const list_t * self)
/*! @cond */ __nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Return the head list_node from the @em list container
 * @details For example,
 * @code
 * ...
 * typedef struct {
 *    ...
 *    list_node_t node;
 *    ...
 * } data_t;
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * list_node_t * node = list_head(&l);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @return non-NULL on success, NULL otherwise
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline list_node_t *list_head(list_t * self)
/*! @cond */ __nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Return the tail list_node from the @em list container
 * @details For example,
 * @code
 * ...
 * typedef struct {
 *    ...
 *    list_node_t node;
 *    ...
 * } data_t;
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * list_node_t * node = list_tail(&l);
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @return non-NULL on success, NULL otherwise
 * @throws UNEXPECTED if @em self pointer is NULL
 */
static inline list_node_t *list_tail(list_t * self)
/*! @cond */ __nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Pretty print and dump the contents of the @em list to output stream @em out
 * @details For example,
 * @code
 * ...
 * typedef struct {
 *    ...
 *    list_node_t node;
 *    ...
 * } data_t;
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * list_dump(&l, stdout) {
 * ...
 * }
 * ...
 * @endcode
 * @memberof list
 * @param self [in] list object @em self pointer
 * @param out [in] output stream
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
extern void list_dump(list_t * self, FILE *)
/*! @cond */ __nonnull((1)) /*! @endcond */ ;

/*!
 * @def list_entry(node_pointer, containing_type, member)
 * @hideinitializer
 * @brief Returns a pointer to the containing structure of a node
 * @param n [in] Pointer to a list_node object
 * @param T [in] Type of the containing structure
 * @param m [in] Name of the list_node member
 */
#define list_entry(n, T, m) container_of(n, T, m)

#if 0
/*!
 * @def list_head(list_pointer, containing_type, member)
 * @hideinitializer
 * @brief Returns a pointer to the head element in the @em list
 * @param n [in] Pointer to a list_node object
 * @param T [in] Type of the containing structure
 * @param m [in] Name of the list_node member
 */
#define list_head(l, T, m) ({								\
    list_node_t * __h = (list_empty(l) ? NULL : list_entry((l)->node.next, T, m))	\
    __h;										\
			    })

/*!
 * @def list_tail(list_pointer, containing_type, member)
 * @hideinitializer
 * @brief Returns a pointer to the tail element in the @em list
 * @param n [in] Pointer to a list_node object
 * @param T [in] Type of the containing structure
 * @param m [in] Name of the list_node member
 */
#define list_tail(h, T, m) ({								\
    list_node_t * __t = (list_empty(l) ? NULL : list_entry((l)->node.prev, T, m))	\
    __t;										\
			    })
#endif

/* ==================================================================== */

/*! @cond */
static inline list_node_t *list_node_next(list_node_t * self)
{
	assert(self != NULL);
	return self->next;
}

static inline list_node_t *list_node_prev(list_node_t * self)
{
	assert(self != NULL);
	return self->prev;
}

static inline void list_init(list_t * self)
{
	assert(self != NULL);
	self->node.prev = self->node.next = &self->node;
}

static inline void list_add_head(list_t * self, list_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	node->next = self->node.next;
	node->prev = &self->node;
	self->node.next->prev = node;
	self->node.next = node;
}

static inline void list_add_tail(list_t * self, list_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	node->next = &self->node;
	node->prev = self->node.prev;
	self->node.prev->next = node;
	self->node.prev = node;
}

static inline list_node_t *list_remove_head(list_t * self)
{
	assert(self != NULL);

	list_node_t *node = self->node.next;

	if (node != NULL) {
		self->node.next = self->node.next->next;
		self->node.next->prev = node->prev;
		node->next = node->prev = NULL;
	}

	return node;
}

static inline list_node_t *list_remove_tail(list_t * self)
{
	assert(self != NULL);

	list_node_t *node = self->node.prev;

	if (node != NULL) {
		self->node.prev = self->node.prev->prev;
		self->node.prev->next = node->next;
		node->prev = node->next = NULL;
	}

	return node;
}

static inline list_node_t *list_remove_node(list_t * self, list_node_t * node)
{
	assert(self != NULL);
	assert(node != NULL);

	node->next->prev = node->prev;
	node->prev->next = node->next;
	node->next = node->prev = NULL;

	return node;
}

static inline bool list_empty(const list_t * self)
{
	return self->node.next == &self->node;
}

static inline list_node_t *list_head(list_t * self)
{
	assert(self != NULL);
	return self->node.next;
}

static inline list_node_t *list_tail(list_t * self)
{
	assert(self != NULL);
	return self->node.prev;
}

/*! @endcond */

/* ==================================================================== */

#endif				/* __LIST_H__ */
