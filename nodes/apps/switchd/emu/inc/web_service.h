#ifndef WEB_SERVICE_H
#define WEB_SERVICE_H

#include <stddef.h>

#include "types.h"

int web_service_handle(EmuModel *model,
                       const char *method,
                       const char *path,
                       const char *body,
                       char *out,
                       size_t outLen,
                       int *status);

#endif /* WEB_SERVICE_H */
