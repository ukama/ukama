/*
 * Core tests for wimc.d.
 * These tests avoid real hub, real agent, and casync.
 */

#include <errno.h>
#include <limits.h>
#include <sqlite3.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include <jansson.h>
#include <uuid/uuid.h>

#include "db.h"
#include "jserdes.h"
#include "package_cache.h"
#include "tasks.h"
#include "wimc.h"

#define TEST_TMP_ROOT "/tmp/wimc-test"
#define TEST_DB_FILE  TEST_TMP_ROOT "/wimc.db"
#define TEST_ASSERT(cond) test_assert((cond), #cond, __FILE__, __LINE__)
#define TEST_EQ_INT(exp, got) test_eq_int((exp), (got), #got, __FILE__, __LINE__)
#define TEST_EQ_STR(exp, got) test_eq_str((exp), (got), #got, __FILE__, __LINE__)
#define TEST_NOT_NULL(ptr) test_assert((ptr) != NULL, #ptr " != NULL", __FILE__, __LINE__)
#define TEST_NULL(ptr) test_assert((ptr) == NULL, #ptr " == NULL", __FILE__, __LINE__)

static int g_tests;
static int g_failures;

static void test_assert(bool ok, const char *expr, const char *file, int line) {
    g_tests++;
    if (!ok) {
        g_failures++;
        fprintf(stderr, "FAIL %s:%d: %s\n", file, line, expr);
    }
}

static void test_eq_int(long exp, long got, const char *expr,
                        const char *file, int line) {
    g_tests++;
    if (exp != got) {
        g_failures++;
        fprintf(stderr, "FAIL %s:%d: %s expected %ld got %ld\n",
                file, line, expr, exp, got);
    }
}

static void test_eq_str(const char *exp, const char *got, const char *expr,
                        const char *file, int line) {
    g_tests++;
    if ((exp == NULL && got != NULL) || (exp != NULL && got == NULL) ||
        (exp != NULL && got != NULL && strcmp(exp, got) != 0)) {
        g_failures++;
        fprintf(stderr, "FAIL %s:%d: %s expected '%s' got '%s'\n",
                file, line, expr, exp ? exp : "(null)", got ? got : "(null)");
    }
}

static void run_cmd(const char *cmd) {
    int rc;

    rc = system(cmd);
    if (rc != 0) {
        fprintf(stderr, "command failed: %s\n", cmd);
        exit(1);
    }
}

static void mkdir_p(const char *path) {
    char tmp[PATH_MAX];
    char *p;

    snprintf(tmp, sizeof(tmp), "%s", path);
    for (p = tmp + 1; *p != '\0'; p++) {
        if (*p == '/') {
            *p = '\0';
            mkdir(tmp, 0755);
            *p = '/';
        }
    }
    mkdir(tmp, 0755);
}

static void write_file(const char *path, const char *body) {
    FILE *fp;

    fp = fopen(path, "w");
    if (fp == NULL) {
        fprintf(stderr, "fopen failed for %s: %s\n", path, strerror(errno));
        exit(1);
    }
    fwrite(body, 1, strlen(body), fp);
    fclose(fp);
}

static void make_pkg_dir(const char *dir, const char *version) {
    char path[PATH_MAX];

    mkdir_p(dir);
    mkdir_p(TEST_TMP_ROOT "/pkg/bin");

    snprintf(path, sizeof(path), "%s/VERSION", dir);
    write_file(path, version);
    snprintf(path, sizeof(path), "%s/manifest.json", dir);
    write_file(path, "{\"name\":\"example\"}\n");
}

static void make_pkg_tar(const char *dir, const char *tarPath,
                         const char *version) {
    char cmd[PATH_MAX * 3];

    make_pkg_dir(dir, version);
    snprintf(cmd, sizeof(cmd), "tar -czf %s -C %s .", tarPath, dir);
    run_cmd(cmd);
}

