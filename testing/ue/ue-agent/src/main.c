/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <arpa/inet.h>
#include <curl/curl.h>
#include <errno.h>
#include <fcntl.h>
#include <linux/if.h>
#include <linux/if_tun.h>
#include <netinet/ip.h>
#include <pthread.h>
#include <signal.h>
#include <stdarg.h>
#include <stdatomic.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/ioctl.h>
#include <sys/select.h>
#include <sys/socket.h>
#include <sys/wait.h>
#include <time.h>
#include <unistd.h>

#include "session.h"

#define UE_AGENT_MAX_STR        256
#define UE_AGENT_PKT_MAX        4096
#define UE_AGENT_DEF_APN        "internet"
#define UE_AGENT_DEF_TUN        "tun0"
#define UE_AGENT_DEF_PREFIX     "22"
#define UE_AGENT_ATTACH_RETRY   3
#define UE_AGENT_ATTACH_CHECK   2

static volatile sig_atomic_t gRun = 1;
static volatile sig_atomic_t gDetach = 0;

typedef struct {
    char imsi[UE_AGENT_MAX_STR];
    char iccid[UE_AGENT_MAX_STR];
    char ip[UE_AGENT_MAX_STR];
    char cidr[UE_AGENT_MAX_STR + 8];
    char apn[UE_AGENT_MAX_STR];
    char epcemuUrl[UE_AGENT_MAX_STR];
    char epcemuDataHost[UE_AGENT_MAX_STR];
    int epcemuDataPort;
    char localDataHost[UE_AGENT_MAX_STR];
    int localDataPort;
    char mediaIp[UE_AGENT_MAX_STR];
    char tunIf[UE_AGENT_MAX_STR];
    int detachOnExit;
} Config;

typedef struct {
    UeSession *session;
    atomic_bool running;
    atomic_bool attached;
} AttachMonitor;

static void on_signal(int sig) {

    if (sig == SIGUSR1) gDetach = 1;
    else gRun = 0;
}

static const char *env_or(const char *name, const char *defv) {

    const char *v;

    v = getenv(name);
    if (v != NULL && v[0] != '\0') return v;

    return defv;
}

static int run_cmd(const char *cmd, ...) {

    va_list args;
    const char *arg;
    char *argv[32];
    int argc;
    int status;
    pid_t pid;

    memset(argv, 0, sizeof(argv));
    argc = 0;
    argv[argc++] = (char *)cmd;

    va_start(args, cmd);
    while (argc < 31) {
        arg = va_arg(args, const char *);
        if (arg == NULL) break;
        argv[argc++] = (char *)arg;
    }
    va_end(args);

    pid = fork();
    if (pid < 0) return -1;

    if (pid == 0) {
        execvp(cmd, argv);
        _exit(127);
    }

    while (waitpid(pid, &status, 0) < 0) {
        if (errno == EINTR) continue;
        return -1;
    }

    if (WIFEXITED(status)) return WEXITSTATUS(status);
    if (WIFSIGNALED(status)) return 128 + WTERMSIG(status);

    return -1;
}

static int tun_create(const char *name) {

    struct ifreq ifr;
    int fd;

    fd = open("/dev/net/tun", O_RDWR);
    if (fd < 0) {
        fprintf(stderr, "open /dev/net/tun failed: %s\n", strerror(errno));
        return -1;
    }

    memset(&ifr, 0, sizeof(ifr));
    ifr.ifr_flags = IFF_TUN | IFF_NO_PI;
    snprintf(ifr.ifr_name, IFNAMSIZ, "%s", name);

    if (ioctl(fd, TUNSETIFF, (void *)&ifr) < 0) {
        fprintf(stderr, "TUNSETIFF %s failed: %s\n", name, strerror(errno));
        close(fd);
        return -1;
    }

    return fd;
}

