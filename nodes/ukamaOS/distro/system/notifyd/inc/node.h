/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_NODE_H_
#define INC_NODE_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

typedef enum {
    TNODE = 1,
    HNODE = 2,
    ANODE = 3,
} NodeType;

typedef enum {
    MOD_COM  = 1,
    MOD_TRX,
    MOD_CNTRL,
    MOD_RFFE,
    MOD_MASK
} ModuleType;

#ifdef __cplusplus
}
#endif
#endif /* INC_NODE_H_ */
