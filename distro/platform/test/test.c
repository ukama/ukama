/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "test.h"

#include "unity.h"
#include "usys_api.h"
#include "usys_dir.h"
#include "usys_error.h"
#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_process.h"
#include "usys_shm.h"
#include "usys_string.h"
#include "usys_sync.h"
#include "usys_thread.h"
#include "usys_timer.h"
#include "usys_types.h"

extern const char *usysErrCodes[];

void setUp(void) {
    // set stuff up here
}

void tearDown(void) {
    // clean stuff up here
}

/* Test error codes */
void test_usys_errors(){

    USysErrCodeIdx idx= USYS_BASE_ERR_CODE;
    for (;idx < ERR_MAX_ERR_CODE; idx++) {
        usys_log_trace("Error Code %d Error string %s", idx, usys_error(idx));
        TEST_ASSERT_EQUAL_STRING(usysErrCodes[idx-USYS_BASE_ERR_CODE], usys_error(idx));
    }
}

/* File Operations test */
void test_usys_fopen_should_return_file_pointer(void) {
    TEST_ASSERT_NOT_EQUAL(NULL, usys_fopen("./test/data/sample.txt", "r"));
}

void test_usys_fopen_should_return_null(void) {
    TEST_ASSERT_EQUAL(NULL, usys_fopen("./test/data/nosuchfile.txt", "r"));
}

void test_usys_fopen_create_new_file(void) {
    TEST_ASSERT_NOT_EQUAL(NULL, usys_fopen("./test/data/newfile.txt", "w+"));
}

void test_usys_fread_fwrite_fseek(void) {
    char data[] = "abcdef";
    char rdata[10];

    FILE* fp = usys_fopen("./test/data/newfile.txt", "w+");
    TEST_ASSERT_MESSAGE(fp!= NULL, "failed to open file.");
    if (fp != NULL) {
        int wbytes = usys_fwrite(data, sizeof(char), sizeof(data), fp);
        TEST_ASSERT_MESSAGE(wbytes == sizeof(data), "Bytes Written not matching");

        int seek = usys_fseek(fp, 0, SEEK_SET);
        TEST_ASSERT_MESSAGE(seek == 0, "Setting file pointer to beginning of file failed.");

        int rbytes = usys_fread(&rdata, sizeof(char), sizeof(data), fp);
        TEST_ASSERT_MESSAGE(rbytes == sizeof(data), "Bytes read not matching");

        TEST_ASSERT_MESSAGE(usys_memcmp(data, rdata, wbytes) == 0, "Read data not matching to written data");
    }

    TEST_ASSERT_MESSAGE(0 == usys_fclose(fp), "Failed to close file.");
}

void test_usys_fgets(void) {
    char data[] = "abcdef";
    char rdata[10];

    FILE* fp = usys_fopen("./test/data/newfile.txt", "r");
    TEST_ASSERT_MESSAGE(fp!= NULL, "failed to open file.");
    if (fp != NULL) {

        usys_fgets(rdata, sizeof(data), fp);

        TEST_ASSERT_MESSAGE(usys_memcmp(data, rdata, sizeof(data)) == 0, "Read data not matching to written data");

    }

    TEST_ASSERT_MESSAGE(0 == usys_fclose(fp), "Failed to close file.");
}

/* Test cases related to thread */

static int counter = 0;
USysThreadId thread[2];

/* Thread sync with mutex */
USysMutex mutex;
void *thread_function_with_mutex(void* arg) {
    USysError err;

    int *tstat = usys_emalloc(sizeof(int));

    USysThreadId tmpid = usys_thread_id();
    usys_log_trace("Thread id %ld waiting for mutex", tmpid);

    err = usys_mutex_lock(&mutex);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    /* Critical section begins */

    counter += 1;

    int temp = counter;

    usys_log_trace("Job %d has started", counter);

    TEST_ASSERT_MESSAGE(0 == usys_sleep(4), "Sleep failed.");

    *tstat = counter;

    usys_log_trace("Job %d has finished", counter);

    TEST_ASSERT_EQUAL_INT_MESSAGE(temp, counter, "critical section access breach");

    /* Critical section begins */

    err = usys_mutex_unlock(&mutex);

    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));


    pthread_exit(tstat);
}

