/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ecc/src/main.h $                                              */
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
 *    File: main.h
 *  Author: Shaun Wetzstein <shaun@us.ibm.com>
 *   Descr: cmdline tool for ECC
 *    Date: 08/30/2012
 */

#ifndef __MAIN_H__
#define __MAIN_H__

#include <stdio.h>

typedef enum {
    c_ERROR = 0,
    c_INJECT = 'I',
    c_REMOVE = 'R',
    c_HEXDUMP = 'H',
} cmd_t;

typedef enum {
    o_ERROR = 0,
    o_OUTPUT = 'o',
} option_t;

typedef enum {
    f_ERROR = 0,
    f_FORCE = 'f',
    f_P8 = 'p',
    f_VERBOSE = 'v',
    f_HELP = 'h',
} flag_t;

typedef struct {
    const char * short_name;

    /* target */
    const char * path;

    /* command */
    cmd_t cmd;

    /* options */
    const char * file;

    /* flags */
    flag_t force, p8, verbose;

    const char ** opt;
    int opt_sz, opt_nr;
} args_t;

extern args_t args;

#define ECC_MAJOR	0x02
#define ECC_MINOR	0x00
#define ECC_PATCH	0x00

#define ECC_EXT		".ecc"

#endif /* __MAIN_H__ */
