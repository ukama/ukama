/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_PROPERTY_H_
#define INC_PROPERTY_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"
#include "noded_macros.h"

/* Permissions*/
#define PERM_EX			0x01
#define PERM_RD 		0x02
#define PERM_WR			0X40
#define PERM_RW         ((PERM_RD)|(PERM_WR))
#define PERM_RWE        ((PERM_RD)|(PERM_WR)|(PERM_EX))

/* Available in device variant or not */
#define PROP_NOTAVAIL	0x00
#define PROP_AVAIL 		0x01

/* Property Type */
#define PROP_TYPE_CONFIG	0x0001
#define PROP_TYPE_STATUS    0x0002
#define PROP_TYPE_ALERT		0x0004
#define PROP_TYPE_EXEC		0x0008
#define PROP_TYPE_MONIT		0x0010 /* Used for LWM2M object definitions like average, min and max */

typedef enum {
    TYPE_NULL = 0,
    TYPE_CHAR = 1,
    TYPE_BOOL,
    TYPE_UINT8,
    TYPE_UINT16,
    TYPE_UINT32,
    TYPE_INT8,
    TYPE_INT16,
    TYPE_INT32,
    TYPE_INT,
    TYPE_FLOAT,
    TYPE_DOUBLE,
    TYPE_ENUM,
    TYPE_STRING,
    TYPE_MAX
} DataType;

typedef enum {
    T1TEMPVALUE = 0,	/* Local Temp Sensor */
    T1MINLIMIT,
    T1MAXLIMIT,
    T1CRITLIMIT, /* Thermal Limit */
    T1MINALARM,
    T1MAXALARM,
    T1CRITALARM,
    T1CRITHYST,
    T1MAXHYST,
    T1OFFSET,
    T2TEMPVALUE,
    T2MINLIMIT,
    T2MAXLIMIT,
    T2CRITLIMIT,
    T2MINALARM,
    T2MAXALARM,
    T2CRITALARM,
    T2CRITHYST,
    T2MAXHYST,
    T2OFFSET,
    T3TEMPVALUE,
    T3MINLIMIT,
    T3MAXLIMIT,
    T3CRITLIMIT,
    T3MINALARM,
    T3MAXALARM,
    T3CRITALARM,
    T3CRITHYST,
    T3MAXHYST,
    T3OFFSET,
    MAXTEMPPROP /* Make sure this is last property.*/
} TempProperty;

typedef enum {
    SHUNTVOLTAGE = 0,
    BUSVOLTAGE,
    CURRENT,
    POWER,
    SHUNTRESISTOR,
    CALIBRATION,
    CRITLOWSHUNTVOLTAGE,
    CRITHIGHSHUNTVOLTAGE,
    SHUNTVOLTAGECRITLOWALARM,
    SHUNTVOLTAGECRITHIGHALARM,
    CRITLOWBUSVOLTAGE,
    CRITHIGHBUSVOLTAGE,
    BUSVOLTAGECRITLOWALARM,
    BUSVOLTAGECRITHIGHALARM,
    CRITHIGHPWR,
    CRITHIGHPWRALARM,
    UPDATEINTERVAL,
    MAXINAPROP /* Make sure this is last property.*/
} INAProperty;

typedef enum {
    ATTVALUE = 1,
    LATCHENABLE,
    MAXATTPROP /* Make sure this is last property.*/
} ATTProperty;

typedef enum {
    VAIN0AIN1 = 0,
    VAIN0AIN3,
    VAIN1AIN3,
    VAIN2AIN3,
    VAIN0GND,
    VAIN1GND,
    VAIN2GND,
    VAIN3GND,
    MAXADCPROP /* Make sure this is last property.*/
} ADCProperty;

typedef enum {
    DIRECTION = 0,
    VALUE,
    EDGE,
    POLARITY,
    MAXGPIOPROP
}GPIOProperty;

typedef enum {
    RBRIGHTNESS = 0,
    RMAX_BRIGHTNESS,
    RTRIGGER,
    GBRIGHTNESS,
    GMAX_BRIGHTNESS,
    GTRIGGER,
    BBRIGHTNESS,
    BMAX_BRIGHTNESS,
    BTRIGGER,
    MAXLEDTRICOLPROP
}LEDProperty;