void test_usys_threads_with_mutex() {
    int idx = 0;
    counter = 0;
    USysError err;
    int ret = 0;
    int arg[2] = {0};
    int *status = 0;
    err = usys_mutex_init(&mutex);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    while (idx < 2) {
        err = usys_thread_create( &thread[idx],
                NULL,
                &thread_function_with_mutex,
                (void*)&arg[idx]);
        TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

        idx++;
    }

    ret = usys_thread_join(thread[0], (void **)(&status));
    TEST_ASSERT_MESSAGE(ret == 0, "Failed to join thread 1.");
    usys_free(status);

    ret = usys_thread_join(thread[1], (void **)(&status));
    TEST_ASSERT_MESSAGE(ret == 0, "Failed to join thread 2.");
    usys_free(status);

    /* Mutex destroy */
    err = usys_mutex_destroy(&mutex);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

}

/* Thread sync with semaphore */

USysSem sem;

void *thread_function_with_semaphore(void* arg) {
    USysError err;

    int *tstat = usys_emalloc(sizeof(int));

    USysThreadId tmpid = usys_thread_id();
    usys_log_trace("Thread id %ld waiting for sem", tmpid);

    err = usys_sem_wait(&sem);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    /* Critical section begins */

    counter += 1;
    int temp = counter;

    usys_log_trace("Job %d has started", counter);

    TEST_ASSERT_MESSAGE(0 == usys_sleep(4), "Sleep failed.");

    *tstat = counter;

    usys_log_trace("Job %d has finished", counter);

    TEST_ASSERT_EQUAL_INT_MESSAGE(temp, counter, "critical section access breach");

    /* Critical section begins */

    err = usys_sem_post(&sem);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));


    pthread_exit(tstat);
}

void test_usys_threads_with_semaphore() {
    int idx = 0;
    counter = 0;
    USysError err;
    int ret = 0;
    int arg[2] = {0};
    int *status = 0;
    err = usys_sem_init(&sem, 1);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

    while (idx < 2) {
        err = usys_thread_create( &thread[idx],
                NULL,
                &thread_function_with_semaphore,
                (void*)&arg[idx]);
        TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

        idx++;
    }

    ret = usys_thread_join(thread[0], (void **)(&status));
    TEST_ASSERT_MESSAGE(ret == 0, "Failed to join thread 1.");
    usys_free(status);

    ret = usys_thread_join(thread[1], (void **)(&status));
    TEST_ASSERT_MESSAGE(ret == 0, "Failed to join thread 2.");
    usys_free(status);

    /* Mutex destroy */
    err = usys_sem_destroy(&sem);
    TEST_ASSERT_MESSAGE(ERR_NONE == err, usys_error(err));

}

/* Strings  test */

/* String copy */
void test_usys_strcpy() {
    char src[12] = "Test Case.";
    char dest[12] = {'\0'};

    /* String compare */
    int ret = usys_strcmp(src, dest);
    TEST_ASSERT_MESSAGE(ret != 0, "Strings are already same.");

    /* String copy */
    usys_strcpy(dest, src);

    /* String compare */
    ret = usys_strcmp(src, dest);
    TEST_ASSERT_MESSAGE(ret == 0, "Strings compare failed after copy");

}

/* string compare */
void test_usys_strncmp() {
    char str1[12] = "Test Case.";
    char str2[5] = "Test";
    char str3[5] = "test";

    int ret = usys_strncmp(str1, str2, usys_strlen(str2));
    TEST_ASSERT_EQUAL_INT(0, ret);

    ret = usys_strncmp(str1, str3, usys_strlen(str3));
    TEST_ASSERT_NOT_EQUAL_INT(0, ret);

}

/* String  length */
void test_usys_strlen() {
    char src[12] = "Test Case.";
    TEST_ASSERT_EQUAL_INT( 10, usys_strlen(src));
}

/* String cat */
void test_usys_strcat() {
    char exp[12] = "Test Case.";
    char src1[5] = "Test";
    char src2[7] = {" Case."};

    TEST_ASSERT_EQUAL_STRING(exp, usys_strcat(src1, src2));
}

