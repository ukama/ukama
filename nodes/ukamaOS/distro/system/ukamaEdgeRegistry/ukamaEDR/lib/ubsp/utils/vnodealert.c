/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "utils/vnodealert.h"

#include "headers/utils/log.h"

int file_watch_fd(const char *filename) {
    static int inot = ERR;
    static int iflags = IN_CLOEXEC | IN_NONBLOCK;
    static uint32_t mask = IN_MODIFY;
    int watch;

    inot = inotify_init1(iflags);
    if (inot == ERR) {
        perror("inotify_init1");
        return ERR;
    }

    watch = inotify_add_watch(inot, filename, mask);
    if (watch == ERR) {
        perror("inotify_add_watch");
        return ERR;
    }
    return inot;
}

int ready_inot(int inot) {
    int ret = 0;
    int i = 0;
    char buffer[EVENT_BUF_LEN]
        __attribute__((aligned(__alignof__(struct inotify_event))));
    int nr;
    char *p;
    struct inotify_event *ev_p;

    nr = read(inot, (char *)buffer, EVENT_BUF_LEN);
    while (i < nr) {
        ev_p = (struct inotify_event *)&buffer[i];
        if (ev_p->mask & IN_OPEN)
            printf("IN_OPEN: ");
        if (ev_p->mask & IN_CLOSE_NOWRITE)
            printf("IN_CLOSE_NOWRITE: ");
        if (ev_p->mask & IN_CLOSE_WRITE)
            printf("IN_CLOSE_WRITE: ");
        if (ev_p->mask & IN_MODIFY) {
            printf("IN_MODIFY: ");
            ret = 1;
        }
        i += EVENT_SIZE + ev_p->len;
    }
    return ret;
}

int poll_file(IRQCfg *cfg) {
    int ret = 0;
    int nfds, num_open_fds;
    struct pollfd pfds[1];
    num_open_fds = nfds = 1;
    int fw_fd = file_watch_fd(cfg->fname);
    if (fw_fd < 0) {
        return fw_fd;
    }
    /* TODO: Better way to close the pfds*/
    memset(&pfds, 0, sizeof(struct pollfd));
    //pfds = calloc(nfds, sizeof(struct pollfd));
    //if (pfds == NULL) {
    //    errExit("malloc");
    //}
    pfds[0].fd = fw_fd;
    pfds[0].events = POLLIN;
    /* Keep calling poll() as long as at least one file descriptor is open */

    while (num_open_fds > 0) {
        int ready = 0;
        log_trace("VNODEALERT:: Started poll() for %s.\n", cfg->fname);
        ready = poll(pfds, nfds, -1);
        if (ready == -1) {
            log_error("VNODEALERT:: poll() error.");
            ret = -1;
        }
        log_trace("VNODEALERT:: poll() received a event: %d\n", ready);
        /* Deal with array returned by poll() */
        for (int j = 0; j < nfds; j++) {
            if (pfds[j].revents != 0) {
                log_debug("VNODEALERT:: poll() fd %d got events: 0x%x",
                          pfds[j].fd, pfds[j].revents);
                if (pfds[j].revents & POLLIN) {
                    if (ready_inot(fw_fd)) {
                        /* Callback to the registered cb */
                        cfg->cb(cfg);
                    }
                } else { /* POLLERR | POLLHUP */
                    log_debug("VNODEALERT:: poll() closing fd %d\n",
                              pfds[j].fd);
                    if (close(pfds[j].fd) == -1)
                        ret = -1;
                    num_open_fds--;
                }
            }
        }
    }

    memset(&pfds, 0, sizeof(struct pollfd));
    log_debug("VNODEALERT:: All file descriptors closed; bye\n");
    return ret;
}
