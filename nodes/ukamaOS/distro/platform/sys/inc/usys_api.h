/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
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
 * @fn    void usys_qexit()
 * @brief terminates the calling process "immediately".  Any open
 *        file descriptors belonging to the process are closed.  Any
 *        children of the process are inherited by init(1) (or by the
 *        nearest "subreaper" process as defined through the use of the
 *        prctl(2) PR_SET_CHILD_SUBREAPER operation).  The process's parent
 *        is sent a SIGCHLD signal.
 * @param status
 *
 */
static inline void usys_qexit(int status) {
    _exit(status);
}

/**
 * @fn    void usys_exit()
 * @brief The exit() function causes normal process termination and the
 *        least significant byte of status (i.e., status & 0xFF) is
 *        returned to the parent
 * @param status
 */
static inline void usys_exit(int status) {
    exit(status);
}

/**
 * @fn     FILE usys_fopen*(const char*, const char*)
 * @brief  Opens the filename pointed to, by filename using the given mode.
 *
 * @param  path
 * @param  mode
 * @return On Success FILE pointer otherwise NULL
 */
static inline FILE *usys_fopen(const char *path, const char *mode) {
    return fopen(path, mode);
}

/**
 * @fn     int usys_fclose(FILE*)
 * @brief  closes the stream pointed by FILE*
 *
 * @param  fp
 * @return On Success 0 otherwise EOF
 */
static inline int usys_fclose(FILE *fp) {
    return fclose(fp);
}

/**
 * @fn     usys_fread
 *
 * @brief  reads nmemb items of data, each size bytes long,
 *         from the stream pointed to by stream, storing them at the
 *         location given by ptr.
 *
 */
static inline int usys_fread(void *ptr, size_t size, size_t nmemb,
                             FILE *stream) {
    return fread(ptr, size, nmemb, stream);
}

/**
 * @fn     int usys_fwrite(const void*, size_t, size_t, FILE*)
 * @brief  Writes an array of count elements,each one with a size of size bytes,
 *         from the block of memory pointed by ptr to the current position
 *         in the stream.
 *
 * @param  ptr
 * @param  size
 * @param  nmemb
 * @param  stream
 * @return The total number of elements successfully written is returned.
 */
static inline int usys_fwrite(const void *ptr, size_t size, size_t nmemb,
                              FILE *stream) {
    return fwrite(ptr, size, nmemb, stream);
}

/**
 * @fn     int usys_fseek(FILE*, long int, int)
 * @brief  Sets the position indicator associated with the stream to a
 *         new position.
 *
 * @param  stream
 * @param  offset
 * @param  origin
 * @return On Success, the function returns zero.
 *         On Failure, it returns non-zero value.
 */
static inline int usys_fseek(FILE *stream, long int offset, int origin) {
    return fseek(stream, offset, origin);
}

/**
 * @fn     int usys_ftell(FILE*)
 * @brief  This function returns the current file position of the stream stream.
 *
 * @param  stream
 * @return On success, the current value of the position indicator is returned.
 *         On failure, -1L is returned.
 */
static inline int usys_ftell(FILE *stream) {
    return ftell(stream);
}

/**
 * @fn     int usys_fputs(char*, FILE*)
 * @brief  The function fputs writes the string s to the stream stream.
 *         The terminating null character is not written.This function does not
 *         add a newline character, either.
 *         It outputs only the characters in the string.
 *
 * @param  str
 * @param  stream
 * @return On success , returns a non-negative value.
 *         On failure, it returns EOF.
 */
static inline int usys_fputs(char* str, FILE *stream) {
    return fputs(str, stream);
}

/**
 * @fn     int usys_rename(const char*, const char*)
 * @brief  renames a file.
 *
 * @param  oldpath
 * @param  newpath
 * @return On success, zero is returned.
 *         On error, -1 is returned,
 */
static inline int usys_rename(const char *oldpath, const char *newpath) {
    return rename(oldpath, newpath);
}

