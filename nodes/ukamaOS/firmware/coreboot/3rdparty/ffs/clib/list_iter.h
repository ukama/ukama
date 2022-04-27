/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/list_iter.h $                                            */
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

/*! @file list_iter.h
 *  @brief List iterator
 *  @details For example,
 *  @code
 *  ...
 *  data_t * d;
 *  ...
 *  list_iter_t it;
 *  list_iter_init(&it, &t, LI_FLAG_FWD);
 *  ...
 *  list_for_each(&it, d, node) {
 *      printf("forward i[%d] f[%f]\n", d->i, d->f);
 *  }
 *  ...
 *  @endcode
 *  @author Shaun Wetzstein <shaun@us.ibm.com>
 *  @date 2010-2011
 */

#ifndef __LIST_ITER_H__
#define __LIST_ITER_H__

#include <stdbool.h>
#include <stdint.h>

#include "clib/builtin.h"
#include "list.h"
#include "nargs.h"

/* ======================================================================= */

typedef struct list_iter list_iter_t;	//!< Alias for the @em list_iter class

/*!
 * @brief list iterator
 * @details list iterator class
*/
struct list_iter {
	list_t *list;		//!< Reference to the target list object
	list_node_t *node;	//!< Current position of the iterator
	uint32_t flags;		//!< Iterator configuration flags
};

/* ======================================================================= */

#define LI_FLAG_NONE		0x00000000	//!< No flag mask
#define LI_FLAG_FWD		0x00000000	//!< Forward (FWD) flag mask
#define LI_FLAG_BWD		0x00000002	//!< Backward (BWD) flag mask
#define LI_FLAG_MASK		0x00000003	//!< All flags mask

/*!
 * @brief Initializes an @em list_iter iterator object
 * @details For example,
 * @code
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * list_iter_t it;
 * list_iter_init(&it, &l, LI_FLAG_FWD);
 * ...
 * @endcode
 * @memberof list_iter
 * @param self [in] list_iter object @em self pointer
 * @param list [in] list container object to iterate
 * @param flags [in] iterator configuration flags
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 * @throws UNEXPECTED if @em list pointer is NULL
 */
extern void list_iter_init(list_iter_t * self, list_t * list, uint32_t flags)
/*! @cond */
__nonnull((1, 2)) /*! @endcond */ ;

/*!
 * @brief Resets a @em list iterator object
 * @details For example,
 * @code
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * list_iter_t it;
 * list_iter_init(&it, &l, l_FLAG_FWD);
 * list_iter_clear(&it);
 * ...
 * @endcode
 * @memberof list_iter
 * @param self [in] list_iter object @em self pointer
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
extern void list_iter_clear(list_iter_t * self)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;

/*!
 * @brief Return a pointer to @em list element bytes at the current iterator position
 * @details For example,
 * @code
 * ...
 * list_t l;
 * list_init(&l, 4, 1024);
 * ...
 * list_iter_t it;
 * list_iter_init(&it, &l, LI_FLAG_FWD);
 * ...
 * data_t * d = (data_t *)list_iter_elem(&it);
 * ...
 * @endcode
 * @memberof list_iter
 * @param self [in] list_iter object @em self pointer
 * @return non-NULL on success, NULL otherwise
 * @throws UNEXPECTED if @em self pointer is NULL
 */
extern list_node_t *list_iter_elem(list_iter_t * self)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;

/*!
 * @fn list_node_t * list_iter_inc(list_iter_t * self, size_t count = 1)
 * @brief Increment the position of an @em list iterator
 * @details If the second (2nd) parameter is omitted, the iterator is incremented by one (1) position.
 * @details For example,
 * @code
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * list_iter_t it;
 * list_iter_init(&it, &l, LI_FLAG_FWD);
 * ...
 * list_iter_inc(&it);
 * ...
 * @endcode
 * @memberof list_iter
 * @param self [in] list_iter object @em self pointer
 * @param count [in] Number of positions to increment (optional)
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
/*! @cond */
#define list_iter_inc(...) STRCAT(list_iter_inc, NARGS(__VA_ARGS__))(__VA_ARGS__)
extern list_node_t *list_iter_inc1(list_iter_t * self)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;
extern list_node_t *list_iter_inc2(list_iter_t * self, size_t count)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;
/*! @endcond */

/*!
 * @fn list_node_td * list_iter_dec(list_iter_t * self, size_t count = 1)
 * @brief decrement the position of an @em list iterator
 * @note If the second (2nd) parameter is omitted, the iterator is decremented by one (1) position.
 * @details For example,
 * @code
 * ...
 * list_t l;
 * list_init(&l);
 * ...
 * list_iter_t it;
 * list_iter_init(&it, &l, LI_FLAG_FWD);
 * ...
 * list_iter_dec(&it, 3);
 * ...
 * @endcode
 * @memberof list_iter
 * @param self [in] list_iter object @em self pointer
 * @param count [in] Number of positions to decrement (optional)
 * @return None
 * @throws UNEXPECTED if @em self pointer is NULL
 */
/*! @cond */
#define list_iter_dec(...) STRCAT(list_iter_dec, NARGS(__VA_ARGS__))(__VA_ARGS__)
extern list_node_t *list_iter_dec1(list_iter_t * self)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;
extern list_node_t *list_iter_dec2(list_iter_t * self, size_t count)
/*! @cond */
__nonnull((1)) /*! @endcond */ ;
/*! @endcond */

/*!
 * @def list_for_each(it, i, m)
 * @hideinitializer
 * @brief List for-each algorithm
 * @param it [in] Tree iterator object
 * @param i [in] Tree element variable
 * @param m [in] Member name
 */
#if 0
#define list_for_each(l, i, m)				\
    for (i = container_of_var(l->node.next, i, m);	\
        &i->m != &(l)->node;				\
        i = container_of_var(i->m.next, i, m))
#define list_for_each_safe(l, i, n, m)			\
    for (i = container_of_var((l)->node.next, i, m),	\
            n = container_of_var(i->m.next, i, m);	\
        &i->m != &(l)->node;				\
        i = n, n = container_of_var(i->m.next, i, m))
#endif

#define list_for_each(it, i, m)					\
    for (i = container_of_var(list_iter_elem(it), i, m);	\
         list_iter_inc(it) != NULL;				\
         i = container_of_var(list_iter_elem(it), i, m))

#endif				/* __LIST_ITER_H__ */
