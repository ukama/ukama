/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "utils/irqdb.h"

#include "headers/errorcode.h"
#include "headers/utils/list.h"
#include "headers/utils/log.h"
#include "utils/vnodealert.h"

ListInfo irqdb;

/* Remove IRQ Config from the irqdb */
static void remove_irq(void *ptr) {
    ListNode *node = (ListNode *)ptr;
    if (node) {
        if (node->data) {
            free(node->data);
            node->data = NULL;
        }
        free(node);
        node = NULL;
    }
}

/* Comparing dev object  with IRQCfg->dev */
int compare_dev(void *lip, void *sd) {
    IRQCfg *ip = (IRQCfg *)lip;
    IRQCfg *op = (IRQCfg *)sd;
    int ret = 0;
    /* If module if  and device name, disc, type matches it means devices is same.*/
    if (!strcmp(ip->obj.mod_UUID, op->obj.mod_UUID) &&
        !strcmp(ip->obj.name, op->obj.name) &&
        !strcmp(ip->obj.disc, op->obj.disc) && (ip->obj.type == op->obj.type) &&
        !strcmp(ip->fname, op->fname)) {
        ret = 1;
    }
    return ret;
}

/* Searching IRQ for a device in the IRQDB */
static IRQCfg *search_device_object(IRQCfg *sdev) {
    IRQCfg *fcfg = NULL;
    /* TODO::Check if it returns proper data.*/
    fcfg = list_search(&irqdb, sdev);
    if (fcfg) {
        log_debug(
            "IRQDB:: IRQ %lu for Device Name %s, Disc: %s Module UUID: %s found in IRQDB.",
            fcfg->pthread, fcfg->obj.name, fcfg->obj.disc, fcfg->obj.mod_UUID);
    } else {
        log_debug(
            "IRQDB:: IRQ %lu for Device Name %s, Disc: %s Module UUID: %s not found in IRQDB.",
            fcfg->pthread, fcfg->obj.name, fcfg->obj.disc, fcfg->obj.mod_UUID);
        if (fcfg) {
            free(fcfg);
            fcfg = NULL;
        }
    }
    return fcfg;
}

/* Cancels IRQ thread and return 1 as a success. */
static int irqdb_unregister_irq(void *data) {
    int ret = 0;
    IRQCfg *cfg = (IRQCfg *)data;
    log_trace(
        "IRQDB:: De-registering IRQ %lu for Device Name %s, Disc: %s Module UUID: %s found.",
        cfg->pthread, cfg->obj.name, cfg->obj.disc, cfg->obj.mod_UUID);
    /* Cancel thread */
    ret = pthread_cancel(cfg->pthread);
    if (ret) {
        ret = 0;
    } else {
        pthread_join(cfg->pthread, NULL);
        ret = 1;
    }
    return ret;
}

void thread_cleanup_handler(void *arg) {
    log_debug("IRQDB:: Cleaning up thread %ld.", pthread_self());
}

/* Thread for IRQ handling.*/
void *threaded_irq_handler(void *args) {
    pthread_cleanup_push(thread_cleanup_handler, NULL);
    pthread_setcancelstate(PTHREAD_CANCEL_ENABLE, NULL);
    IRQCfg *cfg = (IRQCfg *)args;
    if (!cfg) {
        log_debug(
            "IRQDB:::: No IRQ Config found. Failed to start thread for IRQ handling.");

    } else {
        log_debug(
            "IRQDB:: IRQ %lu for Device Name %s, Disc: %s Module UUID: %s started thread %lu for IRQ handling.",
            cfg->pthread, cfg->obj.name, cfg->obj.disc, cfg->obj.mod_UUID,
            pthread_self());

#ifdef TARGET
        //TODO
        while (true) {
            sem_wait(&(cfg->irq_sem));
            log_debug(
                "IRQDB::::Received IRQ %d for Device Name %s, Disc: %s Module UUID: %s found.",
                cfg->irq, cfg->obj.name, cfg->obj.disc, cfg->obj.mod_UUID);
            if (cfg->cb) {
                cfg->cb(&cfg->obj);
            }
        }
#else
        /* Virtual node files are regular files. So alerts are implemented using inotify syscall.*/
        poll_file(cfg);
#endif
    }
    pthread_cleanup_pop(0);
    pthread_exit(NULL);
}

/* Search interrupt for device.*/
IRQCfg *irqdb_search_for_device_irq(Device *dev) {
    IRQCfg scfg = { .pthread = 0, .cb = 0 };
    memcpy(&scfg.obj, &dev->obj, sizeof(DevObj));
    memset(scfg.fname, '\0', 64);
    memcpy(&scfg.fname, dev->sysfile, strlen(dev->sysfile));
    return search_device_object(&scfg);
}

