/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * init -- UkamaOS minimal init
 */

/*
 * basic init flow is:
 * fs.init -> fs initialization
 * exec into init:
 *   setup signal
 *   setup console
 * fork and exec sub.init
 * fork and exec lxce.d
 * loop forever:
 *   wait for lxce.d to exit
 *   process signals
 */

#include <signal.h>
#include <unistd.h>
#include <sys/reboot.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <termios.h>
#include <stdlib.h>
#include <errno.h>
#include <sys/wait.h>
#include <string.h>
#include <stdio.h>
#include <stdarg.h>

#include "init.h"

/* 
 * setup_console -- redirect the stdin, stdout and stderr to console
 *
 */
static void setup_console(void) {

  char *c;
  int fd;

  c = getenv("CONSOLE");
  if (!c) {
    c = getenv("console");
  }

  if (c) {
    fd = open(c, O_RDWR | O_NONBLOCK | O_NOCTTY);
    if (fd >= 0) {
      dup2(fd, STDIN_FILENO);
      dup2(fd, STDOUT_FILENO);
      dup2(fd, STDERR_FILENO);
    }
  }
}

/*
 * setup_term --
 *
 */
static void setup_term(void) {

  struct termios tty;

  if (tcgetattr(STDIN_FILENO, &tty) != 0) {
    return;
  }

  /* special characters */
  tty.c_cc[VINTR]  = 3;	  /* C-c */
  tty.c_cc[VQUIT]  = 28;  /* C-\ */
  tty.c_cc[VERASE] = 127; /* C-? */
  tty.c_cc[VKILL]  = 21;  /* C-u */
  tty.c_cc[VEOF]   = 4;	  /* C-d */
  tty.c_cc[VSTART] = 17;  /* C-q */
  tty.c_cc[VSTOP]  = 19;  /* C-s */
  tty.c_cc[VSUSP]  = 26;  /* C-z */

  /* control modes */
  tty.c_cflag &= CBAUD | CBAUDEX | CSIZE | CSTOPB | PARENB | PARODD | CRTSCTS;
  tty.c_cflag |= CREAD | HUPCL | CLOCAL;

  /* input modes */
  tty.c_iflag = ICRNL | IXON | IXOFF;

  /* output modes */
  tty.c_oflag = OPOST | ONLCR;

  /* local modes */
  tty.c_lflag = ISIG | ICANON | ECHO | ECHOE | ECHOK | ECHOCTL | ECHOKE
    | IEXTEN;

  tcsetattr(STDIN_FILENO, TCSANOW, &tty);
}

/*
 * log_message -- log message to default stdout
 *
 */

ssize_t log_message(int fd, const char *fmt, ...) {

  ssize_t count;
  va_list args;
  char msg[128] = {0};
  int len;

  msg[0] = '\r';
  va_start(args, fmt);
  len = 1 + vsnprintf(msg + 1, sizeof(msg) - 2, fmt, args);
  if (len > sizeof(msg) - 2) {
    len = sizeof(msg) - 2;
  }
  msg[len++] = '\n';
  msg[len]   = '\0';
  va_end(args);

  for (;;) {
  
    count = write(fd, msg, len);
    if (count >= 0 || errno != EINTR)
      break;

    errno = 0; /* reset and repeat */
  }

  return count;
}

/*
 * run_task -- fork/exec to run the task. Wait for the child process to
 *             complete execution, if needed
 *
 */
