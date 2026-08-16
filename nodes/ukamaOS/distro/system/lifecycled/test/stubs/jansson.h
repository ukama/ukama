#pragma once

#include <stddef.h>
#include <stdint.h>

typedef struct json_t json_t;
typedef int64_t json_int_t;
typedef struct {
    int line;
} json_error_t;

#define JSON_COMPACT 0

json_t *json_object(void);
json_t *json_array(void);
json_t *json_string(const char *value);
json_t *json_integer(json_int_t value);
json_t *json_true(void);
json_t *json_null(void);
json_t *json_pack(const char *format, ...);
json_t *json_loads(const char *input, size_t flags, json_error_t *error);
json_t *json_loadb(const char *input,
                   size_t length,
                   size_t flags,
                   json_error_t *error);
json_t *json_object_get(const json_t *object, const char *key);
json_t *json_array_get(const json_t *array, size_t index);
size_t json_array_size(const json_t *array);
const char *json_string_value(const json_t *string);
int json_is_object(const json_t *json);
int json_is_array(const json_t *json);
int json_is_string(const json_t *json);
int json_object_set_new(json_t *object, const char *key, json_t *value);
void json_decref(json_t *json);
char *json_dumps(const json_t *json, size_t flags);

#define json_boolean(value) ((value) ? json_true() : json_pack("b", 0))
#define json_array_foreach(array, index, value)                           \
    for ((index) = 0;                                                     \
         (index) < json_array_size(array) &&                              \
         (((value) = json_array_get((array), (index))) || 1);             \
         (index)++)

