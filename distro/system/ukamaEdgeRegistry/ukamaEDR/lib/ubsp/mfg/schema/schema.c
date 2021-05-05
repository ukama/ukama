
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

#include <getopt.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#define NUMBERS_IN_ID			4
#define NUMBER_OFFSET_IN_ID		2

#define MAX_BOARDS				5
#define VERSION					"0.0.1"
#define MAX_JSON_TAGS 			3

#define MODULE_OFFSET			12
#define UNIT_OFFSET				19

static char* jsontags[MAX_JSON_TAGS] = {
		"unit_info",
		"unit_config",
		"module_info"
};

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

	fp = fopen(filename,"w");
	if(fp == NULL)
	{
		fprintf(stderr,"open file failed\n");
		exit(-1);
	}
	//ret= fprintf(fp,"%s",out);
	ret= fputs(out, fp);

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

	return (ret+1);
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

	EXIT:
	free(data);
	return ret;
}



/* Modify UUID field in JSON */
int modify_uuid(cJSON * obj, char* value, unsigned short int offset) {
	int ret = 0;

	cJSON *juuid = cJSON_GetObjectItem(obj,"UUID");
	if (juuid) {
		char *currval = cJSON_GetStringValue(juuid);
		if (currval) {
			memcpy(currval+offset, value, NUMBERS_IN_ID);
			if (!cJSON_SetValuestring(juuid, currval)) {
				log_error("Schema:: Error:: Setting UUID to %s failed.", currval);
				ret = -1;
			} else {
				log_info("Schema:: UUID set to %s.", currval);
			}
		}else {
			log_info("Schema:: UUID set to %s.", value);
		}

	} else {
		log_error("Schema:: Error:: UUID tag not found.");
		ret = -1;
	}

	return ret;
}

/* Modify JSON */
int modify_json(char* fileName, char* value)
{
	cJSON *root,*obj;
	char *out;

	/* Parse the JSON file */
	root = dofile(fileName);
	if (!root) {
		return -1;
	}

	/* Debug Info */
	out = cJSON_Print(root);
	if (out) {
		log_debug("Before modification file %s is::\n %s\n",fileName, out);
		free(out);
	}

	for (int tag = 0; tag < MAX_JSON_TAGS; tag++) {

		obj = cJSON_GetObjectItem(root, jsontags[tag]);
		if (obj) {

			/* For Unit config which is array */
			cJSON *arrobj = NULL;
			if (!strcmp(jsontags[tag], "unit_config")) {

				/* unit config is supposed to be an array of modules */
				if (cJSON_IsArray(obj)) {

					int elm = cJSON_GetArraySize(obj);
					for (int idx = 0; idx < elm; idx++) {

						/* Update every entry of the UUID */
						arrobj = cJSON_GetArrayItem(obj, idx);
						if (modify_uuid(arrobj, value, MODULE_OFFSET)) {
							log_error("Schema copying uuid failed for Unit/Module Info.");
							return -1;
						}

					}

				} else {
					log_error("Schema:: Expected unit_config array but unknown JSON tag found.");
				}

			} else {
				/* Offset */
				unsigned short int offset = 0;
				if (!strcmp(jsontags[tag], "module_info")) {
					offset = MODULE_OFFSET;
				} else if (!strcmp(jsontags[tag], "unit_info")) {
					offset = UNIT_OFFSET;
				}

				/* For Unit Info and Module info  */
				if (modify_uuid(obj, value, offset)) {
					log_error("Schema copying uuid failed for Unit/Module Info.");
					return -1;
				}
			}

		}
	}

	/* Debug Info */
	out = cJSON_Print(root);
	log_debug("After modification file %s is::\n %s\n",fileName, out);

	/* Update the JSON file */
	if(write_file(fileName,out) > 0 ) {
		log_info("File %s updated successfully.", fileName );
		free(out);
	} else {
		log_error("Write to file %s failed.", fileName );
	}

	/* Clean the cJSON root  */
	cJSON_Delete(root);

	return 0;
}

static struct option long_options[] = {
    { "id", required_argument, 0, 'i' },
    { "file", required_argument, 0, 'f' },
    { "logs", required_argument, 0, 'l' },
    { "help", no_argument, 0, 'h' },
    { "version", no_argument, 0, 'v' },
    { 0, 0, 0, 0 }
};

/* Usage options for the ukamaEDR */
void usage() {
    printf("Usage: ukamaEDR [options] \n");
    printf("Options:\n");
    printf("--h, --help                             Help menu.\n");
    printf(
        "--l, --logs <TRACE> <DEBUG> <INFO>       Log level for the process.\n");
    printf(
        "--i, --id <Numbers in Id>                 ID for file.\n");
    printf(
        "--f, --file <Files>                  Schema files.\n");
    printf("--v, --version                       Software Version.\n");
}

int main(int argc, char** argv) {
    char *uuid = {"\0"};
    char *file[MAX_BOARDS] = {"\0"};
    char *debug = "TRACE";
    char *ip = "";
    set_log_level(debug);

    if (argc < 2 ) {
    	log_error("Not enough arguments.");
    	exit(1);
    }

    int fidx = 0;

    /* Parsing command line args. */
    while (true) {
        int opt = 0;
        int opdidx = 0;

        opt = getopt_long(argc, argv, "h:v:i:f:l:", long_options, &opdidx);
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

        case 'i':
        	uuid = optarg;
        	break;

        case 'f':
        	file[fidx] = optarg;
        	fidx++;
        	break;

        case 'l':
        	debug = optarg;
        	set_log_level(debug);
        	break;

        default:
        	usage();
        	exit(0);
        }
    }
    /* modify every file provided in input.*/
    for(int idx = 0; idx < fidx;idx++) {
    	log_debug("Files[%d] = %s ID %s Length %d \n", idx, file[idx], uuid, strlen(uuid));

    	/* Update JSON */
    	int ret = modify_json(file[idx], uuid);
    	if (ret) {
    		log_error("Schema:: Error:: Failed to update schema %s.", file[idx]);
    	} else {
    		log_info("Schema:: Updated schema %s.", file[idx]);
    	}

    }

    return 0;
}
