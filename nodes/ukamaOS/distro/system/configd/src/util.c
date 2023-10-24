/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <stdlib.h>

#include "configd.h"
#include "jansson.h"
#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"


#define CONFIG_TMP_PATH "/tmp/"
#define CONFIG_STORE_PATH "/etc/config"
#define CONFIG_RUNNING "/etc/config/running"
#define CONFIG_BACKUP "/etc/config/backup"
#define CONFIG_OLD "/etc/config/old"

int is_valid_json(const char *json_string) {
    json_error_t error;
    json_t *json = json_loads(json_string, 0, &error);

    if (json != NULL) {
        json_decref(json); // Release the JSON object
        return 1; // Valid JSON
    } else {
        // Invalid JSON
        fprintf(stderr, "Error: JSON parsing error at line %d, column %d: %s\n",
                error.line, error.column, error.text);
        return 0;
    }
}

int make_path(const char* path) {
    char* p = NULL;
    char* token = NULL;
    char pathCopy[256]; // Adjust the buffer size as needed

    // Create a copy of the path to avoid modifying the original
    usys_strncpy(pathCopy, path, sizeof(pathCopy));

    // Tokenize the path by "/"
    p = pathCopy;
    while ((token = strsep(&p, "/")) != NULL) {
        if (usys_strlen(token) == 0) {
            continue;  // Skip empty tokens
        }

        // Append the token to the current path
        usys_strcat(pathCopy, "/");
        usys_strcat(pathCopy, token);

        // Check if the directory already exists
        struct stat st;
        if (stat(pathCopy, &st) != 0) {
            // If it doesn't exist, create it
            if (mkdir(pathCopy, 0777) != 0) {
            	usys_log_error("Failed to create directory: %s\n", pathCopy);
                return 0; // Return 0 to indicate failure
            }
        }
    }

    return 1; // Return 1 to indicate success
}

int move_dir(const char *source, const char *destination) {

	struct stat st;

	// Check if the source directory exists
	if (stat(source, &st) != 0) {
		printf("Source directory does not exist.\n");
		return -1;
	}

	// Create the destination directory
	if (mkdir(destination, 0777) != 0) {
		perror("Error creating destination directory");
		return -1;
	}

	// Move the source directory to the destination
	if (rename(source, destination) != 0) {
		perror("Error moving directory");
		return -1;
	}

	return 0;
}

int remove_directory(const char *path) {
    struct stat st;
    if (stat(path, &st) != 0) {
        return 0; // Directory doesn't exist
    }

    if (S_ISDIR(st.st_mode)) {
        if (rmdir(path) != 0) {
            perror("Error removing directory");
            return -1; // Error removing directory
        }
    }
    return 0; // Directory removed successfully
}

int create_config(ConfigData* c) {
	char path[512] = {'\0'};
	char fpath[512] = {'\0'};
	sprintf(path,"%s/%s/%s", CONFIG_TMP_PATH, c.version, c.app);

    if (make_path(path) == 0) {
    	usys_log_debug("Directory %s created successfully.\n", path);

        sprint(fpath, "%s/%s",path, c.fileName)
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
        usys_log_debug("File %s created successfully.\n", fpath);

    } else {
        printf("Failed to create directory.\n");
        perror("Error");
        return -1;
    }

    return 0;
}

int create_backup_config(){

	remove_dir(CONFIG_OLD);
	if (move_dir(CONFIG_BACKUP, CONFIG_OLD ) == 0) {
		usys_log_debug("Moved backup config to old config.\n");
	} else {
        usys_log_error("failed to create old config.\n");
        return -1;
    }

    if (move_dir(CONFIG_RUNNING, CONFIG_BACKUP ) == 0) {
		usys_log_debug("Created a backup config to old config.\n");
	}else {
        usys_log_error("failed to create backup config.\n");
        return -1;
    }

	return 0;
}

int restore_config() {

	if (move_dir(CONFIG_BACKUP, CONFIG_RUNNING ) == 0) {
		usys_log_debug("Restore running config done.\n");
	}else {
        usys_log_error("failed to restore running config.\n");
        return -1;
    }

    if (move_dir(CONFIG_OLD, CONFIG_BACKUP ) == 0) {
		usys_log_debug("Restore backup config done.\n");
	}else {
         usys_log_error("failed to restore backup config.\n");
        return -1;
    }

	return 0;

}

int store_config(string version) {
	char sPath[512] = {'\0'};
	sprintf(sPath,"%s/%s", CONFIG_DIR_PATH, version);
	
	/* Create a backup */
	if (create_backup_config() != 0) {
		usys_log_error("Failed to move old config for backup. \n");
		return -1;
	}

	/* Create a config */
	if (move_dir(sPath, CONFIG_RUNNING) != 0) {
		usys_log_error("Failed to create config for backup. Restoring config..\n");
		restore_config();
		return -1;
	}

	return 0;
}