static pid_t run_task(char *taskExe, int wait) {

  sigset_t set;
  pid_t childPid;
  char *cmd[10];
  int status=0, exitStatus;

  switch (childPid = fork()) {
  case -1:
    log_message(STDOUT_FILENO, "Can not fork");
    return childPid;
    break;

  case 0: /* Child process */
    sigfillset(&set);
    sigprocmask(SIG_UNBLOCK, &set, NULL);
    if (strcmp(taskExe, "sub.init")==0) {
      cmd[0] = (char *)"/sbin/sub.init";
      cmd[1] = (char *)"-c";
      cmd[2] = NULL;
    } else if (strcmp(taskExe, "lxce")==0) {
      cmd[0] = (char *)"/sbin/lxce.d";
      cmd[1] = (char *)"-c";
      cmd[2] = (char *)"/conf/lxce_config.toml";
      cmd[3] = (char *)"-m";
      cmd[4] = (char *)"/conf/manifest.json";
      cmd[5] = NULL;
    }

    log_message(STDOUT_FILENO, "Running %s ...", taskExe);
    execv(cmd[0], cmd); 
    _exit(127); /* something went wrong. */
    break; /* never executed. End of child process */

  default: /* Parent */
    if (wait) { /* We will block the main thread for sub.init only */
      waitpid(childPid, &status, 0);
    } else {
      return childPid;
    }
    break;
  }

  if (WIFEXITED(status)) { /* proper termination by calling exit */
    exitStatus = WEXITSTATUS(status);
    switch(exitStatus) {
    case 127:
      log_message(STDOUT_FILENO, "Task execution failed: %s", taskExe);
      break;
    case 0:
      log_message(STDOUT_FILENO, "Task execution success: %s", taskExe);
      break;
    default:
      log_message(STDOUT_FILENO, "Task return invalid code: %d. Continue: %d",
		  taskExe, wait);
      break;
    }
  }

  return childPid;
}

/*
 * process_signals --
 *
 */
static void process_signals(sigset_t *sigset, struct timespec *tspec) {

  int sig;

  sig = sigtimedwait(sigset, NULL, tspec);
  if (sig <= 0) {
    return;
  }

  /* Ignore SIGHUP and SIGCHLD */
  if (sig == SIGHUP || sig == SIGCHLD) { 
    return;
  }

  if (sig == SIGQUIT) {
    /* restart */
  }

  if ((1 << sig) &(0 | (1 << SIGUSR1)
		   | (1 << SIGUSR2)
		   | (1 << SIGTERM))) {
    /* reboot/poweroff */
  }
}

int main(int argc, char **argv) {

  sigset_t sigset;
  pid_t lxcePid, pid;
  int status, exitStatus;
  char *taskExe;

  if (getpid() != 1) {
    log_message(STDERR_FILENO, "init must be run as PID 1");
    _exit(1);
  }
  
  sigemptyset(&sigset);
  sigaddset(&sigset, SIGINT);
  sigaddset(&sigset, SIGPWR);
  sigaddset(&sigset, SIGUSR1);
  sigaddset(&sigset, SIGTERM);
  sigaddset(&sigset, SIGUSR2);
  sigaddset(&sigset, SIGCHLD);

  /* block them all */
  sigprocmask(SIG_BLOCK, &sigset, NULL);

  /* disable CTL-ALT-DEL */
  reboot(RB_DISABLE_CAD);

  /* setup console and term */
  setup_console();
  setup_term();
  chdir("/");
  setsid();

  /* setup environment variables */
  putenv((char *) "HOME=/");
  putenv((char *) "PATH=/sbin:/usr/sbin");
  putenv((char *) "SHELL=/bin/sh");
  putenv((char *) "USER=root");

  log_message(STDOUT_FILENO, "init started: %s", UKAMA_BANNER);

  /* sub.init setup mount, clock, loopback, limits, hostname and resolv.conf */
  run_task("/sbin/sub.init", TRUE);

  /* run lxce.d - the minimal contained apps engine */
  lxcePid = run_task("/sbin/lxce", FALSE);

  /* reached my nirvana state, forever */
  while (TRUE) {

    process_signals(&sigset, NULL);

    /* Wait for any child process(es) to exit */
    while (TRUE) {
      pid = waitpid(-1, &status, WNOHANG);
      /* if error or timeout, repeat again */
      if (pid <= 0) break;

      if (pid == lxcePid) { /* restart the lxce.d? */
	taskExe = "lxce";
	if (WIFEXITED(status)) { /* proper termination by calling exit */
	  exitStatus = WEXITSTATUS(status);
	  switch(exitStatus) {
	  case 127:
	    log_message(STDOUT_FILENO, "Task execution failed: %s", taskExe);
	    break;
	  case 0:
	    log_message(STDOUT_FILENO, "Task execution success: %s", taskExe);
	    break;
	  default:
	    log_message(STDOUT_FILENO, "Task return invalid code: %d Cont: %d",
			taskExe, wait);
	    break;
	  }
	}
      } else {
	log_message(STDOUT_FILENO, "Process terminated: %d", (int)pid);
	break;
      }
    } /* inner while loop */
    sleep(1);
  }
}
