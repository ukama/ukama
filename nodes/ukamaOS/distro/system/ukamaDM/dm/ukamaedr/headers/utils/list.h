/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef UTILS_LIST_H_
#define UTILS_LIST_H_

#include <stdbool.h>
#include <stdint.h>

// a common function used to free malloc'd objects
typedef void (*FreeFxn)(void *);
typedef int (*CompareFxn)(void*, void*);
typedef void* (*CopyDataFxn)(void*);

typedef int (*listIterator)(void *);


typedef struct  __attribute__((__packed__)) listNode {
  void *data;
  struct listNode *next;
} ListNode;

typedef struct {
  uint16_t logicalLength;
  uint16_t elementSize;
  ListNode *head;
  ListNode *tail;
  FreeFxn free;
  CompareFxn cmp;
  CopyDataFxn datacpy;
} ListInfo;

int list_size(ListInfo *list);
int list_if_element_found(ListInfo *list, void *element);
int list_remove(ListInfo *list, void *element);
int list_update(ListInfo *list, void *element);
void list_append(ListInfo *list, void *element);
void list_next(ListInfo *list, ListNode **node);
void list_copy(ListInfo *list, void *data);
void list_for_each(ListInfo *list, listIterator iterator);
void list_new(ListInfo *list, uint16_t elementSize, FreeFxn freeFn,
		CompareFxn cmpFn, CopyDataFxn copyFn);
void list_destroy(ListInfo *list);
void list_prepend(ListInfo *list, void *element);
void* list_search(ListInfo *list, void *element);
void* list_tail(ListInfo *list );
void* list_head(ListInfo *list,  uint8_t removeFromList);
#endif /* UTILS_LIST_H_ */
