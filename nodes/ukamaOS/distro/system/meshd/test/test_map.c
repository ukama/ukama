#include "unity.h"
#include "map.h"
#include "mesh.h"

MapTable *table = NULL;

// Test for init_map_table function
void test_init_map_table(void) {
    TEST_ASSERT_NULL(table->first);
    TEST_ASSERT_NULL(table->last);
}

// Test for create_map_item function
void test_create_map_item(void) {
    MapItem *item = create_map_item("service", "port");
    TEST_ASSERT_NOT_NULL(item);
    TEST_ASSERT_EQUAL_STRING("service", item->serviceInfo->name);
    TEST_ASSERT_EQUAL_STRING("port", item->serviceInfo->port);

    free_map_item(item);
}

// Test for is_existing_item function
void test_is_existing_item(void) {
    add_map_to_table(&table, "service", "port");
    MapItem *item = is_existing_item(table, "service", "port");
    TEST_ASSERT_NOT_NULL(item);

    item = is_existing_item(table, "nonexistent", "port");
    TEST_ASSERT_NULL(item);
}

// Test for add_map_to_table function
void test_add_map_to_table(void) {
    MapItem *item = add_map_to_table(&table, "service", "port");
    TEST_ASSERT_NOT_NULL(item);
    TEST_ASSERT_EQUAL_STRING("service", item->serviceInfo->name);
    TEST_ASSERT_EQUAL_STRING("port", item->serviceInfo->port);
}

// Test for remove_map_item_from_table function
void test_remove_map_item_from_table(void) {
    add_map_to_table(&table, "service", "port");
    remove_map_item_from_table(table, "service", "port");

    MapItem *item = is_existing_item(table, "service", "port");
    TEST_ASSERT_NULL(item);
}

void run_all_tests_map(void) {

    table = (MapTable *)calloc(1, sizeof(MapTable));
    RUN_TEST(test_init_map_table);
    RUN_TEST(test_create_map_item);
    RUN_TEST(test_is_existing_item);
    RUN_TEST(test_add_map_to_table);
    RUN_TEST(test_remove_map_item_from_table);
    free(table);
}

