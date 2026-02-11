/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#ifndef JSON_TYPES_H_
#define JSON_TYPES_H_

#define JTAG_NODE_INFO          "nodeInfo"
#define JTAG_NODE_ID            "UUID"

#define JTAG_TIME               "time"
#define JTAG_MODULE             "module"
#define JTAG_SEVERITY           "severity"
#define JTAG_DETAILS            "details"

#define JTAG_ERROR              "error"
#define JTAG_ERROR_CODE         "code"
#define JTAG_ERROR_CSTRING      "string"

#define JTAG_NAME               "name"
#define JTAG_TYPE               "type"
#define JTAG_VALUE              "value"
#define JTAG_STATUS             "status"
#define JTAG_MESSAGE            "message"
#define JTAG_ENABLED            "enabled"
#define JTAG_UNITS              "units"
#define JTAG_TIMESTAMP          "timestamp"

#define JTAG_FEM_UNIT           "fem_unit"

#define JTAG_GPIO_STATUS        "gpio_status"
#define JTAG_TX_RF_ENABLE       "tx_rf_enable"
#define JTAG_RX_RF_ENABLE       "rx_rf_enable"
#define JTAG_PA_VDS_ENABLE      "pa_vds_enable"
#define JTAG_RF_PAL_ENABLE      "rf_pal_enable"
#define JTAG_VDS_28V_ENABLE     "vds_28v_enable"
#define JTAG_PSU_PGOOD          "psu_pgood"

#define JTAG_DAC_CONFIG         "dac_config"
#define JTAG_CARRIER_VOLTAGE    "carrier_voltage"
#define JTAG_PEAK_VOLTAGE       "peak_voltage"

#define JTAG_ADC_READING        "adc_reading"
#define JTAG_REVERSE_POWER      "reverse_power"
#define JTAG_FORWARD_POWER      "forward_power"
#define JTAG_PA_CURRENT         "pa_current"

#define JTAG_EEPROM_DATA        "eeprom_data"
#define JTAG_SERIAL             "serial"

#define JTAG_SERVICE_INFO       "service_info"
#define JTAG_SERVICE_NAME       "service_name"
#define JTAG_VERSION            "version"
#define JTAG_UPTIME             "uptime"

#define JTAG_OP                 "op"
#define JTAG_OP_ID              "op_id"
#define JTAG_OP_STATE           "state"
#define JTAG_RC                 "rc"
#define JTAG_UPDATED            "updated_ts_ms"

#endif /* JSON_TYPES_H_ */
