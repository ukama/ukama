/* Copyright 2012 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* Task scheduling / events module for Chrome EC operating system */

#ifndef __CROS_EC_TASK_H
#define __CROS_EC_TASK_H

#include "common.h"
#include "compile_time_macros.h"
#include "task_id.h"

/* Task event bitmasks */
/* Tasks may use the bits in TASK_EVENT_CUSTOM_BIT for their own events */
#define TASK_EVENT_CUSTOM_BIT(x) BUILD_CHECK_INLINE(BIT(x), BIT(x) & 0x0ffff)

/* Used to signal that sysjump preparation has completed */
#define TASK_EVENT_SYSJUMP_READY BIT(16)

/* Used to signal that IPC layer is available for sending new data */
#define TASK_EVENT_IPC_READY	BIT(17)

#define TASK_EVENT_PD_AWAKE	BIT(18)

/* npcx peci event */
#define TASK_EVENT_PECI_DONE	BIT(19)

/* I2C tx/rx interrupt handler completion event. */
#ifdef CHIP_STM32
#define TASK_EVENT_I2C_COMPLETION(port) \
				(1 << ((port) + 20))
#define TASK_EVENT_I2C_IDLE	(TASK_EVENT_I2C_COMPLETION(0))
#define TASK_EVENT_MAX_I2C	6
#ifdef I2C_PORT_COUNT
#if (I2C_PORT_COUNT > TASK_EVENT_MAX_I2C)
#error "Too many i2c ports for i2c events"
#endif
#endif
#else
#define TASK_EVENT_I2C_IDLE	BIT(20)
#endif

/* DMA transmit complete event */
#define TASK_EVENT_DMA_TC       BIT(26)
/* ADC interrupt handler event */
#define TASK_EVENT_ADC_DONE	BIT(27)
/* task_reset() that was requested has been completed */
#define TASK_EVENT_RESET_DONE   BIT(28)
/* task_wake() called on task */
#define TASK_EVENT_WAKE		BIT(29)
/* Mutex unlocking */
#define TASK_EVENT_MUTEX	BIT(30)
/*
 * Timer expired.  For example, task_wait_event() timed out before receiving
 * another event.
 */
#define TASK_EVENT_TIMER	(1U << 31)

/* Maximum time for task_wait_event() */
#define TASK_MAX_WAIT_US 0x7fffffff

/**
 * Disable CPU interrupt bit.
 *
 * This might break the system so think really hard before using these. There
 * are usually better ways of accomplishing this.
 */
void interrupt_disable(void);

/**
 * Enable CPU interrupt bit.
 */
void interrupt_enable(void);

/**
 * Return true if we are in interrupt context.
 */
int in_interrupt_context(void);

/**
 * Return current interrupt mask. Meaning is chip-specific and
 * should not be examined; just pass it to set_int_mask() to
 * restore a previous interrupt state after interrupt_disable().
 */
uint32_t get_int_mask(void);

/**
 * Set interrupt mask. As with interrupt_disable(), use with care.
 */
void set_int_mask(uint32_t val);

/**
 * Set a task event.
 *
 * If the task is higher priority than the current task, this will cause an
 * immediate context switch to the new task.
 *
 * Can be called both in interrupt context and task context.
 *
 * @param tskid		Task to set event for
 * @param event		Event bitmap to set (TASK_EVENT_*)
 * @param wait		If non-zero, after setting the event, de-schedule the
 *			calling task to wait for a response event.  Ignored in
 *			interrupt context.
 * @return		The bitmap of events which occurred if wait!=0, else 0.
 */
uint32_t task_set_event(task_id_t tskid, uint32_t event, int wait);

/**
 * Wake a task.  This sends it the TASK_EVENT_WAKE event.
 *
 * @param tskid		Task to wake
 */
static inline void task_wake(task_id_t tskid)
{
	task_set_event(tskid, TASK_EVENT_WAKE, 0);
}

/**
 * Return the identifier of the task currently running.
 */
task_id_t task_get_current(void);

/**
 * Return a pointer to the bitmap of events of the task.
 */
uint32_t *task_get_event_bitmap(task_id_t tskid);

/**
 * Wait for the next event.
 *
 * If one or more events are already pending, returns immediately.  Otherwise,
 * it de-schedules the calling task and wakes up the next one in the priority
 * order.  Automatically clears the bitmap of received events before returning
 * the events which are set.
 *
 * @param timeout_us	If > 0, sets a timer to produce the TASK_EVENT_TIMER
 *			event after the specified micro-second duration.
 *
 * @return The bitmap of received events.
 */
uint32_t task_wait_event(int timeout_us);

/**
 * Wait for any event included in an event mask.
 *
 * If one or more events are already pending, returns immediately.  Otherwise,
 * it de-schedules the calling task and wakes up the next one in the priority
 * order.  Automatically clears the bitmap of received events before returning
 * the events which are set.
 *
 * @param event_mask	Bitmap of task events to wait for.
 *
 * @param timeout_us	If > 0, sets a timer to produce the TASK_EVENT_TIMER
 *			event after the specified micro-second duration.
 *
 * @return		The bitmap of received events. Includes
 *			TASK_EVENT_TIMER if the timeout is reached.
 */
uint32_t task_wait_event_mask(uint32_t event_mask, int timeout_us);

/**
 * Prints the list of tasks.
 *
 * Uses the command output channel.  May be called from interrupt level.
 */
void task_print_list(void);

/**
 * Returns the name of the task.
 */
const char *task_get_name(task_id_t tskid);

