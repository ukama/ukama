/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_API_H_
#define USYS_API_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

/**
 * @fn     FILE usys_fopen*(const char*, const char*)
 * @brief  Opens the filename pointed to, by filename using the given mode.
 *
 * @param  path
 * @param  mode
 * @return On Success FILE pointer otherwise NULL
 */
static inline FILE* usys_fopen(const char* path, const char* mode) {
    return fopen(path, mode);
}

/**
 * @fn     int usys_fclose(FILE*)
 * @brief  closes the stream pointed by FILE*
 *
 * @param  fp
 * @return On Success 0 otherwise EOF
 */
static inline int usys_fclose(FILE* fp) {
    return fclose(fp);
}

/**
 * @fn     usys_fread
 *
 * @brief  This function is a wrapper over glibc fread. Refer fread man page for details.
 *
 */
static inline int usys_fread(void* ptr, size_t size, size_t nmemb, FILE *stream) {
    return fread(ptr, size, nmemb, stream);
}

/**
 * @fn     int usys_fwrite(const void*, size_t, size_t, FILE*)
 * @brief  Writes an array of count elements, each one with a size of size bytes,
 *         from the block of memory pointed by ptr to the current position
 *         in the stream.
 *
 * @param  ptr
 * @param  size
 * @param  nmemb
 * @param  stream
 * @return The total number of elements successfully written is returned.
 */
static inline int usys_fwrite(const void* ptr, size_t size, size_t nmemb, FILE *stream) {
    return fwrite(ptr, size, nmemb, stream);
}

/**
 * @fn     int usys_fseek(FILE*, long int, int)
 * @brief  Sets the position indicator associated with the stream to a new position.
 *
 * @param  stream
 * @param  offset
 * @param  origin
 * @return On Success, the function returns zero.
 *         On Failure, it returns non-zero value.
 */
static inline int usys_fseek(FILE * stream, long int offset, int origin) {
    return fseek(stream, offset, origin);
}

/**
 * @fn     char usys_fgets*(char*, int, FILE*)
 * @brief  Reads characters from stream and stores them as a C string into str until (num-1)
 *         characters have been read or either a newline or the end-of-file is reached,
 *         whichever happens first.
 *
 * @param  s
 * @param  size
 * @param  stream
 * @return On success, the function returns str.
 *         If the end-of-file is encountered while attempting to read a character, the eof indicator is set (feof).
 *         If this happens before any characters could be read, the pointer returned is a null pointer.
 *         If a read error occurs, the error indicator (ferror) is set and a null pointer is also returned
 */
static inline char* usys_fgets(char* s, int size, FILE* stream) {
    return fgets(s, size, stream);
}

/**
 * @def    usys_snprintf
 * @brief  Wraps snprintf function.
 *         The content is stored as a C string in the buffer pointed by str
 *         with size as max length.
 */
#define usys_snprintf(str, size, format, ...) snprintf(str, size, format, ##__VA_ARGS__);

/**
 * @def   usys_fprintf
 * @brief Wraps fprintf function.
 *        Writes the C string pointed by format to the stream.
 */