/* Substring string */
void test_usys_strstr() {
    char str[] ="This is a simple string";
    char * pch;
    int ret = -1;
    pch = usys_strstr (str,"simple");
    if (pch != NULL)
        ret = usys_strncmp (pch,"simple", 6);
    TEST_ASSERT_EQUAL_INT(0, ret );
}

/* String memcmp, memcpy and memset */
void test_usys_memcmp_memset_memcpy() {
    char str[11] = "memory set";
    char mset[11] = "----------";
    char test[5] = "test";

    int ret = usys_memcmp(str, mset, usys_strlen(str));
    TEST_ASSERT_NOT_EQUAL_INT(0, ret );

    usys_memset (str,'-', usys_strlen(str));
    TEST_ASSERT_EQUAL_STRING(mset, str);

    usys_memcpy(str, test, usys_strlen(test));
    ret = usys_memcmp(str, test, usys_strlen(test));
    TEST_ASSERT_EQUAL_INT(0, ret );

}

/* String to token */
void test_usys_strtok() {
    char str[] = "memory - set . test";
    char *tok;
    char *tokens[] = {
            "memory",
            "set",
            "test"
    };
    tok = strtok(str, " -.");
    int ret = usys_memcmp(tok, tokens[0], usys_strlen(tokens[0]));
    TEST_ASSERT_EQUAL_INT(0, ret);
    int i = 0;
    while (tok != NULL) {
        i++;
        tok = strtok(NULL," -.");
        if (tok) {
            ret = usys_memcmp(tok, tokens[i], usys_strlen(tokens[i]));
            TEST_ASSERT_EQUAL_INT(0, ret);
        }
    }

}

/* Shared memory */

#define BLOCKSIZE           256
#define SHMFILE             "/shmfd"
#define PERM                0644
#define SEMAPHORE           "shmSem"
#define BLOCKDATA           "Testing shared memory."

void test_shm_writer() {

    int fd = usys_shm_open(SHMFILE, O_RDWR | O_CREAT, PERM);
    if (fd < 0) {
        usys_log_error("%s : Failed to create a shared memory", __FUNCTION__ );
        return;
    }

    usys_ftruncate(fd, BLOCKSIZE);

    caddr_t mem = usys_mmap(NULL, BLOCKSIZE, PROT_READ | PROT_WRITE, MAP_SHARED,
            fd, 0);
    if ((caddr_t) -1  == mem) {
        usys_log_error("%s : Failed to map a shared memory", __FUNCTION__ );
        return;
    }

    usys_log_trace("%s : shared mem address: %p [0..%d]\n", __FUNCTION__ , mem, BLOCKSIZE - 1);
    usys_log_trace("%s : shared mem file:       /dev/shm%s\n", __FUNCTION__, SHMFILE);

    /* Semaphore code to lock the shared mem */
    USysSem* sem = usys_sem_open(SEMAPHORE ,O_CREAT, PERM, 0);
    if (!sem)  {
        usys_log_error("%s : failed to create open semaphore", __FUNCTION__ );
        return;
    }

    /* Copy data to shared memory */
    usys_strcpy(mem, BLOCKDATA);

    /* Semaphore post */
    if (usys_sem_post(sem) < 0) {
        usys_log_error("%s : failed to post semaphore", __FUNCTION__ );
        return;
    }

    usys_sleep(10);

    /* Clean up */
    if (usys_munmap(mem, BLOCKSIZE) != 0) {
        usys_log_error("%s: failed to unmap shared memory", __FUNCTION__ );
        return;
    }

    close(fd);

    if(usys_sem_close(sem) != 0) {
        usys_log_error("%s: failed to close semaphore", __FUNCTION__ );
        return;
    }

    /* Unlink from the shared memory file */
    usys_shm_unlink(SHMFILE);

    usys_log_error("[%d] %s : completed.", usys_getpid(), __FUNCTION__ );
}

