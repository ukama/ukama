#include "unity.h"
#include "config.h"
#include "mesh.h"
#include "toml.h"

/* defined in config.c as static */
extern int parse_config_entries(Config *config, toml_table_t *configData);
extern int read_line(char *buffer, int size, FILE *fp);
extern void split_strings(char *input, char **str1, char **str2, char *delimiter);
extern int read_nodeid(char **nodeID);
extern int read_hostname_and_nodeid(char *fileName, char **hostname, char **subnetMask);
extern int parse_config_entries(Config *config, toml_table_t *configData);

// Test for split_strings function
void test_split_strings(void) {
    char input[] = "hostname/subnet";
    char *hostname = NULL;
    char *subnet = NULL;

    split_strings(input, &hostname, &subnet, "/");

    TEST_ASSERT_NOT_NULL(hostname);
    TEST_ASSERT_NOT_NULL(subnet);
    TEST_ASSERT_EQUAL_STRING("hostname", hostname);
    TEST_ASSERT_EQUAL_STRING("subnet", subnet);

    free(hostname);
    free(subnet);
}

// Test for read_hostname_and_nodeid function
void test_read_hostname_and_nodeid(void) {
    char *hostname = NULL;
    char *subnetMask = NULL;
    int ret;

    FILE *fp = fopen(".testfile", "w+");
    fprintf(fp, "hostname/subnet;nodeid\n");
    fclose(fp);

    ret = read_hostname_and_nodeid(".testfile", &hostname, &subnetMask);

    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_EQUAL_STRING("hostname", hostname);
    TEST_ASSERT_EQUAL_STRING("subnet", subnetMask);

    free(hostname);
    free(subnetMask);
    remove(".testfile");
}

// Test for read_line function
void test_read_line(void) {
    char buffer[128];
    FILE *fp = fopen("test/testfile.txt", "w+");
    fprintf(fp, "test line\n");
    rewind(fp);

    int ret = read_line(buffer, sizeof(buffer), fp);
    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_EQUAL_STRING("test line", buffer);

    fclose(fp);
    remove("test/testfile.txt");
}

// Test for read_nodeid function
void test_read_nodeid(void) {
    char *nodeID = NULL;
    int ret;

    FILE *fp = fopen(".nodeid", "w+");
    fprintf(fp, "nodeid_value");
    fclose(fp);

    ret = read_nodeid(&nodeID);
    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_EQUAL_STRING("nodeid_value", nodeID);

    free(nodeID);
    remove(".nodeid");
}

// Test for process_config_file function
void test_process_config_file(void) {
    Config config;
    memset(&config, 0, sizeof(config));

    int ret = process_config_file(&config, "../../../../configs/capps/mesh/config.toml");
    TEST_ASSERT_EQUAL_INT(TRUE, ret);
}

// Test for clear_config function
void test_clear_config(void) {
    Config *config;
    config = (Config *)calloc(1, sizeof(Config));

    config->remoteConnect = strdup("remoteConnect");
    config->localHostname = strdup("localHostname");
    config->certFile      = strdup("certFile");
    config->keyFile       = strdup("keyFile");

    clear_config(config);

    TEST_ASSERT_NULL(config->remoteConnect);
    TEST_ASSERT_NULL(config->localHostname);
    TEST_ASSERT_NULL(config->certFile);
    TEST_ASSERT_NULL(config->keyFile);
}

void run_all_tests_config(void) {

    RUN_TEST(test_read_line);
    RUN_TEST(test_split_strings);
    RUN_TEST(test_read_nodeid);
    RUN_TEST(test_process_config_file);
    //    RUN_TEST(test_clear_config);
}

