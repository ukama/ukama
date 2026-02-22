#pragma once

#include <ulfius.h>
#include <pthread.h>
#include <stddef.h>

typedef struct {
    struct _u_instance inst;
    int port;
    pthread_t thread;
    int started;

    pthread_mutex_t lock;
    unsigned notify_posts;
    unsigned fem_puts;
    unsigned noded_gets;

    char last_path[256];
    char last_body[2048];
} MockService;

int mock_noded_start(MockService *ms, int port);
int mock_notifyd_start(MockService *ms, int port);
int mock_femd_start(MockService *ms, int port);
void mock_service_stop(MockService *ms);
void mock_service_reset(MockService *ms);
unsigned mock_notify_posts(MockService *ms);
unsigned mock_fem_puts(MockService *ms);
unsigned mock_noded_gets(MockService *ms);
void mock_last_body(MockService *ms, char *buf, size_t buflen);
