/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/* Mocking cloud-based service provider for testing only. */

#include <stdio.h>
#include <stdlib.h>
#include <sqlite3.h>
#include <string.h>
#include <jansson.h>
#include <ulfius.h>

#define FALSE   0
#define TRUE    1
#define MAX_ENT 128

#define EP_CONTAINERS "/v1/capps/"

typedef struct _u_request  req_t;
typedef struct _u_response resp_t;

typedef struct {
    char *name;
	char *version;
	char *type;
	char *url;
	char *created_at;
	int  size_bytes;
	char *chunk_url;
	char *chunks_url;
} dbEntry;

typedef struct {

	sqlite3 *dbPtr;
	int     numEnt;
	dbEntry ent[MAX_ENT];
} DB;

int table_exists(sqlite3 *db) {

	int ret;
	char *sql = NULL, *errMsg = NULL;

	sql = "SELECT 1 FROM sqlite_master where type='table' and name='Containers'";
	ret = sqlite3_exec(db, sql, NULL, NULL, &errMsg);

	if (ret != SQLITE_OK) {
		fprintf(stdout, "Table does not exist. Query returned error: %s\n", errMsg);
		sqlite3_free(errMsg);
		return FALSE;
	} else {
		return TRUE;
	}
}

static sqlite3 *open_db(char *dbFile) {

	sqlite3 *db = NULL;
	int ret;

	ret = sqlite3_open(dbFile, &db);

	if (ret) {
		fprintf(stderr, "Error opening the dbFile: %s\n", dbFile);
		exit(1);
	} else {
		fprintf(stdout, "db opened: %s\n", dbFile);
	}

	if (!table_exists(db)) {
		sqlite3_close(db);
		return NULL;
	} else {
		fprintf(stdout, "Found right table in db: %s\n", dbFile);
	}

	return db;
}

int callback_get_containers(req_t *request, resp_t *response, void *user_data) {

	int resCode = 200, i, j;
	char *name = NULL, *str = NULL, *response_body = NULL;
	DB *db;
	json_t *json, *artifacts, *artifact, *formats, *format;

	json = json_object();
	db = (DB *)user_data;

	name = (char *)u_map_get(request->map_url, "name");

	if (!name) {
		fprintf(stderr, "Invalid name in GET response for EP: %s.\n", EP_CONTAINERS);
		response_body = strdup("Invalid container name.");
		resCode = 400;
		goto reply;
	}

	fprintf(stdout, "Valid GET request for %s name:%s\n", EP_CONTAINERS, name);

	json_object_set_new(json, "name", json_string(name));
	json_object_set_new(json, "artifacts", json_array());
	artifacts = json_object_get(json, "artifacts");

	for (i = 0; i < db->numEnt; i++) {
		if (strcmp(db->ent[i].name, name) == 0) {

			int artifact_found = 0;
			for (j = 0; j < json_array_size(artifacts); j++) {
				artifact = json_array_get(artifacts, j);
				if (strcmp(json_string_value(json_object_get(artifact, "version")),
                           db->ent[i].version) == 0) {
					formats = json_object_get(artifact, "formats");
					artifact_found = 1;
					break;
				}
			}

			if (!artifact_found) {

				artifact = json_object();
				json_object_set_new(artifact, "version",
                                    json_string(db->ent[i].version));
				json_object_set_new(artifact, "formats",
                                    json_array());
				formats = json_object_get(artifact, "formats");
				json_array_append_new(artifacts, artifact);
			}

			format = json_object();
			json_object_set_new(format, "type", json_string(db->ent[i].type));
			json_object_set_new(format, "url", json_string(db->ent[i].url));
			json_object_set_new(format, "created_at", json_string(db->ent[i].created_at));

			if (db->ent[i].size_bytes) {
				json_object_set_new(format, "size_bytes",
                                    json_integer(db->ent[i].size_bytes));
			}

			if (strcmp(db->ent[i].type, "chunk") == 0) {
				json_t *extra_info = json_object();
				json_object_set_new(extra_info, "chunks",
                                    json_string(db->ent[i].chunks_url));
				json_object_set_new(format, "extra_info", extra_info);
			}
			json_array_append_new(formats, format);
		}
	}

	str = json_dumps(json, 0);
	fprintf(stdout, "JSON Object: %s\n", str);
    free(str);

reply:
	if (resCode == 200) {
		ulfius_set_json_body_response(response, resCode, json);
	} else if (resCode == 400) {
		ulfius_set_string_body_response(response, resCode, response_body);
	}

	json_decref(json);
	free(response_body);

	return U_CALLBACK_CONTINUE;
}

