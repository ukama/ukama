#include <stdio.h>
#include <stdarg.h>
#include <time.h>
#include <string.h>
#include "log.h"
static int gLevel = 2;
void log_set_level(int level) { gLevel = level; }
static void vlog_out(int lvl, const char *tag, const char *fmt, va_list ap) {
    time_t now = time(NULL); struct tm tmv; struct tm *tmp = localtime(&now); if (tmp) tmv = *tmp; else memset(&tmv, 0, sizeof(tmv));
    if (lvl > gLevel) return;
    fprintf(stderr, "%04d-%02d-%02d %02d:%02d:%02d [%s] ", tmv.tm_year+1900, tmv.tm_mon+1, tmv.tm_mday, tmv.tm_hour, tmv.tm_min, tmv.tm_sec, tag);
    vfprintf(stderr, fmt, ap); fputc('\n', stderr);
}
void log_debug(const char *fmt, ...) { va_list ap; va_start(ap, fmt); vlog_out(3, "DBG", fmt, ap); va_end(ap);} 
void log_info(const char *fmt, ...) { va_list ap; va_start(ap, fmt); vlog_out(2, "INF", fmt, ap); va_end(ap);} 
void log_warn(const char *fmt, ...) { va_list ap; va_start(ap, fmt); vlog_out(1, "WRN", fmt, ap); va_end(ap);} 
void log_error(const char *fmt, ...) { va_list ap; va_start(ap, fmt); vlog_out(0, "ERR", fmt, ap); va_end(ap);} 
