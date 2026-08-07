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
#include <signal.h>
#include <stdarg.h>
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

#define UE_AGENT_MAX_STR        256
#define UE_AGENT_PKT_MAX        4096
#define UE_AGENT_DEF_APN        "internet"
#define UE_AGENT_DEF_TUN        "tun0"
#define UE_AGENT_DEF_PREFIX     "22"
#define UE_AGENT_ATTACH_RETRY   3
#define UE_AGENT_ATTACH_CHECK   2

static volatile bool gRun = true;

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

static void on_signal(int sig) {

    (void)sig;
    gRun = false;
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

static size_t discard_body(void *ptr, size_t size, size_t nmemb, void *data) {

    (void)ptr;
    (void)data;

    return size * nmemb;
}

static long http_request(const char *method,
                         const char *url,
                         const char *json) {

    CURL *curl;
    CURLcode res;
    struct curl_slist *headers;
    long code;

    curl = curl_easy_init();
    if (curl == NULL) return 0;

    headers = NULL;
    code = 0;

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 10L);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, discard_body);

    if (json != NULL) {
        headers = curl_slist_append(headers, "Content-Type: application/json");
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json);
    }

    res = curl_easy_perform(curl);
    if (res == CURLE_OK) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    }

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);

    return code;
}

static long attach_ue(Config *cfg) {

    char url[UE_AGENT_MAX_STR * 3];
    char body[UE_AGENT_MAX_STR * 6];

    snprintf(url, sizeof(url), "%s/v1/ue/attach", cfg->epcemuUrl);
    snprintf(body, sizeof(body),
             "{\"imsi\":\"%s\",\"iccid\":\"%s\","
             "\"ip\":\"%s\",\"apn\":\"%s\","
             "\"userPlaneHost\":\"%s\",\"userPlanePort\":%d}",
             cfg->imsi, cfg->iccid, cfg->ip, cfg->apn,
             cfg->localDataHost, cfg->localDataPort);

    return http_request("POST", url, body);
}

static long detach_ue(Config *cfg) {

    char url[UE_AGENT_MAX_STR * 3];
    char body[UE_AGENT_MAX_STR * 2];

    snprintf(url, sizeof(url), "%s/v1/ue/detach", cfg->epcemuUrl);
    snprintf(body, sizeof(body), "{\"imsi\":\"%s\"}", cfg->imsi);

    return http_request("DELETE", url, body);
}

static bool ue_is_attached(Config *cfg) {

    char url[UE_AGENT_MAX_STR * 3];
    long code;

    snprintf(url, sizeof(url), "%s/v1/ue/%s", cfg->epcemuUrl, cfg->imsi);
    code = http_request("GET", url, NULL);

    return code >= 200 && code < 300;
}

static bool wait_for_attach(Config *cfg) {

    long code;

    while (gRun) {
        code = attach_ue(cfg);
        if (code >= 200 && code < 300) return true;

        fprintf(stderr, "attach pending imsi=%s ip=%s http=%ld\n",
                cfg->imsi, cfg->ip, code);
        sleep(UE_AGENT_ATTACH_RETRY);
    }

    return false;
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

static void packet_loop(Config *cfg,
                        int tunFd,
                        int udpFd,
                        struct sockaddr_in *remote) {

    unsigned char buf[UE_AGENT_PKT_MAX];
    struct timeval timeout;
    time_t nextCheck;
    fd_set rfds;
    int maxFd;
    int rc;
    ssize_t n;

    nextCheck = time(NULL) + UE_AGENT_ATTACH_CHECK;

    while (gRun) {
        FD_ZERO(&rfds);
        FD_SET(tunFd, &rfds);
        FD_SET(udpFd, &rfds);
        maxFd = (tunFd > udpFd) ? tunFd : udpFd;

        timeout.tv_sec = 1;
        timeout.tv_usec = 0;

        rc = select(maxFd + 1, &rfds, NULL, NULL, &timeout);
        if (rc < 0) {
            if (errno == EINTR) continue;
            return;
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

        if (time(NULL) >= nextCheck) {
            if (!ue_is_attached(cfg)) return;
            nextCheck = time(NULL) + UE_AGENT_ATTACH_CHECK;
        }
    }
}

static void config_load(Config *cfg) {

    char *slash;

    memset(cfg, 0, sizeof(*cfg));

    snprintf(cfg->imsi, sizeof(cfg->imsi), "%s", env_or("UE_IMSI", ""));
    snprintf(cfg->iccid, sizeof(cfg->iccid), "%s", env_or("UE_ICCID", ""));
    snprintf(cfg->ip, sizeof(cfg->ip), "%s", env_or("UE_IP", ""));
    snprintf(cfg->apn, sizeof(cfg->apn), "%s",
             env_or("UE_APN", UE_AGENT_DEF_APN));

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
}

int main(int argc, char **argv) {

    Config cfg;
    struct sockaddr_in remote;
    bool attached;
    int tunFd;
    int udpFd;

    (void)argc;
    (void)argv;

    signal(SIGINT, on_signal);
    signal(SIGTERM, on_signal);
    curl_global_init(CURL_GLOBAL_DEFAULT);

    config_load(&cfg);

    if (cfg.imsi[0] == '\0' || cfg.ip[0] == '\0') {
        fprintf(stderr, "UE_IMSI and UE_IP are required\n");
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

    attached = false;

    while (gRun) {
        if (!wait_for_attach(&cfg)) break;

        attached = true;
        fprintf(stdout, "UE attached imsi=%s iccid=%s ip=%s data=%s:%d\n",
                cfg.imsi,
                cfg.iccid,
                cfg.ip,
                cfg.epcemuDataHost,
                cfg.epcemuDataPort);
        fflush(stdout);

        packet_loop(&cfg, tunFd, udpFd, &remote);
        if (!gRun) break;

        attached = false;
        fprintf(stdout, "UE detached imsi=%s ip=%s; retrying attach\n",
                cfg.imsi, cfg.ip);
        fflush(stdout);
    }

    if (attached && cfg.detachOnExit) {
        (void)detach_ue(&cfg);
    }

    close(udpFd);
    close(tunFd);
    curl_global_cleanup();

    return 0;
}
