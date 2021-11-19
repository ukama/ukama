/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Functions related to the cgroup and namespace setup for containers.
 * Note: these are Ukama version of the light-weight containers which are
 * **not** OCI compatiable.
 */

#define _GNU_SOURCE
#include <string.h>
#include <errno.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/capability.h>
#include <sys/prctl.h>
#include <sched.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <sys/mount.h>
#include <sys/syscall.h>
#include <fcntl.h>
#include <grp.h>
#include <signal.h>
#include <uuid/uuid.h>

#include "cspace.h"
#include "manifest.h"
#include "log.h"
#include "capp.h"
#include "capp_runtime.h"
#include "utils.h"

static int adjust_capabilities(char *name, int *cap, int size);
static int setup_capabilities(CSpace *space);
static int setup_cspace_security_profile(CSpace *space);
static int setup_user_namespace(CSpace *space);
static int prepare_child_map_files(pid_t pid, CSpace *space);
static int cspace_init_clone(void *arg);
static int handle_create_request(CSpace *space, int seqno, char *params);
static int send_response_packet(CSpace *space, int seqno, char *resp);
static CApp *cspace_capp_init(char *name, char *tag, char *path, uuid_t uuid);
static int cspace_capps_init(CApps **capps);
static int handle_crud_requests(CSpace *space);


/*
 * adjust_capabilities -- Adjust capabilties for the cspace
 *
 */
static int adjust_capabilities(char *name, int *cap, int size) {

  int i, ret=FALSE;
  cap_t capState = NULL;

  /* Drop the capiblities */
  for (i=0; i < size; i++) {
    if (prctl(PR_CAPBSET_DROP, cap[i], 0, 0, 0)) {
      log_error("Space: %s Error dropping capabilities. Error: %s", name,
		strerror(errno));
      return FALSE;
    }
  }

  /* Setting inheritable caps */
  capState = cap_get_proc();
  if (capState == NULL) {
    log_error("Space: %s Failed to get cap state.", name);
    goto done;
  }

  /* Clear the CAP options */
  if (cap_set_flag(capState, CAP_INHERITABLE, size, cap, CAP_CLEAR)!=0) {
    log_error("Space: %s Failed to set the cap flags for cap state. Error: %s",
	      name, strerror(errno));
    goto done;
  }

  /* Commit */
  if (cap_set_proc(capState)==-1) {
    log_error("Space: %s Failed to commit the cap flags. Error: %s", name,
	      strerror(errno));
    goto done;
  }

  ret = TRUE;

 done:
  if (capState) cap_free(capState);
  return ret;
}

/*
 * setup_capabilities -- setup privilage users capabilities. Basic idea is that
 *                       if one of our container go rouge, its ability to
 *                       damage the system should be less than if the same
 *                       program is running as root.
 *
 * Note: see capabilities(7) for details. 
 *
 */
static int setup_capabilities(CSpace *space) {

  /* Ambient capabilties are preserved across an execve() of an un-priviliged
   * program. As per man page: 'The ambient capability set obeys the
   * invariant that no capability can ever be ambient if it is not both
   * permitted and inheritable', hence we need to have both.
   */

  /* The algorithm to set the capabilities for the new program, via execv():
   * From capabilities(7):
   *
   * 'During an execve(2), the kernel calculates the new capabilities of the
   *  process using the following algorithm:
   *
   *  P'(ambient)     = (file is privileged) ? 0 : P(ambient)
   *  P'(permitted)   = (P(inheritable) & F(inheritable)) |
   *                    (F(permitted) & cap_bset) | P'(ambient)
   *  P'(effective)   = F(effective) ? P'(permitted) : P'(ambient)
   *  P'(inheritable) = P(inheritable)    [i.e., unchanged]
   *
   *   where:
   *
   *       P         denotes  the  value of a thread capability set before the
   *                 execve(2)
   *       P'        denotes the value of a thread capability  set  after  the
   *                 execve(2)
   *       F         denotes a file capability set
   *       cap_bset  is  the  value  of the capability bounding set (described
   *                 below).
   */

  int ret;

  ret = adjust_capabilities(space->name, space->cap, space->capCount);
  if (ret == FALSE) {
    log_error("Space: %s Error adjusting capabilties", space->name);
    return ret;
  }

  return TRUE;
}

