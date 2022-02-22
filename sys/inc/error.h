#ifndef USYS_SYS_ERROR_H
#define USYS_SYS_ERROR_H

#ifdef __cplusplus
extern "C" {
#endif

#include "sys_types.h"
#define   USYS_ERROR_CODE_IDX_BASE      (1000) 
#define   USYS_ERROR_CODE_IDX(code)     ((code) - USYS_ERROR_CODE_IDX_BASE)

/*
*  
* Error codes = USYS_ERROR_CODE_IDX_BASE + USysErrorCodeIdx 
*/
typedef enum {
    /* Sample error code */
    USYS_ERR_SOCK_CREATION = (USYS_ERROR_CODE_IDX_BASE+1),
    USYS_ERR_SOCK_CONNECT,
    USYS_ERR_SOCK_SEND,
    USYS_ERR_SOCK_RECV,
    USYS_ERR_MAX_IDX
} USysErrorCodeIdx;


char* usys_error(int err);

#ifdef __cplusplus
extern "C" {
#endif
#endif /* USYS_SYS_ERROR_H */