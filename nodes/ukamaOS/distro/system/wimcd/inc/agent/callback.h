#ifndef WIMC_AGENT_CALLBACK_H_
#define WIMC_AGENT_CALLBACK_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <ulfius.h>

#include "wimc.h"

int agent_web_service_cb_default(const URequest *request,
                                 UResponse *response,
                                 void *userData);

int agent_web_service_cb_post_capp(const URequest *request,
                                   UResponse *response,
                                   void *userData);

int agent_web_service_cb_ping(const URequest *request,
                              UResponse *response,
                              void *userData);

int agent_web_service_cb_version(const URequest *request,
                                 UResponse *response,
                                 void *userData);

int web_service_cb_not_allowed(const URequest *request,
                               UResponse *response,
                               void *userData);

#ifdef __cplusplus
}
#endif

#endif /* WIMC_AGENT_CALLBACK_H_ */
