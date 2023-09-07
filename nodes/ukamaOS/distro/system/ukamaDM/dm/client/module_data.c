

#include "liblwm2m.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>

#include "module_data.h"

extern lwm2m_object_t * objArray[13];

int get_temperature_sensor_value(uint8_t instance, double *sensor_value)
{
    //todo: Add BSP call here to get the temperature sensor value from module.

    /* todo: based on the module info instances, it needs to be determined which module and which sensor to read.
        e.g if instance passed is 5, then we will read module info instances and determine how many temperature sensors does it have
         say,
         first instace has 2,
         next has 1,
         and next 3
        then we will read the temperature sensor of the 3rd module. since 3rd module has 3 sensors and the instace passed is 5, we will
        reads value of 2nd temperature sensor.
     */

    if (sensor_value)
        *sensor_value = 17.5;
    else
        return EXIT_FAILURE;

    return EXIT_SUCCESS;
}