/* Register IRQ for device */
int irqdb_register_for_device_irq(IRQSrcInfo *rsrc, SensorCallbackFxn cb,
                                  ThreadedIRQEnable IRQ_enable) {
    int ret = 0;
    /* Search in list if IRQ is already created for device.*/
    /*TODO: Based on the alert sysfs: we need to decide if one alert file per device or
	 *  for each property which can be alert we need a alert sysfs file.
	 *  In single sysfs file case we need to compare just object other wise  we need to compare sysfs also.
	 *  For now multiple per sensor.
	 */
    IRQCfg scfg = { .pthread = 0, .cb = cb };
    memcpy(&scfg.obj, &rsrc->obj, sizeof(DevObj));
    memset(scfg.fname, '\0', 64);
    memcpy(&scfg.fname, rsrc->src.sysfs_name, strlen(rsrc->src.sysfs_name));

    ret = list_if_element_found(&irqdb, &scfg);
    if (!ret) {
        /* if no IRQ exist for device object create and add one to list.*/
        list_append(&irqdb, &scfg);
    } else {
        /*This will print the info of the IRQ config.*/
        IRQCfg *cfg = search_device_object(&scfg);
        if (cfg) {
            free(cfg);
            cfg = NULL;
        }
    }

    log_trace(
        "IRQDB:: Registering IRQ for Device Name %s, Disc: %s Module UUID: %s found.",
        rsrc->obj.name, rsrc->obj.disc, rsrc->obj.mod_UUID);

    /* Create thread */
    void *args = irqdb.tail->data;
    IRQCfg *tcfg = (IRQCfg *)irqdb.tail->data;
    if (pthread_create(&(tcfg->pthread), NULL, threaded_irq_handler, args)) {
        /* Remove from list */
        if (list_remove(&irqdb, &rsrc->obj)) {
            ret = ERR_UBSP_LIST_DEL_FAILED;
        }
        ret = ERR_UBSP_THREAD_CREATE_FAIL;
    }
    log_trace(
        "IRQDB:: Registered IRQ %lu for Device Name %s, Disc: %s Module UUID: %s found.",
        tcfg->pthread, rsrc->obj.name, rsrc->obj.disc, rsrc->obj.mod_UUID);

    irqdb_print_list();

    /* TODO Try doing this in higher layer after call to irqdb_register_for_device_IRQ
     * on success enable IRQ.
     */

    /*TODO::Enable interrupt here.*/
    if (rsrc->type == IRQ_SYSFS) {
        //IRQ_enable(rsrc);
        //IRQ_enable(rsrc->src.sysfs_name);
    } else {
        //IRQ_enable(rsrc);
    }
    return ret;
}

/* De-register IRQ for device */
int irqdb_deregister_for_device_irq(IRQSrcInfo *rsrc,
                                    ThreadedIRQEnable IRQ_disable) {
    int ret = 0;
    irqdb_print_list();
    /*Creating dummy cfg*/
    IRQCfg scfg = {
        .pthread = 0,
        .cb = NULL,
    };
    memcpy(&scfg.obj, &rsrc->obj, sizeof(DevObj));
    memset(scfg.fname, '\0', 64);
    memcpy(&scfg.fname, rsrc->src.sysfs_name, strlen(rsrc->src.sysfs_name));

    /*This will print the required info.*/
    IRQCfg *cfg = search_device_object(&scfg);
    if (cfg) {
        /* Return 1 if cancel is success */
        ret = irqdb_unregister_irq(cfg);
        if (ret) {
            ret = 0;
            /* Remove from list */
            if (list_remove(&irqdb, cfg)) {
                ret = ERR_UBSP_LIST_DEL_FAILED;
            }
        } else {
            ret = ERR_UBSP_THREAD_CANCEL_FAIL;
            log_error(
                "IRQDB(%d):: Failed to cancel IRQ thread for Device Name %s, Disc: %s Module UUID: %s found.",
                ret, cfg->obj.name, cfg->obj.disc, cfg->obj.mod_UUID);
        }
        free(cfg);
        cfg = NULL;
    } else {
        ret = ERR_UBSP_DEV_IRQ_NOT_REG;
        log_error(
            "IRQDB(%d):: De-registering failed No IRQ for Device Name %s, Disc: %s Module UUID: %s found.",
            ret, rsrc->obj.name, rsrc->obj.disc, rsrc->obj.mod_UUID);
    }
    /* Just for confirmation */
    irqdb_print_list();
    return ret;
}

/* IRDB init. */
void irqdb_init() {
    list_new(&irqdb, sizeof(IRQCfg), remove_irq, compare_dev, NULL);
    log_trace("IRQDB:: IRQDB initialized.");
}

void irqdb_exit() {
    log_trace("IRQDB:: Current IRQDB list.");
    irqdb_print_list();
    list_for_each(&irqdb, irqdb_unregister_irq);
    list_destroy(&irqdb);
    log_trace("IRQDB:: Removing IRQDB.");
}

int irqdb_print_irq(void *data) {
    uint8_t ret = 0;
    if (!data) {
        return ret;
    }
    IRQCfg *cfg = (IRQCfg *)data;
    log_trace(
        "********************************************************************************");
    log_trace("* ThreadId                  : %lu", cfg->pthread);
    log_trace("* Name                      : %s", cfg->obj.name);
    log_trace("* Disc               	   : %s", cfg->obj.disc);
    log_trace("* Module UUID               : %s", cfg->obj.mod_UUID);
    log_trace("* Type                      : %d", cfg->obj.type);
    log_trace("* SysFile Name:             : %s", cfg->fname);
    log_trace(
        "********************************************************************************");
    ret = 1;
    return ret;
}

void irqdb_print_list() {
    log_trace(
        "********************************************************************************");
    log_trace(
        "****************************** IRQDB LIST *************************************");
    if (irqdb.logicalLength > 0) {
        log_trace("Logical length is                          %d.",
                  irqdb.logicalLength);
        log_trace("Element size is                            %d.",
                  irqdb.elementSize);
        log_trace("Head is at                                 %p.", irqdb.head);
        log_trace("Tail is at                            	  %p.", irqdb.tail);
        log_trace(
            "********************************************************************************");
        list_for_each(&irqdb, irqdb_print_irq);
        log_trace(
            "********************************************************************************");
    } else {
        log_trace("IRQDB is empty.");
        return;
    }
}