typedef enum {
    STRICTLYLESSTHEN = 0,
    LESSTHENEQUALTO,
    STRICTLYGREATERTHEN,
    GREATERTHENEQUALTO
}AlertCondition;

/* If any property is depending on other properties it needs to hold this structure.*/
typedef struct  __attribute__((__packed__)) {
    int currIdx;
    int lmtIdx;
    AlertCondition cond;
} DepProperty;

/* For each device we have set of properties which can read, configured and queried.*/
/* This could be use to store user config also*/
typedef struct  __attribute__((__packed__)) {
    uint16_t id;
    char name[32];
    DataType dataType;
    uint8_t perm;
    uint8_t available;  //for ADT we have three configs Low Level , high level and critical level for SE98 we have only low and high.
    uint16_t propType;
    char units[24];
    char sysFname[64];
    DepProperty *depProp; //For Alerts: this may hold location of current and default values of the alert property.*/
} Property;

typedef struct  __attribute__((__packed__)) {
    DataType type;
    uint16_t size;
} MapDataType;

/**
 * @fn      int get_alert_cond(char*)
 * @brief   Converts string to AlertCondition
 *
 * @param   cond
 * @return  On success, enum AlertCondition i.e Integer value 0 or greater
 *          On failure, -1
 */
int get_alert_cond(char* cond);

/**
 * @fn      int get_prop_perm(char*)
 * @brief   Converts string to Permission values
 *
 * @param   perm
 * @return  On success, positive integer, one of permission macros values
 *          On failure, -1
 */
int get_prop_perm(char* perm);

/**
 * @fn      int get_prop_type(char*)
 * @brief   Converts string to property type values
 *
 * @param   type
 * @return  On success, positive integer, one of property type values
 *          On failure, -1
 */
int get_prop_type(char* type);

/**
 * @fn      int get_prop_avail(char*)
 * @brief   Converts string to property available values
 *
 * @param   avail
 * @return  On success, positive integer, one of property available values
 *          On failure, -1
 */
int get_prop_avail(char* avail);

/**
 * @fn      int get_prop_data_type(char*)
 * @brief   Converts string to DataType
 *
 * @param   type
 * @return  On success, enum DataType i.e Integer value 0 or greater
 *          On failure, -1
 */
int get_prop_data_type(char *type);

/**
 * @fn      int get_property_count(char*)
 * @brief   Read the property count for sensor from parsed data.
 *
 * @param   dev
 * @return  On success, 0
 *          On failure, non zero value
 */
int get_property_count(char* dev);

/**
 * @fn      int validate_irq_limits(double, double, int)
 * @brief   Compares the current sensor value cur to limits
 *          using condition cond
 *
 * @param   cur
 * @param   lmt
 * @param   cond
 * @return  On success, 1
 *          On Invalid condition, -1
 *          On failure, 0
 */
int validate_irq_limits(double cur, double lmt, int cond);

/**
 * @fn      void print_properties(Property*, uint16_t)
 * @brief   list the sensor properties
 *
 * @param   prop
 * @param   count
 */
void print_properties(Property* prop, uint16_t count);

/**
 * @fn      char get_sysfs_name*(char*)
 * @brief   read the filename from the path.
 *
 * @param   fpath
 * @return  On success, filename
 *          On failure, NULL
 */
char* get_sysfs_name(char* fpath);

/**
 * @fn      uint16_t get_sizeof(DataType)
 * @brief   return the size of the data type.
 *
 * @param   type
 * @return  On success, size of the data type
 *          On failure, zero value
 */
uint16_t get_sizeof(DataType type);

/**
 * @fn      Property get_property_table*(char*)
 * @brief   Reads the property table for sensor dev from the parsed data.
 *
 * @param   dev
 * @return  On success, property table
 *          On failure, NULL
 */
Property* get_property_table(char* dev);

#ifdef __cplusplus
}
#endif

#endif /* INCLUDE_PROPERTY_H_ */
