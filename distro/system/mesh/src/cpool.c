/*
 * Connection pool thread.
 *
 */

#include "cpool.h"

/*
 * create work.
 * add work.
 * delete work - 
 * get work - get work packet out
 * create cpool.- create the connection thread.
 * destroy cpool - destroy the connection thread.
 */

extern CPool *cpoolTX, *cpoolRX;
extern int maxCpoolTh;

/*
 * find_assigned_thread --
 *
 */
CPool *find_assigned_thread(unsigned char *ipAddress, int flag) {

  int i;
  CPool *ptr;

  if (ipAddress == NULL) {
    return NULL;
  }

  if (flag == TX) {
    ptr = &cpoolTX[0];
  } else if (flag == RX) {
    ptr = &cpoolRX[0];
  } else {
    return NULL;
  }

  for (i=0; i<maxCpoolTh; i++) {

    if (ptr[i].state != THREAD_READY)
      continue;

    if (strcmp(ptr[i].clientIP, ipAddress) == 0) {
      return &ptr[i];
    }
  }

  return NULL;
}

/*
 * create_work --
 *
 */

static CPoolWork *create_work(Packet data, thread_func_t pre, void *preArgs,
			      thread_func_t post, void *postArgs) {

  CPoolWork *work;

  if (pre == NULL || post == NULL)
    return NULL;

  work = (CPoolWork *)calloc(1, sizeof(CPoolWork));
  if (!work) {
    log_error("Error allocating memory: %d", sizeof(CPoolWork));
    return NULL;
  }

  work->preFunc  = pre;
  work->postFunc = post;
  work->preArgs  = preArgs;
  work->postArgs = postArgs;
  
  work->data = data;
  work->next = NULL;

  return work;
}

/*
 * destroy_work --
 *
 */

static void destroy_work(CPoolWork *work) {

  if (!work) {
    return;
  }

  free(work->data);
  free(work);
}

/*
 * get_work -- get work from specific FIFO queues (TX or RX) a in FIFO manner.
 *
 */
static CPoolWork *get_work(CPool *cp) {

  CPoolWork *work, *first, *last;
  
  if (cp == NULL) {
    return NULL;
  }
    
  if (cp->tddFlag == TX_QUEUE) {
    first = cp->firstTX;
    last = cp->lastTX;
  } else if (cp->tddFlag == RX_QUEUE) {
    first = cp->firstRX;
    last = cp->lastRX;
  } else {
    return NULL;
  }
     
  work = first;
  if (work == NULL) {
    return NULL;
  }

  if (work->next == NULL) { /* Removed the only item from queue. */
    first = NULL;
    last = NULL;
  } else {
    first = work->next;
  }

  return work;
}

/*
 * Connection pool thread worker. 
 *
 */

static void *cpool_worker(void *arg) {

  CPool *cp = arg;
  CPoolWork *work=NULL;

  CPoolWork *first, *last;
  pthread_mutex_t *mutex;
  pthread_cond_t  *cond;
  int stop, exit, other, ret;

  pid_t tID;

  if (cp->tddFlag == RX) {
    first  = cp->firstRX;
    last  = cp->lastRX;
    mutex = &(cp->rxMutex);
    cond = &(cp->rxCondWait);
    stop  = cp->stopRX;
    exit  = cp->exitRX;
    other = cp->exitTX;
  } else if (cp->tddFlag == TX) {
    first  = cp->firstTX;
    last  = cp->lastTX;
    mutex = &(cp->txMutex);
    cond = &(cp->txCondWait);
    stop  = cp->stopTX;
    exit  = cp->exitTX;
    other = cp->exitRX;
  }

  tID = syscall(__NR_gettid);

  /* Thread wait for the parent to assign it the SSL connection handler and ask
   * for it to execute on it. It can be long wait. Thread is on stand-by mode.
   * Once it has the lock, it will let it go only when it is asked to "stop"
   */
  log_debug("TID-%s: standby mode. Waiting for SSL connection handler.", tID);
  pthread_mutex_lock(&cp->active);

  while (TRUE) {

    log_debug("TID-%d: Acquring lock", tID);
    
    pthread_mutex_lock(mutex);
    
    /* Exit the thread if told. */
    if (exit) {
      break; 
    }

    /* Don't process any packet. */
    if (stop) {
      pthread_mutex_unlock(&cp->active); /* back in stand-by mode. */
      pthread_mutex_unlock(mutex);
      continue;
    }

    /* There is no work in the queue, conditional wait */
    if (first == NULL) {
      log_debug("TID-%d: Waiting on work cond", tID);
      pthread_cond_wait(cond, mutex);
    }

    /* We have some work to do. */
    work = get_work(cp);

    if (work != NULL) {
      if (work->preFunc) {
	work->preFunc(work->preArgs);
      }

      /* Send the actual packet over SSL/TLS to remote client. */
      while ((ret = mbedtls_ssl_write(cp->ssl, work, strlen(work))) <= 0) {
        if( ret != MBEDTLS_ERR_SSL_WANT_READ &&
	    ret != MBEDTLS_ERR_SSL_WANT_WRITE ) {
	  log_error("ssl_write failed with error code: %d", ret);
	  goto post_func;
        }
      }

    post_func:
      if (work->postFunc) {
	work->postFunc(work->postArgs);
      }
      destroy_work(work);
    }
    
    pthread_mutex_unlock(mutex); /* pre/post func could delay the unlock. */
  }

  /* Thread is done. */

  /* check if we need to close the socket connection. */
  if (other) {
    /* Close the secure socket connection. */
    //    close_connection(cp);
  }
  
  return NULL;
}