/*
 * setup_cspace_security_profile -- setup security profile of the cspace
 *
 */
static int setup_cspace_security_profile(CSpace *space) {

  return setup_capabilities(space);
}

/*
 * setup_user_namespace -- setup the user/group on the child
 *
 */
static int setup_user_namespace(CSpace *space) {

  int ns=FALSE, resp=FALSE;
  int size, ret;

  /* unshare the user namespace. */
  if (!unshare(CLONE_NEWUSER)) {
    ns = TRUE;
  }

  /* Write on the socket connection with parent. Parent need to setup values
   * in the map files under /proc
   */
  if (write(space->sockets[CHILD_SOCKET], &ns, sizeof(ns)) != sizeof(ns)) {
    log_error("Space: %s Error writing to parent socket. Size: %d Value: %d",
	      space->name, sizeof(ns), ns);
    return FALSE;
  }

  /* Read response back from the parent. */
  size = read(space->sockets[CHILD_SOCKET], &resp, sizeof(resp));
  if (size != sizeof(resp)) {
    log_error("Space: %s Error reading from parent socket. size: %d Got: %d",
	      space->name, sizeof(resp), size);
    return FALSE;
  }

  if (!resp) {
    log_error("Space: %s Parent failed to setup map file", space->name);
    return FALSE;
  }

  /* Switch over the uid and gid */
  log_debug("Space: %s Switching to uid: %d and gid: %d", space->name,
	    space->uid, space->gid);

  ret = setgroups(1, &space->gid);
  if (ret != 0) {
    log_error("Space: %s error setting groups. gid: %d Error: %s", space->name,
	      space->gid, strerror(errno));
    return FALSE;
  }

  ret = setresgid(space->gid, space->gid, space->gid);
  if (ret != 0) {
    log_error("Space: %s Error setting group id to: %d Error: %s", space->name,
	      space->gid, strerror(errno));
    return FALSE;
  }

  ret = setresuid(space->uid, space->uid, space->uid);
  if (ret != 0) {
    log_error("Space: %s Error setting user id to: %d Error: %s", space->name,
	      space->uid, strerror(errno));
    return FALSE;
  }

  return TRUE;
}

/*
 * prepare_child_map_files -- setup map files (uid_map, gid_map and setgroups)
 *
 */
static int prepare_child_map_files(pid_t pid, CSpace *space) {

  int ns=0, size;
  int mapFd = 0;
  char *mapFiles[] = {"uid_map", "gid_map"};
  char mapPath[LXCE_MAX_PATH] = {0};
  char **file;

  if (space==NULL) return FALSE;

  size = read(space->sockets[PARENT_SOCKET], &ns, sizeof(ns));
  if (size != sizeof(ns)) {
    log_error("Error reading from client socket. Expected size: %d Got: %d",
	      sizeof(ns), size);
    return FALSE;
  }

  if (!ns) {
    log_error("Space: %s Child user namespace issue. unshare.", space->name);
    return FALSE;
  }

  for (file=&mapFiles[0]; *file; file++) {
    sprintf(mapPath, "/proc/%d/%s", pid, *file);

    if ((mapFd = open(mapPath, O_WRONLY)) == -1) {
      log_error("Space: %s Error opening map file: %s Error: %s", space->name,
		mapPath, strerror(errno));
      return FALSE;
    }

    if (dprintf(mapFd, "0 %d 1\n", USER_NS_OFFSET) == -1) {
      log_error("Space: %s Error writing to map file: %s Error: %s",
		space->name, mapPath, strerror(errno));
      close(mapFd);
      return FALSE;
    }
  }

  /* Inform child, it can proceed. */
  size = write(space->sockets[PARENT_SOCKET], &(int){TRUE}, sizeof(int));
  if (size != sizeof(int)) {
    log_error("Space: %s Error writing to child socket", space->name);
    log_error("Expected size: %d Wrote: %d", sizeof(int), size);
    return FALSE;
  }

  return TRUE;
}

