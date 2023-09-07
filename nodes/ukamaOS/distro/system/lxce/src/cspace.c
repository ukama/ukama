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
#include <dirent.h>

#include "space.h"
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
static int cspace_init_clone(void *arg);
static int handle_create_request(CSpace *space, int seqno, char *params);
static int send_response_packet(CSpace *space, int seqno, char *resp);
static CApp *cspace_capp_init(char *name, char *tag, char *path, uuid_t uuid);
static int cspace_capps_init(CApps **capps);
static int handle_crud_requests(CSpace *space);
static int valid_cspace_rootfs_pkg(char *fileName);

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
  } else {
    log_error("Space: [%s] Unable to unshare the user namespace. %s",
	      space->name, strerror(errno));
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

  ret = setgroups(1, & (gid_t) {space->gid});
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
  //setup_cspace_security_profile(space);

  /* Step-3: Setup mounts */
  log_debug("Setting up mounts for space: %s at rootfs: %s", space->name,
	    space->rootfs);
  if (!setup_mounts(AREA_TYPE_CSPACE, space->rootfs, space->name)) {
    log_error("Space: %s Error setting up rootfs mount: %s", space->name,
	      space->rootfs);
    return FALSE;
  }

  /* Step-4: setup user namespace */
  setup_user_namespace(space);

  /* Step-5: cSpace stays in this state forever
   * Accept capp CRUD calls from the parent process.
   */
  handle_crud_requests(space);

  return TRUE;
}

/*
 * create_cspace --
 *
 */
int create_cspace(CSpace *space, pid_t *pid) {

  if (space == NULL) return FALSE;

  return create_space(AREA_TYPE_CSPACE,
		      space->sockets, space->nameSpaces,
		      space->name, pid,
		      cspace_init_clone, (void *)space);
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

  if (!set_str_object_value(json, &(space->version), JSON_VERSION, TRUE, NULL)) {
    return FALSE;
  }

  if (!set_str_object_value(json, &(space->target), JSON_TARGET, TRUE, NULL)) {
    return FALSE;
  }

  if (strcmp(space->target, LXCE_SERIAL)==0) {
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

  /* setup veth IP address */
  set_str_object_value(json, &(space->vethIP), JSON_VETH_IP, FALSE,
		       CSPACE_DEFAULT_VETH_IP);

  /* 3: 2 for / and 1 for NULL */
  space->rootfs = (char *)malloc(strlen(DEF_CSPACE_ROOTFS_PATH) +
				 strlen(space->name) + 3);
  sprintf(space->rootfs, "%s/%s/", DEF_CSPACE_ROOTFS_PATH, space->name);

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

    memset(buffer, 0, CSPACE_MAX_BUFFER);

    count = read(space->sockets[CHILD_SOCKET], buffer, CSPACE_MAX_BUFFER);

    if (count <=0  && errno != EAGAIN) {
      log_error("Error reading packet from cspace socket. Name: %s",
		space->name);
      return CSPACE_READ_ERROR;
    }

    if (count <= 0 && errno == EAGAIN) {
      continue;
    }

    cmd    = (char *)malloc(count+1);
    params = (char *)malloc(count+1);

    if (!cmd || !params) {
      log_error("Memory allocation error. Size: %d", 2*count);
      return CSPACE_MEMORY_ERROR;
    }
    memset(cmd, 0, count+1);
    memset(params, 0, count+1);

    /* we have some packet. Let's see what we got
     * packet format is [cmd seq_id some_text]
     */
    sscanf(buffer, "%s %d %s", cmd, &seqno, params);

    log_debug("%d %s %d %s", count, cmd, seqno, params);
    
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

  add_to_apps(space->apps, capp, PEND_LIST, 0);

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

  capp->params  = (CAppParams *)calloc(1, sizeof(CAppParams));
  capp->state   = (CAppState *)calloc(1, sizeof(CAppState));
  capp->config  = (CAppConfig *)calloc(1, sizeof(CAppConfig));
  capp->runtime = (CAppRuntime *)calloc(1, sizeof(CAppRuntime));

  if (capp->params == NULL || capp->state == NULL || capp->config == NULL ||
      capp->runtime == NULL ) {
    log_error("Memory allocation error of sizes: %d %d %d %d ",
	      sizeof(CAppParams), sizeof(CAppState), sizeof(CAppConfig),
	      sizeof(CAppRuntime));
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

/*
 * valid_cspace_rootfs_pkg --
 *
 */
static int valid_cspace_rootfs_pkg(char *fileName) {

  struct stat stats;

  if (fileName == NULL) return FALSE;

  stat(fileName, &stats);
  if (S_ISREG(stats.st_mode)) {
    return TRUE;
  }

  return FALSE;
}

/*
 * cspace_unpack_rootfs --
 *
 */
int cspace_unpack_rootfs(char *destDir) {

  char pkg[CSPACE_MAX_BUFFER] = {0};
  char runMe[CSPACE_MAX_BUFFER] = {0};
  struct stat stats;
  int ret;

  /* Steps are as follow:
   * 1. check the existance of rootfs pkg "cspace_rootfs.tar.gz". Currently,
   *    default location is /capps/pkgs
   * 2. remove existing rootfs at 'destDir'
   * 3. untar pkg to 'destDir'
   */

  sprintf(pkg, "%s/%s", DEF_CSPACE_ROOTFS_PKG_PATH, DEF_CSPACE_ROOTFS_PKG_NAME);
  if (valid_cspace_rootfs_pkg(pkg) == FALSE) {
    log_error("Unable to find cspace rootfs pkg at: %s", pkg);
    return FALSE;
  }

  /* Check if directory exist */
  stat(destDir, &stats);
  if (S_ISDIR(stats.st_mode)) {
    sprintf(runMe, "/bin/rm -rf %s", destDir);
    if ((ret = system(runMe)) < 0) {
      log_error("Error removing existing cspace rootfs path at: %s Error: %s",
		destDir, strerror(errno));
      return FALSE;
    }
  }

  /* re-create the directory */
  if(mkdir(destDir, 0700) < 0) {
    log_error("Error creating cspsace rootfs dir: %s Error: %s", destDir,
	      strerror(errno));
    return FALSE;
  }

  /* untar to destDir */
  sprintf(runMe, "/bin/tar xfz %s/%s -C %s", DEF_CSPACE_ROOTFS_PKG_PATH,
	  DEF_CSPACE_ROOTFS_PKG_NAME, destDir);
  if ((ret = system(runMe)) < 0) {
    log_error("Unable to unpack the cspace rootfs: %s/%s to %s Code: %d",
	      DEF_CSPACE_ROOTFS_PKG_PATH, DEF_CSPACE_ROOTFS_PKG_NAME,
	      destDir, ret);
    return FALSE;
  }

  log_debug("cspace rootfs successfully unpack at: %s", destDir);

  return TRUE;
}