static int configure_tun(Config *cfg) {

    int fd;

    run_cmd("ip", "link", "delete", cfg->tunIf, NULL);

    fd = tun_create(cfg->tunIf);
    if (fd < 0) return -1;

    if (run_cmd("ip", "addr", "replace", cfg->cidr, "dev", cfg->tunIf,
                NULL) != 0) {
        fprintf(stderr, "failed to assign %s to %s\n", cfg->cidr,
                cfg->tunIf);
        close(fd);
        return -1;
    }

    if (run_cmd("ip", "link", "set", cfg->tunIf, "up", NULL) != 0) {
        fprintf(stderr, "failed to bring up %s\n", cfg->tunIf);
        close(fd);
        return -1;
    }

    if (cfg->mediaIp[0] != '\0') {
        run_cmd("ip", "route", "replace", cfg->mediaIp, "dev", cfg->tunIf,
                NULL);
    }

    return fd;
}

static void *monitor_attachment(void *data) {

    AttachMonitor *monitor = data;

    while (atomic_load(&monitor->running)) {
        sleep(UE_AGENT_ATTACH_CHECK);
        if (!atomic_load(&monitor->running)) break;

        if (!session_poll(monitor->session)) {
            atomic_store(&monitor->attached, false);
            break;
        }
    }

    return NULL;
}

static int udp_setup(Config *cfg, struct sockaddr_in *remote) {

    struct sockaddr_in local;
    int fd;
    int opt;

    fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (fd < 0) return -1;

    opt = 1;
    setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    memset(&local, 0, sizeof(local));
    local.sin_family = AF_INET;
    local.sin_port = htons(cfg->localDataPort);
    local.sin_addr.s_addr = htonl(INADDR_ANY);

    if (bind(fd, (struct sockaddr *)&local, sizeof(local)) < 0) {
        fprintf(stderr, "bind local data port %d failed: %s\n",
                cfg->localDataPort, strerror(errno));
        close(fd);
        return -1;
    }

    memset(remote, 0, sizeof(*remote));
    remote->sin_family = AF_INET;
    remote->sin_port = htons(cfg->epcemuDataPort);
    if (inet_pton(AF_INET, cfg->epcemuDataHost, &remote->sin_addr) != 1) {
        close(fd);
        return -1;
    }

    return fd;
}

static void packet_loop(UeSession *session,
                        int tunFd,
                        int udpFd,
                        struct sockaddr_in *remote) {

    unsigned char buf[UE_AGENT_PKT_MAX];
    struct timeval timeout;
    AttachMonitor monitor;
    pthread_t monitorThread;
    fd_set rfds;
    int maxFd;
    int rc;
    ssize_t n;

    monitor.session = session;
    atomic_init(&monitor.running, true);
    atomic_init(&monitor.attached, true);

    rc = pthread_create(&monitorThread, NULL, monitor_attachment, &monitor);
    if (rc != 0) {
        fprintf(stderr, "failed to start attach monitor: %s\n", strerror(rc));
        return;
    }

    while (gRun && !gDetach && atomic_load(&monitor.attached)) {
        FD_ZERO(&rfds);
        FD_SET(tunFd, &rfds);
        FD_SET(udpFd, &rfds);
        maxFd = (tunFd > udpFd) ? tunFd : udpFd;

        timeout.tv_sec = 1;
        timeout.tv_usec = 0;

        rc = select(maxFd + 1, &rfds, NULL, NULL, &timeout);
        if (rc < 0) {
            if (errno == EINTR) continue;
            break;
        }

        if (rc > 0 && FD_ISSET(tunFd, &rfds)) {
            n = read(tunFd, buf, sizeof(buf));
            if (n > 0) {
                sendto(udpFd, buf, n, 0,
                       (struct sockaddr *)remote, sizeof(*remote));
            }
        }

        if (rc > 0 && FD_ISSET(udpFd, &rfds)) {
            n = recvfrom(udpFd, buf, sizeof(buf), 0, NULL, NULL);
            if (n > 0) write(tunFd, buf, n);
        }
    }

    atomic_store(&monitor.running, false);
    pthread_join(monitorThread, NULL);
}

