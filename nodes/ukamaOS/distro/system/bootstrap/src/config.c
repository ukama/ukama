/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>
#include <ctype.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#include "usys_log.h"

#include "config.h"
#include "toml.h"

int parse_config(Config *config, toml_table_t *configData) {

	toml_datum_t nodedHost;
	toml_datum_t meshConfig, remoteIPFile;

	/* sanity check */
	if (config == NULL) return FALSE;
	if (configData == NULL) return FALSE;

	/* Read the config data from the config.toml and load into Config. */
	/* noded-host */
	nodedHost = toml_string_in(configData, NODED_HOST);
	if (!nodedHost.ok) {
		usys_log_debug("[%s] is missing, setting to default: %s", NODED_HOST,
				  DEF_NODED_HOST);
		config->nodedHost = strdup(DEF_NODED_HOST);
	} else {
		config->nodedHost = strdup(nodedHost.u.s);
	}

	/* mesh-config */
	meshConfig = toml_string_in(configData, MESH_CONFIG);
	if (!meshConfig.ok) {
		usys_log_debug("[%s] is missing, setting to default: %s", MESH_CONFIG,
				  DEF_MESH_CONFIG);
		config->meshConfig = strdup(DEF_MESH_CONFIG);
	} else {
		config->meshConfig = strdup(meshConfig.u.s);
	}

	/* remote-IP-file */
	remoteIPFile = toml_string_in(configData, REMOTE_IP_FILE);
	if (!remoteIPFile.ok) {
		usys_log_debug("[%s] is missing, setting to default: %s", REMOTE_IP_FILE,
				  DEF_REMOTE_IP_FILE);
		config->remoteIPFile = strdup(DEF_REMOTE_IP_FILE);
	} else {
		config->remoteIPFile = strdup(remoteIPFile.u.s);
	}

	if (nodedHost.ok)       free(nodedHost.u.s);
	if (meshConfig.ok)      free(meshConfig.u.s);
	if (remoteIPFile.ok)    free(remoteIPFile.u.s);

    return TRUE;
}

static bool is_valid_hostname_strict(const char *s) {

    if (!s || !*s) return false;

    size_t n = strlen(s);
    if (n > 253) return false;                  // DNS max hostname length (practical)
    if (s[0] == '.' || s[0] == '-') return false;
    if (s[n - 1] == '.' || s[n - 1] == '-') return false;

    for (size_t i = 0; i < n; i++) {
        unsigned char c = (unsigned char)s[i];

        // Absolutely forbid any whitespace or control chars
        if (isspace(c) || iscntrl(c)) return false;

        // Allow only letters, digits, dot, hyphen
        if (!(isalnum(c) || c == '.' || c == '-')) return false;

        // Forbid consecutive dots
        if (c == '.' && i + 1 < n && s[i + 1] == '.') return false;
    }

    // Label rules: each label 1..63, no leading/trailing hyphen per label
    size_t label_len = 0;
    for (size_t i = 0; i <= n; i++) {
        char c = s[i];
        if (c == '.' || c == '\0') {
            if (label_len == 0 || label_len > 63) return false;

            // check label doesn't end with '-'
            if (i > 0 && s[i - 1] == '-') return false;

            // check label doesn't start with '-'
            size_t start = i - label_len;
            if (s[start] == '-') return false;

            label_len = 0;
        } else {
            label_len++;
        }
    }

    return true;
}

bool read_bootstrap_server_info(char **out) {
    if (!out) return false;
    *out = NULL;

    FILE *fp = fopen(DEF_BOOTSTRAP_FILE, "rb");
    if (!fp) {
        usys_log_error("Error opening bootstrap file '%s': %s",
                       DEF_BOOTSTRAP_FILE, strerror(errno));
        return false;
    }

    /* Read first line only; if line is longer than this, we reject */
    char line[512];
    if (!fgets(line, sizeof(line), fp)) {
        usys_log_error("Bootstrap file '%s' is empty/unreadable", DEF_BOOTSTRAP_FILE);
        fclose(fp);
        return false;
    }

    /* If there's a second line or extra bytes, we don't care, but we DO reject 
     * newline chars in the token.
     */
    fclose(fp);

    /* Strip ONLY trailing newline chars produced by fgets (we do not "trim" spaces/tabs/etc.)
     * This allows the file to be a normal text file with newline at end, 
     * but the hostname token itself will still be strict (no whitespace allowed inside).
     */
    size_t len = strcspn(line, "\r\n");
    line[len] = '\0';

    /* Reject if there was any whitespace/control in the token (strict) */
    if (!is_valid_hostname_strict(line)) {
        /* Helpful logging: show escaped bytes around the token */
        usys_log_error("Invalid bootstrap server value in '%s' (must be strict hostname [A-Za-z0-9.-], no whitespace/control chars): '%s'",
                       DEF_BOOTSTRAP_FILE, line);
        return false;
    }

    char *dup = strdup(line);
    if (!dup) {
        usys_log_error("Memory allocation failed duplicating bootstrap server");
        return false;
    }

    *out = dup;
    return true;
}

int process_config_file(char *fileName, Config *config) {

	FILE *fp;
	toml_table_t *fileData, *configData;
	char errBuf[MAX_BUFFER];

	if ((fp = fopen(fileName, "r")) == NULL) {
		usys_log_error("Error opening config file: %s: %s", fileName,
				  strerror(errno));
		return FALSE;
	}

	/* Prase the TOML file entries. */
	fileData = toml_parse_file(fp, errBuf, sizeof(errBuf));
	fclose(fp);
	if (!fileData) {
		usys_log_error("Error parsing the config file %s: %s", fileName, errBuf);
		return FALSE;
	}

	/* Parse the config. */
	configData = toml_table_in(fileData, CONFIG);

	if (configData == NULL) {
		usys_log_error("[Config] section parsing error in file: %s", fileName);
		toml_free(fileData);
		return FALSE;
	}

	if (!parse_config(config, configData)) {
        usys_log_error("Unable to prase config file: %s", fileName);
        toml_free(fileData);
        return FALSE;
    }

    if (!read_bootstrap_server_info(&config->bootstrapRemoteServer)) {
        usys_log_debug("Unable to read bootstrap server info");
        config->bootstrapRemoteServer = strdup(DEF_BOOTSTRAP_SERVER);
    }

	toml_free(fileData);
	return TRUE;
}

void print_config(Config *config) {

	if (config == NULL) return;

	if (config->nodedHost) {
		usys_log_debug("noded host: %s", config->nodedHost);
	}

    usys_log_debug("noded port: %d", config->nodedPort);

	if (config->meshConfig) {
		usys_log_debug("mesh config: %s", config->meshConfig);
	}

	if (config->remoteIPFile) {
		usys_log_debug("remote IP file: %s", config->remoteIPFile);
	}

	if (config->bootstrapRemoteServer) {
	    usys_log_debug("bootstrap remote server: %s", config->bootstrapRemoteServer);
	}

    if (config->bootstrapRemotePort) {
        usys_log_debug("bootstrap remote port: %d", config->bootstrapRemotePort);
    }
}

void clear_config(Config *config) {

	if (!config) return;

	free(config->nodedHost);
	free(config->meshConfig);
	free(config->remoteIPFile);
	free(config->bootstrapRemoteServer);
}
