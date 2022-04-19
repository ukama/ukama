/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/type.h $                                                 */
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

/*! @file type.h
 *  @brief Type helpers
 *  @author Shaun Wetzstein <shaun@us.ibm.com>
 *  @date 2010-2011
 */

#ifndef __TYPE_H__
#define __TYPE_H__

#include "align.h"

/*!
 * @def cast_type(e,T)
 * @hideinitializer
 * @brief Casts @em p to type @em T
 * @param e [in] Expression
 * @param T [in] Type name
 */
#define cast_type(e,T)						\
    ((T)(e))

/*!
 * @def check_type(e,T)
 * @hideinitializer
 * @brief Check the type of expression @em e is equal to type @em T
 * @param e [in] Expression
 * @param T [in] Type name
 */
#define check_type(e,T)						\
    ((typeof(e) *)0 != (T *)0)

/*!
 * @def check_types_match(e,T)
 * @hideinitializer
 * @brief Check the types of expressions @em e1 and @em e2 are equal
 * @param e1 [in] Expression
 * @param e2 [in] Expression
 */
#define check_types_match(e1,e2)				\
    ((typeof(e1) *) 0 != (typeof(e2) *)0)

/*!
 * @def container_of(member_ptr, containing_type, member)
 * @hideinitializer
 * @brief Returns a pointer to the containing structure of a member
 * @param p [in] Pointer to a strucure member
 * @param T [in] Type of the containing structure
 * @param m [in] Name of the member @em p references
 */
#define container_of(p,T,m)					\
    ((T *)((char *)(p) - offsetof(T,m)) - check_types_match(*(p),((T*)0)->m))

/*!
 * @def container_of_var(member_ptr, var, member)
 * @hideinitializer
 * @brief Returns a pointer to the containing structure of a member
 * @param p [in] Pointer to a strucure member
 * @param v [in] Pointer to a variable with the same of type as the containing structure
 * @param m [in] Name of the member @em p references
 */
#define container_of_var(p, v, m)				\
    container_of(p, typeof(*v), m)

#endif				/* __TYPE_H__ */
