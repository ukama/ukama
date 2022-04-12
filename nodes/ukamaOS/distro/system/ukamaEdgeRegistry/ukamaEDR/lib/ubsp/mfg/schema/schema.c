
/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Tool for modifying the JSON schema for the device.
 * This binary is only meant to create a dummy device config by reading a sample JSON's
 * This tool would take JSON file and new  digits of serial number (LAST FOUR digit XXXX ) as command line arguments.
 *
 * Example usage is:
 * schema -i 1101 --f mfgdata/schema/com.json --f mfgdata/schema/lte.json --f mfgdata/schema/mask.json
 * schema -i 1101 --f mfgdata/schema/rfctrl.json --f mfgdata/schema/rffe.json
 *
 */
#include "utils/cJSON.h"
#include "headers/utils/log.h"
#include "mfg/common/mfg_helper.h"

#include <getopt.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define VERSION					"0.0.1"

#define MAX_JSON_TAGS 			3

/* Schema struct */
typedef struct {
	char *name;
	char *uuid;
	char *muuid;
	char *fileName;
} UnitSchema;

/* JSON TAGS */
static char* jsontags[MAX_JSON_TAGS] = {
		"unit_info",
		"unit_config",
		"module_info"
};

UnitSchema unitschema[MAX_BOARDS] = {'\0'};

/* Set the verbosity level for logs. */
void set_log_level(char *slevel) {
	int ilevel = LOG_TRACE;
	if (!strcmp(slevel, "TRACE")) {
		ilevel = LOG_TRACE;
	} else if (!strcmp(slevel, "DEBUG")) {
		ilevel = LOG_DEBUG;
	} else if (!strcmp(slevel, "INFO")) {
		ilevel = LOG_INFO;
	}
	log_set_level(ilevel);
}

/* Write to JSON file */
int write_file(char *filename, char *out)
{
	int ret = 0;
	FILE *fp = NULL;

	/* Open file */
	fp = fopen(filename,"w");
	if(fp == NULL)
	{
		fprintf(stderr,"open file failed\n");
		exit(-1);
	}

	/* Write to a file */
	ret= fputs(out, fp);

	/* Close file */
	if(fp != NULL)
		fclose(fp);

	return ret;
}

/* Read a file */
int read_file(char *filename, char **data) {
	FILE *f;
	size_t len = 0;
	int ret = 0;

	/* Open file */
	f=fopen(filename,"rb");

	/* Length of file */
	fseek(f,0,SEEK_END);
	len=ftell(f);
	fseek(f,0,SEEK_SET);

	/* Memory allocation*/
	char *fdata=(char*)malloc(len+1);
	if (fdata) {
		ret = fread(fdata,1,len,f);
		fdata[len]='\0';
		*data = fdata;
	}

	/* Close file */
	fclose(f);

	return ret;
}

/* Perform the parsing operation on the JSON file */
cJSON *dofile(char *filename)
{
	cJSON *json,*ret;
	char *data = NULL;

	/* Read file */
	if (read_file( filename, &data)  <= 0) {
		log_error("Schema:: Error:: Failed to read file %s.",filename);
		ret =  NULL;
		goto EXIT;
	}

	/* Parse file */
	json=cJSON_Parse(data);
	if (!json)
	{
		log_error("Schema:: Error before: [%s]\n",cJSON_GetErrorPtr());
		ret = NULL;
		goto EXIT;
	}
	else
	{
		ret = json;
	}

	/* Cleanup */
	EXIT:
	if (data) {
		free(data);
	}
	return ret;
}



/* Modify UUID field in JSON */
int modify_uuid(cJSON * obj, char* value) {
	int ret = 0;

	/* Get UUID value */
	cJSON *juuid = cJSON_GetObjectItem(obj,"UUID");
	if (juuid) {

		/* Assumption is we have some UUID present in file. Read current UUID string */
		char *currval = cJSON_GetStringValue(juuid);
		if (currval) {

			/* Update the UUID */
			if (!cJSON_SetValuestring(juuid, value)) {
				log_error("Schema:: Error:: Setting UUID to %s failed.", value);
				ret = -1;
			} else {
				log_debug("Schema:: UUID %s is updated to %s.", currval, value);
			}

		}else {
			log_debug("Schema:: UUID set to %s.", value);
		}

	} else {
		log_error("Schema:: Error:: UUID tag not found.");
		ret = -1;
	}

	return ret;
}

