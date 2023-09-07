/* Copyright 2018 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Common test code to test lid angle calculation.
 */

#include "accelgyro.h"
#include "host_command.h"
#include "motion_common.h"
#include "motion_sense.h"
#include "task.h"
#include "timer.h"

/*****************************************************************************/
/* Mock functions */
static int accel_init(const struct motion_sensor_t *s)
{
	return EC_SUCCESS;
}

static int accel_read(const struct motion_sensor_t *s, intv3_t v)
{
	rotate(s->xyz, *s->rot_standard_ref, v);
	return EC_SUCCESS;
}

static int accel_get_range(const struct motion_sensor_t *s)
{
	return s->default_range;
}

static int accel_get_resolution(const struct motion_sensor_t *s)
{
	return 0;
}

int test_data_rate[2] = { 0 };

static int accel_set_data_rate(const struct motion_sensor_t *s,
			      const int rate,
			      const int rnd)
{
	test_data_rate[s - motion_sensors] = rate | (rnd ? ROUND_UP_FLAG : 0);
	return EC_SUCCESS;
}

static int accel_get_data_rate(const struct motion_sensor_t *s)
{
	return test_data_rate[s - motion_sensors];
}

const struct accelgyro_drv test_motion_sense = {
	.init = accel_init,
	.read = accel_read,
	.get_range = accel_get_range,
	.get_resolution = accel_get_resolution,
	.set_data_rate = accel_set_data_rate,
	.get_data_rate = accel_get_data_rate,
};

struct motion_sensor_t motion_sensors[] = {
	[BASE] = {
		.name = "base",
		.active_mask = SENSOR_ACTIVE_S0_S3_S5,
		.chip = MOTIONSENSE_CHIP_LSM6DS0,
		.type = MOTIONSENSE_TYPE_ACCEL,
		.location = MOTIONSENSE_LOC_BASE,
		.drv = &test_motion_sense,
		.rot_standard_ref = NULL,
		.default_range = 2,  /* g, enough for laptop. */
		.config = {
			/* EC use accel for angle detection */
			[SENSOR_CONFIG_EC_S0] = {
				.odr = TEST_LID_FREQUENCY,
			},
		},
	},
	[LID] = {
		.name = "lid",
		.active_mask = SENSOR_ACTIVE_S0,
		.chip = MOTIONSENSE_CHIP_KXCJ9,
		.type = MOTIONSENSE_TYPE_ACCEL,
		.location = MOTIONSENSE_LOC_LID,
		.drv = &test_motion_sense,
		.rot_standard_ref = NULL,
		.default_range = 2,  /* g, enough for laptop. */
		.config = {
			/* EC use accel for angle detection */
			[SENSOR_CONFIG_EC_S0] = {
				.odr = TEST_LID_FREQUENCY,
			},
		},
	},
};
const unsigned int motion_sensor_count = ARRAY_SIZE(motion_sensors);

/* Read 6 samples from array to sensor vectors, convert units if necessary. */
void feed_accel_data(const float *array, int *idx,
		int (filler)(const struct motion_sensor_t*, const float))
{
	int i, j;

	for (i = 0; i < motion_sensor_count; i++) {
		struct motion_sensor_t *s = &motion_sensors[i];

		for (j = X; j <= Z; j++)
			s->xyz[j] = filler(s, array[*idx + i * 3 + j]);
	}
	*idx += 6;
}

void wait_for_valid_sample(void)
{
	uint8_t sample;
	uint8_t *lpc_status = host_get_memmap(EC_MEMMAP_ACC_STATUS);

	sample = *lpc_status & EC_MEMMAP_ACC_STATUS_SAMPLE_ID_MASK;
	usleep(TEST_LID_EC_RATE);
	task_wake(TASK_ID_MOTIONSENSE);
	while ((*lpc_status & EC_MEMMAP_ACC_STATUS_SAMPLE_ID_MASK) == sample)
		usleep(TEST_LID_SLEEP_RATE);
}


