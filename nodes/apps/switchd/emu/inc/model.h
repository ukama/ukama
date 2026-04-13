#ifndef MODEL_H
#define MODEL_H

#include "types.h"

void model_init(EmuModel *model, const EmuConfig *cfg);
void model_recompute(EmuModel *model);
void model_set_reachable(EmuModel *model, int reachable);
int model_set_port_link(EmuModel *model, unsigned int portId, int up);
int model_set_port_admin(EmuModel *model, unsigned int portId, int up);
int model_set_port_poe(EmuModel *model, unsigned int portId, int on);
void model_stage_firmware(EmuModel *model, const char *path,
                          const char *filename);

#endif /* MODEL_H */
