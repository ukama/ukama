#ifndef SNMP_AGENT_H
#define SNMP_AGENT_H

#include "types.h"

int snmp_agent_start(EmuModel *model);
void snmp_agent_stop(EmuModel *model);

#endif /* SNMP_AGENT_H */