/* Search unit name and return the UUID for it */
char* read_uuid_for_module_name(char* name) {
	for(unsigned short int idx = 0; idx < MAX_BOARDS; idx++) {
		if ( (unitschema[idx].name) && (!strcasecmp(name, unitschema[idx].name))) {
			return unitschema[idx].muuid;
		}
	}
	return NULL;
}

/* Update Unit config */
int  update_unit_config(cJSON **obj) {
	int ret = -1;
	cJSON *unitcfgobj = *obj;

	/* unit config is supposed to be an array of modules */
	if (cJSON_IsArray(unitcfgobj)) {
		const cJSON *module = NULL;
		cJSON_ArrayForEach(module, unitcfgobj)
		{
			ret = -1;
			cJSON *jname = cJSON_GetObjectItem(module, "name");
			if (jname) {
				char *currval = cJSON_GetStringValue(jname);
				if (currval) {
					char* muuid = read_uuid_for_module_name(currval);
					if (muuid) {
						if (modify_uuid(module, muuid) == 0) {
							ret = 0;
						}
					}
				}
			}

			/* ret should be 0 after every loop if not something is wrong with schema */
			if(ret) {
				break;
			}

		}

	} else {
		log_error("Schema:: Expected unit_config array but unknown JSON tag found.");
	}
	return ret;
}

/* Modify JSON */
int modify_json(unsigned short idx)
{
	int ret = 0;
	cJSON *root,*obj;
	char *out;
	char* value = NULL;
	/* Parse the JSON file */
	root = dofile(unitschema[idx].fileName);
	if (!root) {
		return -1;
	}

	/* Debug Info */
	out = cJSON_Print(root);
	if (out) {
		log_trace("Before modification file %s is::\n %s\n", unitschema[idx].fileName, out);
		free(out);
	}

	/* Update all available tags in JSON */
	for (int tag = 0; tag < MAX_JSON_TAGS; tag++) {

		obj = cJSON_GetObjectItem(root, jsontags[tag]);
		if (obj) {

			/* For Unit Config which is array */
			if ( unitschema[idx].muuid && (!strcmp(jsontags[tag], "unit_config"))) {

				/* Update Unit Config */
				ret  = update_unit_config(&obj);
				if (ret) {
					log_error("Schema:: Failed to update unit config for %s file.", unitschema[idx].fileName);
					return ret;
				}

			} else {

				if (unitschema[idx].muuid && (!strcmp(jsontags[tag], "module_info"))) {
					/* Module Info */
					value = unitschema[idx].muuid;

				} else if (unitschema[idx].uuid && (!strcmp(jsontags[tag], "unit_info"))) {
					/* Unit Info */
					value = unitschema[idx].uuid;
				}

				/* For Unit Info and Module info  */
				if (modify_uuid(obj, value )) {
					log_error("Schema copying uuid failed for Unit/Module Info.");
					return -1;
				}
			}

		}
	}

	/* Debug Info */
	out = cJSON_Print(root);
	log_trace("After modification file %s is::\n %s\n",unitschema[idx].fileName, out);

	/* Update the JSON file */
	if(write_file(unitschema[idx].fileName,out) > 0 ) {
		log_info("File %s updated successfully.", unitschema[idx].fileName );
		free(out);
	} else {
		log_error("Write to file %s failed.", unitschema[idx].fileName );
	}

	/* Clean the cJSON root  */
	cJSON_Delete(root);

	return 0;
}

/* Command line args */
static struct option long_options[] = {
		{ "name", required_argument, 0, 'n' },
		{ "uuid", required_argument, 0, 'u' },
		{ "muuid", required_argument, 0, 'm' },
		{ "file", required_argument, 0, 'f' },
		{ "logs", required_argument, 0, 'l' },
		{ "help", no_argument, 0, 'h' },
		{ "version", no_argument, 0, 'v' },
		{ 0, 0, 0, 0 }
};

