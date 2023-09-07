/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/test/exception.c $                                       */
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

#define EXCEPTION_FOO   1
#define EXCEPTION_BAR   2
#define EXCEPTION_TAZ   3

int bar(int i) {
    int * data = MALLOC(sizeof *data);
    *data = i;

    throw_bytes(EXCEPTION_BAR, data, sizeof *data);

    return 0;
}

int foo(int i) {

    exception_t ex;
    try {
	printf("%s: %d\n", __func__, __LINE__);
        bar(i);
    } catch (EXCEPTION_BAR, ex) {
        printf("foo: CAUGHT %s(%d) EXCEPTION_BAR data[%d]\n",
               ex.file, ex.line, *(int *)ex.data);
        throw_bytes(EXCEPTION_FOO, ex.data, ex.size);
    } end_try

    throw_bytes(EXCEPTION_FOO, "this is a test", 14);

	printf("%s: %d\n", __func__, __LINE__);

    return 0;
}

int main(void) {
    exception_t ex;

    try {
	printf("%s: %d\n", __func__, __LINE__);
        foo(1);
        printf("try block: AFTER foo <-- should not get here\n");
    } catch (EXCEPTION_FOO, ex) {
        printf("main: CAUGHT %s(%d) EXCEPTION_FOO data[%d]\n",
               ex.file, ex.line,        *(int *)ex.data);
    } else (ex) {
        printf("main: CAUGHT %s(%d) data[%d]\n",
               ex.file, ex.line, *(int *)ex.data);
    } end_try

    return 0;
}
