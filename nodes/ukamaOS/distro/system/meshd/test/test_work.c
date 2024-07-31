#include "unity.h"
#include "work.h"
#include "mesh.h"

/* work.c */
extern WorkItem *create_work_item(char *data, thread_func_t pre,
                                  void *preArgs, thread_func_t post,
                                  void *postArgs);

// Mock functions and global variables
WorkList *list = NULL;

// Test for init_work_list function
void test_init_work_list(void) {
    TEST_ASSERT_NULL(list->first);
    TEST_ASSERT_NULL(list->last);
}

// Test for create_work_item function
void test_create_work_item(void) {
    WorkItem *item = create_work_item("data", NULL, NULL, NULL, NULL);
    TEST_ASSERT_NOT_NULL(item);
    TEST_ASSERT_EQUAL_STRING("data", item->data);

    destroy_work_item(item);
}

// Test for add_work_to_queue function
void test_add_work_to_queue(void) {
    int ret = add_work_to_queue(&list, "data", NULL, NULL, NULL, NULL);
    TEST_ASSERT_EQUAL_INT(TRUE, ret);
    TEST_ASSERT_NOT_NULL(list->first);
    TEST_ASSERT_NOT_NULL(list->last);
}

// Test for get_work_to_transmit function
void test_get_work_to_transmit(void) {
    add_work_to_queue(&list, "data", NULL, NULL, NULL, NULL);
    WorkItem *item = get_work_to_transmit(list);
    TEST_ASSERT_NOT_NULL(item);
    TEST_ASSERT_EQUAL_STRING("data", item->data);

    destroy_work_item(item);
}

void run_all_tests_work(void) {

    list = (WorkList *)malloc(sizeof(WorkList));
    init_work_list(&list);

    RUN_TEST(test_init_work_list);
    RUN_TEST(test_create_work_item);
    RUN_TEST(test_add_work_to_queue);
    RUN_TEST(test_get_work_to_transmit);

    free(list);
}

