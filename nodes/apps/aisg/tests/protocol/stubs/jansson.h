#ifndef AISG_TEST_STUB_JANSSON_H_
#define AISG_TEST_STUB_JANSSON_H_

#include <stddef.h>

typedef struct json_t json_t;
typedef long long json_int_t;

static inline json_t *json_object(void) { return (json_t *)0; }
static inline json_t *json_array(void) { return (json_t *)0; }
static inline json_t *json_string(const char *s) { (void)s; return (json_t *)0; }
static inline json_t *json_integer(json_int_t v) { (void)v; return (json_t *)0; }
static inline json_t *json_real(double v) { (void)v; return (json_t *)0; }
static inline json_t *json_boolean(int v) { (void)v; return (json_t *)0; }
static inline json_t *json_true(void) { return (json_t *)0; }
static inline json_t *json_false(void) { return (json_t *)0; }
static inline json_t *json_null(void) { return (json_t *)0; }
static inline int json_object_set_new(json_t *o, const char *k, json_t *v)
{ (void)o; (void)k; (void)v; return 0; }
static inline int json_array_append_new(json_t *a, json_t *v)
{ (void)a; (void)v; return 0; }
static inline void json_decref(json_t *v) { (void)v; }
static inline json_t *json_object_get(const json_t *o, const char *k)
{ (void)o; (void)k; return (json_t *)0; }
static inline int json_is_string(const json_t *v) { (void)v; return 0; }
static inline int json_is_integer(const json_t *v) { (void)v; return 0; }
static inline int json_is_number(const json_t *v) { (void)v; return 0; }
static inline const char *json_string_value(const json_t *v)
{ (void)v; return (const char *)0; }
static inline json_int_t json_integer_value(const json_t *v)
{ (void)v; return 0; }
static inline double json_number_value(const json_t *v) { (void)v; return 0.0; }

#endif /* AISG_TEST_STUB_JANSSON_H_ */
