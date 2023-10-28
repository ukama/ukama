#ifndef FORWARD_H
#define FORWARD_H

#include <pthread.h>

#include "mesh.h"
#include "work.h"

typedef struct forward_ {

    char *uuid;

    int  httpCode;
    int  size;
    void *data;

    pthread_mutex_t mutex;
    pthread_cond_t  hasData;

    struct forward_ *next;
} Forward;

typedef struct {

    Forward *first;
    Forward *last;

    pthread_mutex_t mutex;
} ForwardList;

void init_forward_list(ForwardList **list);
void free_forward_item(Forward *item);

void remove_item_from_list(ForwardList *list, char *uuid);
Forward *is_existing_item_in_list(ForwardList *list, char *uuid);
Forward *add_client_to_list(ForwardList **list, char *uuid);
#endif /* FORWARD_H */
