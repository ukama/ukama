/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/list.c $                                             */
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

#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <errno.h>

#include "assert.h"
#include "misc.h"
#include "list.h"

/* ======================================================================= */

void list_dump(list_t * self, FILE * out)
{
	assert(self != NULL);

	fprintf(out,
		"===================================================================\n");
	fprintf(out, "head: %8lx node: %8lx\n", (ulong) self,
		(ulong) & self->node);

	list_node_t *node = &self->node;
	do {
		fprintf(out, "    node: %8lx - prev: %8lx - next: %8lx\n",
			(ulong) node, (ulong) node->prev,
			(ulong) node->next);
		node = node->next;
	} while (node != &self->node);

	fprintf(out,
		"===================================================================\n");
}
