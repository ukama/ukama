/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "util.h"
#include "config_macros.h"
#include <dirent.h>
#include <sys/stat.h>
#include <unistd.h>

#define MOVE_DIR true
#define COPY_DIR false

int is_valid_json(const char *json_string) {
	json_error_t error;
	json_t *json = json_loads(json_string, 0, &error);

	if (json != NULL) {
		//usys_log_debug("Json data is : \n %s", json_dumps(json, JSON_INDENT(4)));
		json_decref(json); // Release the JSON object
		return 1; // Valid JSON
	} else {
		// Invalid JSON
		fprintf(stderr, "Error: JSON parsing error at line %d, column %d: %s\n",
				error.line, error.column, error.text);
		return 0;
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

int clean_empty_dir(char* path) {
	struct dirent *entry;
	struct stat st;
	DIR *dir = opendir(path);
	usys_log_debug("Removing empty directories from %s", path);

	while ((entry = readdir(dir))) {
		// Skip . and ..
		if (strcmp(entry->d_name, ".") == 0 || strcmp(entry->d_name, "..") == 0)
			continue;

		char sourcePath[512];
		snprintf(sourcePath, sizeof(sourcePath), "%s/%s", path, entry->d_name);

		if (lstat(sourcePath, &st) == -1) {
			usys_log_error("Failed getting file status for %s", sourcePath);
			perror("Error getting file status");
			return -1;
		}

		if (S_ISDIR(st.st_mode)) {
			// If it's a directory, recursively move it
			if (is_dir_empty(sourcePath)) {
				remove_dir(sourcePath);
			} else {
				clean_empty_dir(sourcePath);
			}
		}
	}

	closedir(dir);
	return 0;

}

int make_path(const char* path) {
	char* p = NULL;
	char* token = NULL;
	char pathCopy[512]; // Adjust the buffer size as needed
	char npath[512]={'\0'};
	// Create a copy of the path to avoid modifying the original
	usys_strncpy(pathCopy, path, sizeof(pathCopy));

	// Tokenize the path by "/"
	p = pathCopy;
	while ((token = strsep(&p, "/")) != NULL) {
		if (usys_strlen(token) == 0) {
			continue;  // Skip empty tokens
		}

		// Append the token to the current path
		usys_strcat(npath, "/");
		usys_strcat(npath, token);

		// Check if the directory already exists
		struct stat st;
		if (stat(npath, &st) != 0) {
			// If it doesn't exist, create it
			if (mkdir(npath, 0777) != 0) {
				usys_log_error("Failed to create directory: %s\n", npath);
				perror("error");
				return 0; // Return 0 to indicate failure
			}
		}
	}

	return 1; // Return 1 to indicate success
}

/* Remove config file */
int remove_config(ConfigData *c) {
	char path[512] = {'\0'};
	sprintf(path,"%s/%s/%s/%s", CONFIG_TMP_PATH, c->version, c->app, c->fileName);
	if (remove(path) != 0) {
		usys_log_error("Failed removing config %s", path);
		perror("Error deleting file");
		return -1; // Error deleting file
	}
	return 0;
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
		if (mkdir(destination, 0777) != 0) {
			usys_log_error("Failed creating dir %s", destination);
			perror("Error creating destination directory");
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
	struct dirent *entry;
	struct stat st;
	DIR *dir = opendir(path);

	if (dir == NULL) {
		usys_log_error("Failed opening dir %s", dir);
		perror("Failed to open directory");
		return -1; // Error opening directory
	}

	while ((entry = readdir(dir))) {
		if (strcmp(entry->d_name, ".") == 0 || strcmp(entry->d_name, "..") == 0)
			continue;

		char entryPath[PATH_MAX];
		snprintf(entryPath, sizeof(entryPath), "%s/%s", path, entry->d_name);

		if (lstat(entryPath, &st) == -1) {
			usys_log_error("Failed getting file status for %s", entryPath);
			perror("Error getting file status");
			continue;
		}

		if (S_ISDIR(st.st_mode)) {
			if (remove_dir(entryPath) != 0) {
				closedir(dir);
				return -1; // Error deleting subdirectory
			}
		} else {
			if (remove(entryPath) != 0) {
				usys_log_error("Failed deleting file %s", entryPath);
				perror("Error deleting file");
				closedir(dir);
				return -1; // Error deleting file
			}
		}
	}

	closedir(dir);

	if (rmdir(path) != 0) {
		usys_log_error("Failed deleting dir %s", path);
		perror("Error deleting directory");
		return -1; // Error deleting directory
	}

	return 0; // Directory and its contents deleted successfully
}

int create_config(ConfigData* c) {
	char path[512] = {'\0'};
	char fpath[512] = {'\0'};
	sprintf(path,"%s/%s/%s", CONFIG_TMP_PATH, c->version, c->app);

	remove(path);
	if (make_path(path)) {
		usys_log_debug("Directory %s created successfully.\n", path);

		sprintf(fpath,"%s/%s", path, c->fileName);
		// Create and write to files in the new directory
		FILE* file = fopen(fpath, "w");
		if (file == NULL) {
			usys_log_error("Failed to create file %s\n", fpath);
			perror(NULL);
			return -1;
		}

		if (c->data != NULL) {
			if(fputs(c->data, file) == EOF) {
				perror("Failed to write to file");
				fclose(file); // Close the file
				return -1; // Return an error code
			}
		}
		fclose(file);
		usys_log_debug("File %s created successfully.\n", fpath);

	} else {
		printf("Failed to create directory.\n");
		perror("error");
		return -1;
	}

	return 0;
}

int create_backup_config(){

	remove_dir(CONFIG_OLD);
	if (clone_dir(CONFIG_BACKUP, CONFIG_OLD, MOVE_DIR) == 0) {
		usys_log_debug("Moved backup config to old config.\n");
	} else {
		usys_log_error("failed to create old config.\n");
		perror("error");
		return -1;
	}

	remove_dir(CONFIG_BACKUP);
	if (clone_dir(CONFIG_RUNNING, CONFIG_BACKUP, MOVE_DIR) == 0) {
		usys_log_debug("Created a backup config to old config.\n");
	}else {
		usys_log_error("failed to create backup config.\n");
		perror("error");
		return -1;
	}

	return 0;
}

int restore_config() {

	remove_dir(CONFIG_RUNNING);
	if (clone_dir(CONFIG_BACKUP, CONFIG_RUNNING, MOVE_DIR) == 0) {
		usys_log_debug("Restore running config done.\n");
	}else {
		usys_log_error("failed to restore running config.\n");
		perror("error");
		return -1;
	}

	remove_dir(CONFIG_BACKUP);
	if (clone_dir(CONFIG_OLD, CONFIG_BACKUP, MOVE_DIR ) == 0) {
		usys_log_debug("Restore backup config done.\n");
	}else {
		usys_log_error("failed to restore backup config.\n");
		perror("error");
		return -1;
	}

	return 0;

}

int store_config(char* version) {
	char sPath[512] = {'\0'};
	sprintf(sPath,"%s/%s", CONFIG_TMP_PATH, version);

	/* Create a backup */
	if (create_backup_config() != 0) {
		usys_log_error("Failed to move old config for backup. \n");
		perror("error");
		return -1;
	}

	/* Create a config */
	remove_dir(CONFIG_RUNNING);
	if (clone_dir(sPath, CONFIG_RUNNING, MOVE_DIR) != 0) {
		usys_log_error("Failed to create config for backup. Restoring config..\n");
		perror("error");
		restore_config();
		return -1;
	}

	return 0;
}

int prepare_for_new_config(ConfigData* c) {
	char path[512] = {'\0'};
	sprintf(path,"%s/%s", CONFIG_TMP_PATH, c->version);

	/* remove old residue */
	remove_dir(path);

	/* Create new folder and copy current running config */
	if (clone_dir(CONFIG_RUNNING, path, COPY_DIR) != 0) {
		usys_log_error("Failed to prepare tmp config for download.\n");
		perror("error");
		return -1;
	}

	return 0;
}


