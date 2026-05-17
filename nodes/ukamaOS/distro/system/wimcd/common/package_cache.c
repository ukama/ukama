/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <ctype.h>
#include <dirent.h>
#include <errno.h>
#include <fcntl.h>
#include <limits.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#include "package_cache.h"
#include "usys_log.h"

#define PKG_SUFFIX      ".tar.gz"
#define PKG_TMP_DIR     DEFAULT_APPS_PKGS_PATH "/.tmp"
#define VERSION_FILE    "VERSION"

static int mkdir_p(const char *path, mode_t mode) {

    char tmp[PATH_MAX];
    char *p;

    if (path == NULL || *path == '\0') {
        return -1;
    }

    if (strlen(path) >= sizeof(tmp)) {
        return -1;
    }

    snprintf(tmp, sizeof(tmp), "%s", path);

    for (p = tmp + 1; *p != '\0'; p++) {
        if (*p == '/') {
            *p = '\0';
            if (mkdir(tmp, mode) != 0 && errno != EEXIST) {
                return -1;
            }
            *p = '/';
        }
    }

    if (mkdir(tmp, mode) != 0 && errno != EEXIST) {
        return -1;
    }

    return 0;
}

static int rm_rf(const char *path) {

    DIR *dir;
    struct dirent *ent;
    struct stat st;
    char child[PATH_MAX];

    if (path == NULL || lstat(path, &st) != 0) {
        return 0;
    }

    if (!S_ISDIR(st.st_mode)) {
        return unlink(path);
    }

    dir = opendir(path);
    if (dir == NULL) {
        return -1;
    }

    while ((ent = readdir(dir)) != NULL) {
        if (strcmp(ent->d_name, ".") == 0 ||
            strcmp(ent->d_name, "..") == 0) {
            continue;
        }

        if (snprintf(child, sizeof(child), "%s/%s",
                     path, ent->d_name) >= (int)sizeof(child)) {
            closedir(dir);
            return -1;
        }

        if (rm_rf(child) != 0) {
            closedir(dir);
            return -1;
        }
    }

    closedir(dir);
    return rmdir(path);
}

static int run_wait(char *const argv[]) {

    pid_t pid;
    int status;

    pid = fork();
    if (pid < 0) {
        return -1;
    }

    if (pid == 0) {
        execvp(argv[0], argv);
        _exit(127);
    }

    if (waitpid(pid, &status, 0) < 0) {
        return -1;
    }

    if (!WIFEXITED(status) || WEXITSTATUS(status) != 0) {
        return -1;
    }

    return 0;
}

static char *trim_line(char *s) {

    char *end;

    if (s == NULL) {
        return NULL;
    }

    while (*s != '\0' && isspace((unsigned char)*s)) {
        s++;
    }

    end = s + strlen(s);
    while (end > s && isspace((unsigned char)*(end - 1))) {
        end--;
    }
    *end = '\0';

    return s;
}

static int file_size_ok(const char *path) {

    struct stat st;

    if (path == NULL || stat(path, &st) != 0) {
        return 0;
    }

    if (!S_ISREG(st.st_mode) && !S_ISLNK(st.st_mode)) {
        return 0;
    }

    return st.st_size > 0;
}
 
bool pkg_is_valid_identifier(const char *value) {

    size_t i;

    if (value == NULL || *value == '\0') {
        return false;
    }

    if (strlen(value) >= WIMC_MAX_NAME_LEN) {
        return false;
    }

    for (i = 0; value[i] != '\0'; i++) {
        if ((value[i] >= 'a' && value[i] <= 'z') ||
            (value[i] >= 'A' && value[i] <= 'Z') ||
            (value[i] >= '0' && value[i] <= '9') ||
            value[i] == '-' || value[i] == '_' || value[i] == '.') {
            continue;
        }
        return false;
    }

    return true;
}

bool pkg_is_alias_tag(const char *tag) {

    return tag != NULL && strcmp(tag, WIMC_ALIAS_LATEST) == 0;
}

int pkg_path_for_tag(const char *name, const char *tag,
                     char *path, size_t pathLen) {

    if (!pkg_is_valid_identifier(name) ||
        !pkg_is_valid_identifier(tag) || path == NULL) {
        return -1;
    }

    if (snprintf(path, pathLen, "%s/%s_%s%s", DEFAULT_APPS_PKGS_PATH,
                 name, tag, PKG_SUFFIX) >= (int)pathLen) {
        return -1;
    }

    return 0;
}

int pkg_tmp_tar_path(const char *uuidStr, const char *name,
                     const char *tag, char *path, size_t pathLen) {

    if (!pkg_is_valid_identifier(uuidStr) ||
        !pkg_is_valid_identifier(name) ||
        !pkg_is_valid_identifier(tag) || path == NULL) {
        return -1;
    }

    if (snprintf(path, pathLen, "%s/%s/%s_%s%s", PKG_TMP_DIR,
                 uuidStr, name, tag, PKG_SUFFIX) >= (int)pathLen) {
        return -1;
    }

    return 0;
}