/*
 * add_work -- Add work to the thread (cp) queue assigned to the remote
 *             client at destIP.
 *
 */
int add_work(unsigned char *destIP, Packet data, thread_func_t pre,
	     void *preArgs, thread_func_t post, void *postArgs, int flag) {


  CPoolWork *work;
  pthread_mutex_t *mutex;
  CPoolWork *fPtr, *lPtr;
  CPool *cp;

  if (destIP == NULL) {
    return FALSE;
  }

  cp = find_assigned_thread(destIP, flag);
  
  if (cp == NULL) {
    return FALSE;
  }

  work = create_work(data, pre, preArgs, post, postArgs);

  if (work == FALSE) {
    return FALSE;
  }

  if (cp->tddFlag == TX) {
    mutex = &(cp->txMutex);
    fPtr = cp->firstTX;
    lPtr = cp->lastTX;
  } else {
    mutex = &(cp->rxMutex);
    fPtr = cp->firstRX;
    lPtr = cp->lastRX;
  }

  pthread_mutex_lock(mutex);

  if (fPtr == NULL) {
    fPtr = work; /* XXX check the ptr. */
    lPtr = fPtr;
  } else {
    lPtr->next = work;
    lPtr = work;
  }

  pthread_cond_broadcast(&(cp->txCondWait));
  pthread_mutex_unlock(mutex);
  
  return TRUE;
}

/*
 * create_cpool -- create x-many connection pool threads and detach them.
 *
 */
int create_cpool(pthread_t *tArray, CPool *cpArray, int num, int flag) {

  int i, count, ret;
  pthread_t *thread; // XXX should be arrayyyyyy.
  
  if (num == 0) {
    return FALSE; /* Default do nothing. */
  }

  /* Loop through. */
  for (i=0; i<num; i++) {
    CPool *cp;

    cp = &cpArray[i];
    thread = &tArray[i];

    /* Initialize mutexs. */
    pthread_mutex_init(&(cp->txMutex), NULL);
    pthread_mutex_init(&(cp->txDataFlag), NULL);
    pthread_mutex_init(&(cp->rxMutex), NULL);
    pthread_mutex_init(&(cp->rxDataFlag), NULL);
    pthread_mutex_init(&(cp->active), NULL);
    pthread_cond_init(&(cp->txCondWait), NULL);
    pthread_cond_init(&(cp->rxCondWait), NULL);

    cp->firstTX = NULL;
    cp->lastTX  = NULL;
    cp->firstRX = NULL;
    cp->lastRX  = NULL;

    cp->ssl = NULL;
    cp->stopTX = FALSE;
    cp->stopRX = FALSE;
    cp->exitTX = FALSE;
    cp->exitRX = FALSE;
    cp->clientIP = calloc(sizeof(unsigned char), 32);

    if (flag != TX || flag != RX) {
      cp->tddFlag = TX; /* default is TX queue. */
    } else {
      cp->tddFlag = flag;
    }

    /* Put thread in the stand-by mode. Thread will be put into ready state
     * once/if remote client shows up.
     */
    pthread_mutex_lock(&cp->active);
    cp->state = THREAD_STANDBY;

    /* Now create the real thread and detach!. */
    ret = pthread_create(thread, NULL, cpool_worker, cp);
    if (ret != 0) {
      log_error("Error creating Connector pool thread: %d with error: %d",
		i, ret);
      goto cleanup;
    }

    ret = pthread_detach(*thread);
    if (ret != 0) {
      log_error("Error creating Connector pool thread: %d with error: %d",
		i, ret);
      goto cleanup;
    }

    count++;
  }

  log_debug("Successfully created %d connection thread pools", num);
  return TRUE;

 cleanup:
  
  /* KILL THREADS XXXX and return. */
  for (i=0; i<count; i++) {
    destroy_cpool(&cpArray[i]);
  }
  
  return FALSE;
}

