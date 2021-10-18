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

#include "contained.h"
#include "manifest.h"
#include "log.h"

static int basic[] = { CAP_BLOCK_SUSPEND, /* block system suspend. */
		       CAP_IPC_LOCK,      /* Lock memory */
		       CAP_MAC_ADMIN,     /* allow MAC config */
		       CAP_MAC_OVERRIDE,  /* override MAC. */
};

static int serviceType[] = { CAP_NET_ADMIN,        /* Network related ops */
			     CAP_NET_BIND_SERVICE, /* privilage ports */
			     CAP_SETFCAP,          /* arbitray cap on file. */
			     CAP_SETUID,           /* Manipulate process UIDs */
			     CAP_SYS_ADMIN,        /* various cap. */
			     CAP_SYS_BOOT,         /* use reboot */
			     CAP_SYS_MODULE,       /* load kernel modules */
			     CAP_SYS_NICE,         /* process nice() */
			     CAP_SYS_RAWIO,      /* I/O ops and various thing */
			     CAP_SYS_TIME,         /* system clock */
			     CAP_SYSLOG,           /* privileged syslog */
			     CAP_WAKE_ALARM,       /* CLOCK_BOOTTIME_ALARM */
			     CAP_SYS_RESOURCE   /* resources related action */ 
};

static int shutdownType[] = { CAP_SYSLOG,    /* privileged syslog */
			      CAP_SYS_TIME,  /* system clock */
			      CAP_NET_ADMIN  /* Network related operations. */
};

static int adjust_capabilities(int *cap, int size);
static int setup_capabilities(char *podType);
static int setup_pod_security_profile(char *type);
static int setup_user_namespace(Pod *pod);
static int prepare_child_map_files(pid_t pid, Pod *pod);
static int cInit_clone(void *arg);

/*
 * adjust_capabilities -- Adjust capabilties for process
 *
 */