int pkg_extract_path(const char *uuidStr, const char *name,
                     const char *tag, char *path, size_t pathLen) {

    if (!pkg_is_valid_identifier(uuidStr) ||
        !pkg_is_valid_identifier(name) ||
        !pkg_is_valid_identifier(tag) || path == NULL) {
        return -1;
    }

    if (snprintf(path, pathLen, "%s/%s/%s_%s", PKG_TMP_DIR,
                 uuidStr, name, tag) >= (int)pathLen) {
        return -1;
    }

    return 0;
}

int pkg_ensure_cache_dirs(void) {

    if (mkdir_p(DEFAULT_APPS_PKGS_PATH, 0755) != 0) {
        usys_log_error("Failed to create package dir: %s",
                       DEFAULT_APPS_PKGS_PATH);
        return -1;
    }

    if (mkdir_p(PKG_TMP_DIR, 0700) != 0) {
        usys_log_error("Failed to create package tmp dir: %s", PKG_TMP_DIR);
        return -1;
    }

    return 0;
}

int pkg_cleanup_tmp(void) {

    DIR *dir;
    struct dirent *ent;
    char path[PATH_MAX];

    if (pkg_ensure_cache_dirs() != 0) {
        return -1;
    }

    dir = opendir(PKG_TMP_DIR);
    if (dir == NULL) {
        return -1;
    }

    while ((ent = readdir(dir)) != NULL) {
        if (strcmp(ent->d_name, ".") == 0 ||
            strcmp(ent->d_name, "..") == 0) {
            continue;
        }

        if (snprintf(path, sizeof(path), "%s/%s", PKG_TMP_DIR,
                     ent->d_name) >= (int)sizeof(path)) {
            closedir(dir);
            return -1;
        }

        if (rm_rf(path) != 0) {
            usys_log_warn("Failed to cleanup tmp package path: %s", path);
        }
    }

    closedir(dir);
    return 0;
}

int pkg_read_version_from_dir(const char *dir, char *version,
                              size_t versionLen) {

    FILE *fp;
    char path[PATH_MAX];
    char line[WIMC_MAX_NAME_LEN];
    char *value;

    if (dir == NULL || version == NULL || versionLen == 0) {
        return -1;
    }

    if (snprintf(path, sizeof(path), "%s/%s", dir,
                 VERSION_FILE) >= (int)sizeof(path)) {
        return -1;
    }

    fp = fopen(path, "r");
    if (fp == NULL) {
        return -1;
    }

    if (fgets(line, sizeof(line), fp) == NULL) {
        fclose(fp);
        return -1;
    }
    fclose(fp);

    value = trim_line(line);
    if (!pkg_is_valid_identifier(value) || pkg_is_alias_tag(value)) {
        return -1;
    }

    snprintf(version, versionLen, "%s", value);
    return 0;
}

/* VERSION
 * ./VERSION
 * some-top-dir/VERSION
 */
int pkg_read_version_from_tar(const char *path, char *version,
                              size_t versionLen) {

    FILE *fp;
    char cmd[PATH_MAX * 2];
    char extractCmd[PATH_MAX * 2];
    char line[PATH_MAX];
    char entry[PATH_MAX];
    char *value;
    size_t len;

    if (path == NULL || version == NULL || versionLen == 0) {
        return -1;
    }

    if (strchr(path, '\'') != NULL) {
        return -1;
    }

    memset(cmd, 0, sizeof(cmd));
    memset(extractCmd, 0, sizeof(extractCmd));
    memset(line, 0, sizeof(line));
    memset(entry, 0, sizeof(entry));

    usys_log_debug("Trying to read VERSION from tgz: %s", path);

    if (snprintf(cmd, sizeof(cmd),
                 "tar -tzf '%s' 2>/dev/null", path) >=
        (int)sizeof(cmd)) {
        return -1;
    }

    fp = popen(cmd, "r");
    if (fp == NULL) {
        return -1;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        value = trim_line(line);
        if (value == NULL || *value == '\0') {
            continue;
        }

        if (strchr(value, '\'') != NULL ||
            strstr(value, "..") != NULL) {
            continue;
        }

        if (strcmp(value, VERSION_FILE) == 0 ||
            strcmp(value, "./" VERSION_FILE) == 0) {
            snprintf(entry, sizeof(entry), "%s", value);
            break;
        }

        len = strlen(value);
        if (len > strlen("/" VERSION_FILE) &&
            strcmp(value + len - strlen("/" VERSION_FILE),
                   "/" VERSION_FILE) == 0) {
            snprintf(entry, sizeof(entry), "%s", value);
            break;
        }
    }

    if (pclose(fp) == -1) {
        return -1;
    }

    if (entry[0] == '\0') {
        return -1;
    }

    if (snprintf(extractCmd, sizeof(extractCmd),
                 "tar -xO -zf '%s' '%s' 2>/dev/null | head -n 1",
                 path, entry) >= (int)sizeof(extractCmd)) {
        return -1;
    }

    fp = popen(extractCmd, "r");
    if (fp == NULL) {
        return -1;
    }

    memset(line, 0, sizeof(line));
    if (fgets(line, sizeof(line), fp) == NULL) {
        pclose(fp);
        return -1;
    }

    if (pclose(fp) == -1) {
        return -1;
    }

    value = trim_line(line);
    if (!pkg_is_valid_identifier(value) || pkg_is_alias_tag(value)) {
        return -1;
    }

    snprintf(version, versionLen, "%s", value);
    return 0;
}

