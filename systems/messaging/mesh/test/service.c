/*
 * Minimal HTTP service for proxy testing (standalone)
 *
 * Build:
 *     gcc -O2 -Wall -Wextra -o service service.c
 *
 * Run:
 *     ./service deviced
 *
 * Port lookup:
 *   - Tries /etc/services for tcp port of <service_name>
 *   - Falls back to <base_port> + (stable hash of service_name % range)
 *
 * Endpoints:
 *   GET /v1/ping     -> 200 OK
 *   GET /v1/version  -> 200 vX.Y.Z
 *
 * Methods:
 *   - Valid EP, wrong method -> 405 Method Not Allowed
 *   - Unknown EP             -> 404 Not Found
 */

#define _POSIX_C_SOURCE 200809L

#include <arpa/inet.h>
#include <ctype.h>
#include <errno.h>
#include <netdb.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <time.h>
#include <unistd.h>

static int send_all(int fd, const char *buf, size_t len) {
    size_t off = 0;

    while (off < len) {
        ssize_t n = send(fd, buf + off, len - off, 0);
        if (n < 0) {
            if (errno == EINTR)
                continue;
            return -1;
        }
        off += (size_t)n;
    }
    return 0;
}

static void http_reply(int fd,
                       int code,
                       const char *status,
                       const char *ctype,
                       const char *body) {
    char hdr[512];
    size_t body_len = body ? strlen(body) : 0;

    int n = snprintf(hdr, sizeof(hdr),
                     "HTTP/1.1 %d %s\r\n"
                     "Connection: close\r\n"
                     "Content-Type: %s\r\n"
                     "Content-Length: %zu\r\n"
                     "\r\n",
                     code, status, ctype, body_len);

    if (n <= 0)
        return;

    (void)send_all(fd, hdr, (size_t)n);

    if (body && body_len)
        (void)send_all(fd, body, body_len);
}

static unsigned long djb2_hash(const char *s) {
    unsigned long h = 5381;
    unsigned char c;

    while ((c = (unsigned char)*s++) != 0) {
        h = ((h << 5) + h) + c; /* h * 33 + c */
    }
    return h;
}

static int lookup_port_etc_services_tcp(const char *service_name) {
    struct servent *se;

    /* getservbyname reads /etc/services (and NSS) */
    se = getservbyname(service_name, "tcp");
    if (se == NULL)
        return -1;

    return (int)ntohs((uint16_t)se->s_port);
}

static int service_port(const char *service_name) {
    int port = lookup_port_etc_services_tcp(service_name);
    if (port > 0 && port <= 65535)
        return port;

    /* Fallback: stable derived port range */
    const int base_port = 20000;
    const int range = 20000; /* [20000..39999] */
    unsigned long h = djb2_hash(service_name);

    port = base_port + (int)(h % (unsigned long)range);
    return port;
}

static int is_http_method_token(const char *m) {
    /* Keep it simple but strict */
    return (strcmp(m, "GET") == 0 ||
            strcmp(m, "POST") == 0 ||
            strcmp(m, "PUT") == 0 ||
            strcmp(m, "DELETE") == 0 ||
            strcmp(m, "HEAD") == 0 ||
            strcmp(m, "OPTIONS") == 0 ||
            strcmp(m, "PATCH") == 0);
}

static void handle_request(int cfd) {
    char req[4096];
    ssize_t r = recv(cfd, req, sizeof(req) - 1, 0);
    if (r <= 0)
        return;

    req[r] = '\0';

    char method[16] = {0};
    char target[256] = {0};

    if (sscanf(req, "%15s %255s", method, target) != 2) {
        http_reply(cfd, 400, "Bad Request", "text/plain", "bad request\n");
        return;
    }

    if (!is_http_method_token(method)) {
        http_reply(cfd, 400, "Bad Request", "text/plain", "bad request\n");
        return;
    }

    /* Strip query string if present */
    char path[256] = {0};
    const char *q = strchr(target, '?');
    if (q) {
        size_t n = (size_t)(q - target);
        if (n >= sizeof(path))
            n = sizeof(path) - 1;
        memcpy(path, target, n);
        path[n] = '\0';
    } else {
        strncpy(path, target, sizeof(path) - 1);
        path[sizeof(path) - 1] = '\0';
    }

    int is_get = (strcmp(method, "GET") == 0);

    if (strcmp(path, "/v1/ping") == 0) {
        if (!is_get) {
            http_reply(cfd, 405, "Method Not Allowed", "text/plain",
                       "method not allowed\n");
            return;
        }
        http_reply(cfd, 200, "OK", "text/plain", "OK\n");
        return;
    }

    if (strcmp(path, "/v1/version") == 0) {
        if (!is_get) {
            http_reply(cfd, 405, "Method Not Allowed", "text/plain",
                       "method not allowed\n");
            return;
        }

        char version[64];
        snprintf(version, sizeof(version),
                 "v%d.%d.%d\n",
                 rand() % 10,
                 rand() % 20,
                 rand() % 100);

        http_reply(cfd, 200, "OK", "text/plain", version);
        return;
    }

    http_reply(cfd, 404, "Not Found", "text/plain", "not found\n");
}

int main(int argc, char **argv) {
    if (argc != 2) {
        fprintf(stderr, "Usage: %s <service_name>\n", argv[0]);
        return 2;
    }

    const char *service_name = argv[1];
    int port = service_port(service_name);

    srand((unsigned)time(NULL));

    int s = socket(AF_INET, SOCK_STREAM, 0);
    if (s < 0) {
        perror("socket");
        return 1;
    }

    int yes = 1;
    setsockopt(s, SOL_SOCKET, SO_REUSEADDR, &yes, sizeof(yes));

    struct sockaddr_in addr;
    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port = htons((uint16_t)port);
    addr.sin_addr.s_addr = htonl(INADDR_ANY);

    if (bind(s, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        perror("bind");
        close(s);
        return 1;
    }

    if (listen(s, 64) < 0) {
        perror("listen");
        close(s);
        return 1;
    }

    printf("Service '%s' listening on port %d\n", service_name, port);

    for (;;) {
        struct sockaddr_in caddr;
        socklen_t clen = sizeof(caddr);

        int cfd = accept(s, (struct sockaddr *)&caddr, &clen);
        if (cfd < 0) {
            if (errno == EINTR)
                continue;
            perror("accept");
            break;
        }

        handle_request(cfd);
        close(cfd);
    }

    close(s);
    return 0;
}
