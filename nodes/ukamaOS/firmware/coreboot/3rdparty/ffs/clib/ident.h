/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/ident.h $                                                */
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
 * @file ident.h
 * @brief Identification object
 * @author Shaun Wetzstein <shaun@us.ibm.com>
 * @date 2008-2011
 */

#ifndef __IDENT_H__
#define __IDENT_H__

#include <stdint.h>

#define IDENT_MAGIC_0		0
#define IDENT_MAGIC_1		1
#define IDENT_MAGIC_2		2
#define IDENT_MAGIC_3		3
#define IDENT_MAJOR		4
#define IDENT_MINOR		5
#define IDENT_PATCH		6
#define IDENT_FLAGS		7
#define IDENT_SIZE		8

#define INIT_IDENT		{0,}

#define MAGIC_CHECK(i, m) ({					\
    bool rc = (((i)[IDENT_MAGIC_0] != (m)[IDENT_MAGIC_0]) ||	\
               ((i)[IDENT_MAGIC_1] != (m)[IDENT_MAGIC_1]) ||	\
               ((i)[IDENT_MAGIC_2] != (m)[IDENT_MAGIC_2]) ||	\
               ((i)[IDENT_MAGIC_3] != (m)[IDENT_MAGIC_3]));	\
    rc;								\
			    })

typedef uint8_t ident_t[IDENT_SIZE];

#endif				/* __IDENT_H__ */