static const char *pkg_version_without_v(const char *version)
{
    if (version == NULL) {
        return NULL;
    }

    if (version[0] == 'v' || version[0] == 'V') {
        return version + 1;
    }

    return version;
}

static bool pkg_versions_match(const char *expected, const char *actual)
{
    const char *expectedNorm;
    const char *actualNorm;

    if (expected == NULL || actual == NULL) {
        return false;
    }

    expectedNorm = pkg_version_without_v(expected);
    actualNorm   = pkg_version_without_v(actual);

    return strcmp(expectedNorm, actualNorm) == 0;
}

int pkg_validate_tar(const char *name, const char *tag, const char *path,
                     PackageInfo *info)
{
    char actualVersion[WIMC_MAX_NAME_LEN];

    if (info != NULL) {
        memset(info, 0, sizeof(PackageInfo));
        if (name != NULL) snprintf(info->name, sizeof(info->name), "%s", name);
        if (tag != NULL) snprintf(info->tag, sizeof(info->tag), "%s", tag);
        if (path != NULL) snprintf(info->path, sizeof(info->path), "%s", path);
        info->alias = pkg_is_alias_tag(tag);
    }

    if (!pkg_is_valid_identifier(name) ||
        !pkg_is_valid_identifier(tag) || path == NULL) {
        if (info != NULL) {
            snprintf(info->error, sizeof(info->error), "invalid input");
        }
        return 0;
    }

    if (!file_size_ok(path)) {
        if (info != NULL) {
            info->exists = false;
            snprintf(info->error, sizeof(info->error), "file missing/empty");
        }
        return 0;
    }

    if (info != NULL) {
        info->exists = true;
    }

    memset(actualVersion, 0, sizeof(actualVersion));
    if (pkg_read_version_from_tar(path, actualVersion,
                                  sizeof(actualVersion)) != 0) {
        if (info != NULL) {
            snprintf(info->error, sizeof(info->error),
                     "missing/invalid VERSION");
        }
        return 0;
    }

    if (!pkg_is_alias_tag(tag) && !pkg_versions_match(tag, actualVersion)) {
        if (info != NULL) {
            snprintf(info->actualVersion, sizeof(info->actualVersion),
                     "%s", actualVersion);
            snprintf(info->error, sizeof(info->error),
                     "VERSION does not match requested tag");
        }

        usys_log_error("Version mismatch: Expected: %s Got: %s",
                       tag, actualVersion);

        return 0;
    }

    if (info != NULL) {
        info->valid = true;
        snprintf(info->actualVersion, sizeof(info->actualVersion),
                 "%s", actualVersion);
    }

    return 1;
}

static int make_immutable_package(const char *dir, const char *tmpTar) {

    char *argv[] = { "tar", "-czf", (char *)tmpTar, "-C", (char *)dir,
                     ".", NULL };

    return run_wait(argv);
}

static int publish_alias(const char *name,
                         const char *tag,
                         const char *actualVersion,
                         const char *actualPath) {

    char aliasPath[WIMC_MAX_PATH_LEN];
    const char *base;

    if (name == NULL || tag == NULL || actualVersion == NULL ||
        actualPath == NULL) {
        return -1;
    }

    /*
     * Create alias for:
     *   latest -> actual version
     *   1.0.1-abcdefgh -> v1.0.1-abcdefgh
     *
     * Do not create alias when requested tag and actual version are identical.
     */
    if (strcmp(tag, WIMC_ALIAS_LATEST) != 0 &&
        strcmp(tag, actualVersion) == 0) {
        return 0;
    }

    if (strcmp(tag, WIMC_ALIAS_LATEST) != 0 &&
        !pkg_versions_match(tag, actualVersion)) {
        return 0;
    }

    if (pkg_path_for_tag(name, tag, aliasPath, sizeof(aliasPath)) != 0) {
        return -1;
    }

    if (strcmp(aliasPath, actualPath) == 0) {
        return 0;
    }

    base = strrchr(actualPath, '/');
    base = base ? base + 1 : actualPath;

    unlink(aliasPath);
    if (symlink(base, aliasPath) != 0) {
        usys_log_error("Failed to create alias %s -> %s: %s",
                       aliasPath, base, strerror(errno));
        return -1;
    }

    return 0;
}

