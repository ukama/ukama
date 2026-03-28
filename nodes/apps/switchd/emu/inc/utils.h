/* ... */
#ifndef UTILS_H
#define UTILS_H

#include <stddef.h>

#include "switchemu.h"

int json_get_string_field(const char *json, const char *key,
                          char *buf, size_t bufLen);
int json_get_int_field(const char *json, const char *key, int *value);
int json_get_bool_field(const char *json, const char *key, int *value);
int write_all(int fd, const void *buf, size_t len);
const char *bool_to_json(int value);

#endif /* UTILS_H */
