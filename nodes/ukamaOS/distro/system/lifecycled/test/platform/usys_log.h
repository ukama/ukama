#pragma once

#define USYS_LOG_DEBUG 0
#define USYS_LOG_INFO  1
#define USYS_LOG_WARN  2
#define USYS_LOG_ERROR 3

void usys_log_set_service(const char *service);
void usys_log_set_level(int level);
void usys_log_debug(const char *format, ...);
void usys_log_info(const char *format, ...);
void usys_log_warn(const char *format, ...);
void usys_log_error(const char *format, ...);