static WimcReq *make_fetch_request(const char *name, const char *tag,
                                   const char *method) {
    WimcReq *req;

    req = calloc(1, sizeof(WimcReq));
    TEST_NOT_NULL(req);

    req->type = WREQ_FETCH;
    req->fetch = calloc(1, sizeof(WFetch));
    TEST_NOT_NULL(req->fetch);

    req->fetch->content = calloc(1, sizeof(WContent));
    TEST_NOT_NULL(req->fetch->content);

    uuid_generate(req->fetch->uuid);
    req->fetch->cbURL = strdup("http://127.0.0.1:19079/v1/apps/example/v1-abc/stats");
    req->fetch->interval = 10;

    req->fetch->content->name = strdup(name);
    req->fetch->content->tag = strdup(tag);
    req->fetch->content->method = strdup(method);
    req->fetch->content->indexURL = strdup("http://127.0.0.1:19080/fake.caidx");
    req->fetch->content->storeURL = strdup("http://127.0.0.1:19080/chunks/");
    req->fetch->content->expectedSizeBytes = 1024;

    return req;
}

static void free_fetch_request(WimcReq *req) {
    if (req == NULL) return;
    if (req->fetch != NULL) {
        if (req->fetch->content != NULL) {
            free(req->fetch->content->name);
            free(req->fetch->content->tag);
            free(req->fetch->content->method);
            free(req->fetch->content->indexURL);
            free(req->fetch->content->storeURL);
            free(req->fetch->content);
        }
        free(req->fetch->cbURL);
        free(req->fetch);
    }
    free(req);
}

static void reset_test_dirs(void) {
    run_cmd("rm -rf " TEST_TMP_ROOT);
    mkdir_p(TEST_TMP_ROOT);
    run_cmd("rm -rf /ukama/apps/pkgs /ukama/apps/db");
    mkdir_p("/ukama/apps/pkgs/.tmp");
    mkdir_p("/ukama/apps/db");
}

static void test_pkg_identifiers(void) {
    TEST_ASSERT(pkg_is_valid_identifier("example"));
    TEST_ASSERT(pkg_is_valid_identifier("v1-abc"));
    TEST_ASSERT(pkg_is_valid_identifier("example.d"));
    TEST_ASSERT(pkg_is_valid_identifier("example_1"));

    TEST_ASSERT(!pkg_is_valid_identifier(NULL));
    TEST_ASSERT(!pkg_is_valid_identifier(""));
    TEST_ASSERT(!pkg_is_valid_identifier("../bad"));
    TEST_ASSERT(!pkg_is_valid_identifier("bad/name"));
    TEST_ASSERT(!pkg_is_valid_identifier("bad name"));
    TEST_ASSERT(!pkg_is_valid_identifier("bad;name"));
}

static void test_pkg_path_for_tag(void) {
    char path[WIMC_MAX_PATH_LEN];

    TEST_EQ_INT(0, pkg_path_for_tag("example", "v1-abc", path, sizeof(path)));
    TEST_EQ_STR("/ukama/apps/pkgs/example_v1-abc.tar.gz", path);

    TEST_EQ_INT(-1, pkg_path_for_tag("../example", "v1-abc", path,
                                      sizeof(path)));
}

static void test_pkg_version_from_dir_and_tar(void) {
    char version[WIMC_MAX_NAME_LEN];

    make_pkg_tar(TEST_TMP_ROOT "/pkg-v1-abc",
                 TEST_TMP_ROOT "/example_v1-abc.tar.gz",
                 "v1-abc\n");

    TEST_EQ_INT(0, pkg_read_version_from_dir(TEST_TMP_ROOT "/pkg-v1-abc",
                                             version, sizeof(version)));
    TEST_EQ_STR("v1-abc", version);

    memset(version, 0, sizeof(version));
    TEST_EQ_INT(0, pkg_read_version_from_tar(TEST_TMP_ROOT
                                             "/example_v1-abc.tar.gz",
                                             version, sizeof(version)));
    TEST_EQ_STR("v1-abc", version);
}

