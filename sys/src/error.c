#include "error.h"

char* usysErrorCodes[USYS_ERR_MAX_IDX] = {
    [USYS_ERROR_CODE_IDX(USYS_ERR_SOCK_CREATION)] = "Failed to create socket",
    [USYS_ERROR_CODE_IDX(USYS_ERR_SOCK_CONNECT)] = "Failed to connect to socket",
    [USYS_ERROR_CODE_IDX(USYS_ERR_SOCK_SEND)] = "Failed to send on socket",
    [USYS_ERROR_CODE_IDX(USYS_ERR_SOCK_RECV)] = "Failed to read from socket"
}

/**
 * @brief Read error description 
 * 
 * @param err 
 * @return char* 
 */
char* usys_error(int err) {
    if (err < USYS_ERROR_CODE_IDX_BASE)) {
        return strerror(err);
    } else {
        if (err < USYS_ERR_MAX_IDX) {
         return usysErrorCodes[err];
        }
    }

}