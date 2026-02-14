/* 
 * - If FEMD_SYSROOT is set, it rewrites any JSON path fields:
 *      invtSysFsFile, devSysFsFile
 *   so they point to the actual filesystem under FEMD_SYSROOT.
 *
 * Example:
 *   FEMD_SYSROOT=/tmp/sys
 *   "/tmp/sys/class/hwmon/hwmon0/temp1_input" stays as-is
 *   "/sys/class/hwmon/..." -> "/tmp/sys/class/hwmon/..."
 *   "/class/hwmon/..."     -> "/tmp/sys/class/hwmon/..."
 *   "/dev/i2c-1"           -> "/tmp/sys/dev/i2c-1"
 *
 * If FEMD_SYSROOT is not set, behavior is unchanged.
 */

/*
 * Tool for modifying the JSON schema for the device.
 * This binary is only meant to create a dummy device config by reading a
 * sample JSON's
 *
 * Example usage is:
 *
 * genSchema --u UK-7001-HNODE-SA03-1102 \
 * --n ComV1 --m UK-7001-COM-1102  --f mfgdata/schema/com.json \
 * --n LTE --m UK-7001-TRX-1102  --f mfgdata/schema/lte.json \
 * --n MASK --m UK-7001-MSK-1102 --f mfgdata/schema/mask.json
 *
 * ANode
 * genSchema --u UK-8001-ANODE-SA03-1102 \
 * --n "RF CTRL BOARD" --m UK-8001-RFC-1102 --f mfgdata/schema/rfctrl.json \
 * --n "RF BOARD" --m UK-8001-RFA-1102 --f mfgdata/schema/rffe.json
 *
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

#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

#define VERSION                 "0.0.1"
#define MAX_JSON_TAGS           3

#define ENV_FEMD_SYSROOT        "FEMD_SYSROOT"

/* Schema struct */
typedef struct {
    char *name;
    char *uuid;
    char *muuid;
    char *fileName;
} NodeSchema;

/* JSON TAGS */
static char* jsonKeyTag[MAX_JSON_TAGS] = {
                JTAG_NODE_INFO,
                JTAG_NODE_CONFIG,
                JTAG_MODULE_INFO
};

NodeSchema nodeSchema[MAX_BOARDS] = {'\0'};

static const char *femd_sysroot(void) {
    const char *v = getenv(ENV_FEMD_SYSROOT);
    return (v && v[0] != '\0') ? v : NULL;
}

static bool starts_with(const char *s, const char *pfx) {
    if (!s || !pfx) return false;
    size_t n = strlen(pfx);
    return strncmp(s, pfx, n) == 0;
}

static int resolve_with_sysroot(char *out, size_t outsz, const char *path) {
    const char *root = femd_sysroot();
    if (!out || outsz == 0 || !path || path[0] == '\0') return -1;

    if (!root) {
        int n = snprintf(out, outsz, "%s", path);
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }

    /* Already absolute under sysroot */
    if (starts_with(path, root)) {
        int n = snprintf(out, outsz, "%s", path);
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }

    /* /sys/... -> <root>/... */
    if (starts_with(path, "/sys/")) {
        int n = snprintf(out, outsz, "%s%s", root, path + 4);
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }

    /* /dev/... -> <root>/dev/... */
    if (starts_with(path, "/dev/")) {
        int n = snprintf(out, outsz, "%s%s", root, path);
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }

    /* "/class/...", "/bus/...", "/devices/..." etc -> <root> + path */
    if (path[0] == '/') {
        int n = snprintf(out, outsz, "%s%s", root, path);
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }

    /* Relative path: leave unchanged */
    {
        int n = snprintf(out, outsz, "%s", path);
        return (n < 0 || (size_t)n >= outsz) ? -1 : 0;
    }
}

static bool is_path_key(const char *k) {
    return (k &&
            (!strcmp(k, "devSysFsFile") ||
             !strcmp(k, "invtSysFsFile")));
}

static void json_apply_sysroot_paths(JsonObj *j) {
    if (!j || !femd_sysroot()) return;

    if (json_is_object(j)) {
        const char *k = NULL;
        JsonObj *v = NULL;

        json_object_foreach(j, k, v) {

            if (is_path_key(k) && json_is_string(v)) {
                const char *oldp = json_string_value(v);
                if (oldp && oldp[0] != '\0') {
                    char newp[512] = {0};
                    if (resolve_with_sysroot(newp, sizeof(newp), oldp) == 0) {
                        if (strcmp(oldp, newp) != 0) {
                            json_object_set_new(j, k, json_string(newp));
                        }
                    }
                }
                continue;
            }

            if (json_is_object(v) || json_is_array(v)) {
                json_apply_sysroot_paths(v);
            }
        }
        return;
    }

    if (json_is_array(j)) {
        size_t i = 0;
        JsonObj *v = NULL;
        json_array_foreach(j, i, v) {
            if (json_is_object(v) || json_is_array(v)) {
                json_apply_sysroot_paths(v);
            }
        }
        return;
    }
}

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

/* Search node name and return the UUID for it */
char* read_uuid_for_module_name(const char* name) {
    for(unsigned short int idx = 0; idx < MAX_BOARDS; idx++) {
        if ( (nodeSchema[idx].name) && (!strcasecmp(name, nodeSchema[idx].name))) {
            return nodeSchema[idx].muuid;
        }
    }
    return NULL;
}

