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
  int stop, exit, other;

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

  while (TRUE) {

    log_debug("TID-%d: Acquring lock", tID);
    
    pthread_mutex_lock(mutex);
    
    /* Exit the thread if told. */
    if (exit) {
      break; 
    }

    /* Don't process any packet. */
    if (stop) {
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
 * add_work -- Add work to the thread (cp) queue.
 *
 */
int add_work(CPool *cp, Packet data, thread_func_t pre, void *preArgs,
	     thread_func_t post, void *postArgs) {

  CPoolWork *work;
  pthread_mutex_t *mutex;
  CPoolWork *fPtr, *lPtr;
  
  if (!cp) {
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

  /* XXXX
  pthread_cond_broadcast(&(tp->work_cond));
  */
  
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
    pthread_mutex_init(&(cp->tddMutex), NULL);
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

    if (flag != TX || flag != RX) {
      cp->tddFlag = TX; /* default is TX queue. */
    } else {
      cp->tddFlag = flag;
    }

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
