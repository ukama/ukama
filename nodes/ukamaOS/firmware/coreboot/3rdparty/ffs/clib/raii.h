/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/raii.h $                                                 */
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

/*
 *   File: raii.h
 * Author:
 *  Descr: RAII macro
 *   Note:
 *   Date: 04/03/13
 */

#ifndef __RAII_H__
#define __RAII_H__

#include "attribute.h"

#define CLEANUP(type,name,func) 			\
    void __cleanup_ ## name (type * __v) { 		\
		if (__v != NULL)			\
			func(__v); 			\
    }							\
    type name __cleanup(__cleanup_##name)

#define RAII(type,name,ctor,dtor) 			\
    void __cleanup_##name (type * __v) { 		\
		if (__v != NULL && *__v != NULL) 	\
			dtor(*__v); 			\
    }							\
    type name __cleanup(__cleanup_##name) = (ctor)

#endif				/* __RAII_H__ */
