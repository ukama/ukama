/* Test-only platform adapters; production uses libusys. */
#include "usys_file.h"
#include "usys_log.h"
void usys_log_set_service(const char *service) { (void)service; }
void usys_log_set_level(int level) { (void)level; }
void usys_log_debug(const char *format, ...) { (void)format; }
void usys_log_info(const char *format, ...) { (void)format; }
void usys_log_warn(const char *format, ...) { (void)format; }
void usys_log_error(const char *format, ...) { (void)format; }
int usys_find_service_port(char *service) { (void)service; return 0; }
