/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/align.h $                                                */
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
 * @file align.h
 * @brief Alignment helpers
 * @author Shaun Wetzstein <shaun@us.ibm.com>
 * @date 2010-2011
 */

#ifndef __ALIGN_H__
#define __ALIGN_H__

/*!
 * @def alignof(t)
 * @hideinitializer
 * @brief Returns the alignment of an object or minimum alignment required by a type
 * @param t [in] Object or type name
 */
#define alignof(t)			__alignof__(t)

#ifndef offsetof

/*!
 * @def offsetof(t,m)
 * @hideinitializer
 * @brief Returns the offset of a member within a structure
 * @param t [in] Structure type name
 * @param m [in] Member name within a structure
 */
#define offsetof(t,m)			__builtin_offsetof(t, m)
#endif

#endif				/* __ALIGN_H__ */
