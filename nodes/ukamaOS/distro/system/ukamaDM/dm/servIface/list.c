/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "headers/utils/list.h"

#include <string.h>
#include <stdlib.h>

void list_new(ListInfo *list, uint16_t elementSize, FreeFxn freeFn,
              CompareFxn cmpFn, CopyDataFxn copyFn) {
    if (list) {
        if (!elementSize) {
            return;
        }
        list->logicalLength = 0;
        list->elementSize = elementSize;
        list->head = list->tail = NULL;
        list->free = freeFn;
        list->cmp = cmpFn;
        list->datacpy = copyFn;
    }
}

void list_destroy(ListInfo *list) {
    if (list) {
        ListNode *current = NULL;
        while (list->head != NULL) {
            current = list->head;
            list->head = current->next;
            if (list->free) {
                list->free(current);
            } else {
                if (current->data) {
                    free(current->data);
                }
                free(current);
            }
            list->logicalLength--;
        }
    }
}

void list_prepend(ListInfo *list, void *element) {
    ListNode *node = malloc(sizeof(ListNode));
    if (list->datacpy) {
        node->data = list->datacpy(element);
    } else {
        node->data = malloc(list->elementSize);
        if (node->data) {
            memcpy(node->data, element, list->elementSize);
        }
    }
    node->next = list->head;
    list->head = node;

    // first node?
    if (!list->tail) {
        list->tail = list->head;
    }

    list->logicalLength++;
}

void list_append(ListInfo *list, void *element) {
    ListNode *node = malloc(sizeof(ListNode));
    if (list->datacpy) {
        node->data = list->datacpy(element);
    } else {
        node->data = malloc(list->elementSize);
        if (node->data) {
            memcpy(node->data, element, list->elementSize);
        }
    }
    node->next = NULL;

    if (list->logicalLength == 0) {
        list->head = list->tail = node;
    } else {
        list->tail->next = node;
        list->tail = node;
    }

    list->logicalLength++;
}

void *list_search(ListInfo *list, void *element) {
    void *data = NULL;
    if (list) {
        ListNode *node = NULL;
        uint8_t count = 0;
        if (list->logicalLength == 0) {
            list->head = list->tail = node;
        } else {
            node = list->head;
            while (node) {
                if (list->cmp) {
                    if (list->cmp(node->data, element)) {
                        if (list->datacpy) {
                            data = list->datacpy(node->data);
                        } else {
                            data = malloc(list->elementSize);
                            if (data) {
                                memcpy(data, node->data, list->elementSize);
                            }
                        }
                        break;
                    }
                }
                node = node->next;
                count++;
            }
        }
    }
    return data;
}

int list_if_element_found(ListInfo *list, void *element) {
    int ret = 0;
    if (list) {
        ListNode *node = NULL;
        uint8_t count = 0;
        if (list->logicalLength == 0) {
            list->head = list->tail = node;
        } else {
            node = list->head;
            while (node) {
                if (list->cmp) {
                    if (list->cmp(node->data, element)) {
                        ret = 1;
                        break;
                    }
                }
                node = node->next;
                count++;
            }
        }
    }
    return ret;
}

int list_remove(ListInfo *list, void *element) {
    int ret = -1;
    ListNode *node = NULL;
    ListNode *prev = NULL;
    if (!list) {
        return ret;
    }
    if (list->logicalLength == 0) {
        list->head = list->tail = node;
    } else {
        node = list->head;
        while (node) {
            if (list->cmp) {
                if (list->cmp(node->data, element)) {
                    if (list->free) {
                        if ((node != list->head) && (node != list->tail)) {
                            if (prev) {
                                prev->next = node->next;
                            }
                        }
                        if (node == list->head) {
                            list->head = node->next;
                        }
                        if (node == list->tail) {
                            list->tail = prev;
                            if (prev) {
                                prev->next = NULL;
                            }
                        }
                        list->free(node);
                        list->logicalLength--;
                        ret = 0;
                    }
                    /* if we don't have free fxn so it would cause memory leaks*/
                    break;
                }
            }
            prev = node;
            node = node->next;
        }
    }
    return ret;
}

/* Only the non-primary member(which are not the part of primary key or used for search) can be updated */
int list_update(ListInfo *list, void *element) {
    int ret = -1;
    if (!list) {
        return ret;
    }
    ListNode *node = NULL;
    ListNode *prev = NULL;
    if (list->logicalLength == 0) {
        list->head = list->tail = node;
    } else {
        node = list->head;
        while (node) {
            if (list->cmp) {
                if (list->cmp(node->data, element)) {
                    if (list->datacpy) {
                        /* Create a node */
                        ListNode *newnode = malloc(sizeof(ListNode));
                        if (newnode) {
                            /* Assign data */
                            newnode->next = node->next;
                            newnode->data = list->datacpy(element);
                        }
                        if ((node != list->head) && (node != list->tail)) {
                            if (prev) {
                                prev->next = newnode;
                            }
                        }
                        if (node == list->head) {
                            list->head = newnode;
                        }
                        if (node == list->tail) {
                            list->tail = newnode;
                        }
                        /* Free the node */
                        if (list->free) {
                            list->free(node);
                        }

                    } else {
                        memcpy(node->data, element, list->elementSize);
                    }
                    ret = 0;
                    break;
                }
            }
            prev = node;
            node = node->next;
        }
    }
    return ret;
}

void list_for_each(ListInfo *list, listIterator iterator) {
    if (!iterator) {
        return;
    }
    if (!list) {
        return;
    }
    ListNode *node = list->head;
    uint8_t result = true;
    while (node != NULL && result) {
        result = iterator(node->data);
        node = node->next;
    }
}

void list_next(ListInfo *list, ListNode **node) {
    if (!list) {
        return;
    }

    /* Starting from head */
    if (!(*node)) {
        *node = list->head;
    } else {
        *node = (*node)->next;
    }
}

void *list_head(ListInfo *list, uint8_t removeFromList) {
    void *element = NULL;
    if (!list) {
        return NULL;
    }

    if (!list->head) {
        return element;
    }

    ListNode *node = list->head;
    if (list->datacpy) {
        element = list->datacpy(node->data);
    } else {
        element = malloc(list->elementSize);
        if (element) {
            memcpy(element, node->data, list->elementSize);
        }
    }

    if (removeFromList) {
        list->head = node->next;
        list->logicalLength--;

        free(node->data);
        free(node);
    }
    return element;
}

void *list_tail(ListInfo *list) {
    void *element = NULL;
    if (!list) {
        return NULL;
    }

    if (!list->tail) {
        return element;
    }

    ListNode *node = list->tail;
    if (list->datacpy) {
        element = list->datacpy(node->data);
    } else {
        element = malloc(list->elementSize);
        if (element) {
            memcpy(element, node->data, list->elementSize);
        }
    }
    return element;
}

int list_size(ListInfo *list) {
    if (list) {
        return list->logicalLength;
    } else {
        return -1;
    }
}

/* This function doesn't support deep copy. */
void list_copy(ListInfo *list, void *data) {
    int ret = 0;
    if (!list) {
        return;
    }
    ListNode *node = NULL;
    uint16_t count = 0;
    uint16_t el_sz = list->elementSize;
    if (list->logicalLength == 0) {
        return;
    }
    if (data) {
        node = list->head;
        while (node) {
            memcpy(data + (count * el_sz), node->data, el_sz);
            node = node->next;
            count++;
        }
    }
}
