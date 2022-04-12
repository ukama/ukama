/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/value.h $                                                */
/*                                                                        */
/* OpenPOWER FFS Project                                                  */
/*                                                                        */
/* Contributors Listed Below - COPYRIGHT 2014,2015                        */
/* [+] International Business Machines Corp.                              */
/*                                                                        */
/*                                                                        */
/* Licensed under the Apache License, Version 2.0 (the "License");        */
/* you may not use this file except in compliance with the License.       */
/* You may obtain a copy of the License at                                */
/*                                                                        */
/*     http://www.apache.org/licenses/LICENSE-2.0                         */
/*                                                                        */
/* Unless required by applicable law or agreed to in writing, software    */
/* distributed under the License is distributed on an "AS IS" BASIS,      */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or        */
/* implied. See the License for the specific language governing           */
/* permissions and limitations under the License.                         */
/*                                                                        */
/* IBM_PROLOG_END_TAG                                                     */

/*!
 * @file value.h
 * @brief Value Container
 * @details Values are a kind of variant structure that can contain
 *          values of different types
 * @author Shaun Wetzstein <shaun@us.ibm.com>
 * @date 2010-2011
 */

#ifndef __VALUE_H__
#define __VALUE_H__

#include <stdint.h>
#include <stdbool.h>

#include "libclib.h"

/* ======================================================================= */

typedef enum value_type_enum value_type_t;	//!< Alias for the @em value_type_enum enum
typedef struct value value_t;	//!< Alias for the @em value  class

/*!
 * @brief Value types
 * @details Supported value type
 */
enum value_type_enum {
	VT_UNKNOWN = 0,		//!< Uninitialized or unknown type
	VT_I8,			//!<  8-bit signed integer type
	VT_I16,			//!< 16-bit signed integer type
	VT_I32,			//!< 32-bit signed integer type
	VT_I64,			//!< 64-bit signed integer type
	VT_U8,			//!<  8-bit unsigned integer type
	VT_U16,			//!< 16-bit unsigned integer type
	VT_U32,			//!< 32-bit unsigned integer type
	VT_U64,			//!< 64-bit unsigned integer type
	VT_REAL32,		//!< 32-bit float type
	VT_REAL64,		//!< 64-bit float type
	VT_STR,			//!< Character string type
	VT_STR_INLINE,		//!< Character string inline type
	VT_STR_OFF,		//!< Character string offset type
	VT_STR_CONST,		//!< Constant character string type
	VT_BLOB,		//!< Blob type
	VT_BLOB_INLINE,		//!< Blob inline type
	VT_BLOB_OFF,		//!< Blob offset type
	VT_BLOB_FILE,		//!< Blob file type
} /*! @cond */ __packed /*! @endcond */ ;

/*! @cond */
#ifndef VALUE_PAD_DEFAULT
#define VALUE_PAD_DEFAULT 12
#endif
/*! @endcond */ ;

/*!
 * @brief value container
 */
struct value {
	uint32_t type:8,	//!< Value type
	 size:24;		//!< Value size (in bytes)

	union {
		int8_t i8;	//!<  8-bit signed integer
		int16_t i16;	//!< 16-bit signed integer
		int32_t i32;	//!< 32-bit signed integer
		int64_t i64;	//!< 64-bit signed integer

		uint8_t u8;	//!<  8-bit unsigned integer
		uint16_t u16;	//!< 16-bit unsigned integer
		uint32_t u32;	//!< 32-bit unsigned integer
		uint64_t u64;	//!< 64-bit unsigned integer

		float r32;	//!< 32-bit float
		double r64;	//!< 64-bit float

		void *ptr;	//!< pointer

		uint8_t data[VALUE_PAD_DEFAULT];
	};			//!< Anonymous variant record
};

/* ==================================================================== */

/*! @cond */
#define value_type(...) STRCAT(value_type, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline value_type_t value_type1(const value_t * self)
{
	return self ? self->type : VT_UNKNOWN;
}

static inline void value_type2(value_t * self, value_type_t type)
{
	if (self != NULL)
		self->type = type;
}

/*! @endcond */

/*! @cond */
#define value_size(...) STRCAT(value_size, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline uint32_t value_size1(const value_t * self)
{
	return self ? self->size : 0;
}

static inline void value_size2(value_t * self, uint32_t size)
{
	if (self != NULL)
		self->size = size;
}

/*! @endcond */

/*!
 * @brief Clear the contents of a @em value
 * @memberof value
 * @param self [in] value object @em self pointer
 */
