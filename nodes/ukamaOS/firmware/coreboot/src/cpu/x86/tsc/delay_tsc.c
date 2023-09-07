/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <arch/early_variables.h>
#include <cpu/x86/tsc.h>
#include <pc80/i8254.h>
#include <smp/spinlock.h>
#include <delay.h>
#include <thread.h>

void init_timer(void)
{
	(void)tsc_freq_mhz();
}

void udelay(unsigned int us)
{
	unsigned long long start;
	unsigned long long current;
	unsigned long long clocks;

	if (!thread_yield_microseconds(us))
		return;

	start = rdtscll();
	clocks = us;
	clocks *= tsc_freq_mhz();
	current = rdtscll();
	while ((current - start) < clocks) {
		cpu_relax();
		current = rdtscll();
	}
}

#if CONFIG(TSC_MONOTONIC_TIMER)
#include <timer.h>

static struct monotonic_counter {
	int initialized;
	struct mono_time time;
	uint64_t last_value;
} mono_counter_g CAR_GLOBAL;

static inline struct monotonic_counter *get_monotonic_context(void)
{
	return car_get_var_ptr(&mono_counter_g);
}

void timer_monotonic_get(struct mono_time *mt)
{
	uint64_t current_tick;
	uint64_t ticks_elapsed;
	unsigned long ticks_per_usec;
	struct monotonic_counter *mono_counter;

	mono_counter = get_monotonic_context();
	if (!mono_counter->initialized) {
		init_timer();
		mono_counter->last_value = rdtscll();
		mono_counter->initialized = 1;
	}

	current_tick = rdtscll();
	ticks_elapsed = current_tick - mono_counter->last_value;
	ticks_per_usec = tsc_freq_mhz();

	/* Update current time and tick values only if a full tick occurred. */
	if (ticks_elapsed >= ticks_per_usec) {
		uint64_t usecs_elapsed;

		usecs_elapsed = ticks_elapsed / ticks_per_usec;
		mono_time_add_usecs(&mono_counter->time, (long)usecs_elapsed);
		mono_counter->last_value = current_tick;
	}

	/* Save result. */
	*mt = mono_counter->time;
}
#endif