/**
 * @fn     int usys_remove(const char*)
 * @brief  deletes the given filename so that it is no longer accessible.
 *
 * @param  filename
 * @return On success, zero is returned.
 *         On error, -1 is returned.
 */
static inline int usys_remove(const char *filename) {
    return remove(filename);
}

/**
 * @fn    int usys_access(const char*, int)
 * @brief The access function checks to see whether the file named by
 *        filename can be accessed in the way specified by the how argument.
 *
 * @param  filename
 * @param  how
 * @return On Success, the function returns zero.
 *         On Failure, it returns -1.
 */
static inline int usys_access(const char *filename, int how) {
    return access(filename, how);
}

/**
 * @fn     int usys_open(const char*, int, mode_t)
 * @brief  opens the file specified by pathname.
 *
 * @param  pathname
 * @param  flags
 * @param  mode
 * @return On success return the new file descriptor.
 *         On error, -1
 */
static inline int usys_open(const char *pathname, int flags, mode_t mode) {
    return open(pathname, flags, mode);
};

/**
 * @fn     int usys_close(int)
 * @brief  closes a file descriptor
 *
 * @param  fd
 * @return On success return 0.
 *         On error, -1
 */
static inline int usys_close(int fd) {
    return close(fd);
}

/**
 * @fn     int usys_fsync(int)
 * @brief  flushes the data of file referred to by the file
 *         descriptor fd to the disk device
 *
 * @param  fd
 * @return On success, return zero.
 *         On error, -1
 */
static inline int usys_fsync(int fd) {
    return fsync(fd);
}

/**
 * @fn     int usys_stat(const char*, struct stat*)
 * @brief  return information about a file, in the buffer
 *         pointed to by statbuf.
 *
 * @param  pathname
 * @param  statbuf
 * @return On success, return zero.
 *         On error, -1
 */
static inline int usys_stat(const char *pathname, struct stat *statbuf) {
    return stat(pathname, statbuf);
}

/**
 * @fn     int usys_lstat(const char*, struct stat*)
 * @brief  return information about a link, in the buffer
 *         pointed to by statbuf.
 *
 * @param  pathname
 * @param  statbuf
 * @return On success, return zero.
 *         On error, -1
 */
static inline int usys_lstat(const char *pathname, struct stat *statbuf) {
    return lstat(pathname, statbuf);
}

/**
 * @fn off_t usys_lseek(int, off_t, int)
 * @brief  repositions the file offset of the open file description
 *         associated with the file descriptor fd to the argument offset
 *         according to the directive whence
 *
 * @param fd
 * @param offset
 * @param whence
 * @return On success resulting offset location as measured in bytes
 *         from the beginning of the file.
 *         On error, -1
 */
static inline off_t usys_lseek(int fd, off_t offset, int whence) {
    return lseek(fd, offset, whence);
}

/**
 * @fn     ssize_t usys_read(int, void*, size_t)
 * @brief  attempts to read up to count bytes from file descriptor fd
 *         into the buffer starting at buf.
 *
 * @param  fd
 * @param  buf
 * @param  count
 * @return On success, the number of bytes read is returned.
 *         On error, -1 is returned
 */
static inline ssize_t usys_read(int fd, void *buf, size_t count) {
    return read(fd, buf, count);
}

/**
 * @fn     ssize_t usys_write(int, const void*, size_t)
 * @brief  writes up to count bytes from the buffer starting at buf
 *         to the file referred to by the file descriptor fd.
 *
 * @param  fd
 * @param  buf
 * @param  count
 * @return On success, the number of bytes written is returned.
 *         On error, -1 is returned.
 */
static inline ssize_t usys_write(int fd, const void *buf, size_t count) {
    return write(fd, buf, count);
}

/**
 * @fn     ssize_t usys_readlink(const char*, char*, size_t)
 * @brief  places the contents of the symbolic link pathname in
 *         the buffer buf, which has size bufsiz.
 *
 * @param  pathname
 * @param  buf
 * @param  bufsiz
 * @return On success, these calls return the number of bytes placed in buf.
 *         On error, -1 is returned
 */
