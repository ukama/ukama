/*
 * Connection Pool threads.
 *
 */

#include <pthread.h>
#include <stdlib.h>

#include "ssl.h"

#define TX 1
#define RX 2

#define TX_QUEUE 1
#define RX_QUEUE 2

#define FALSE 0
#define TRUE  1

typedef void (*thread_func_t)(void *arg);
typedef char* Packet;

/*
 * Link list of object connection thread waiting to process.
 */
typedef struct cpool_work_t {
  
  thread_func_t preFunc;  /* data packet pre-process function */
  thread_func_t postFunc; /* data packet post-processing function. */
  void          *preArgs; /* Args for pre-function. */
  void          *postArgs;/* Args for post-funciion. */

  Packet data;             /* Data packet to process. */

  struct cpool_work_t *next; /* Link to next item in the queue. */
  
}CPoolWork;

/*
 * Connection threads. 
 * Each thread is handling single client. It maintains two Queues: 1. TX, for
 * outgoing packets and  2.) RX, for all incoming packets.
 */
typedef struct cpool_t {

  /* TX */
  CPoolWork       *firstTX;   /* First item in the TX queue. */
  CPoolWork       *lastTX;    /* Pointer to last item in the TX queue. */
  pthread_mutex_t txMutex;    /* Mutex for insert and remove */
  pthread_mutex_t txDataFlag; /* Flag to identify there are data packets. */
  int             stopTX;     /* Stop transmiting the packets. */
  int             exitTX;     /* if TX thread is to exit or exited. */
  
  /* RX */
  CPoolWork       *firstRX;   /* First item in the RX queue. */
  CPoolWork       *lastRX;    /* Pointer to last item in the RX queue. */
  pthread_mutex_t rxMutex;    /* Mutex for insert and remove */
  pthread_mutex_t rxDataFlag; /* Flag to identify there are data packets. */
  int             stopRX;     /* Stop receiving of data packets. */
  int             exitRX;     /* if RX thread is to exit or exited. */

  /* Misc. */
  pthread_mutex_t     tddMutex; /* Mutex for the SSL connection handler. */
  mbedtls_ssl_context *ssl;     /* SSL connection handler. */
  int tddFlag;                  /* Thread is TX or RX. */
}CPool;