void test_shm_reader(char *readdata) {
    int fd = usys_shm_open(SHMFILE, O_RDWR, PERM);
    if (fd < 0) {
        usys_log_error("%s : Failed to create a shared memory", __FUNCTION__ );
        return;
    }

    caddr_t mem = usys_mmap(NULL, BLOCKSIZE, PROT_READ | PROT_WRITE, MAP_SHARED,
            fd, 0);
    if ((caddr_t) -1  == mem) {
        usys_log_error("%s : Failed to map a shared memory", __FUNCTION__ );
        return;
    }

    usys_log_trace("%s : shared mem address: %p [0..%d]\n", __FUNCTION__ , mem, BLOCKSIZE - 1);
    usys_log_trace("%s : shared mem file:       /dev/shm%s\n", __FUNCTION__, SHMFILE);

    /* Semaphore code to lock the shared mem */
    USysSem* sem = usys_sem_open(SEMAPHORE ,O_CREAT, PERM, 0);
    if (!sem)  {
        usys_log_error("%s : failed to create open semaphore", __FUNCTION__ );
        return;
    }

    /* Wait to acquire sem  */
    if (!usys_sem_wait(sem)) {

        /* Copy data from shared memory */
        usys_memcpy(readdata, mem,  usys_strlen(BLOCKDATA));

        usys_sem_post(sem);
    }

    /* Clean up */
    if (usys_munmap(mem, BLOCKSIZE) != 0) {
        usys_log_error("%s: failed to unmap shared memory", __FUNCTION__ );
        return;
    }

    close(fd);

    if(usys_sem_close(sem) != 0) {
        usys_log_error("%s: failed to close semaphore", __FUNCTION__ );
        return;
    }

    /* Unlink from the shared memory file */
    usys_shm_unlink(SHMFILE);

    usys_log_trace("[%d] %s : completed.", usys_getpid(), __FUNCTION__ );
}

/* Test function for shared memory */
void test_usys_shm(void) {
	USysPid child_pid;
    char readdata[BLOCKSIZE] = {'\0'};
    child_pid = usys_fork ();
    if ( child_pid == 0) {
    	/* Child process */
        usys_log_trace("[%d] child process for shm writer successfully created", usys_getpid());
        usys_log_trace ("[%d] child_PID = %d,parent_PID = %d\n",
        		usys_getpid(), usys_getpid(), usys_getppid( ) );
        test_shm_writer();
        usys_log_trace("[%d] child process for shm writer completed", usys_getpid());
        usys_qexit(0);
    } else if (child_pid > 0) {
    	/* Parent process */
        usys_sleep(2);
        usys_log_trace("[%d] shm reader successfully created!", usys_getpid());
        test_shm_reader(readdata);
        usys_wait(NULL);

        int ret = usys_memcmp(readdata, BLOCKDATA, usys_strlen(BLOCKDATA));
        TEST_ASSERT_EQUAL_INT(ret, 0);
    }
    usys_log_trace("[%d] test shm completed", usys_getpid());
}

/* Timer callback */
void timer_stat(int signum){
    struct timeval current_time;
    usys_gettimeofday(&current_time);
    usys_log_info("Timer Callback seconds : %ld micro seconds : %ld",
            current_time.tv_sec, current_time.tv_usec);
}

/* Test timer  Just check the creation of timer
 * Not a high resolution timer */
void test_usys_timer(void) {

    struct timeval current_time;
    usys_gettimeofday(&current_time);
    usys_log_info("seconds : %ld micro seconds : %ld",
        current_time.tv_sec, current_time.tv_usec);

    /* Create a timer */
    int ret = usys_timer(100000, timer_stat);
    TEST_ASSERT_EQUAL_INT( 1, ret );

    /* Come time for timer to run */
    unsigned int num = 0xFFFFFFFF;
    while(num >0) {
        num--;
    };

    /* Stop timer */
    ret = usys_timer(0, timer_stat);
    TEST_ASSERT_EQUAL_INT( 1, ret );

    usys_log_trace("Waiting to stop timer.");
    /* Some time for timer to stop */
    num = 0xFFFFFF;
    while(num >0) {
        num--;
    };
}

