#ifndef OID_MAP_H
#define OID_MAP_H

#include <stddef.h>

#include "types.h"

int oid_get_string(EmuModel *model, const char *oid, char *buf, size_t len);
int oid_get_int(EmuModel *model, const char *oid, int *value);
int oid_set_int(EmuModel *model, const char *oid, int value);
int oid_get_next(EmuModel *model, const char *oid, char *nextOid, size_t nextLen);

#endif /* OID_MAP_H */
