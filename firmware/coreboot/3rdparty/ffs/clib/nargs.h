/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/nargs.h $                                                */
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
 * @file nargs.h
 * @brief Macro utilities
 * @author Shaun Wetzstein <shaun@us.ibm.com>
 * @date 2010-2011
 */

#ifndef __NARGS_H__
#define __NARGS_H__

/*!
 * @def STRCAT(x,y)
 * @hideinitializer
 * @brief C string concatination of @em x and @em y
 * @param x [in] C-style string
 * @param y [in] C-style string
 */
/*! @cond */
#define STRCAT(x,y)	__C__(x, y)
#define __C__(x,y)	x ## y
/*! @endcond */

/*!
 * @def NARGS(...)
 * @hideinitializer
 * @brief Return the number of pre-process macro parameters
 */
/*! @cond */
#define __NARGS__(junk, _10, _9, _8, _7, _6, _5, _4, _3, _2, _1, _, ...) _
#define NARGS(...) __NARGS__(junk, ##__VA_ARGS__, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0)
/*! @endcond */

#endif				/* __NARGS_H__ */