/*
 * destroy_cpool -- gracefully destroy connection pool thread worker.
 *
 */

void destroy_cpool(CPool *cp) {

  CPoolWork *ptr1, *ptr2;
  pthread_mutex_t *mutex;
  int stop, exit, other;
  
  /* Destroying a connection thread means following steps:
   *
   * 1. drop/free all packets on the queues.
   * 2. Close the SSL/TLS connection.
   * 3. Tell the thread to STOP.
   * 4. Destroy all mutex.
   * 5. Free the allocated memory.
   */
  
  if (cp == NULL) {
    return;
  }

  if (cp->tddFlag == RX) {
    ptr1  = cp->firstRX;
    ptr2  = cp->lastRX;
    mutex = &(cp->rxMutex);
    stop  = cp->stopRX;
    exit  = cp->exitRX;
    other = cp->exitTX;
  } else if (cp->tddFlag == TX) {
    ptr1  = cp->firstTX;
    ptr2  = cp->lastTX;
    mutex = &(cp->txMutex);
    stop  = cp->stopTX;
    exit  = cp->exitTX;
    other = cp->exitRX;
  }
      
  pthread_mutex_lock(mutex);

  while (ptr1) {
    ptr2 = ptr1->next;
    destroy_work(ptr1);
    ptr1 = ptr2;
  }

  /* Set the flags. */
  if (cp->tddFlag == RX) {
    cp->stopRX = TRUE;
    cp->exitRX = TRUE;
  } else if (cp->tddFlag == TX) {
    cp->stopTX = TRUE;
    cp->exitTX = TRUE;
  }

  /* If we are the last thread on this connection, close the socket. */
  if (other == TRUE) {
    /* XXXX */
  }
  
  pthread_mutex_unlock(mutex);

  /* Wait for thread to exit. */
  /* XXX */

  if (cp->tddFlag == RX) {
    pthread_mutex_destroy(&(cp->rxMutex));
    pthread_mutex_destroy(&(cp->rxDataFlag));
  } else {
    pthread_mutex_destroy(&(cp->txMutex));
    pthread_mutex_destroy(&(cp->txDataFlag));
  }

  free(cp);
}

/*
 * assign_thread -- assign cpool 'stand-by' thread to the remote client.
 *
 */

int assign_thread(unsigned char *clientIP, CPool *cpArrayTX,
		  CPool *cpArrayRX, int num, mbedtls_ssl_context *ssl) {

  int i;

  if (num == 0) {
    return FALSE; /* Default do nothing. */
  }

  /* Loop through. */
  for (i=0; i<num; i++) {

    CPool *cpTX, *cpRX;

    cpTX = &cpArrayTX[i];
    cpTX = &cpArrayRX[i];

    if (cpTX->state == THREAD_STANDBY) {

      cpTX->state = THREAD_READY;
      cpRX->state = THREAD_READY;

      cpTX->ssl = ssl;
      cpRX->ssl = ssl;

      sprintf(cpTX->clientIP, "%d.%d.%d.%d", clientIP[0], clientIP[1],
	      clientIP[2], clientIP[3]);
      sprintf(cpRX->clientIP, "%d.%d.%d.%d", clientIP[0], clientIP[1],
	      clientIP[2], clientIP[3]);

      pthread_mutex_unlock(&cpTX->active);
      pthread_mutex_unlock(&cpRX->active);

      return i;
    }
  }

  log_debug("No available thread to assign");

  return FALSE;
}

/* unassign. de-active, etc. etc. */