static int adjust_capabilities(int *cap, int size) {

  int i, ret=FALSE;
  cap_t capState = NULL;

  /* Drop the basic capiblities */
  for (i=0; i < size; i++) {
    if (prctl(PR_CAPBSET_DROP, cap[i], 0, 0, 0)) {
      log_error("Error dropping capabilities. Error: %s", strerror(errno));
      return FALSE;
    }
  }

  /* Setting inheritable caps */
  capState = cap_get_proc();
  if (capState == NULL) {
    log_error("Failed to get cap state.");
    goto done;
  }

  /* Clear the CAP options */
  if (cap_set_flag(capState, CAP_INHERITABLE, size, cap, CAP_CLEAR)!=0) {
    log_error("Failed to set the cap flags for cap state. Error: %s",
	      strerror(errno));
    goto done;
  }

  /* Commit */
  if (cap_set_proc(capState)==-1) {
    log_error("Failed to commit the cap flags. Error: %s", strerror(errno));
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
static int setup_capabilities(char *podType) {

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

  /* For onboot:   basic
   *     service:  basic & service
   *     shutdown: basic & shutdown
   */

  ret = adjust_capabilities(basic, sizeof(basic)/sizeof(*basic));
  if (ret == FALSE) {
    log_error("Error adjusting basic capabilties");
    return ret;
  }

  if (strcmp(podType, POD_TYPE_SERVICE) == 0) {
    ret = adjust_capabilities(serviceType,
			      sizeof(serviceType)/sizeof(*serviceType));
    if (ret == FALSE) {
      log_error("Error adjusting service capabilities");
      /* Note: we don't need to undo the basic capabilities as upon error the
       * child process (pod setup) terminates. Otherwise, we should clean up.
       */
      return FALSE;
    }
  } else if  (strcmp(podType, POD_TYPE_SHUTDOWN) == 0) {
    ret = adjust_capabilities(shutdownType,
			      sizeof(shutdownType)/sizeof(*shutdownType));
    if (ret == FALSE) {
      log_error("Error adjusting shutdown capabilities");
      return FALSE;
    }
  }

  return TRUE;
}

/*
 * setup_secure_pod -- setup security profile of the pod
 *
 */
static int setup_pod_security_profile(char *type) {

  return setup_capabilities(type);
}

/*
 * setup_user_namespace -- setup the user/group on the child
 *
 */
static int setup_user_namespace(Pod *pod) {

  int ns=FALSE, resp=FALSE;
  int size, ret;

  /* unshare the user namespace. */
  if (!unshare(CLONE_NEWUSER)) {
    ns = TRUE;
  }

  /* Write on the socket connection with parent. Parent need to setup values
   * in the map files under /proc
   */
  if (write(pod->sockets[0], &ns, sizeof(ns)) != sizeof(ns)) {
    log_error("Error writing to parent socket. Size: %d Value: %d", sizeof(ns),
	      ns);
    return FALSE;
  }

  /* Read response back from the parent. */
  size = read(pod->sockets[0], &resp, sizeof(resp));
  if (size != sizeof(resp)) {
    log_error("Error reading from parent socket. Expected size: %d Got: %d",
	      sizeof(resp), size);
    return FALSE;
  }

  if (!resp) {
    log_error("Parent failed to setup map file");
    return FALSE;
  }

  /* Switch over the uid and gid */
  log_debug("Switching to uid: %d and gid: %d", pod->uid, pod->gid);

  ret = setgroups(1, &pod->gid);
  if (ret != 0) {
    log_error("Error setting groups. gid: %d Error: %s", pod->gid,
	      strerror(errno));
    return FALSE;
  }

  ret = setresgid(pod->gid, pod->gid, pod->gid);
  if (ret != 0) {
    log_error("Error setting group id to: %d Error: %s", pod->gid,
	      strerror(errno));
    return FALSE;
  }

  ret = setresuid(pod->uid, pod->uid, pod->uid);
  if (ret != 0) {
    log_error("Error setting user id to: %d Error: %s", pod->uid,
	      strerror(errno));
    return FALSE;
  }

  return TRUE;
}

/*
 * prepare_child_map_files -- setup map files (uid_map, gid_map and setgroups)
 *
 */
static int prepare_child_map_files(pid_t pid, Pod *pod) {

  int ns=0, size;
  int mapFd = 0;
  char *mapFiles[] = {"uid_map", "gid_map"};
  char mapPath[LXCE_MAX_PATH] = {0};
  char **file;

  if (pod==NULL) return FALSE;

  size = read(pod->sockets[0], &ns, sizeof(ns));
  if (size != sizeof(ns)) {
    log_error("Error reading from client socket. Expected size: %d Got: %d",
	      sizeof(ns), size);
    return FALSE;
  }

  if (!ns) {
    log_error("Child user namespace issue. unshare.");
    return FALSE;
  }

  for (file=&mapFiles[0]; *file; file++) {
    sprintf(mapPath, "/proc/%d/%s", pid, *file);

    if ((mapFd = open(mapPath, O_WRONLY)) == -1) {
      log_error("Error opening map file: %s Error: %s", mapPath,
		strerror(errno));
      return FALSE;
    }

    if (dprintf(mapFd, "0 %d 1\n", USER_NS_OFFSET) == -1) {
      log_error("Error writing to map file: %s Error: %s", mapPath,
		strerror(errno));
      close(mapFd);
      return FALSE;
    }
  }

  /* Inform child, it can proceed. */
  size = write(pod->sockets[0], &(int){TRUE}, sizeof(int));
  if (size != sizeof(int)) {
    log_error("Error writing to child socket. Expected size: %d Wrote: %d",
	      sizeof(int), size);
    return FALSE;
  }

  return TRUE;
}

/*
 * setup_mounts --
 *
 */
static int setup_mounts(Pod *pod) {

  int ret=FALSE;
  char tempMount[] = "/tmp/tmp.ukama.XXXXXX"; /* last 6 char needs to be X */
  char oldRoot[]   = "/tmp/tmp.ukama.XXXXXX/oldroot.XXXXXX"; /* same */

  if (mount(NULL, "/", NULL, MS_REC | MS_PRIVATE, NULL)) {
    log_error("Failed to remount as MS_PRIVATE. Error: %s", strerror(errno));
    return ret;
  }

  /* make temp and bind mount */
  if (!mkdtemp(tempMount)) {
    log_error("Failed to make temp dir: %s. Error: %s", tempMount,
	      strerror(errno));
    return FALSE;
  }

  if (mount(pod->mountDir, tempMount, NULL, MS_BIND | MS_PRIVATE, NULL)) {
    log_error("Failed to do bind mount. %s %s Error: %s", pod->mountDir,
	      tempMount, strerror(errno));
    return FALSE;
  }

  if (!mkdtemp(oldRoot)) {
    log_error("Failed to create old Root directory. %s Error :%s", oldRoot,
	      strerror(errno));
    return FALSE;
  }

  /* pivot root */
  if (syscall(SYS_pivot_root, pod->mountDir, oldRoot)) {
    log_error("Failed to pivot_root from %s to %s", pod->mountDir, oldRoot);
    return FALSE;
  }

  log_debug("Pivot root sucessfully done. from %s to %s", pod->mountDir,
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
 * cInit_clone --
 *
 */
static int cInit_clone(void *arg) {

  Pod *pod = (Pod *)arg;
  char *hostName;

  if (pod->hostName) {
    hostName = pod->hostName;
  } else {
    hostName = CONTD_DEFAULT_HOSTNAME;
  }

  /* Step-1: setup hostname. */
  if (sethostname(hostName, strlen(hostName))) {
    log_error("Error setting host name: %s", hostName);
    return FALSE;
  }

  /* Step-2: setup security profile (cap and seccomp) */
  setup_pod_security_profile(pod->type);

  /* Step-3: setup mounts*/
  setup_mounts(pod);
  
  /* Step-4: setup user namespace */
  setup_user_namespace(pod);
  
  return TRUE;
}

/*
 * create_ukama_pod -- Create POD for the three type of containers:
 *                     boot, service and shutdown. Each POD has its own
 *                     namespace (PID, UID, NET, MOUNT) and cgroups.
 *                     A simple process PID-1 cInit.d is running within
 *                     each POD responsible to process request and act as
 *                     parent process of every container running within the
 *                     POD.
 */
int create_ukama_pod(Pod *pod, Manifest *manifest, char *type) {

  int cloneFlags=0;
  int childStatus;
  pid_t pid;
  char *stack=NULL;
  
  /* logic is as follows, for each pod type in the manifest:
   *
   * 1. Create socketpair - this will be used to communicate between lxce and
   *                        cInit.d
   * 2. Setup proper cgroups.
   * 3. Clone with proper flags for namespaces
   * 4. setup mount
   * 5. setup user namespace
   * 6. Limit capabilities
   * 7. Restrict system calls
   * 8. execv
   */

  /* Sanity check */
  if (pod == NULL || manifest == NULL) return FALSE;

  pod->type = strdup(type);
  
  /* Create socket pairs.
   * Re: SOCK_SEQPACKET:
   * http://urchin.earth.li/~twic/Sequenced_Packets_Over_Ordinary_TCP.html
   */
  if (socketpair(AF_LOCAL, SOCK_SEQPACKET, 0, pod->sockets)) {
    log_error("Error creating socket pair for pod type: %s", pod->type);
    return FALSE;
  }

  /* child only access one. */
  if (fcntl(pod->sockets[0], F_SETFD, FD_CLOEXEC)) {
    fprintf(stderr, "Failed to close socket via fcntl for type: %s",
	    pod->type);
    if (pod->sockets[0]) close(pod->sockets[0]);
    if (pod->sockets[1]) close(pod->sockets[1]);
    
    return FALSE;
  }

  /* clone with proper flags for namespaces */
  cloneFlags = SIGCHLD |
    CLONE_NEWNS     |
    CLONE_NEWPID    |
    CLONE_NEWUTS;

  pid = clone(cInit_clone, stack + STACK_SIZE, cloneFlags, pod);
  if (pid == -1) {
    log_error("Unable to clone cInit for type: %s", pod->type);
    return FALSE;
  }

  close(pod->sockets[1]);
  pod->sockets[1] = 0;

  /* prepare child process gid/uid map files. */
  if (prepare_child_map_files(pid, pod) == FALSE) {
    log_error("Error preparing map files for child process. Terminating it");
    kill(pid, SIGKILL);
    return FALSE;
  }

  /* Wait on child. XXX - fix me.*/
  waitpid(pid, &childStatus, 0);

  if (pod->sockets[0]) close(pod->sockets[0]);
  if (pod->sockets[1]) close(pod->sockets[1]);
  
  return TRUE;
}

/*
 * set_integer_object_value --
 *
 */
static int set_integer_object_value(json_t *json, int *param, char *objName,
				    int mandatory, int defValue) {

  json_t *obj;

  obj = json_object_get(json, objName);
  if (obj==NULL) {
    if (mandatory) {
      log_error("Missing Mandatory JSON field: %s Setting to default: %d",
		objName, defValue);
      if (defValue)  {
	*param = defValue;
      } else {
	return FALSE;
      }
    } else {
      log_debug("Missing JSON field: %s. Ignored.", objName);
      *param = 0;
    }
  } else {
    *param = json_integer_value(obj);
  }

  return TRUE;
}

/*
 * set_str_object_value --
 *
 */
static int set_str_object_value(json_t *json, char *param, char *objName,
				int mandatory, char *defValue) {

  json_t *obj;

  obj = json_object_get(json, objName);
  if (obj==NULL) {
    if (mandatory) {
      log_error("Missing Mandatory JSON field: %s Setting to default: %s",
		objName, defValue);
      if (defValue)  {
	param = strdup(defValue);
      } else {
	return FALSE;
      }
    } else {
      log_debug("Missing JSON field: %s. Ignored.", objName);
      param = NULL;
    }
  } else {
    param = strdup(json_string_value(obj));
  }

  return TRUE;
}

/*
 * namespace_flag --
 *
 */
static int namespaces_flag(char *ns) {

  if (strcmp(ns, "pid")==0) {
    return CLONE_NEWPID;
  } else if (strcmp(ns, "uts")==0) {
    return CLONE_NEWUTS;
  } else if (strcmp(ns, "net")==0) {
    return CLONE_NEWNET;
  } else if (strcmp(ns, "mount")==0) {
    return CLONE_NEWNS;
  } else if (strcmp(ns, "user")==0) {
    return CLONE_NEWUSER;
  } else {
    log_error("Unsupported namespace type detecetd: %s", ns);
    return 0;
  }

  return 0;
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
 * deserialize_contdSpace_file -- convert the json into internal struct
 *
 */
static int deserialize_contdSpace_file(ContdSpace *space, json_t *json) {

  int j=0, size=0;
  json_t *obj;
  json_t *jArray, *jElem;

  if (space == NULL) return FALSE;
  if (json == NULL) return FALSE;

  if (!set_str_object_value(json, space->version, JSON_VERSION, TRUE, NULL)) {
    return FALSE;
  }

  if (!set_str_object_value(json, space->target, JSON_TARGET, TRUE, NULL)) {
    return FALSE;
  }

  if (strcmp(space->target, LXCE_SERIAL)==0) {
    if (!set_str_object_value(json, space->serial, JSON_SERIAL, TRUE, NULL)) {
      return FALSE;
    }
  } else {
    set_str_object_value(json, space->serial, JSON_SERIAL, FALSE, NULL);
  }

  if (!set_str_object_value(json, space->name, JSON_NAME, TRUE, NULL)) {
    return FALSE;
  }

  set_str_object_value(json, space->hostName, JSON_HOSTNAME, FALSE,
		       CONTD_DEFAULT_HOSTNAME);

  set_integer_object_value(json, &space->uid, JSON_UID, FALSE, 0);
  set_integer_object_value(json, &space->gid, JSON_GID, FALSE, 0);

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
 * process_contdSpace_config --
 *
 */
int process_contdSpace_config(char *fileName, ContdSpace *contdSpace) {

  int ret=FALSE;
  FILE *fp=NULL;
  char *buffer=NULL;
  long size=0;
  json_t *json;
  json_error_t jerror;

  /* Sanity check */
  if (fileName==NULL) return FALSE;
  if (contdSpace==NULL) return FALSE;

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
  ret = deserialize_contdSpace_file(contdSpace, json);

 done:
  if (buffer) free(buffer);
  fclose(fp);

  json_decref(json);
  return ret;
}