static int config_load(Config *cfg) {

    char *slash;
    const char *apn;

    memset(cfg, 0, sizeof(*cfg));

    snprintf(cfg->imsi, sizeof(cfg->imsi), "%s", env_or("UE_IMSI", ""));
    snprintf(cfg->iccid, sizeof(cfg->iccid), "%s", env_or("UE_ICCID", ""));
    snprintf(cfg->ip, sizeof(cfg->ip), "%s", env_or("UE_IP", ""));
    apn = getenv("UE_APN");
    if (apn == NULL) apn = UE_AGENT_DEF_APN;
    if (strlen(apn) >= sizeof(cfg->apn)) {
        fprintf(stderr, "UE_APN is too long\n");
        return -1;
    }
    snprintf(cfg->apn, sizeof(cfg->apn), "%s", apn);

    snprintf(cfg->epcemuUrl, sizeof(cfg->epcemuUrl), "%s",
             env_or("EPCEMU_URL", "http://127.0.0.1:18092"));
    snprintf(cfg->epcemuDataHost, sizeof(cfg->epcemuDataHost), "%s",
             env_or("EPCEMU_DATA_HOST", "127.0.0.1"));
    snprintf(cfg->localDataHost, sizeof(cfg->localDataHost), "%s",
             env_or("UE_DATA_HOST", "127.0.0.1"));
    snprintf(cfg->mediaIp, sizeof(cfg->mediaIp), "%s",
             env_or("MEDIA_IP", ""));
    snprintf(cfg->tunIf, sizeof(cfg->tunIf), "%s",
             env_or("UE_TUN", UE_AGENT_DEF_TUN));

    cfg->epcemuDataPort = atoi(env_or("EPCEMU_DATA_PORT", "18110"));
    cfg->localDataPort = atoi(env_or("UE_DATA_PORT", "41001"));
    cfg->detachOnExit = atoi(env_or("UE_DETACH_ON_EXIT", "1"));

    slash = strchr(cfg->ip, '/');
    if (slash != NULL) {
        snprintf(cfg->cidr, sizeof(cfg->cidr), "%s", cfg->ip);
        *slash = '\0';
    } else {
        snprintf(cfg->cidr, sizeof(cfg->cidr), "%s/%s", cfg->ip,
                 UE_AGENT_DEF_PREFIX);
    }

    return 0;
}

int main(int argc, char **argv) {

    Config cfg;
    struct sockaddr_in remote;
    UeSession session;
    int attempts;
    int result;
    int tunFd;
    int udpFd;

    (void)argc;
    (void)argv;

    signal(SIGINT, on_signal);
    signal(SIGTERM, on_signal);
    signal(SIGUSR1, on_signal);
    curl_global_init(CURL_GLOBAL_DEFAULT);

    if (config_load(&cfg) != 0) {
        curl_global_cleanup();
        return 1;
    }

    if (cfg.imsi[0] == '\0' || cfg.ip[0] == '\0') {
        fprintf(stderr, "UE_IMSI and UE_IP are required\n");
        curl_global_cleanup();
        return 1;
    }

    memset(&session, 0, sizeof(session));
    session.url = cfg.epcemuUrl;
    session.imsi = cfg.imsi;
    session.iccid = cfg.iccid;
    session.ip = cfg.ip;
    session.requestedApn = cfg.apn;
    session.dataHost = cfg.localDataHost;
    session.dataPort = cfg.localDataPort;
    if (session_init(&session) != 0) {
        fprintf(stderr, "failed to initialize UE status\n");
        curl_global_cleanup();
        return 1;
    }

    tunFd = configure_tun(&cfg);
    if (tunFd < 0) {
        curl_global_cleanup();
        return 1;
    }

    udpFd = udp_setup(&cfg, &remote);
    if (udpFd < 0) {
        close(tunFd);
        curl_global_cleanup();
        return 1;
    }

    result = 0;
    while (gRun) {
        if (gDetach) {
            (void)session_detach(&session);
        } else if (session_attach(&session)) {
            fprintf(stdout, "UE attached imsi=%s apn=%s pdn=%s "
                    "data_state=%s\n", cfg.imsi, session.apn,
                    session.pdnId, session.dataState);
            fflush(stdout);
            packet_loop(&session, tunFd, udpFd, &remote);
        }
        if (gRun) sleep(UE_AGENT_ATTACH_RETRY);
    }

    if (cfg.detachOnExit || gDetach) {
        result = 1;
        for (attempts = 0; attempts < 3; attempts++) {
            if (session_detach(&session)) {
                result = 0;
                break;
            }
        }
        if (result) fprintf(stderr, "UE cleanup remains pending\n");
    }

    close(udpFd);
    close(tunFd);
    curl_global_cleanup();

    return result;
}