static inline void value_clear(value_t * self)
{
	if (self != NULL) {
		if (self->type == VT_BLOB || self->type == VT_STR)
			if (self->ptr != NULL)
				free(self->ptr);
		memset(self, 0, sizeof *self);
	}
}

#define value_i8(...) STRCAT(value_i8_, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline int8_t value_i8_1(const value_t * v)
{
	assert(value_type(v) == VT_I8);
	return v->i8;
}

static inline value_t *value_i8_2(value_t * v, const int8_t i8)
{
	value_clear(v), v->type = VT_I8, v->i8 = i8, v->size = sizeof(i8);
	return v;
}

#define value_i16(...) STRCAT(value_i16_, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline int16_t value_i16_1(const value_t * v)
{
	assert(value_type(v) == VT_I16);
	return v->i16;
}

static inline value_t *value_i16_2(value_t * v, const int16_t i16)
{
	value_clear(v), v->type = VT_I16, v->i16 = i16, v->size = sizeof(i16);
	return v;
}

#define value_i32(...) STRCAT(value_i32_, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline int32_t value_i32_1(const value_t * v)
{
	assert(value_type(v) == VT_I32);
	return v->i32;
}

static inline value_t *value_i32_2(value_t * v, const int32_t i32)
{
	value_clear(v), v->type = VT_I32, v->i32 = i32, v->size = sizeof(i32);
	return v;
}

#define value_i64(...) STRCAT(value_i64_, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline int64_t value_i64_1(const value_t * v)
{
	assert(value_type(v) == VT_I64);
	return v->i64;
}

static inline value_t *value_i64_2(value_t * v, const int64_t i64)
{
	value_clear(v), v->type = VT_I64, v->i64 = i64, v->size = sizeof(i64);
	return v;
}

#define value_u8(...) STRCAT(value_u8_, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline uint8_t value_u8_1(const value_t * v)
{
	assert(value_type(v) == VT_U8);
	return v->u8;
}

static inline value_t *value_u8_2(value_t * v, const uint8_t u8)
{
	value_clear(v), v->type = VT_U8, v->u8 = u8, v->size = sizeof(u8);
	return v;
}

#define value_u16(...) STRCAT(value_u16_, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline uint16_t value_u16_1(const value_t * v)
{
	assert(value_type(v) == VT_U16);
	return v->u16;
}

static inline value_t *value_u16_2(value_t * v, const uint16_t u16)
{
	value_clear(v), v->type = VT_U16, v->u16 = u16, v->size = sizeof(u16);
	return v;
}

#define value_u32(...) STRCAT(value_u32_, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline uint32_t value_u32_1(const value_t * v)
{
	assert(value_type(v) == VT_U32);
	return v->u32;
}

static inline value_t *value_u32_2(value_t * v, const uint32_t u32)
{
	value_clear(v), v->type = VT_U32, v->u32 = u32, v->size = sizeof(u32);
	return v;
}

#define value_u64(...) STRCAT(value_u64_, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline uint64_t value_u64_1(const value_t * v)
{
	assert(value_type(v) == VT_U64);
	return v->u64;
}

static inline value_t *value_u64_2(value_t * v, const uint64_t u64)
{
	value_clear(v), v->type = VT_U64, v->u64 = u64, v->size = sizeof(u64);
	return v;
}

#define value_float(...) STRCAT(value_float_, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline float value_float1(const value_t * v)
{
	assert(value_type(v) == VT_REAL32);
	return v->r32;
}

static inline value_t *value_float2(value_t * v, const float r32)
{
	value_clear(v), v->type = VT_REAL32, v->r32 = r32, v->size =
	    sizeof(r32);
	return v;
}

#define value_double(...) STRCAT(value_double, NARGS(__VA_ARGS__))(__VA_ARGS__)
static inline double value_double1(const value_t * v)
{
	assert(value_type(v) == VT_REAL64);
	return v->r64;
}

static inline value_t *value_double2(value_t * v, const double r64)
{
	value_clear(v), v->type = VT_REAL64, v->r64 = r64, v->size =
	    sizeof(r64);
	return v;
}

#define value_blob(...) STRCAT(value_blob, NARGS(__VA_ARGS__))(__VA_ARGS__)
/*! @cond */
static inline void *value_blob1(const value_t * v)
{
	if (value_type(v) == VT_BLOB)
		return v->ptr;
	else if (value_type(v) == VT_BLOB_INLINE)
		return (void *)v->data;
	else if (value_type(v) == VT_BLOB_OFF)
		return NULL;
	else if (value_type(v) == VT_BLOB_FILE) ;
	else {
		UNEXPECTED("invalid value_t blob type '%d'", value_type(v));
	}

	return NULL;
}