int read_entries(void *arg, int argc, char **argv, char **colName) {

	int i, ent;
	DB *db = (DB *)arg;

	if (db == NULL) {
		fprintf(stderr, "Memory failure\n");
		exit(1);
	}

	ent = db->numEnt;

	for (i = 0; i < argc; i++) {
		if (strcmp(colName[i], "Name") == 0) {
			db->ent[ent].name = strdup(argv[i]);
		} else if (strcmp(colName[i], "Version") == 0) {
			db->ent[ent].version = strdup(argv[i]);
		} else if (strcmp(colName[i], "Type") == 0) {
			db->ent[ent].type = strdup(argv[i]);
		} else if (strcmp(colName[i], "URL") == 0) {
			db->ent[ent].url = strdup(argv[i]);
		} else if (strcmp(colName[i], "CreatedAt") == 0) {
			db->ent[ent].created_at = strdup(argv[i]);
		} else if (strcmp(colName[i], "SizeBytes") == 0 && argv[i]) {
			db->ent[ent].size_bytes = atoi(argv[i]);
		} else if (strcmp(colName[i], "ChunkURL") == 0) {
			db->ent[ent].chunk_url = strdup(argv[i]);
		} else if (strcmp(colName[i], "ChunksURL") == 0) {
			db->ent[ent].chunks_url = strdup(argv[i]);
		}
	}

	if (db->ent[db->numEnt].name)
		db->numEnt++;

	return 0;
}

int read_all_db_entries(DB *db) {

	int val = FALSE;
	char buf[128];
	char *err = NULL;

	if (db == NULL || db->dbPtr == NULL) {
		return FALSE;
	}

	sprintf(buf, "SELECT * FROM Containers;");

	val = sqlite3_exec(db->dbPtr, buf, read_entries, db, &err);

	if (val != SQLITE_OK) {
		fprintf(stderr, "SQL read error, query failure: %s\n", err);
		sqlite3_free(err);
        return FALSE;
	} else {
		fprintf(stdout, " Query: %s\n Response ok\n", buf);
		return TRUE;
	}

	return FALSE;
}

void free_db_entries(DB *db) {

	int i;

	for (i = 0; i < db->numEnt; i++) {
		free(db->ent[i].name);
		free(db->ent[i].version);
		free(db->ent[i].type);
		free(db->ent[i].url);
		free(db->ent[i].created_at);
		free(db->ent[i].chunk_url);
		free(db->ent[i].chunks_url);
	}

	free(db);
}

int main(int argc, char **argv) {

	int port;
	struct _u_instance inst;
	DB *db = NULL;

	if (argc < 2) {
		fprintf(stderr, "USAGE: %s port dbFile\n", argv[0]);
		return 0;
	}

	db = (DB *)calloc(sizeof(DB), 1);
	if (!db) {
		fprintf(stderr, "Memory allocation issue of size: %ld\n", sizeof(DB));
		exit(1);
	}

	port      = atoi(argv[1]);
	db->dbPtr = open_db(argv[2]);

	if (!db->dbPtr) {
		fprintf(stderr, "Error opening db file: %s\n", argv[2]);
		return 1;
	}

	if (ulfius_init_instance(&inst, port, NULL, NULL) != U_OK) {
		fprintf(stderr, "Error ulfius_init_instance, abort\n");
		return 1;
	}

	ulfius_add_endpoint_by_val(&inst, "GET", "/v1/hub/apps/", ":name", 0,
                               &callback_get_containers, db);

	if (read_all_db_entries(db) == FALSE) {
        fprintf(stderr, "Error reading the db\n");
        return 1;
	}

	if (ulfius_start_framework(&inst) == U_OK) {
		fprintf(stdout, "Framework started on port %d\n", inst.port);
		getchar();
	} else {
		fprintf(stderr, "Error starting framework\n");
	}

	fprintf(stdout, "End framework\n");

	ulfius_stop_framework(&inst);
	ulfius_clean_instance(&inst);

	free_db_entries(db);

	return 0;
}





