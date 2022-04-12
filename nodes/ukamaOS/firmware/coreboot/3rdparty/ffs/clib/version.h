/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/version.h $                                              */
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

#ifndef __VERSION__H__
#define __VERSION__H__

#define VER_TO_MAJOR(x)	(((x) & 0xFF000000) >> 24)
#define VER_TO_MINOR(x)	(((x) & 0x00FF0000) >> 16)
#define VER_TO_PATCH(x)	(((x) & 0x0000FF00) >> 8)

#define MAJOR_TO_VER(x)	((0xFF & (x)) << 24)
#define MINOR_TO_VER(x)	((0xFF & (x)) << 16)
#define PATCH_TO_VER(x)	((0xFF & (x)) << 8)

#define VER(x,y,z) (MAJOR_TO_VER(x) | MINOR_TO_VER(y) | PATCH_TO_VER(z))

#endif				/* __VERSION__H__ */
