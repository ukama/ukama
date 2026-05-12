#ifndef CONFIG_H_
#define CONFIG_H_

#include <stdbool.h>

typedef enum {
    LOOKOUT_APP_MANAGER_STARTERD = 0,
    LOOKOUT_APP_MANAGER_SUPERVISORD
} LookoutAppManager;

typedef enum {
    LOOKOUT_NODE_UNKNOWN = 0,
    LOOKOUT_NODE_TOWER,
    LOOKOUT_NODE_AMPLIFIER,
    LOOKOUT_NODE_CONTROL
} LookoutNodeType;

typedef struct {

    int servicePort;
    int nodedPort;
    int starterdPort;

    char *nodeID;

    LookoutAppManager appManager;
    LookoutNodeType   nodeType;

    bool isTowerNode;
    bool isAmplifierNode;
    bool isControlNode;
} Config;

#endif /* CONFIG_H_ */