static void test_pkg_validate_tar(void) {
    PackageInfo info;

    make_pkg_tar(TEST_TMP_ROOT "/pkg-v1-abc",
                 TEST_TMP_ROOT "/example_v1-abc.tar.gz",
                 "v1-abc\n");

    TEST_ASSERT(pkg_validate_tar("example", "v1-abc",
                                 TEST_TMP_ROOT "/example_v1-abc.tar.gz",
                                 &info));
    TEST_ASSERT(info.exists);
    TEST_ASSERT(info.valid);
    TEST_EQ_STR("v1-abc", info.actualVersion);

    TEST_ASSERT(!pkg_validate_tar("example", "v1-xyz",
                                  TEST_TMP_ROOT "/example_v1-abc.tar.gz",
                                  &info));
    TEST_ASSERT(info.exists);
    TEST_ASSERT(!info.valid);
    TEST_EQ_STR("v1-abc", info.actualVersion);

    TEST_ASSERT(pkg_validate_tar("example", "latest",
                                 TEST_TMP_ROOT "/example_v1-abc.tar.gz",
                                 &info));
    TEST_ASSERT(info.alias);
    TEST_ASSERT(info.valid);
}

static void test_pkg_missing_version_fails(void) {
    PackageInfo info;

    mkdir_p(TEST_TMP_ROOT "/missing-version");
    write_file(TEST_TMP_ROOT "/missing-version/manifest.json", "{}\n");
    run_cmd("tar -czf " TEST_TMP_ROOT "/missing.tar.gz -C "
            TEST_TMP_ROOT "/missing-version .");

    TEST_ASSERT(!pkg_validate_tar("example", "v1-abc",
                                  TEST_TMP_ROOT "/missing.tar.gz",
                                  &info));
    TEST_ASSERT(info.exists);
    TEST_ASSERT(!info.valid);
}

static void test_pkg_publish_from_dir(void) {
    char publishedPath[WIMC_MAX_PATH_LEN];
    char actualVersion[WIMC_MAX_NAME_LEN];
    PackageInfo info;

    make_pkg_dir(TEST_TMP_ROOT "/pkg-v1-abc", "v1-abc\n");

    TEST_EQ_INT(0, pkg_publish_from_dir("example", "latest", "uuid-test",
                                        TEST_TMP_ROOT "/pkg-v1-abc",
                                        publishedPath, sizeof(publishedPath),
                                        actualVersion,
                                        sizeof(actualVersion)));
    TEST_EQ_STR("v1-abc", actualVersion);
    TEST_EQ_STR("/ukama/apps/pkgs/example_latest.tar.gz", publishedPath);
    TEST_ASSERT(pkg_validate_tar("example", "latest", publishedPath, &info));
}

static void test_db_lifecycle(void) {
    sqlite3 *db = NULL;
    char *status = NULL;
    char *path = NULL;
    char *actualVersion = NULL;
    char *error = NULL;

    TEST_ASSERT(db_open_or_create(TEST_DB_FILE, &db));
    TEST_NOT_NULL(db);

    TEST_ASSERT(db_insert_entry(db, "example", "latest", WIMC_STATUS_QUEUED));
    TEST_ASSERT(db_update_package_status(db, "example", "latest",
                                         "/ukama/apps/pkgs/example_latest.tar.gz",
                                         WIMC_STATUS_AVAILABLE,
                                         "v1-abc", NULL));

    TEST_ASSERT(db_read_package(db, "example", "latest", &status, &path,
                                &actualVersion, &error));
    TEST_EQ_STR(WIMC_STATUS_AVAILABLE, status);
    TEST_EQ_STR("/ukama/apps/pkgs/example_latest.tar.gz", path);
    TEST_EQ_STR("v1-abc", actualVersion);
    TEST_NULL(error);

    free(status);
    free(path);
    free(actualVersion);
    free(error);
    sqlite3_close(db);
}

static void test_db_stale_downloads_failed(void) {
    sqlite3 *db = NULL;
    char *status = NULL;

    TEST_ASSERT(db_open_or_create(TEST_DB_FILE, &db));
    TEST_ASSERT(db_insert_entry(db, "example", "v1-abc",
                                WIMC_STATUS_DOWNLOADING));
    TEST_ASSERT(db_mark_old_downloads_failed(db));
    TEST_ASSERT(db_read_status(db, "example", "v1-abc", &status));
    TEST_EQ_STR(WIMC_STATUS_FAILED, status);

    free(status);
    sqlite3_close(db);
}

