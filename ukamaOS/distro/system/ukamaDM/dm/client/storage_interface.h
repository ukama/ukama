#ifndef STORAGE_INTERFACE_H_
#define STORAGE_INTERFACE_H_

#include "stdio.h"
#include "stdlib.h"
#include "stdint.h"
#include "string.h"

/* Error codes */
#define STORAGE_INTERFACE_SUCCESS                       0
#define STORAGE_INTERFACE_INVALID_PARAM                 -1001
#define STORAGE_INTERFACE_UNKNOWN_MODULE_CLASS          -1002

/* String literals. */
#define STORAGE_INTERFACE_MODULE_CLASS_MASK             "Mask"
#define STORAGE_INTERFACE_MODULE_CLASS_RADIO            "Radio"
#define STORAGE_INTERFACE_MODULE_CLASS_COM              "Compute"
#define STORAGE_INTERFACE_MODULE_CLASS_POWER            "Power"
#define STORAGE_INTERFACE_MODULE_CLASS_CTRL             "Controller"
#define STORAGE_INTERFACE_MODULE_CLASS_FEM              "Frontend"

/* Macros for constants. */
#define STORAGE_INTERFACE_TEMP_SENSORS_MASK             2
#define STORAGE_INTERFACE_TEMP_SENSORS_RADIO            2
#define STORAGE_INTERFACE_TEMP_SENSORS_COM              2
#define STORAGE_INTERFACE_TEMP_SENSORS_POWER            2
#define STORAGE_INTERFACE_TEMP_SENSORS_CTRL             1
#define STORAGE_INTERFACE_TEMP_SENSORS_FEM              2

/* Macros for constants. */
#define STORAGE_INTERFACE_DIGITAL_INPUT_MASK             2
#define STORAGE_INTERFACE_DIGITAL_INPUT_RADIO            2
#define STORAGE_INTERFACE_DIGITAL_INPUT_COM              2
#define STORAGE_INTERFACE_DIGITAL_INPUT_POWER            2
#define STORAGE_INTERFACE_DIGITAL_INPUT_CTRL             1
#define STORAGE_INTERFACE_DIGITAL_INPUT_FEM              2

/* Macros for constants. */
#define STORAGE_INTERFACE_DIGITAL_OUTPUT_MASK             2
#define STORAGE_INTERFACE_DIGITAL_OUTPUT_RADIO            2
#define STORAGE_INTERFACE_DIGITAL_OUTPUT_COM              2
#define STORAGE_INTERFACE_DIGITAL_OUTPUT_POWER            2
#define STORAGE_INTERFACE_DIGITAL_OUTPUT_CTRL             1
#define STORAGE_INTERFACE_DIGITAL_OUTPUT_FEM              2

/* Reads the number of modules attached to the device. */
int Read_Module_Count(uint8_t *module_count);

/* Read the Module uuid of each module. */
int Read_Each_Module_UUID(char **uuid, uint8_t module_number);

/* Read the Module manufacturer of each module. */
int Read_Each_Module_Manufacturer(char **manufacturer, uint8_t module_number);

/* Read the Module model of each module. */
int Read_Each_Module_Model(char **model, uint8_t module_number);

/* Read the Module part number of each module. */
int Read_Each_Module_PartNumber(char **partnum, uint8_t module_number);

/* Read the Module manufacturing date of each module. */
int Read_Each_Module_Mfgdate(char **mfgdate, uint8_t module_number);

/* Read the Module module class of each module. */
int Read_Each_Module_Moduleclass(char **moduleclass, uint8_t module_number);

/* Read the Module software version of each module. */
int Read_Each_Module_Swversion(char **swversion, uint8_t module_number);

/* Read the Module hardware version of each module. */
int Read_Each_Module_Hwversion(char **hwversion, uint8_t module_number);

/* Read the number of temperature sensors in the module. */
int Read_Temperature_Sensor_Count(char *model_class, uint8_t *count);

/* Read the number of digital inputs in the module. */
int Read_Digital_Input_Count(char *model_class, uint8_t *count);

/* Read the number of digital outputs in the module. */
int Read_Digital_Output_Count(char *model_class, uint8_t *count);

#endif