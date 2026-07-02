#ifndef AISG_TEST_STUB_USYS_LOG_H_
#define AISG_TEST_STUB_USYS_LOG_H_

#define LOG_TRACE 0
#define LOG_DEBUG 1
#define LOG_INFO  2
#define USYS_LOG_TRACE 0
#define USYS_LOG_DEBUG 1
#define USYS_LOG_INFO  2

static inline void usys_log_set_level(int level) { (void)level; }
static inline void usys_log_set_service(const char *svc) { (void)svc; }

#endif /* AISG_TEST_STUB_USYS_LOG_H_ */
