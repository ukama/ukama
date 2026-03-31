#ifndef SWITCHEMU_LOG_H
#define SWITCHEMU_LOG_H
void log_set_level(int level);
void log_debug(const char *fmt, ...);
void log_info(const char *fmt, ...);
void log_warn(const char *fmt, ...);
void log_error(const char *fmt, ...);
#endif