void test_usys_read_write_strings_to_file() {
    char* fileName = "./test/data/testfile.txt";
    char buff[32] = "Ukama test file created.";
    char testBuff[32] = { '\0' };
    int size = usys_strlen(buff);

    /* Create file */
    usys_file_init(fileName);

    /* Write to file */
    usys_file_write(fileName, buff, 0, size);

    /* Read from file */
    usys_file_read(fileName, testBuff, 0, size);

    /* Compare read and written string */
    TEST_ASSERT_EQUAL_STRING(buff, testBuff);

    TEST_ASSERT_EQUAL_INT( 0 , usys_remove(fileName));
}

void test_usys_read_write_numbers_to_file() {
    int type = 0xabcd;
    int testType;
    char ty1[4];
    char ty[4];
    char* fileName = "./test/data/testfile.txt";
    usys_memcpy(ty, &type, 4);

    /* Create file */
    usys_file_init(fileName);

    /* Write a number */
    usys_file_write(fileName, ty, 0, 4);

    /* Read a number */
    usys_file_read(fileName, ty1, 0, 4);

    usys_memcpy(&testType, ty1, 4);

    /* Compare read and written number */
    TEST_ASSERT_EQUAL_CHAR_ARRAY(ty, ty1, 4);

    /* Remove file */
    TEST_ASSERT_EQUAL_INT( 0 , usys_remove(fileName));
}

void test_usys_read_write_arrays_to_file() {
    char* fileName = "./test/data/testfile.txt";
    uint16_t write[3] = { 455, 35, 6335 };
    uint16_t read[3] = { 0 };

    /* Create a file */
    usys_file_init(fileName);

    /* Write numbers */
    usys_file_write_number(fileName, write, 22, 3, sizeof(uint16_t));

    /* Read */
    usys_file_read_number(fileName, read, 22, 3, sizeof(uint16_t));

    /* Compare read and written bytes */
    TEST_ASSERT_EQUAL_UINT16_ARRAY(write, read, 3);

    /* Remove file */
    TEST_ASSERT_EQUAL_INT( 0 , usys_remove(fileName));
}

void test_usys_write_failure_no_file_exist() {
    char* fileName = "./test/data/testfile.txt";
    char buff[32] = "Ukama test file created.";
    char testBuff[32] = { '\0' };

    int size = usys_strlen(buff);

    /* Write to file */
    TEST_ASSERT_GREATER_THAN( size ,
    		usys_file_write(fileName, buff, 0, size));

    /* Remove file */
    TEST_ASSERT_EQUAL_INT( 0 , usys_remove(fileName));
}

/* Directory related test cases */

void test_usys_open_read_close_dir() {
	struct dirent *de;

	/* Open directory */
	DIR *dr = usys_opendir(".");
	if (dr == NULL) {
		TEST_ASSERT_MESSAGE(dr != NULL, usys_error(errno));
	}

	/* Read directory */
	while ((de = usys_readdir(dr)) != NULL) {
		usys_log_trace("%s\n", de->d_name);
	}

	TEST_ASSERT_EQUAL_INT( 0 , usys_closedir(dr));
}

void test_usys_opendir_fail() {
	struct dirent *de;

	/* Open directory */
	DIR *dr = usys_opendir("/xyz");
	usys_log_warn("Error: %s", usys_error(errno));
	TEST_ASSERT_MESSAGE(dr == NULL, "Directory wasn't expected");

}

void test_usys_getcwd_ch_dir() {
	struct dirent *de;
	/* cuurent working directory */
	char cwd[256];
	TEST_ASSERT_MESSAGE(getcwd(cwd, sizeof(cwd)) != NULL, "Couldn't get "
			"current working dir.");
	usys_log_trace("Current working dir is %s", cwd);

	/* open directory */
	DIR *dr = usys_opendir(".");
	if (dr == NULL) {
		TEST_ASSERT_MESSAGE(dr != NULL, usys_error(errno));
	}

	while ((de = usys_readdir(dr)) != NULL) {
		usys_log_trace("%s\n", de->d_name);
	}

	/* Close directory */
	TEST_ASSERT_EQUAL_INT( 0 , usys_closedir(dr));

	/* Change directory */
	TEST_ASSERT_MESSAGE( usys_chdir("/") == 0, "Change directory failed");

	/* Open directory */
	dr = usys_opendir(".");
	if (dr == NULL) {
		TEST_ASSERT_MESSAGE(dr != NULL, usys_error(errno));
	}

	/* Read directory */
	while ((de = usys_readdir(dr)) != NULL) {
		usys_log_trace("%s\n", de->d_name);
	}

	/* Close directory */
	TEST_ASSERT_EQUAL_INT( 0 , usys_closedir(dr));

	/* Change directory */
	TEST_ASSERT_MESSAGE( usys_chdir(cwd) == 0, "Change directory failed");
}