int pkg_publish_from_dir(const char *name, const char *tag,
                         const char *uuidStr, const char *dir,
                         char *publishedPath, size_t publishedPathLen,
                         char *actualVersion,
                         size_t actualVersionLen) {

    char version[WIMC_MAX_NAME_LEN];
    char actualPath[WIMC_MAX_PATH_LEN];
    char tmpTar[WIMC_MAX_PATH_LEN];
    char tmpParent[WIMC_MAX_PATH_LEN];
    char *slash;
    PackageInfo info;

    if (!pkg_is_valid_identifier(name) ||
        !pkg_is_valid_identifier(tag) ||
        !pkg_is_valid_identifier(uuidStr) || dir == NULL) {
        return -1;
    }

    if (pkg_read_version_from_dir(dir, version, sizeof(version)) != 0) {
        usys_log_error("Package %s:%s has missing/invalid VERSION",
                       name, tag);
        return -1;
    }

    if (!pkg_is_alias_tag(tag) && strcmp(tag, version) != 0) {
        usys_log_error("Package %s:%s VERSION mismatch: %s",
                       name, tag, version);
        return -1;
    }

    if (pkg_path_for_tag(name, version, actualPath, sizeof(actualPath)) != 0) {
        return -1;
    }

    if (pkg_tmp_tar_path(uuidStr, name, version, tmpTar,
                         sizeof(tmpTar)) != 0) {
        return -1;
    }

    snprintf(tmpParent, sizeof(tmpParent), "%s", tmpTar);
    slash = strrchr(tmpParent, '/');
    if (slash == NULL) {
        return -1;
    }
    *slash = '\0';

    if (mkdir_p(tmpParent, 0700) != 0) {
        return -1;
    }

    unlink(tmpTar);
    if (make_immutable_package(dir, tmpTar) != 0) {
        usys_log_error("Failed to create package tarball: %s", tmpTar);
        return -1;
    }

    if (!pkg_validate_tar(name, version, tmpTar, &info)) {
        usys_log_error("Created package failed validation: %s",
                       info.error[0] ? info.error : "unknown");
        unlink(tmpTar);
        return -1;
    }

    if (rename(tmpTar, actualPath) != 0) {
        usys_log_error("Failed to publish package %s: %s",
                       actualPath, strerror(errno));
        unlink(tmpTar);
        return -1;
    }

    if (publish_alias(name, tag, version, actualPath) != 0) {
        return -1;
    }

    if (publishedPath != NULL && publishedPathLen > 0) {
        if (pkg_is_alias_tag(tag)) {
            pkg_path_for_tag(name, tag, publishedPath, publishedPathLen);
        } else {
            snprintf(publishedPath, publishedPathLen, "%s", actualPath);
        }
    }

    if (actualVersion != NULL && actualVersionLen > 0) {
        snprintf(actualVersion, actualVersionLen, "%s", version);
    }

    return 0;
}

int pkg_publish_tar(const char *name,
                    const char *tag,
                    const char *tmpTar,
                    char *publishedPath, size_t publishedPathLen,
                    char *actualVersion, size_t actualVersionLen) {

    char actualPath[WIMC_MAX_PATH_LEN];
    PackageInfo info;

    if (!pkg_is_valid_identifier(name) ||
        !pkg_is_valid_identifier(tag) || tmpTar == NULL) {
        return -1;
    }

    memset(&info, 0, sizeof(info));
    if (!pkg_validate_tar(name, tag, tmpTar, &info)) {
        usys_log_error("Downloaded package failed validation: %s",
                       info.error[0] ? info.error : "unknown");
        return -1;
    }

    if (pkg_path_for_tag(name, info.actualVersion,
                         actualPath, sizeof(actualPath)) != 0) {
        return -1;
    }

    if (rename(tmpTar, actualPath) != 0) {
        usys_log_error("Failed to publish package %s: %s",
                       actualPath, strerror(errno));
        return -1;
    }

    if (publish_alias(name, tag, version, actualPath) != 0) {
        return -1;
    }

    if (publishedPath != NULL && publishedPathLen > 0) {
        if (pkg_is_alias_tag(tag)) {
            pkg_path_for_tag(name, tag, publishedPath, publishedPathLen);
        } else {
            snprintf(publishedPath, publishedPathLen, "%s", actualPath);
        }
    }

    if (actualVersion != NULL && actualVersionLen > 0) {
        snprintf(actualVersion, actualVersionLen, "%s", info.actualVersion);
    }

    return 0;
}
