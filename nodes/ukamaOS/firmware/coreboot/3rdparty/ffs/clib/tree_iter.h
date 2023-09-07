/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/tree_iter.h $                                            */
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

/*! @file tree_iter.h
 *  @brief Binary Tree Iterator
 *  @details For example,
 *  @code
 *  ...
 *  data_t * d;
 *  ...
 *  tree_iter_t it;
 *  tree_iter_init(&it, &t, TI_FLAG_DFT_FWD);
 *  ...
 *  tree_for_each(&it, d, node) {
 *      printf("depth first (FWD) i[%d] f[%f]\n", d->i, d->f);
 *  }
 *  ...
 *  @endcode
 *  @author Shaun Wetzstein <shaun@us.ibm.com>
 *  @date 2010-2012
 */

#ifndef __TREE_ITER_H__
#define __TREE_ITER_H__

#include <stdint.h>
#include <stdbool.h>

#include "clib/builtin.h"
#include "compare.h"
#include "type.h"

#include "tree.h"

/* ======================================================== */

typedef struct tree_iter tree_iter_t;	//!< Alias for the @em tree_iter class

/*!
 * @brief tree iterator
 * @details Binary tree container iterator
 */
struct tree_iter {
	tree_t *tree;		//!< Tree container to iterate

	tree_node_t *node;	//!< Current position in the iteration
	tree_node_t *safe;

	uint32_t flags;		//!< Iterator configuration flags
};

/* ======================================================== */

#define TI_FLAG_NONE		0x00000000	//!< All flag mask
#define TI_FLAG_FWD		0x00000001	//!< Depth-first-traversal Forward (FWD) flag mask
#define TI_FLAG_BWD		0x00000002	//!< Depth-first-traversal Backward (BWD) flag mask
#define TI_FLAG_MASK		0x00000003	//!< Depth-first-traversal All flag mask

/*!
 * @brief Initializes an @em tree_iter iterator object
 * @memberof tree_iter
 * @param self [in] tree_iter object @em self pointer
 * @param tree [in] tree container object to iterate
 * @param flags [in] iterator configuration flags
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 * @throws UNEXPECTED if @em tree pointer is NULL
 */
extern int tree_iter_init(tree_iter_t * self, tree_t * tree, uint32_t flags)
/*! @cond */
__nonnull((1, 2)) /*! @endcond */ ;

/*!
 * @brief Resets an @em tree iterator object
 * @memberof tree_iter
 * @param self [in] tree_iter object @em self pointer
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
extern int tree_iter_clear(tree_iter_t * self)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Return a pointer to a @em tree_node element at the current
 * iterator position
 * @memberof tree_iter
 * @param self [in] tree_iter object @em self pointer
 * @return non-NULL on success, NULL otherwise
 * @throws UNEXPECTED if @em self pointer is NULL
 */
extern tree_node_t *tree_iter_elem(tree_iter_t * self)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;

/*!
 * @fn tree_node_t * tree_iter_inc(tree_iter_t * self, size_t count = 1)
 * @brief Increment the position of an @em tree iterator
 * @details If the second (2nd) parameter is omitted, the iterator is
 * incremented by one (1) position.
 * @details For example,
 * @code
 * ...
 * tree_t l;
 * tree_init(&l);
 * ...
 * tree_iter_t it;
 * tree_iter_init(&it, &l, LI_FLAG_FWD);
 * ...
 * tree_iter_inc(&it);
 * ...
 * @endcode
 * @memberof tree_iter
 * @param self [in] tree_iter object @em self pointer
 * @param count [in] Number of positions to increment (optional)
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
#define tree_iter_inc(...) STRCAT(tree_iter_inc, NARGS(__VA_ARGS__))(__VA_ARGS__)
extern tree_node_t *tree_iter_inc1(tree_iter_t * self)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;
extern tree_node_t *tree_iter_inc2(tree_iter_t * self, size_t count)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;

/*!
 * @fn tree_node_td * tree_iter_dec(tree_iter_t * self, size_t count = 1)
 * @brief decrement the position of an @em list iterator
 * @note If the second (2nd) parameter is omitted, the iterator is decremented
 * by one (1) position.
 * @details For example,
 * @code
 * ...
 * tree_t l;
 * tree_init(&l);
 * ...
 * tree_iter_t it;
 * tree_iter_init(&it, &l, LI_FLAG_FWD);
 * ...
 * tree_iter_dec(&it, 3);
 * ...
 * @endcode
 * @memberof tree_iter
 * @param self [in] tree_iter object @em self pointer
 * @param count [in] Number of positions to decrement (optional)
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
#define tree_iter_dec(...) STRCAT(tree_iter_dec, NARGS(__VA_ARGS__))(__VA_ARGS__)
extern tree_node_t *tree_iter_dec1(tree_iter_t * self)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;
extern tree_node_t *tree_iter_dec2(tree_iter_t * self, size_t count)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;

/*!
 * @def tree_for_each(it, i, m)
 * @hideinitializer
 * @brief Tree for-each algorithm
 * @param it [in] Tree iterator object
 * @param i [in] Tree element variable
 * @param m [in] Member name
 */
#define tree_for_each(it, i, m)						      \
    for (tree_iter_clear(it), i = container_of_var(tree_iter_elem(it), i, m); \
         tree_iter_elem(it) != NULL;					      \
         i = container_of_var(tree_iter_inc(it), i, m))

/* ======================================================== */

#endif				/* __TREE_ITER_H__ */
