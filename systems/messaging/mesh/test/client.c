/*
 * Minimal HTTP client for proxy testing
 *
 * Build:
 *     gcc -O2 -Wall -Wextra -o client client.c -lcurl
 *
 * Run:
 *     ./client 127.0.0.1 8080 deviced v1/ping
 *     ./client 127.0.0.1 8080 deviced v1/version
 *
 * URL constructed as:
 *     http://<host>:<port>/<service_name>/<endpoint>
 */

#define _POSIX_C_SOURCE 200809L

#include <curl/curl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct {
    char   *data;
    size_t  len;
} Buffer;

static size_t write_cb(void *ptr, size_t size, size_t nmemb, void *userdata) {
    size_t n = size * nmemb;
    Buffer *b = (Buffer *)userdata;

    char *p = realloc(b->data, b->len + n + 1);
    if (!p)
        return 0;

    b->data = p;
    memcpy(b->data + b->len, ptr, n);
    b->len += n;
    b->data[b->len] = '\0';
    return n;
}

int main(int argc, char **argv) {
    if (argc != 5) {
        fprintf(stderr,
                "Usage: %s <host> <port> <service_name> <endpoint>\n",
                argv[0]);
        return 2;
    }

    const char *host = argv[1];
    int port = atoi(argv[2]);
    const char *service = argv[3];
    const char *endpoint = argv[4];

    if (port <= 0 || port > 65535) {
        fprintf(stderr, "Invalid port\n");
        return 2;
    }

    char url[2048];
    snprintf(url, sizeof(url),
             "http://%s:%d/%s/%s",
             host, port, service, endpoint);

    CURL *curl;
    CURLcode rc;
    long status = 0;

    Buffer body = {0};

    if (curl_global_init(CURL_GLOBAL_DEFAULT) != 0) {
        fprintf(stderr, "curl_global_init failed\n");
        return 1;
    }

    curl = curl_easy_init();
    if (!curl) {
        fprintf(stderr, "curl_easy_init failed\n");
        curl_global_cleanup();
        return 1;
    }

    struct curl_slist *hdrs = NULL;
    hdrs = curl_slist_append(hdrs, "User-Agent: mesh-client/1.0");
    hdrs = curl_slist_append(hdrs, "Accept: */*");

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, hdrs);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &body);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 10L);

    rc = curl_easy_perform(curl);
    if (rc != CURLE_OK) {
        fprintf(stderr, "REQUEST FAILED: %s\n",
                curl_easy_strerror(rc));
        goto out;
    }

    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &status);

    printf("URL:    %s\n", url);
    printf("STATUS: %ld\n", status);

    if (body.data && body.len) {
        printf("BODY:\n");
        fwrite(body.data, 1, body.len, stdout);
        if (body.data[body.len - 1] != '\n')
            putchar('\n');
    }

out:
    free(body.data);
    curl_slist_free_all(hdrs);
    curl_easy_cleanup(curl);
    curl_global_cleanup();

    return (rc == CURLE_OK) ? 0 : 1;
}
