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
 * genSchema --n ComV1 --u UK-7001-HNODE-SA03-1102 --m UK-7001-COM-1102 --f mfgdata/schema/com.json --m UK-7001-LTE-1102 --n LTE --f mfgdata/schema/lte.json --n MASK --m UK-7001-MSK-1102 --f mfgdata/schema/mask.json
 * schema -i 1101 --f mfgdata/schema/rfctrl.json --f mfgdata/schema/rffe.json
 *
 */

#include "device.h"
#include "errorcode.h"
#include "inventory.h"
#include "json_types.h"
#include "ledger.h"
#include "noded_macros.h"
#include "property.h"

#include "utils/crc32.h"
#include "utils/mfg_helper.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_getopt.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

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
static char* jsonKeyTag[MAX_JSON_TAGS] = {
                JTAG_UNIT_INFO,
                JTAG_UNIT_CONFIG,
                JTAG_MODULE_INFO
};

UnitSchema unitSchema[MAX_BOARDS] = {'\0'};

/* Set the verbosity level for logs. */
void set_log_level(char *slevel) {
    int ilevel = LOG_TRACE;
    if (!usys_strcmp(slevel, "TRACE")) {
        ilevel = LOG_TRACE;
    } else if (!usys_strcmp(slevel, "DEBUG")) {
        ilevel = LOG_DEBUG;
    } else if (!usys_strcmp(slevel, "INFO")) {
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
    fp = usys_fopen(filename,"w");
    if(fp == NULL)
    {
        usys_fprintf(stderr,"open file failed\n");
        usys_exit(-1);
    }

    /* Write to a file */
    ret= usys_fputs(out, fp);

    /* Close file */
    if(fp != NULL)
        usys_fclose(fp);

    return ret;
}

/* Read a file */
int read_file(char *filename, char **data) {
    FILE *f;
    size_t len = 0;
    int ret = 0;

    /* Open file */
    f=usys_fopen(filename,"rb");

    /* Length of file */
    usys_fseek(f,0,SEEK_END);
    len=usys_ftell(f);
    usys_fseek(f,0,SEEK_SET);

    /* Memory allocation*/
    char *fdata=(char*)usys_zmalloc(len+1);
    if (fdata) {
        ret = usys_fread(fdata,1,len,f);
        fdata[len]='\0';
        *data = fdata;
    }

    /* Close file */
    usys_fclose(f);

    return ret;
}

/* Perform the parsing operation on the JSON file */
JsonObj* dofile(char *filename) {
    JsonObj *jSchema = NULL;
    JsonErrObj *jErr = NULL;
    char *data = NULL;

    /* Read file */
    if (read_file( filename, &data)  <= 0) {
        usys_log_error("Schema:: Error:: Failed to read file %s.",filename);
        goto EXIT;
    }

    /* Parse file */
    jSchema = json_loads(data, JSON_DECODE_ANY, jErr);
    if (!jSchema) {
        parser_error(jErr, "Failed to parse schema");
        goto EXIT;
    }

    /* Cleanup */
    EXIT:
    if (data) {
        usys_free(data);
    }
    return jSchema;
}



/* Modify UUID field in JSON */
int modify_uuid(JsonObj * jObj, char* value) {
    int ret = -1;

    /* Get UUID value */
    char *uuid = NULL;
    if(!parser_read_string_object((const JsonObj*)jObj, JTAG_UUID, &uuid)) {
        usys_log_error("Schema:: Error:: UUID tag not found.");
        ret = -1;
    } else {
        /* Updating new value */
        ret = json_object_set_new(jObj, JTAG_UUID, json_string(value));
        if(ret) {
            usys_log_error("Schema:: Error:: Setting UUID to %s failed.", value);
            ret = -1;
        } else {
            log_debug("Schema:: UUID %s is updated to %s.", uuid, value);
            usys_free(uuid);
            ret = 0;
        }
    }

    return ret;
}

/* Search unit name and return the UUID for it */
char* read_uuid_for_module_name(const char* name) {
    for(unsigned short int idx = 0; idx < MAX_BOARDS; idx++) {
        if ( (unitSchema[idx].name) && (!strcasecmp(name, unitSchema[idx].name))) {
            return unitSchema[idx].muuid;
        }
    }
    return NULL;
}

/* Update Unit config */
int  update_unit_config( const JsonObj **obj) {
    int ret = -1;
   const JsonObj *unitCfgObj = *obj;

    /* unit config is supposed to be an array of modules */
    if (json_is_array(unitCfgObj)) {
        JsonObj *module = NULL;
        int iter = 0;
        json_array_foreach(unitCfgObj, iter, module) {
            ret = -1;
            JsonObj *jName = json_object_get(module, JTAG_NAME);
            if (jName) {
                const char *currVal = json_string_value(jName);
                if (currVal) {
                    char* mUuid = read_uuid_for_module_name(currVal);
                    if (mUuid) {
                        if (modify_uuid(module, mUuid) == 0) {
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
        usys_log_error("Schema:: Expected unit_config array but unknown JSON tag found.");
    }
    return ret;
}

/* Modify JSON */
int modify_json(unsigned short idx)
{
    int ret = 0;
    JsonObj *root = NULL;
    JsonObj *obj = NULL;
    char *out;
    char* value = NULL;
    /* Parse the JSON file */
    root = dofile(unitSchema[idx].fileName);
    if (!root) {
        return -1;
    }

    /* Debug Info */
    out = json_dumps(root, (JSON_INDENT(4)|JSON_COMPACT|JSON_ENCODE_ANY) );
    if (out) {
        usys_log_trace("Before modification file %s is::\n %s\n", unitSchema[idx].fileName, out);
        usys_free(out);
    }

    /* Update all available tags in JSON */
    for (int tag = 0; tag < MAX_JSON_TAGS; tag++) {

        obj = json_object_get(root, jsonKeyTag[tag]);
        if (obj) {

            /* For Unit Config which is array */
            if ( unitSchema[idx].muuid && (!usys_strcmp(jsonKeyTag[tag], JTAG_UNIT_CONFIG))) {

                /* Update Unit Config */
                ret  = update_unit_config((const JsonObj**)&obj);
                if (ret) {
                    usys_log_error("Schema:: Failed to update unit config for %s file.", unitSchema[idx].fileName);
                    return ret;
                }

            } else {

                if (unitSchema[idx].muuid && (!usys_strcmp(jsonKeyTag[tag], JTAG_MODULE_INFO))) {
                    /* Module Info */
                    value = unitSchema[idx].muuid;

                } else if (unitSchema[idx].uuid && (!usys_strcmp(jsonKeyTag[tag], JTAG_UNIT_INFO))) {
                    /* Unit Info */
                    value = unitSchema[idx].uuid;
                }

                /* For Unit Info and Module info  */
                if (modify_uuid(obj, value )) {
                    usys_log_error("Schema copying uuid failed for Unit/Module Info.");
                    return -1;
                }
            }

        }
    }

    /* Debug Info */
    out = json_dumps(root, (JSON_INDENT(4)|JSON_COMPACT|JSON_ENCODE_ANY) );
    usys_log_trace("After modification file %s is::\n %s\n",unitSchema[idx].fileName, out);

    /* Update the JSON file */
    if(write_file(unitSchema[idx].fileName,out) > 0 ) {
        usys_log_info("File %s updated successfully.", unitSchema[idx].fileName );
        usys_free(out);
    } else {
        usys_log_error("Write to file %s failed.", unitSchema[idx].fileName );
    }

    /* Clean the cJSON root  */
    json_decref(root);

    return 0;
}

/* Command line args */
static struct option longOptions[] = {
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
        usys_log_error("Not enough arguments.");
        usage();
        usys_exit(1);
    }


    int fidx = 0;
    int midx = 0;
    int nidx = 0;
    int uuidCount = 0;

    /* Parsing command line args. */
    while (true) {
        int opt = 0;
        int opdIdx = 0;

        opt = usys_getopt_long(argc, argv, "h:v:u:f:l:m:n:", longOptions, &opdIdx);
        if (opt == -1) {
            break;
        }

        switch (opt) {
            case 'h':
                usage();
                usys_exit(0);
                break;

            case 'v':
                usys_puts(VERSION);
                break;

            case 'n':
                name[nidx] = optarg;
                nidx++;
                break;

            case 'u':
                uuid = optarg;
                uuidCount++;
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
                usys_exit(0);
        }
    }

    /* Check for arguments */
    if (uuidCount != 1) {

        usys_log_error("Schema:: Error:: Schema expects one uuid argument which is must and you provided %d.", uuidCount);
        usage();
        usys_exit(0);

    } else if ((fidx != nidx) || (fidx != midx) || (fidx < 1)) {

        usys_log_error("Schema:: Error:: Schema expects module uuid, name and file for each module "
                        "and you provided %d module name %d modules uuid and %d module files.", nidx, midx, fidx);
        usage();
        usys_exit(0);

    }

    /* Verify UUID */
    if (verify_uuid(uuid)) {
        usys_log_error("Schema:: Error:: Check the Unit UUID %s", uuid);
        usage();
        exit(0);
    }

    /* Update Unit Schema */
    for(int idx = 0; idx < fidx;idx++) {
        usys_log_trace("Files[%d] = %s Module UUID %s\n", idx, file[idx], mid[idx]);

        /* Verify module uuid and name */
        if (verify_uuid(mid[idx]) || verify_board_name(name[idx]) ) {
            usage();
            exit(0);
        }

        unitSchema[idx].name = name[idx];
        if (idx==0) {
            unitSchema[idx].uuid = uuid;
        } else {
            unitSchema[idx].uuid = NULL;
        }
        unitSchema[idx].muuid = mid[idx];
        unitSchema[idx].fileName = file[idx];
    }

    /* Modify every file provided in input.*/
    for(int idx = 0; idx < fidx;idx++) {

        /* Update JSON */
        int ret = modify_json(idx);
        if (ret) {
            usys_log_error("Schema:: Error:: Failed to update schema %s.", file[idx]);
            usys_exit(0);
        } else {
            usys_log_info("Schema:: Updated schema %s.", file[idx]);
        }

    }

    return 0;
}
