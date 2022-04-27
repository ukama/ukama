/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/test/exc_throw.c $                                       */
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


#include <stdio.h>

#include <clib/exception.h>

#define EXCEPTION_FOO   10

void foo() {
    throw(EXCEPTION_FOO, ex.data, ex.size);
}

int main() {
    exception_t ex;

    try {
        printf("try block: BEFORE foo\n");
        foo();
        printf("try block: AFTER foo <-- should not get here\n");
    } catch (EXCEPTION_FOO, ex) {
        printf("main: CAUGHT %s(%d) EXCEPTION_FOO data[%d]\n",
               ex.file, ex.line,        *(int *)ex.data);
    } else (ex) {
        printf("main: CAUGHT %s(%d) data[%d]\n",
               ex.file, ex.line, *(int *)ex.data);
    } end

    return 0;
}