static inline ssize_t usys_readlink(const char *pathname, char *buf,
                                    size_t bufsiz) {
    return readlink(pathname, buf, bufsiz);
}

/**
 * @fn     char usys_fgets*(char*, int, FILE*)
 * @brief  Reads characters from stream and stores them as a C string into str
 *         until (num-1) characters have been read or either a newline or
 *         the end-of-file is reached, whichever happens first.
 *
 * @param  s
 * @param  size
 * @param  stream
 * @return On success, the function returns str.
 *         If the end-of-file is encountered while attempting to read a
 *         character, the eof indicator is set (feof). If this happens before
 *         any characters could be read, the pointer returned is a null pointer.
 *         If a read error occurs, the error indicator (ferror)
 *         is set and a null pointer is also returned
 */
static inline char *usys_fgets(char *s, int size, FILE *stream) {
    return fgets(s, size, stream);
}

/**
 * @def    usys_snprintf
 * @brief  Wraps snprintf function.
 *         The content is stored as a C string in the buffer pointed by str
 *         with size as max length.
 */
#define usys_snprintf(str, size, format, ...) \
    snprintf(str, size, format, ##__VA_ARGS__);

/**
 * @def   usys_fprintf
 * @brief Wraps fprintf function.
 *        Writes the C string pointed by format to the stream.
 */
#define usys_fprintf(stream, format, ...) \
    { fprintf(stream, format, ##__VA_ARGS__); }

/**
 * @def   usys_printf
 * @brief Wraps printf function
 *        Writes the C string pointed by format to the standard output (stdout).
 */
#define usys_printf(format, ...) \
    { printf(format, ##__VA_ARGS__); }

/**
 * @def   usys_sprintf
 * @brief Wraps sprintf function.
 *        The content is stored as a C string in the buffer pointed by str
 */
#define usys_sprintf(str, format, ...) \
    { sprintf(str, format, ##__VA_ARGS__); }

/**
 * @def   usys_sscanf
 * @brief Reads data from str and stores them according to parameter
 *        format into the locations given by the additional arguments.
 */
#define usys_sscanf(str, format, ...) \
    { sscanf(str, format, ##__VA_ARGS__); }

/**
 * @fn    struct timeval usys_gettimeofday()
 * @brief function gets the system’s clock time
 *
 * @param  tv
 * @return void
 *
 */
static inline void usys_gettimeofday(struct timeval *tv) {
    gettimeofday(tv, NULL);
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
static inline size_t usys_strftime(char *s, size_t max, const char *format,
                                   const struct tm *tm) {
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
 *         with the values that correspond to the local time representation
 *         of timer.
 */
static inline struct tm *usys_localtime(const time_t *timep) {
    return localtime(timep);
}

/**
 * @fn     uint16_t usys_ntohs(uint16_t)
 * @brief  Converts the unsigned short integer netshort from network byte order
 *         to host byte order.
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
static inline uint32_t usys_inet_addr(const char *data) {
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
static inline int usys_inet_pton_ipv4(const char *data, uint32_t *ip_addr) {
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
static inline int usys_inet_pton_ipv6(const char *data, uint8_t *ip_addr) {
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
inline static const char* usys_inet_ntop_ipv4(const char *data, char *ip_addr,
                                              socklen_t size) {
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
inline static const char* usys_inet_ntop_ipv6(const char *data, char *ip_addr,
                                              socklen_t size) {
    return inet_ntop(AF_INET6, data, ip_addr, size);
}

/**
 * @fn     double usys_sqrt(double)
 * @brief  This function calculates the nonnegative value of the square
 *         root of x.
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

/**
 * @fn     int usys_puts(char*)
 * @brief  The puts function writes the string s to the stream stdout
 *         followed by a newline. The terminating null character of the
 *         string is not written.
 *
 * @param  str
 * @return On success, non-negative value is returned.
 *         On error, the function returns EOF
 */
static inline int usys_puts(char* str) {
  return puts(str);
}

/**
 * @fn     int usys_putchar(char)
 * @brief  Writes a character to the standard output stdout
 *
 * @param  c
 * @return On success, the character written is returned.
 *         On failure, EOF is returned and the error indicator (ferror) is set.
 */
static inline int usys_putchar(char c) {
  return putchar(c);
}

/**
 * @fn     int usys_fflush(FILE*)
 * @brief  If the stream was open for writing any unwritten data in its output
 *         buffer is written to the file.
 *
 * @param  stream
 * @return On success, zero is returned.
 *         On failure, EOF is returned and the error indicator (ferror) is set.
 */
static inline int usys_fflush(FILE * stream) {
  return fflush(stream);
}

/**
 * @fn     int usys_atoi(const char*)
 * @brief  Parses the str interpreting its content as an integral number.
 *
 * @param  str
 * @return On success, the function returns the converted integral number
 *         as an int value.If the int value returned is out of the range of int,
 *         it causes undefined behavior.
 */
static inline int usys_atoi(const char * str) {
    return atoi(str);
}

/**
 * @fn     double usys_atof(const char*)
 * @brief  Parses the str interpreting its content as an floating point.
 *
 * @param  str
 * @return On success, the function returns the converted floating point number
 *         as a double value.
 *         If no valid conversion could be performed,
 *         the function returns zero (0.0).
 *         If the converted value would be out of the range of representable
 *         values by a double, it causes undefined behavior.
 */
static inline double usys_atof(const char* str) {
    return atof(str);
}

/**
 * @fn     long int usys_strtol(const char*, char**, int)
 * @brief  Parses the C-string str interpreting its content as an integral
 *         number of the specified base, which is returned as a long int value.
 *         If endptr is not a null pointer, the function also sets the value of
 *         endptr to point to the first character after the number.
 *
 * @param  str
 * @param  endptr
 * @param  base
 * @return On success, the function returns the converted integral number
 *         as a long int value.
 *         If no valid conversion could be performed, a zero value is returned.
 *         If the value read is out of the range of representable values by
 *         a long int, the function returns LONG_MAX or LONG_MIN
 */
static inline long int usys_strtol(const char* str, char** endptr, int base){
    return strtol(str, endptr, base);
}

/**
 * @fn     double usys_strtod(const char*, char**)
 * @brief  Parses the C string str interpreting its content as a floating point
 *         number (according to the current locale) and returns its value as a
 *         long double. If endptr is not a null pointer, the function also sets
 *         the value of endptr to point to the first character after the number.
 *
 * @param  str
 * @param  endptr
 * @return On success, the function returns the converted floating
 *         point number as a value of type long double.
 *         If no valid conversion could be performed,
 *         the function returns zero (0.0L).
 *         If the correct value is out of the range of representable values
 *         for the type, a positive or negative HUGE_VALL is returned,
 */
static inline double usys_strtod(const char* str, char** endptr){
    return strtod(str, endptr);
}

/**
 * @fn     double usys_round(double)
 * @brief  Returns the integral value that is nearest to x,
 *         with halfway cases rounded away from zero.
 *
 * @param  x
 * @return The value of x rounded to the nearest integral
 *         (as a floating-point value).
 */
static inline double usys_round(double x) {
    return round(x);
}

/**
 * @fn      double usys_time(time_t*)
 * @brief   Returns the time as the number of seconds since the Epoch,
 *          1970-01-01 00:00:00 +0000 (UTC).
 *          If t is non-NULL, the return value is also stored in the
 *          memory pointed to by t.
 *
 * @param   t
 * @return  On success, the value of time in seconds since
 *          the Epoch is returned.
 *          On error, ((time_t) -1) is returned, and errno is set appropriately.
 */
static inline double usys_time(time_t *t) {
    return time(t);
}

#ifdef __cplusplus
}
#endif

#endif /* USYS_API_H_ */
