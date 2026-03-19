#ifndef PERSISTENCE_H
#define PERSISTENCE_H

#include "types.h"

int persistence_save(const EmuModel *model, const char *path);
int persistence_load(EmuModel *model, const char *path);

#endif /* PERSISTENCE_H */