static inline value_t *value_blob4(value_t * self, const void *blob, size_t len,
				   size_t pad)
{
	value_clear(self);

	if (pad <= len) {
		self->ptr = malloc(len + 1);
		self->type = VT_BLOB;

		memcpy(self->ptr, blob, len);
		memset(self->ptr + len, 0, 1);
	} else {
		self->type = VT_BLOB_INLINE;

		memcpy(self->data, blob, len);
		self->data[len] = '\0';
	}

	self->size = len;

	return self;
}

static inline value_t *value_blob3(value_t * self, const void *blob, size_t len)
{
	return value_blob4(self, blob, len, VALUE_PAD_DEFAULT);
}

static inline value_t *value_blob2(value_t * self, const char *path)
{
	value_t *rc = value_blob4(self, path, strlen(path), VALUE_PAD_DEFAULT);
	if (rc != NULL)
		value_type(rc, VT_BLOB_FILE);
	return rc;
}

/*! @endcond */

#define value_string(...) STRCAT(value_string, NARGS(__VA_ARGS__))(__VA_ARGS__)
/*! @cond */
static inline char *value_string1(const value_t * v)
{
	if (value_type(v) == VT_STR)
		return (char *)v->ptr;
	else if (value_type(v) == VT_STR_INLINE)
		return (char *)v->data;
	else if (value_type(v) == VT_STR_OFF) ;
	else {
		UNEXPECTED("invalid value string type '%d'", value_type(v));
	}

	return NULL;
}

static inline value_t *value_string4(value_t * self, const char *str,
				     size_t len, size_t pad)
{
	value_blob(self, (const char *)str, len, pad);

	if (self->type == VT_BLOB)
		self->type = VT_STR;
	else if (self->type == VT_BLOB_INLINE)
		self->type = VT_STR_INLINE;

	return self;
}

static inline value_t *value_string3(value_t * self, const char *str,
				     size_t len)
{
	return value_string4(self, str, len, VALUE_PAD_DEFAULT);
}

static inline value_t *value_string2(value_t * self, const char *str)
{
	return value_string4(self, str, strlen(str), VALUE_PAD_DEFAULT);
}

static inline value_t *value_set_string_const(value_t * self, const char *str,
					      size_t len)
{
	value_clear(self);

	self->type = VT_STR_CONST;
	self->ptr = (void *)str;
	self->size = len;

	return self;
}

static inline value_t *value_set_string_offset(value_t * self, size_t off,
					       size_t len)
{
	value_clear(self);

	self->type = VT_STR_OFF;
	self->u32 = off;
	self->size = len;

	return self;
}

/*! @endcond */

/*!
 * @brief Pretty print the contents of a @em value to stdout
 * @memberof value
 * @param self [in] value object @em self pointer
 * @param out [in] output stream
 */
extern void value_dump(const value_t * self, FILE * out)
/*! @cond */
__nonnull((1, 2)) /*! @endcond */ ;

/*! @cond */
#define __choose__	__builtin_choose_expr
#define __compat__	__builtin_types_compatible_p
#define __const__	__builtin_constant_p
/*! @endcond */

/*!
 * @fn T value_get(const value * self, T)
 * @brief Return a @em value as primitive type @em T
 * @note @em T can be one fo the following types:
 *       - @em char, @em int8_t
 *       - @em short, @em int16_t
 *       - @em int, @em long @em int32_t
 *       - @em long @em long, @em int64_t
 *       - @em unsigned @em char, @em uint8_t
 *       - @em unsigned @em short, @em uint16_t
 *       - @em unsigned @em int, @em long @em uint32_t
 *       - @em unsigned @em long @em long, @em uint64_t
 *       - @em float, @em double
 *       - @em char @em *, @em unsigned @em char @em *
 * @memberof value
 * @param self [in] value object @em self pointer
 * @param T [in] Type to cast
 * @throws ASSERTION if T is not equal to the object's type
 */