/*
 * setup_mounts --
 *
 */
static int setup_mounts(CSpace *space) {

  int ret=FALSE;
  char tempMount[] = "/tmp/tmp.ukama.XXXXXX"; /* last 6 char needs to be X */
  char oldRoot[]   = "/tmp/tmp.ukama.XXXXXX/oldroot.XXXXXX"; /* same */

  if (mount(NULL, "/", NULL, MS_REC | MS_PRIVATE, NULL)) {
    log_error("Space: %s Failed to remount as MS_PRIVATE. Error: %s",
	      space->name, strerror(errno));
    return ret;
  }

  /* make temp and bind mount */
  if (!mkdtemp(tempMount)) {
    log_error("Space: %s Failed to make temp dir: %s. Error: %s", space->name,
	      tempMount, strerror(errno));
    return FALSE;
  }

  if (mount(space->mountDir, tempMount, NULL, MS_BIND | MS_PRIVATE, NULL)) {
    log_error("Space: %s Failed to do bind mount. %s %s Error: %s",
	      space->name, space->mountDir, tempMount, strerror(errno));
    return FALSE;
  }

  if (!mkdtemp(oldRoot)) {
    log_error("Failed to create old Root directory. %s Error :%s", oldRoot,
	      strerror(errno));
    return FALSE;
  }

  /* pivot root */
  if (syscall(SYS_pivot_root, space->mountDir, oldRoot)) {
    log_error("Failed to pivot_root from %s to %s", space->mountDir, oldRoot);
    return FALSE;
  }

  log_debug("Pivot root sucessfully done. from %s to %s", space->mountDir,
	    oldRoot);

  /* clean up */

  if (chdir("/")) {
    log_error("Error changing director to / after pivot");
    return FALSE;
  }

  if (umount2(oldRoot, MNT_DETACH) || rmdir(oldRoot)) {
    log_error("Failed to umount/rm old root: %s. Error: %s", oldRoot,
	      strerror(errno));
    return FALSE;
  }

  return TRUE;
}

/*
 * cspace_init_clone --
 *
 */
static int cspace_init_clone(void *arg) {

  CSpace *space = (CSpace *)arg;
  char *hostName=NULL;

  /* Close parent socket */
  close(space->sockets[PARENT_SOCKET]);

  if (space->hostName) {
    hostName = space->hostName;
  } else {
    hostName = CSPACE_DEFAULT_HOSTNAME;
  }

  /* Step-1: setup hostname. */
  if (sethostname(hostName, strlen(hostName))) {
    log_error("Space: %s Error setting host name: %s", space->name,
	      hostName);
    return FALSE;
  }

  /* Step-2: setup security profile (cap and seccomp) */
  setup_cspace_security_profile(space);

  /* Step-3: setup user namespace */
  setup_user_namespace(space);

  /* Step-4: cSpace stays in this state forever
   * Accept capp CRUD calls from the parent process.
   */
  handle_crud_requests(space);

  return TRUE;
}

/*
 * create_cspace -- create contained spaces
 */
