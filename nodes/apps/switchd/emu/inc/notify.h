#ifndef NOTIFY_H
#define NOTIFY_H

#include "types.h"

int notify_send_alarm(const EmuConfig *cfg,
                      const EmuAlarm *alarm,
                      const EmuSwitchInfo *info);

#endif /* NOTIFY_H */
