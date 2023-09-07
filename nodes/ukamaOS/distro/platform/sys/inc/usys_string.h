/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_STRING_H_
#define USYS_STRING_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

/**
 * @fn     char usys_strtok*(char*, const char*)
 * @brief  Abstraction to strtok
 *
 * @param  str
 * @param  delim
 * @return If a token is found, a pointer to the beginning of the token.
 *         Otherwise, a null pointer.
 */
static inline char *usys_strtok(char *str, const char *delim) {
    return strtok(str, delim);
}

/**
 * @fn     char usys_strcpy*(char*, const char*)
 * @brief  Copies the C string pointed by source into the array pointed by
 *         destination, including the terminating null character
 *
 * @param  dest
 * @param  src
 * @return destination
 */
static inline char *usys_strcpy(char *dest, const char *src) {
    return strcpy(dest, src);
}

/**
 * @fn     char usys_strncpy*(char*, const char*, size_t)
 * @brief  Copies the first num characters of source to destination.
 *         If the end of the source C string is found before num characters
 *         have been copied, destination is padded with zeros until a total of
 *         num characters have been written to it.
 *
 * @param  dest
 * @param  src
 * @param  n
 * @return destination string
 */
static inline char *usys_strncpy(char *dest, const char *src, size_t n) {
    return strncpy(dest, src, n);
}

/**
 * @fn     int usys_strcmp(const char*, const char*)
 * @brief  Compares the C string s1 to the C string s2.
 *
 * @param  s1
 * @param  s2
 * @return Returns an integral value indicating the difference in strings
 */
static inline int usys_strcmp(const char *s1, const char *s2) {
    return strcmp(s1, s2);
}

/**
 * @fn     int usys_strncmp(const char*, const char*, size_t)
 * @brief  Compares the first n bytes of C string s1 to the C string s2.
 *
 * @param  s1
 * @param  s2
 * @param  n
 * @return Returns an integral value indicating the difference in strings
 */
static inline int usys_strncmp(const char *s1, const char *s2, size_t n) {
    return strncmp(s1, s2, n);
}

/**
 * @fn     int usys_strcasecmp(const char*, const char*)
 * @brief  Compares two strings irrespective of the case of characters
 *         by converting strings to lower case and then comparing them
 *
 * @param  s1
 * @param  s2
 * @return Returns an integral value indicating the difference in strings
 */
static inline int usys_strcasecmp(const char *s1, const char *s2) {
    return strcasecmp(s1, s2);
}

/**
 * @fn     int usys_strncasecmp(const char*, const char*, size_t)
 * @brief  Compares first n bytes of two strings irrespective of the case of
 *         characters by converting strings to lower case and
 *         then comparing them
 *
 * @param  s1
 * @param  s2
 * @param  n
 * @return
 */
static inline int usys_strncasecmp(const char *s1, const char *s2, size_t n) {
    return strncasecmp(s1, s2, n);
}

/**
 * @fn     size_t usys_strlen(const char*)
 * @brief  The length of a C string is determined by the terminating
 *         null-character
 *
 * @param  s
 * @return integer value indicating length of string.
 */
static inline size_t usys_strlen(const char *s) {
    return strlen(s);
}

/**
 * @fn     char usys_strcat*(char*, const char*)
 * @brief  Appends a copy of the source string to the destination string.
 *
 * @param  dest
 * @param  src
 * @return destination string
 */
static inline char *usys_strcat(char *dest, const char *src) {
    return strcat(dest, src);
}

/**
 * @fn     char usys_strncat*(char*, const char*, size_t)
 * @brief  Appends first n  of the source string to the destination string.
 *
 * @param  destination
 * @param  source
 * @param  num
 * @return destination string
 */
static inline char *usys_strncat(char *dest, const char *src, size_t num) {
    return strcat(dest, src);
}

/**
 * @fn     char usys_strstr*(const char*, const char*)
 * @brief  a pointer to the first occurrence of str2 in str1, or a null
 *         pointer if str2 is not part of str1.
 *
 * @param  str1
 * @param  str2
 * @return A pointer to the first occurrence in str1 of the entire sequence
 *         of characters specified in str2, or a null pointer if the sequence
 *         is not present in str1.
 */
static inline char *usys_strstr(const char *str1, const char *str2) {
    return strstr((char *)str1, (char *)str2);
}

/**
 * @fn     int usys_strspn(const char*, const char*)
 * @brief  Returns the length of the initial portion of str1 which consists
 *         only of characters that are part of str2.
 *
 * @param  str1
 * @param  str2
 * @return
 */
static inline int usys_strspn(const char *str1, const char *str2) {
    return strspn((char *)str1, (char *)str2);
}

/**
 * @fn     void usys_memset*(void*, int, size_t)
 * @brief  Sets the first num bytes of the block of memory pointed by ptr to
 *         the specified value
 *
 * @param  ptr
 * @param  value
 * @param  num
 * @return ptr
 */
static inline void *usys_memset(void *ptr, int value, size_t num) {
    return memset(ptr, value, num);
}

/**
 * @fn     void usys_memcpy*(void*, void*, size_t)
 * @brief  Copies the values of num bytes from the location pointed to by
 *         source directly to the memory block pointed to by destination.
 *
 * @param  dest
 * @param  src
 * @param  num
 * @return dest
 */
static inline char *usys_memcpy(void *dest, const void *src, size_t num) {
    return memcpy(dest, src, num);
}

/**
 * @fn      int usys_memcmp(void*, void*, size_t)
 * @brief   Compares the first num bytes of the block of memory pointed by ptr1
 *          to the first num bytes pointed by ptr2, returning zero if they all
 *          match or a value different from zero representing which is greater
 *          if they do not
 *
 * @param  ptr1
 * @param  ptr2
 * @param  num
 * @return
 */
static inline int usys_memcmp(void *ptr1, void *ptr2, size_t num) {
    return memcmp(ptr1, ptr2, num);
}

/**
 * @fn      char* usys_strdup(char*)
 * @brief   Returns a pointer to a null-terminated byte string,
 *          which is a duplicate of the string pointed to by str
 *
 * @param   str
 * @return  On Success, Address of a null terminated duplicate string
 *          On Failure, NULL
 */
static inline char* usys_strdup(char* str) {
    return  strdup(str);
}

/**
 * @fn      int usys_strndup(char*)
 * @brief   Returns a pointer to a null-terminated byte string,
 *          which is a duplicate of the string pointed to by str till nth
 *          character
 *
 * @param   str
 * @param   n
 * @return  On Success, Address of a null terminated duplicate string
 *          On Failure, NULL
 */
static inline char* usys_strndup(char* str, size_t n) {
    return  strndup(str, n);
}

#ifdef __cplusplus
}
#endif

#endif /* USYS_STRING_H_ */