#define usys_fprintf(stream, format, ...)\
{\
    fprintf(stream, format, ##__VA_ARGS__);\
}

/**
 * @def   usys_printf
 * @brief Wraps printf function
 *        Writes the C string pointed by format to the standard output (stdout).
 */
#define usys_printf(format, ...)\
{\
    printf(format, ##__VA_ARGS__);\
}

/**
 * @def   usys_sprintf
 * @brief Wraps sprintf function.
 *        The content is stored as a C string in the buffer pointed by str
 */
#define usys_sprintf(str, format, ...)\
{\
    sprintf(str, format, ##__VA_ARGS__);\
}

/**
 * @def   usys_sscanf
 * @brief Reads data from str and stores them according to parameter
 *        format into the locations given by the additional arguments.
 */
#define usys_sscanf(str, format, ...)\
{\
    sscanf(str, format, ##__VA_ARGS__);\
}

/**
 * @fn    struct timeval usys_gettimeofday()
 * @brief function gets the system’s clock time
 *
 * @return On success,return time value
 *         On failure the returns 0 initilaized values.
 */
static inline struct timeval usys_gettimeofday() {
    struct timeval tv = {0};
    gettimeofday(&tv, NULL);
    return tv;
}

/**
 * @fn     size_t usys_strftime(char*, size_t, const char*, const struct tm*)
 * @brief  This function formats the broken-down time tm according
 *         to the format specification format and places the result in the
 *         character array s of size max.
 *
 * @param  s
 * @param  max
 * @param  format
 * @param  tm
 * @return returns the number of bytes placed in the array s
 */
static inline size_t usys_strftime(char *s, size_t max, const char *format, const struct tm *tm) {
    return strftime(s, max, format, tm);
}

/**
 * @fn     struct tm usys_localtime*(const time_t*)
 * @brief  Uses the value pointed by timer to fill a tm structure with the
 *         values that represent the corresponding time,
 *         expressed for the local timezone.
 *
 * @param  timep
 * @return A pointer to a tm structure with its members filled
 *         with the values that correspond to the local time representation of timer.
 */
static inline struct tm *usys_localtime(const time_t* timep) {
    return localtime(timep);
}

/**
 * @fn     uint16_t usys_ntohs(uint16_t)
 * @brief  Converts the unsigned short integer netshort from network byte order to
 *         host byte order.
 *
 * @param  netshort
 * @return uint16_t
 */
static inline uint16_t usys_ntohs(uint16_t netshort) {
    return ntohs(netshort);
}

/**
 * @fn     uint16_t usys_htons(uint16_t)
 * @brief  Converts the unsigned short integer hostshort from host byte order to
 *         network byte order.
 *
 * @param  hostshort
 * @return uint16_t
 */
static inline uint16_t usys_htons(uint16_t hostshort) {
    return htons(hostshort);
}

/**
 * @fn     uint32_t usys_ntohl(uint32_t)
 * @brief  Converts the unsigned integer netlong from network byte order to
 *         host byte order.
 *
 * @param  netlong
 * @return uint32_T
 */
static inline uint32_t usys_ntohl(uint32_t netlong) {
    return ntohl(netlong);
}

/**
 * @fn     uint32_t usys_htonl(uint32_t)
 * @brief  Converts the unsigned integer hostlong from host byte order to
 *         network byte order.
 *
 * @param  hostlong
 * @return uint32_t
 */
static inline uint32_t usys_htonl(uint32_t hostlong) {
    return htonl(hostlong);
}

/**
 * @fn     uint32_t usys_inet_addr(const char*)
 * @brief  converts the Internet host address cp from IPv4
 *         numbers-and-dots notation into binary data in network byte order.
 *
 * @param  data
 * @return uint32_t
 */
static inline uint32_t usys_inet_addr(const char* data) {
    return inet_addr(data);
}

/**
 * @fn     uint64_t usys_le64toh(uint64_t)
 * @brief  Convert from little-endian order to host byte order.
 *
 * @param  netlonglong
 * @return uint64_t
 */
static inline uint64_t usys_le64toh(uint64_t netlonglong) {
    return le64toh(netlonglong);
}

/**
 * @fn     uint64_t usys_htobe64(uint64_t)
 * @brief  Convert from host byte order to big-endian order.
 *
 * @param  hostlonglong
 * @return uint64_t
 */
static inline uint64_t usys_htobe64(uint64_t hostlonglong) {
    return htobe64(hostlonglong);
}

/**
 * @fn     uint32_t usys_rand()
 * @brief  Generates pseudo-random integral number in the range
 *
 * @return An integer value between 0 and RAND_MAX.
 */
static inline uint32_t usys_rand() {
    return rand();
}

/**
 * @fn     int usys_inet_pton_ipv4(const char*, uint32_t*)
 * @brief  This function converts the character string src into a network
 *         address structure in the IPv4 af address family, then copies the
 *         network address structure to dst.
 *
 * @param  data
 * @param  ip_addr
 * @return 1 if success
 *         0 if src does not contain a character string representing a valid
 *         network address.
 */
static inline int usys_inet_pton_ipv4(const char* data, uint32_t* ip_addr) {
    return inet_pton(AF_INET, data, ip_addr);
}

/**
 * @fn     int usys_inet_pton_ipv6(const char*, uint8_t*)
 * @brief  This function converts the character string src into a network
 *         address structure in the IPv6 af address family, then copies the
 *         network address structure to dst.
 *
 * @param  data
 * @param  ip_addr
 * @return 1 if success
 *         0 if src does not contain a character string representing a valid
 *         network address.
 */
static inline int usys_inet_pton_ipv6(const char* data, uint8_t* ip_addr) {
    return inet_pton(AF_INET6, data, ip_addr);
}


/**
 * @fn     const char usys_inet_ntop_ipv4*(const char*, char*, socklen_t)
 * @brief  This function converts the network address structure src in the IPv4
 *         af address family into a character string.
 *
 * @param  data
 * @param  ip_addr
 * @param  size
 * @return On success, inet_ntop() returns a non-null pointer to dst.
 *         NULL is returned if there was an error
 */
static inline const char* usys_inet_ntop_ipv4(const char* data, char* ip_addr, socklen_t size) {
    return inet_ntop(AF_INET, data, ip_addr, size);
}

/**
 * @fn     const char usys_inet_ntop_ipv6*(const char*, char*, socklen_t)
 * @brief  This function converts the network address structure src in the IPv6
 *         af address family into a character string.
 *
 * @param  data
 * @param  ip_addr
 * @param  size
 * @return On success, inet_ntop() returns a non-null pointer to dst.
 *         NULL is returned if there was an error
 */
static inline const char* usys_inet_ntop_ipv6(const char* data, char* ip_addr, socklen_t size) {
    return inet_ntop(AF_INET6, data, ip_addr, size);
}

/**
 * @fn     double usys_sqrt(double)
 * @brief  This function calculates the nonnegative value of the square root of x.
 *
 * @param  num
 * @return double
 */
static inline double usys_sqrt(double num) {
    return sqrt(num);
}


/**
 * @fn     double usys_cos(double)
 * @brief  function calculates the cosine of num.
 *
 * @param  num
 * @return double
 */
static inline double usys_cos(double num) {
    return cos(num);
}


/**
 * @fn     double usys_sin(double)
 * @brief  function calculates the sine of num.
 *
 * @param  num
 * @return double
 */
static inline double usys_sin(double num) {
    return sin(num);
}

/**
 * @fn    double usys_atan(double)
 * @brief function calculates the arc tangent of num.
 *
 * @param  num
 * @return double
 */
static inline double usys_atan(double num) {
    return atan(num);
}

/**
 * @fn     uint32_t sleep(uint32_t)
 * @brief  causes the calling thread to sleep either until the
 *         number of real-time seconds specified in seconds have elapsed or
 *         until a signal arrives which is not ignored.
 *
 * @param  sec
 * @return Zero if the requested time has elapsed, or the number of seconds
 *         left to sleep, if the call was interrupted by a signal handler.
 */
static inline uint32_t usys_sleep(uint32_t sec) {
    return sleep(sec);
}

/**
 * @fn     int usleep(uint32_t)
 * @brief  The usleep() function suspends execution of the calling thread
 *         for (at least) usec microseconds.  The sleep may be lengthened
 *         slightly by any system activity or by the time spent processing
 *         the call or by the granularity of system timers.
 *
 * @param  usec
 * @return The usleep() function returns 0 on success.  On error, -1 is
 *         returned, with errno set to indicate the error.
 */
static inline int usys_usleep(useconds_t usec) {
    return usleep(usec);
}

#ifdef __cplusplus
}
#endif

#endif /* USYS_API_H_ */
