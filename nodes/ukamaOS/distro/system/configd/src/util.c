/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <dirent.h>
#include <sys/stat.h>
#include <unistd.h>

#include "util.h"

#define MOVE_DIR true
#define COPY_DIR false

bool is_valid_json(const char *json_string) {

	json_error_t error;
	json_t *json = NULL;

    json = json_loads(json_string, 0, &error);
	if (json != NULL) {
		json_decref(json);
		return USYS_TRUE;
	} else {
		usys_log_error("Error: JSON parsing error at line %d, column %d: %s\n",
                       error.line, error.column, error.text);
		return USYS_FALSE;
	}
}

int is_dir_empty(const char *directoryPath) {
    DIR *dir = opendir(directoryPath);

    if (dir == NULL) {
        perror("Unable to open directory");
        return -1;  // Return -1 to indicate an error
    }

    struct dirent *entry;
    int entriesCount = 0;

    while ((entry = readdir(dir)) != NULL) {
        if (strcmp(entry->d_name, ".") == 0 || strcmp(entry->d_name, "..") == 0) {
            continue; // Skip current directory and parent directory entries
        }
        entriesCount++;
    }

    closedir(dir);

    // If there are no entries other than "." and "..", the directory is empty
    return entriesCount == 0;
}

static int mkdir_p(const char *path, mode_t mode) {

    char temp[PATH_MAX];
    char *p = NULL;
    size_t len;

    snprintf(temp, sizeof(temp), "%s", path);
    len = strlen(temp);
    if (len > 0 && temp[len - 1] == '/') {
        temp[len - 1] = 0;
    }

    for (p = temp + 1; *p; p++) {
        if (*p == '/') {
            *p = 0;

            if (mkdir(temp, mode) != 0) {
                if (errno != EEXIST) {
                    return -1;
                }
            }
            *p = '/';
        }
    }

    if (mkdir(temp, mode) != 0) {
        if (errno != EEXIST) {
            return -1;
        }
    }

    return 0;
}

bool remove_config_file_from_staging_area(SessionData *s) {

	char path[MAX_PATH] = {0};

	sprintf(path,"%s/%d/%s/%s", CONFIG_TMP_PATH, s->timestamp, s->app, s->fileName);
	if (remove(path) != 0) {
		usys_log_error("Failed removing config: %s Error: %s",
                       path, strerror(errno));
		return USYS_FALSE;
	}

	return USYS_TRUE;
}

/* Copy or move file.
 * Flag true => move
 * Flag false => copy
 */
int clone_file(const char *source, const char *destination, bool flag) {
	FILE *src, *dest;
	char ch;

	src = fopen(source, "rb");
	if (src == NULL) {
		usys_log_error("Failed creating %s", source);
		perror("Error opening source file");
		return 1;
	}

	dest = fopen(destination, "wb");
	if (dest == NULL) {
		usys_log_error("Failed creating %s", destination);
		perror("Error creating destination file");
		fclose(src);
		return 1;
	}

	while ((ch = fgetc(src)) != EOF) {
		fputc(ch, dest);
	}

	fclose(src);
	fclose(dest);
	if (flag) {
		if (remove(source) != 0) {
			usys_log_error("Failed removing %s", source);
			perror("Error removing source file");
			return 1;
		}
	}

	return 0;

}

/* Copy or move directory.
 * Flag true => move
 * Flag false => copy
 */
int clone_dir(const char *source, const char *destination, bool flag) {

	struct dirent *entry;
	struct stat st;
	DIR *dir = opendir(source);

	// Check if the source directory exists
	if (dir == NULL) {
		usys_log_error("Failed opening %s", source);
		perror("Failed to open source directory");
		return -1;
	}

	// Create the destination directory
	if (stat(destination, &st) != 0) {
		if (mkdir_p(destination, 0777) != 0) {
			usys_log_error("Failed creating dir %s", destination);
			closedir(dir);
			return -1;
		}
	}

	while ((entry = readdir(dir))) {
		// Skip . and ..
		if (strcmp(entry->d_name, ".") == 0 || strcmp(entry->d_name, "..") == 0)
			continue;

		char sourcePath[512];
		char destPath[512];
		snprintf(sourcePath, sizeof(sourcePath), "%s/%s", source, entry->d_name);
		snprintf(destPath, sizeof(destPath), "%s/%s", destination, entry->d_name);

		if (lstat(sourcePath, &st) == -1) {
			usys_log_error("Failed getting file status for %s", sourcePath);
			perror("Error getting file status");
			return -1;
		}

		if (S_ISDIR(st.st_mode)) {

			if (is_dir_empty(sourcePath)) {
				int rs = remove(sourcePath);
				usys_log_trace("Empty dir %s is %s removed", sourcePath, ((rs==0)? "successfully":"not"));
				continue;
			}

			// If it's a directory, recursively move it
			if (clone_dir(sourcePath, destPath, flag)==0) {
				if (flag) {
					remove_dir(sourcePath);
				}
			} else {
				return -1;
			}
		} else {
			// If it's a file, move it
			if (clone_file(sourcePath, destPath, flag) != 0) {
				usys_log_error("Failed moving file from %s to %s", sourcePath, destPath);
				perror("Error moving file");
				return -1;
			}
		}
	}

	closedir(dir);
	return 0;
}

int remove_dir(const char *path) {

	struct dirent *entry = NULL;
	struct stat st;
	DIR *dir = NULL;

    dir = opendir(path);
	if (dir == NULL) {
		usys_log_error("Failed opening dir %s", dir);
		usys_log_error("Failed to open directory: %s. Error: %s",
                       path, strerror(errno));
		return -1;
	}

	while ((entry = readdir(dir))) {
		if (strcmp(entry->d_name, ".") == 0 || strcmp(entry->d_name, "..") == 0)
			continue;

		char entryPath[PATH_MAX];
		snprintf(entryPath, sizeof(entryPath), "%s/%s", path, entry->d_name);

		if (lstat(entryPath, &st) == -1) {
			usys_log_error("Failed getting file status for %s Error: %s",
                           entryPath, strerror(errno));
			continue;
		}

		if (S_ISDIR(st.st_mode)) {
			if (remove_dir(entryPath) != 0) {
				closedir(dir);
				return -1;
			}
		} else {
			if (remove(entryPath) != 0) {
				usys_log_error("Failed deleting file %s Error: %s",
                               entryPath, strerror(errno));
				closedir(dir);
				return -1;
			}
		}
	}

	closedir(dir);

	if (rmdir(path) != 0) {
		usys_log_error("Failed deleting dir %s Error: %s",
                       path, strerror(errno));
		return -1;
	}

	return 0; // Directory and its contents deleted successfully
}

bool create_config_file_in_staging_area(SessionData *s) {

    char fpath[MAX_FILE_PATH] = {0};
    FILE *file = NULL;

    if (s->data == NULL) return USYS_FALSE;

    snprintf(fpath, sizeof(fpath), "%s/%d/%s/%s",
             CONFIG_TMP_PATH, s->timestamp, s->app, s->fileName);

    /* overwrite an existing or create new file */
    file = fopen(fpath, "w+");
    if (file == NULL) {
        usys_log_error("Failed to create file %s Error: %s",
                       fpath, strerror(errno));
        return USYS_FALSE;
    }

    if (fputs(s->data, file) == EOF) {
        usys_log_error("Failed to write to file %s Error: %s",
                       fpath, strerror(errno));
        fclose(file);
        return USYS_FALSE;
    }

    fclose(file);
    usys_log_debug("File %s created/updated successfully", fpath);

    return USYS_TRUE;
}
