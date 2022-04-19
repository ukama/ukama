/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_list.h"

#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"


void usys_list_new(ListInfo *list, uint16_t elementSize, FreeFxn freeFn,
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

void usys_list_destroy(ListInfo *list) {
    if (list) {
        ListNode *current = NULL;
        while (list->head != NULL) {
            current = list->head;
            list->head = current->next;
            if (list->free) {
                list->free(current);
            } else {
                if (current->data) {
                    usys_free(current->data);
                }
                usys_free(current);
            }
            list->logicalLength--;
        }
    }
}

void usys_list_prepend(ListInfo *list, void *element) {
    ListNode *node = usys_zmalloc(sizeof(ListNode));
    if (list->datacpy) {
        node->data = list->datacpy(element);
    } else {
        node->data = usys_zmalloc(list->elementSize);
        if (node->data) {
            usys_memcpy(node->data, element, list->elementSize);
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

void usys_list_append(ListInfo *list, void *element) {
    ListNode *node = usys_zmalloc(sizeof(ListNode));
    if (list->datacpy) {
        node->data = list->datacpy(element);
    } else {
        node->data = usys_zmalloc(list->elementSize);
        if (node->data) {
            usys_memcpy(node->data, element, list->elementSize);
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

void *usys_list_search(ListInfo *list, void *element) {
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
                            data = usys_zmalloc(list->elementSize);
                            if (data) {
                                usys_memcpy(data, node->data, list->elementSize);
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

int usys_list_if_element_found(ListInfo *list, void *element) {
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

int usys_list_remove(ListInfo *list, void *element) {
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

/* Only the non-primary member
 * (which are not the part of primary key or used for search) can be updated
*/
int usys_list_update(ListInfo *list, void *element) {
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
                        ListNode *newnode = usys_zmalloc(sizeof(ListNode));
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
                        usys_memcpy(node->data, element, list->elementSize);
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

void usys_list_for_each(ListInfo *list, listIterator iterator) {
    if (!iterator) {
        return;
    }
    if (!list) {
        return;
    }
    ListNode *node = list->head;
    int result = true;
    while (node != NULL && result) {
        result = iterator(node->data);
        node = node->next;
    }
}

void usys_list_next(ListInfo *list, ListNode **node) {
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

void *usys_list_head(ListInfo *list, uint8_t removeFromList) {
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
        element = usys_zmalloc(list->elementSize);
        if (element) {
            usys_memcpy(element, node->data, list->elementSize);
        }
    }

    if (removeFromList) {
        list->head = node->next;
        list->logicalLength--;

        usys_free(node->data);
        usys_free(node);
    }
    return element;
}

void *usys_list_tail(ListInfo *list) {
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
        element = usys_zmalloc(list->elementSize);
        if (element) {
            usys_memcpy(element, node->data, list->elementSize);
        }
    }
    return element;
}

int usys_list_size(ListInfo *list) {
    if (list) {
        return list->logicalLength;
    } else {
        return -1;
    }
}

/* This function doesn't support deep copy. */
void usys_list_copy(ListInfo *list, void *data) {
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
            usys_memcpy(data + (count * el_sz), node->data, el_sz);
            node = node->next;
            count++;
        }
    }
}