/* Usage options for the ukamaEDR */
void usage() {
	printf("Usage: schema [options] \n");
	printf("Options:\n");
	printf(
			"--h, --help                                                          Help menu.\n");
	printf(
			"--l, --logs <TRACE>|<DEBUG>|<INFO>                                   Log level for the process.\n");
	printf(
			"--n, --name <ComV1>|<LTE>|<MASK>|<RF CTRL BOARD>,<RF BOARD>          Name of module.\n");
	printf(
			"--u, --uuid <string 24 character long>                               UUIID for file.\n");
	printf(
			"--m, --muuid <string 24 character long>                              Module UIID for file.\n");
	printf(
			"--f, --file <Files>                                                  Schema files.\n");
	printf(
			"--v, --version                                                       Software Version.\n");
}


/* JSON Schema UUID Update utility */
int main(int argc, char** argv) {
	char *uuid = {"\0"};
	char *name[MAX_BOARDS] = {"\0"};
	char *mid[MAX_BOARDS] = {"\0"};
	char *file[MAX_BOARDS] = {"\0"};
	char *debug = "TRACE";
	char *ip = "";
	set_log_level(debug);

	if (argc < 2 ) {
		log_error("Not enough arguments.");
		usage();
		exit(1);
	}


	int fidx = 0;
	int midx = 0;
	int nidx = 0;
	int uuidcount = 0;

	/* Parsing command line args. */
	while (true) {
		int opt = 0;
		int opdidx = 0;

		opt = getopt_long(argc, argv, "h:v:u:f:l:m:n:", long_options, &opdidx);
		if (opt == -1) {
			break;
		}

		switch (opt) {
		case 'h':
		usage();
		exit(0);
		break;

		case 'v':
			puts(VERSION);
			break;

		case 'n':
			name[nidx] = optarg;
			nidx++;
			break;

		case 'u':
			uuid = optarg;
			uuidcount++;
			break;

		case 'f':
			file[fidx] = optarg;
			fidx++;
			break;

		case 'l':
			debug = optarg;
			set_log_level(debug);
			break;

		case 'm':
			mid[midx] = optarg;
			midx++;
			break;

		default:
			usage();
			exit(0);
		}
	}

	/* Check for arguments */
	if (uuidcount != 1) {

		log_error("Schema:: Error:: Schema expects one uuid argument which is must and you provided %d.", uuidcount);
		usage();
		exit(0);

	} else if ((fidx != nidx) || (fidx != midx) || (fidx < 1)) {

		log_error("Schema:: Error:: Schema expects module uuid, name and file for each module "
				"and you provided %d module name %d modules uuid and %d module files.", nidx, midx, fidx);
		usage();
		exit(0);

	}

	/* Verify UUID */
	if (verify_uuid(uuid)) {
		log_error("Schema:: Error:: Check the Unit UUID %s", uuid);
		usage();
		exit(0);
	}

	/* Update Unit Schema */
	for(int idx = 0; idx < fidx;idx++) {
		log_trace("Files[%d] = %s Module UUID %s\n", idx, file[idx], mid[idx]);

		/* Verify module uuid and name */
		if (verify_uuid(mid[idx]) || verify_boardname(name[idx]) ) {
			usage();
			exit(0);
		}

		unitschema[idx].name = name[idx];
		if (idx==0) {
			unitschema[idx].uuid = uuid;
		} else {
			unitschema[idx].uuid = NULL;
		}
		unitschema[idx].muuid = mid[idx];
		unitschema[idx].fileName = file[idx];
	}

	/* Modify every file provided in input.*/
	for(int idx = 0; idx < fidx;idx++) {

		/* Update JSON */
		int ret = modify_json(idx);
		if (ret) {
			log_error("Schema:: Error:: Failed to update schema %s.", file[idx]);
			exit(0);
		} else {
			log_info("Schema:: Updated schema %s.", file[idx]);
		}

	}

	return 0;
}