void test_usys_seek_tell_dir() {
	struct dirent *de ;
	DIR *dir;
	long int offset = 0;
	long int loc= 0;
	int idx = 0;

	/* Open dir */
	dir = usys_opendir(".");

	/* Read directory */
	while ((de = usys_readdir(dir))!= NULL) {

		/* Tell dir */
		offset = usys_telldir(dir);
		TEST_ASSERT_GREATER_THAN_INT64( 0 , offset);
		if (idx++ == 2) {
			loc = offset;
		}

		usys_log_debug("DName:%s Offset: %ld ", de->d_name, offset) ;
	}

	/* Seek dir */
	usys_seekdir(dir, loc);

	usys_log_debug("reading again");
	while ((de = usys_readdir(dir))!= NULL) {

		/* Tell dir */
		offset = usys_telldir(dir);
		TEST_ASSERT_GREATER_THAN_INT64( 0 , offset);

		usys_log_debug("DName:%s Offset: %ld", de->d_name, offset);
	}

	/* Close directory */
	TEST_ASSERT_EQUAL_INT( 0 , usys_closedir(dir));
}

/* Test process related API's */

void test_usys_fork_wait_pid_ppid_prgp() {
	USysPid child_pid;
	int waitstatus = 0;

	USysPid c_pid = usys_getpid();
	usys_log_debug("Current Process id %d, parent_PID = %d groupid %d",
			c_pid, usys_getppid(), usys_getpgrp());

    child_pid = usys_fork ();
    if ( child_pid == 0) {
    	/* Child process */
        usys_log_trace(" [%d] Child process created", usys_getpid());
        usys_log_trace ("[%d] child_PID = %d, parent_PID = %d groupid %d",
        		usys_getpid(), usys_getpid(), usys_getppid(), usys_getpgrp());

        usys_log_trace("[%d] child process completed", usys_getpid());
        usys_qexit(10);

    } else if (child_pid > 0) {

    	/* Parent process */
    	usys_log_trace("[%d] Back in parent process !", usys_getpid());

    	/* Waitng for child */
    	usys_wait(&waitstatus);
    	if (waitstatus == -1) {
    		usys_log_error("Failed to wait for child exit status.");
    		TEST_ASSERT_MESSAGE(waitstatus, "Failed to wait for "
    				"child exit status.");
    	}

    	if (WIFEXITED(waitstatus)) {
    		int exitstatus = WEXITSTATUS(waitstatus);
    		usys_log_error("Exit status of a child process %d",
    				exitstatus);
    		TEST_ASSERT_MESSAGE(exitstatus == 10 , "Failed to get proper "
    				"exit status from child.");
    	}
    	else {
    		usys_log_error("FAILURE: unknown child exited");
    	}
    }
}

void test_usys_get_set_rlimit() {
	  struct rlimit limit;

	  /* Get the file size resource limit. */
	  TEST_ASSERT_MESSAGE(usys_getrlimit(RLIMIT_FSIZE, &limit) == 0,
			  "Reading resource limit failed.");
	  usys_log_debug("Soft limit is %llu hard limit is %llu", limit.rlim_cur, limit.rlim_max);

	  /* Set the file size resource limit. */
	  limit.rlim_cur = 65535;
	  limit.rlim_max = 65535;
	  TEST_ASSERT_MESSAGE(usys_setrlimit(RLIMIT_FSIZE, &limit) == 0,
			  "Setting resource limit failed.");

	  /* Get the file size resource limit. */
	  TEST_ASSERT_MESSAGE(usys_getrlimit(RLIMIT_FSIZE, &limit) == 0,
			  "Reading resource limit failed.");
	  usys_log_debug("Soft limit is %llu hard limit is %llu", limit.rlim_cur, limit.rlim_max);

}