#if DELETE
/*! @cond */
#define value_get(v, T)	({					\
    __choose__(__compat__(T, char),				\
           value_get_int8((v)),					\
    __choose__(__compat__(T, int8_t),				\
           value_get_int8((v)),					\
    __choose__(__compat__(T, short),				\
           value_get_int16((v)),				\
    __choose__(__compat__(T, int16_t),				\
           value_get_int16((v)),				\
    __choose__(__compat__(T, int),				\
           value_get_int32((v)),				\
    __choose__(__compat__(T, long),				\
           value_get_int32((v)),				\
    __choose__(__compat__(T, int32_t),				\
           value_get_int32((v)),				\
    __choose__(__compat__(T, long long),			\
           value_get_int64((v)),				\
    __choose__(__compat__(T, int64_t),				\
           value_get_int64((v)),				\
    __choose__(__compat__(T, unsigned char),			\
           value_get_uint8((v)),				\
    __choose__(__compat__(T, uint8_t),				\
           value_get_uint8((v)),				\
    __choose__(__compat__(T, unsigned short),			\
           value_get_uint16((v)),				\
    __choose__(__compat__(T, uint16_t),				\
           value_get_uint16((v)),				\
    __choose__(__compat__(T, unsigned int),			\
           value_get_uint32((v)),				\
    __choose__(__compat__(T, unsigned long),			\
           value_get_uint32((v)),				\
    __choose__(__compat__(T, uint32_t),				\
           value_get_uint32((v)),				\
    __choose__(__compat__(T, unsigned long long),		\
           value_get_uint64((v)),				\
    __choose__(__compat__(T, uint64_t),				\
           value_get_uint64((v)),				\
    __choose__(__compat__(T, float),				\
           value_get_float((v)),				\
    __choose__(__compat__(T, double),				\
           value_get_double((v)),				\
    __choose__(__compat__(T, char *),				\
           value_get_string((v)),				\
    __choose__(__compat__(T, unsigned char *),			\
           value_get_string((v)),				\
    0xDEADBEEF))))))))))))))))))))));				\
})
/*! @endcond */
#endif
/*!
 * @fn void value_set(const value * self, T v)
 * @brief Assign a @em value to @em T v
 * @memberof value
 * @param self [in] value object @em self pointer
 * @param v [in] value to assign
 */

#if 0
/*! @cond */
#define value_set(v, x)	({						\
     __choose__(__compat__(typeof (x), char),				\
           value_set_int8((v),(int8_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), int8_t),				\
           value_set_int8((v),(int8_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), short),				\
           value_set_int16((v),(int16_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), int16_t),				\
           value_set_int16((v),(int16_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), int),				\
           value_set_int32((v),(int32_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), long),				\
           value_set_int32((v),(int32_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), int32_t),				\
           value_set_int32((v),(int32_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), long long),			\
           value_set_int64((v),(int64_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), int64_t),				\
           value_set_int64((v),(int64_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), unsigned char),			\
           value_set_uint8((v),(uint8_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), uint8_t),				\
           value_set_uint8((v),(uint8_t)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), unsigned short),			\
           value_set_uint16((v),(uint16_t)(intptr_t)(x)),		\
    __choose__(__compat__(typeof (x), uint16_t),			\
           value_set_uint16((v),(uint16_t)(intptr_t)(x)),		\
    __choose__(__compat__(typeof (x), unsigned int),			\
           value_set_uint32((v),(uint32_t)(intptr_t)(x)),		\
    __choose__(__compat__(typeof (x), unsigned long),			\
           value_set_uint32((v),(uint32_t)(intptr_t)(x)),		\
    __choose__(__compat__(typeof (x), uint32_t),			\
           value_set_uint32((v),(uint32_t)(intptr_t)(x)),		\
    __choose__(__compat__(typeof (x), unsigned long long),		\
           value_set_uint64((v),(uint64_t)(intptr_t)(x)),		\
    __choose__(__compat__(typeof (x), uint64_t),			\
           value_set_uint64((v),(uint64_t)(intptr_t)(x)),		\
    __choose__(__compat__(typeof (x), float),				\
           value_set_float((v),(float)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), double),				\
           value_set_double((v),(double)(intptr_t)(x)),			\
    __choose__(__compat__(typeof (x), char *),				\
        __choose__(__const__((x)),					\
            value_set_string_const((v),(const char *)(intptr_t)(x),	\
                                   strlen((const char *)(intptr_t)(x))),\
           value_set_string((v),(char *)(intptr_t)(x),			\
                            strlen((const char *)(intptr_t)(x)))),	\
    __choose__(__compat__(typeof (x), unsigned char *),			\
        __choose__(__const__((x)),					\
            value_set_string_const((v),(const char *)(intptr_t)(x),	\
                                   strlen((const char *)(intptr_t)(x))),\
            value_set_string((v),(const char *)(intptr_t)(x),		\
                             strlen((const char *)(intptr_t)(x)))),	\
   ((void)0)))))))))))))))))))))));					\
})
/*! @endcond */
#endif

/* ==================================================================== */

#endif				/* __VALUE_H__ */
