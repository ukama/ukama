/*
 * Copyright (c) 2012-2017 Roberto E. Vargas Caballero
 *
 * SPDX-License-Identifier: BSD-3-Clause
 */
/*
 * Portions copyright (c) 2018, ARM Limited and Contributors.
 * All rights reserved.
 */

#ifndef STDLIB_H
#define STDLIB_H

#include <stdlib_.h>

#ifndef NULL
#define NULL ((void *) 0)
#endif

#define _ATEXIT_MAX 1

extern void abort(void);
extern int atexit(void (*func)(void));
extern void exit(int status);

#endif /* STDLIB_H */