static void test_tasks_add_find_delete(void) {
    WimcReq *req1;
    WimcReq *req2;
    WTasks *tasks = NULL;
    WTasks *found;

    req1 = make_fetch_request("example", "v1-abc", "chunk");
    req2 = make_fetch_request("example", "v1-xyz", "chunk");

    add_to_tasks(&tasks, req1);
    add_to_tasks(&tasks, req2);

    TEST_NOT_NULL(tasks);
    TEST_NOT_NULL(tasks->next);

    found = find_task_by_uuid(tasks, req1->fetch->uuid);
    TEST_NOT_NULL(found);
    TEST_EQ_STR("example", found->content->name);
    TEST_EQ_STR("v1-abc", found->content->tag);
    TEST_EQ_INT(1024, found->content->expectedSizeBytes);

    delete_from_tasks(&tasks, found);
    TEST_NOT_NULL(tasks);
    TEST_EQ_STR("v1-xyz", tasks->content->tag);

    clear_tasks(&tasks);
    TEST_NULL(tasks);

    free_fetch_request(req1);
    free_fetch_request(req2);
}

static void test_hub_response_deserialize(void) {
    const char *body =
        "{"
        "\"name\":\"example\","
        "\"artifacts\":[{"
        "\"version\":\"v1-abc\","
        "\"formats\":[{"
        "\"type\":\"chunk\","
        "\"url\":\"/fake/example/v1-abc.caidx\","
        "\"created_at\":\"2026-01-01T00:00:00Z\","
        "\"extra_info\":{\"chunks\":\"/fake/chunks/\"}"
        "}]"
        "}]"
        "}";
    json_error_t err;
    json_t *json;
    Artifact **artifacts = NULL;
    int count = 0;

    json = json_loads(body, 0, &err);
    TEST_NOT_NULL(json);

    TEST_ASSERT(deserialize_hub_response(&artifacts, &count, json));
    TEST_EQ_INT(1, count);
    TEST_EQ_STR("example", artifacts[0]->name);
    TEST_EQ_STR("v1-abc", artifacts[0]->version);
    TEST_EQ_STR("chunk", artifacts[0]->formats[0]->type);
    TEST_EQ_STR("/fake/chunks/", artifacts[0]->formats[0]->extraInfo);

    json_decref(json);
}

static void test_wimc_request_serialize(void) {
    WimcReq *req;
    json_t *json = NULL;

    req = make_fetch_request("example", "v1-abc", "chunk");
    TEST_ASSERT(serialize_wimc_request(req, &json));
    TEST_NOT_NULL(json);

    TEST_EQ_STR("example", json_string_value(json_object_get(json,
                                                              JSON_NAME)));
    TEST_EQ_STR("v1-abc", json_string_value(json_object_get(json,
                                                             JSON_TAG)));
    TEST_EQ_INT(1024, json_integer_value(json_object_get(json,
                                                         JSON_EXPECTED_SIZE)));

    json_decref(json);
    free_fetch_request(req);
}

static void run_test(const char *name, void (*fn)(void)) {
    printf("RUN %s\n", name);
    reset_test_dirs();
    fn();
}

int main(void) {
    run_test("test_pkg_identifiers", test_pkg_identifiers);
    run_test("test_pkg_path_for_tag", test_pkg_path_for_tag);
    run_test("test_pkg_version_from_dir_and_tar", test_pkg_version_from_dir_and_tar);
    run_test("test_pkg_validate_tar", test_pkg_validate_tar);
    run_test("test_pkg_missing_version_fails", test_pkg_missing_version_fails);
    run_test("test_pkg_publish_from_dir", test_pkg_publish_from_dir);
    run_test("test_db_lifecycle", test_db_lifecycle);
    run_test("test_db_stale_downloads_failed", test_db_stale_downloads_failed);
    run_test("test_tasks_add_find_delete", test_tasks_add_find_delete);
    run_test("test_hub_response_deserialize", test_hub_response_deserialize);
    run_test("test_wimc_request_serialize", test_wimc_request_serialize);

    if (g_failures != 0) {
        fprintf(stderr, "FAILED: %d failures across %d checks\n",
                g_failures, g_tests);
        return 1;
    }

    printf("PASS: %d checks\n", g_tests);
    return 0;
}
