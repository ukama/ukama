#ifndef TFTP_SERVER_H
#define TFTP_SERVER_H

#include "types.h"

int tftp_server_start(EmuModel *model);
void tftp_server_stop(EmuModel *model);

#endif /* TFTP_SERVER_H */