int create_cspace(CSpace *space, pid_t *pid) {

  char *stack=NULL;
  
  /* logic is as follows:
   *
   * 1. Create socketpair - this will be used to communicate between lxce and
   *                        cSpace
   * 2. Setup proper cgroups.
   * 3. Clone with proper flags for namespaces
   * 4. setup mount
   * 5. setup user namespace
   * 6. Limit capabilities
   * 7. execv into cSpace.init
   */

  /* Sanity check */
  if (space == NULL) return FALSE;
  
  /* Create socket pairs.
   * Re: SOCK_SEQPACKET:
   * http://urchin.earth.li/~twic/Sequenced_Packets_Over_Ordinary_TCP.html
   */
  if (socketpair(AF_LOCAL, SOCK_SEQPACKET, 0, space->sockets)) {
    log_error("Space: %s Error creating socket pair", space->name);
    return FALSE;
  }

  if (!(stack = malloc(STACK_SIZE))) {
    log_error("Error allocating stack of size: %d", STACK_SIZE);
    return FALSE;
  }

  /* clone with proper flags for namespaces */
  *pid = clone(cspace_init_clone, stack + STACK_SIZE,
	      SIGCHLD | space->nameSpaces, space);
  if (*pid == -1) {
    log_error("Space: %s Unable to clone cInit", space->name);
    return FALSE;
  }

  /* Close child socket */
  close(space->sockets[CHILD_SOCKET]);

  /* prepare child process gid/uid map files. */
  if (prepare_child_map_files(*pid, space) == FALSE) {
    log_error("Error preparing map files for child process. Terminating it");
    kill(pid, SIGKILL); /* Kill child process */
    close(space->sockets[PARENT_SOCKET]);
    close(space->sockets[CHILD_SOCKET]);
    return FALSE;
  }

  return TRUE;
}

/*
 * str_to_cap --
 *
 */
static int str_to_cap(char *str) {

  if (strcmp(str, "CAP_BLOCK_SUSPEND")==0) {
    return CAP_BLOCK_SUSPEND;
  } else if (strcmp(str, "CAP_IPC_LOCK")==0) {
    return CAP_IPC_LOCK;
  } else if (strcmp(str, "CAP_MAC_ADMIN")==0) {
    return CAP_MAC_ADMIN;
  } else if (strcmp(str, "CAP_MAC_OVERRIDE")==0) {
    return CAP_MAC_OVERRIDE;
  }

  log_error("Invalid capabilities: %s", str);
  return 0;
}

/*
 * deserialize_cspace_file -- convert the json into internal struct
 *
 */
static int deserialize_cspace_file(CSpace *space, json_t *json) {

  int j=0, size=0;
  json_t *obj;
  json_t *jArray, *jElem;

  if (space == NULL) return FALSE;
  if (json == NULL) return FALSE;

  if (!set_str_object_value(json, &(space->version), JSON_VERSION, TRUE,
			    NULL)) {
    return FALSE;
  }

  if (!set_str_object_value(json, &(space->target), JSON_TARGET, TRUE, NULL)) {
    return FALSE;
  }

  if (space->target == LXCE_SERIAL) {
    if (!set_str_object_value(json, &(space->serial), JSON_SERIAL, TRUE, NULL)) {
      return FALSE;
    }
  } else {
    set_str_object_value(json, &(space->serial), JSON_SERIAL, FALSE, NULL);
  }

  if (!set_str_object_value(json, &(space->name), JSON_NAME, TRUE, NULL)) {
    return FALSE;
  }

  set_str_object_value(json, &(space->hostName), JSON_HOSTNAME, FALSE,
		       CSPACE_DEFAULT_HOSTNAME);

  set_integer_object_value(json, &(space->uid), JSON_UID, FALSE, 0);
  set_integer_object_value(json, &(space->gid), JSON_GID, FALSE, 0);

  /* Look for namespaces. */
  space->nameSpaces = 0;
  jArray = json_object_get(json, JSON_NAMESPACES);
  if (jArray != NULL) {
    size = json_array_size(jArray);

    for (j=0; j<size; j++) {
      jElem = json_array_get(jArray, j);
      if (jElem) {
	obj = json_object_get(jElem, JSON_TYPE);
	if (obj)
	  space->nameSpaces |= namespaces_flag(json_string_value(obj));
      }
    }
  } else {
    log_debug("No valid namespaces found.");
  }

  /* Look for capabilities */
  jArray = json_object_get(json, JSON_CAPABILITIES);
  if (jArray != NULL) {
    size = json_array_size(jArray);
    space->capCount = size;

    if (size > CONTD_MAX_CAPS) {
      log_error("%d many more Capabilities are defined than supported: 5d",
		(size-CONTD_MAX_CAPS), CONTD_MAX_CAPS);
      return FALSE;
    }

    for (j=0; j<size; j++) {
      jElem = json_array_get(jArray, j);
      if (jElem) {
	obj = json_object_get(jElem, JSON_TYPE);
	if (obj)
	  space->cap[j] = str_to_cap(json_string_value(obj));
      }
    }
  } else {
    log_debug("No valid capabilities found.");
  }

  return TRUE;
}