/* Update Node Config */
int  update_node_config( const JsonObj **obj) {
    int ret = -1;
    const JsonObj *nodeCfgObj = *obj;

    /* node config is supposed to be an array of modules */
    if (json_is_array(nodeCfgObj)) {
        JsonObj *module = NULL;
        int iter = 0;
        json_array_foreach(nodeCfgObj, iter, module) {
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
        usys_log_error("Schema:: Expected node_config array but unknown JSON tag found.");
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
    root = dofile(nodeSchema[idx].fileName);
    if (!root) {
        return -1;
    }

    /* Debug Info */
    out = json_dumps(root, (JSON_INDENT(4)|JSON_COMPACT|JSON_ENCODE_ANY) );
    if (out) {
        usys_log_trace("Before modification file %s is::\n %s\n", nodeSchema[idx].fileName, out);
        usys_free(out);
    }

    /* Update all available tags in JSON */
    for (int tag = 0; tag < MAX_JSON_TAGS; tag++) {

        obj = json_object_get(root, jsonKeyTag[tag]);
        if (obj) {

            /* For  Node Config which is array */
            if ( nodeSchema[idx].muuid && (!usys_strcmp(jsonKeyTag[tag], JTAG_NODE_CONFIG))) {

                /* Update  Node Config */
                ret  = update_node_config((const JsonObj**)&obj);
                if (ret) {
                    usys_log_error("Schema:: Failed to update node config for %s file.",
                                   nodeSchema[idx].fileName);
                    json_decref(root);
                    return ret;
                }

            } else {

                if (nodeSchema[idx].muuid && (!usys_strcmp(jsonKeyTag[tag], JTAG_MODULE_INFO))) {
                    /* Module Info */
                    value = nodeSchema[idx].muuid;

                } else if (nodeSchema[idx].uuid && (!usys_strcmp(jsonKeyTag[tag], JTAG_NODE_INFO))) {
                    /* Node Info */
                    value = nodeSchema[idx].uuid;
                }

                /* For Node Info and Module info  */
                if (modify_uuid(obj, value )) {
                    usys_log_error("Schema copying uuid failed for Unit/Module Info.");
                    json_decref(root);
                    return -1;
                }
            }

        }
    }

    /* Apply FEMD_SYSROOT path rewrite (devSysFsFile / invtSysFsFile) */
    json_apply_sysroot_paths(root);

    /* Debug Info */
    out = json_dumps(root, (JSON_INDENT(4)|JSON_COMPACT|JSON_ENCODE_ANY) );
    usys_log_trace("After modification file %s is::\n %s\n",nodeSchema[idx].fileName, out);

    /* Update the JSON file */
    if(write_file(nodeSchema[idx].fileName,out) > 0 ) {
        usys_log_info("File %s updated successfully.", nodeSchema[idx].fileName );
        usys_free(out);
    } else {
        usys_log_error("Write to file %s failed.", nodeSchema[idx].fileName );
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

/* Usage options */
void usage() {
    printf("Usage: schema [options] \n");
    printf("Options:\n");
    printf("--h, --help                                                          Help menu.\n");
    printf("--l, --logs <TRACE>|<DEBUG>|<INFO>                                   Log level for the process.\n");
    printf("--n, --name <ComV1>|<LTE>|<MASK>|<RF CTRL BOARD>,<RF BOARD>          Name of module.\n");
    printf("--u, --uuid <string 24 character long>                               UUIID for file.\n");
    printf("--m, --muuid <string 24 character long>                              Module UIID for file.\n");
    printf("--f, --file <Files>                                                  Schema files.\n");
    printf("--v, --version                                                       Software Version.\n");
}

/* JSON Schema UUID Update utility */
int main(int argc, char** argv) {
    char *uuid = {"\0"};
    char *name[MAX_BOARDS] = {"\0"};
    char *mid[MAX_BOARDS] = {"\0"};
    char *file[MAX_BOARDS] = {"\0"};
    char *debug = "TRACE";

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
                        "and you provided %d module name %d modules uuid and %d module files.",
                       nidx, midx, fidx);
        usage();
        usys_exit(0);
    }

    /* Verify UUID */
    if (verify_uuid(uuid)) {
        usys_log_error("Schema:: Error:: Check the Node UUID %s", uuid);
        usage();
        exit(0);
    }

    /* Update Node Schema */
    for(int idx = 0; idx < fidx;idx++) {
        usys_log_trace("Files[%d] = %s Module UUID %s\n", idx, file[idx], mid[idx]);

        /* Verify module uuid and name */
        if (verify_uuid(mid[idx]) || verify_board_name(name[idx]) ) {
            usage();
            exit(0);
        }

        nodeSchema[idx].name = name[idx];
        if (idx==0) {
            nodeSchema[idx].uuid = uuid;
        } else {
            nodeSchema[idx].uuid = NULL;
        }
        nodeSchema[idx].muuid = mid[idx];
        nodeSchema[idx].fileName = file[idx];
    }

    /* Modify every file provided in input.*/
    for(int idx = 0; idx < fidx;idx++) {
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
