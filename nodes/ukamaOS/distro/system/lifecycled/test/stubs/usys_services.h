#pragma once

#define SERVICE_CONFIG  "config"
#define SERVICE_STARTER "starter"
#define SERVICE_NOTIFY  "notify"

int usys_find_service_port(const char *service);