#ifdef CONFIG_TASK_PROFILING
/**
 * Start tracking an interrupt.
 *
 * This must be called from interrupt context (!) before the interrupt routine
 * is called.
 */
void task_start_irq_handler(void *excep_return);
void task_end_irq_handler(void *excep_return);
#else
#define task_start_irq_handler(excep_return)
#endif

/**
 * Change the task scheduled to run after returning from the exception.
 *
 * If task_send_event() has been called and has set need_resched flag,
 * re-computes which task is running and eventually swaps the context
 * saved on the process stack to restore the new one at exception exit.
 *
 * This must be called from interrupt context (!) and is designed to be the
 * last call of the interrupt handler.
 */
void task_resched_if_needed(void *excep_return);

/**
 * Initialize tasks and interrupt controller.
 */
void task_pre_init(void);

/**
 * Start task scheduling.  Does not normally return.
 */
int task_start(void);

/**
 * Return non-zero if task_start() has been called and task scheduling has
 * started.
 */
int task_start_called(void);

#ifdef CONFIG_FPU
/**
 * Clear floating-point used flag for currently executing task. This means the
 * FPU regs will not be stored on context switches until the next time floating
 * point is used for currently executing task.
 */
void task_clear_fp_used(void);
#endif

/**
 * Mark all tasks as ready to run and reschedule the highest priority task.
 */
void task_enable_all_tasks(void);

/**
 * Enable a task.
 */
void task_enable_task(task_id_t tskid);

/**
 * Disable a task.
 *
 * If the task disable itself, this will cause an immediate reschedule.
 */
void task_disable_task(task_id_t tskid);

/**
 * Enable an interrupt.
 */
void task_enable_irq(int irq);

/**
 * Disable an interrupt.
 */
void task_disable_irq(int irq);

/**
 * Software-trigger an interrupt.
 */
void task_trigger_irq(int irq);

/*
 * A task that supports resets may call this to indicate that it may be reset
 * at any point between this call and the next call to task_disable_resets().
 *
 * Calling this function will trigger any resets that were requested while
 * resets were disabled.
 *
 * It is not expected for this to be called if resets are already enabled.
 */
void task_enable_resets(void);

/*
 * A task that supports resets may call this to indicate that it may not be
 * reset until the next call to task_enable_resets(). Any calls to task_reset()
 * during this time will cause a reset request to be queued, and executed
 * the next time task_enable_resets() is called.
 *
 * Must not be called if resets are already disabled.
 */
void task_disable_resets(void);

/*
 * If the current task was reset, completes the reset operation.
 *
 * Returns a non-zero value if the task was reset; tasks with state outside
 * of the stack should perform any necessary cleanup immediately after calling
 * this function.
 *
 * Tasks that support reset must call this function once at startup before
 * doing anything else.
 *
 * Must only be called once at task startup.
 */
int task_reset_cleanup(void);

/*
 * Resets the specified task, which must not be the current task,
 * to initial state.
 *
 * Returns EC_SUCCESS, or EC_ERROR_INVAL if the specified task does
 * not support resets.
 *
 * If wait is true, blocks until the task has been reset. Otherwise,
 * returns immediately - in this case the task reset may be delayed until
 * that task can be safely reset. The duration of this delay depends on the
 * task implementation.
 */
int task_reset(task_id_t id, int wait);

/**
 * Clear a pending interrupt.
 *
 * Note that most interrupts can be removed from the pending state simply by
 * handling whatever caused the interrupt in the first place.  This only needs
 * to be called if an interrupt handler disables itself without clearing the
 * reason for the interrupt, and then the interrupt is re-enabled from a
 * different context.
 */
void task_clear_pending_irq(int irq);

struct mutex {
	uint32_t lock;
	uint32_t waiters;
};

/**
 * Lock a mutex.
 *
 * This tries to lock the mutex mtx.  If the mutex is already locked by another
 * task, de-schedules the current task until the mutex is again unlocked.
 *
 * Must not be used in interrupt context!
 */
void mutex_lock(struct mutex *mtx);

/**
 * Release a mutex previously locked by the same task.
 */
void mutex_unlock(struct mutex *mtx);

struct irq_priority {
	uint8_t irq;
	uint8_t priority;
};

/*
 * Some cores may make use of this struct for mapping irqs to handlers
 * for DECLARE_IRQ in a linker-script defined section.
 */
struct irq_def {
	int irq;

	/* The routine which was declared as an IRQ */
	void (*routine)(void);

	/*
	 * The routine usually needs wrapped so the core can handle it
	 * as an IRQ.
	 */
	void (*handler)(void);
};

/*
 * Implement the DECLARE_IRQ(irq, routine, priority) macro which is
 * a core specific helper macro to declare an interrupt handler "routine".
 */
#ifdef CONFIG_COMMON_RUNTIME
#include "irq_handler.h"
#else
#define IRQ_HANDLER(irqname) CONCAT3(irq_, irqname, _handler)
#define IRQ_HANDLER_OPT(irqname) CONCAT3(irq_, irqname, _handler_optional)
#define DECLARE_IRQ(irq, routine, priority) DECLARE_IRQ_(irq, routine, priority)
#define DECLARE_IRQ_(irq, routine, priority) \
	void IRQ_HANDLER_OPT(irq)(void) __attribute__((alias(#routine)));

/* Include ec.irqlist here for compilation dependency */
#define ENABLE_IRQ(x)
#include "ec.irqlist"
#endif

#endif  /* __CROS_EC_TASK_H */