/*
 * process_cspace_config --
 *
 */
int process_cspace_config(char *fileName, CSpace *space) {

  int ret=FALSE;
  FILE *fp=NULL;
  char *buffer=NULL;
  long size=0;
  json_t *json;
  json_error_t jerror;

  /* Sanity check */
  if (fileName==NULL) return FALSE;
  if (space==NULL) return FALSE;

  if ((fp = fopen(fileName, "rb")) == NULL) {
    log_error("Error opening file: %s Error %s", fileName, strerror(errno));
    return FALSE;
  }

  /* Read everything into buffer */
  fseek(fp, 0, SEEK_END);
  size = ftell(fp);
  fseek(fp, 0, SEEK_SET);

  if (size > CONFIG_MAX_SIZE) {
    log_error("Error opening file: %s Error: File size too big: %ld",
	      fileName, size);
    fclose(fp);
    return FALSE;
  }

  buffer = (char *)malloc(size+1);
  if (buffer==NULL) {
    log_error("Error allocating memory of size: %ld", size+1);
    fclose(fp);
    return FALSE;
  }
  memset(buffer, 0, size+1);
  fread(buffer, 1, size, fp); /* Read everything into buffer */

  /* Trying loading it as JSON */
  json = json_loads(buffer, 0, &jerror);
  if (json==NULL) {
    log_error("Error loading contd config into JSON format. File: %s Size: %ld",
	      fileName, size);
    log_error("JSON error on line: %d: %s", jerror.line, jerror.text);
    goto done;
  }

  /* Now convert JSON into internal struct */
  ret = deserialize_cspace_file(space, json);

  if (space) {
    (space)->configFile = strdup(fileName);
  }

 done:
  if (buffer) free(buffer);
  fclose(fp);

  json_decref(json);
  return ret;
}

/*
 * handle_crud_requests --
 *
 */
static int handle_crud_requests(CSpace *space) {

  struct timeval tv;
  int count, seqno;
  char buffer[CSPACE_MAX_BUFFER] = {0};
  char *cmd=NULL, *params=NULL;

  /* time-out socket */
  tv.tv_sec  = 5; /* XXX - check on this. */
  tv.tv_usec = 0;
  setsockopt(space->sockets[CHILD_SOCKET], SOL_SOCKET, SO_RCVTIMEO,
	     (const char*)&tv, sizeof tv);

  while (TRUE) {

    count = recv(space->sockets[CHILD_SOCKET], buffer, CSPACE_MAX_BUFFER, 0);

    if (count <=0  && errno != EAGAIN) {
      log_error("Error reading packet from cspace socket. Name: %s",
		space->name);
      return CSPACE_READ_ERROR;
    }

    if (count == 0 && errno == EAGAIN) {
      continue;
    }

    cmd    = (char *)malloc(count+1);
    params = (char *)malloc(count+1);

    if (!cmd || !params) {
      log_error("Memory allocation error. Size: %d", 2*count);
      return CSPACE_MEMORY_ERROR;
    }

    /* we have some packet. Let's see what we got
     * packet format is [cmd seq_id some_text]
     */
    sscanf(buffer, "%s %d %s", cmd, &seqno, params);

    if (strcmp(cmd, CAPP_CMD_CREATE)==0) {
      handle_create_request(space, seqno, params);
      create_and_run_capps(space->apps);
    } else {
      log_error("Invalid command recevied: %s", cmd);
    }

    free(cmd);
    free(params);

  } /* forever loop */
}

/*
 * handle_create_request --
 *
 */
static int handle_create_request(CSpace *space, int seqno, char *params) {

  uuid_t uuid;
  char idStr[36+1];
  CApp *capp=NULL;
  char *name, *tag, *path, *resp;

  if (!space || !seqno) return FALSE;

  /* Generate ID */
  uuid_generate(uuid);
  uuid_unparse(uuid, &idStr[0]);

  name = strtok(params, ":");
  tag  = strtok(NULL, ":");
  path = strtok(NULL, ":");

  if (!space->apps) {
    if (!cspace_capps_init(&(space->apps))) {
      resp = "ERROR";
      goto reply;
    }
  }

  capp = cspace_capp_init(name, tag, path, uuid);
  if (capp == NULL) {
    resp = "ERROR";
    goto reply;
  }

  add_to_apps(space->apps, capp, PEND_LIST, NULL);

  resp = &idStr[0];

 reply:
  if (!send_response_packet(space, seqno, resp)) {
    log_error("Error sending response packet. %s", space->name);
    return FALSE;
  }

  return TRUE;
}

/*
 * send_response_packet --
 *
 */
static int send_response_packet(CSpace *space, int seqno, char *resp) {

  char *data=NULL;
  int size;

  if (!space || !resp) return FALSE;

  if (space->sockets[CHILD_SOCKET] <= 0) {
    log_error("Socket pair is closed between thread and cspace. Name: %s",
	      space->name);
    return FALSE;
  }

  size = (3*sizeof(int)+2) + (strlen(resp)+1);

  data = (char *)malloc(size);
  if (data == NULL) {
    log_error("Memory allocation error. Size: %d", size);
    return FALSE;
  }

  sprintf(data, "%d %s", seqno, resp);

  if (send(space->sockets[CHILD_SOCKET], data, strlen(data), 0) <0) {
    log_error("Sending response packet to thread over socket failed. %s",
	      space->name);
    free(data);
    return FALSE;
  }

  log_debug("Response sent. Resp: %s", data);

  free(data);
  return TRUE;
}

/*
 * cspace_capps_init --
 *
 */
static int cspace_capps_init(CApps **capps) {

  if (*capps != NULL) return TRUE;

  *capps = (CApps *)calloc(1, sizeof(CApps));
  if (*capps == NULL) {
    log_error("Memory allocation error of size: %d", sizeof(CApps));
    return FALSE;
  }

  return TRUE;
}

/*
 * cspace_capp_init --
 *
 */
static CApp *cspace_capp_init(char *name, char *tag, char *path, uuid_t uuid) {

  CApp *capp=NULL;

  if (!name || !tag || !path) return NULL;

  /* We have valid name, tag and path. */
  capp = (CApp *)calloc(1, sizeof(CApp));
  if (!capp) {
    log_error("Memory allocation error of size: %d", sizeof(CApp));
    return NULL;
  }

  capp->params = (CAppParams *)malloc(sizeof(CAppParams));
  capp->state = (CAppState *)malloc(sizeof(CAppState));

  if (capp->params == NULL || capp->state == NULL) {
    log_error("Memory allocation error of sizes: %d %d", sizeof(CAppParams),
	      sizeof(CAppState));
    goto failure;
  }

  capp->params->name  = strdup(name);
  capp->params->tag   = strdup(tag);
  capp->params->path  = strdup(path);
  capp->params->space = NULL;
  uuid_copy(capp->params->uuid, uuid);

  capp->state->state       = CAPP_STATE_PENDING;
  capp->state->exit_status = CAPP_STATE_INVALID;

  capp->policy = NULL;
  capp->space  = NULL;

  return capp;

 failure:
  clear_capp(capp);
  return NULL;
}
