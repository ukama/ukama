/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "retap.h"

void retap_request_init(RetapRequest *request, uint8_t procedure)
{
    if (request == NULL) {
        return;
    }

    memset(request, 0, sizeof(RetapRequest));
    request->procedure = procedure;
}

void retap_response_init(RetapResponse *response)
{
    if (response == NULL) {
        return;
    }

    memset(response, 0, sizeof(RetapResponse));
}

bool retap_encode_request(RetapRequest *request,
                          uint8_t *buf,
                          size_t size,
                          size_t *len)
{
    if (request == NULL || buf == NULL || len == NULL) {
        return false;
    }

    if (request->dataLen > RETAP_MAX_PAYLOAD) {
        return false;
    }

    if (size < request->dataLen + RETAP_HEADER_LEN) {
        return false;
    }

    buf[0] = request->procedure;
    buf[1] = (uint8_t)(request->dataLen & 0xFF);
    buf[2] = (uint8_t)((request->dataLen >> 8) & 0xFF);

    if (request->dataLen > 0) {
        memcpy(&buf[RETAP_HEADER_LEN], request->data, request->dataLen);
    }

    *len = request->dataLen + RETAP_HEADER_LEN;

    return true;
}

bool retap_decode_response(const uint8_t *buf,
                           size_t len,
                           RetapResponse *response)
{
    uint16_t dataLen;
    const uint8_t *data;
    size_t appLen;

    if (buf == NULL || response == NULL || len < RETAP_HEADER_LEN) {
        return false;
    }

    dataLen = (uint16_t)buf[1] | ((uint16_t)buf[2] << 8);
    if ((size_t)dataLen + RETAP_HEADER_LEN != len) {
        return false;
    }

    /* Class-1 responses from a single-antenna RET carry at least OK/FAIL. */
    if (dataLen < 1) {
        return false;
    }

    data = &buf[RETAP_HEADER_LEN];

    retap_response_init(response);
    response->procedure = buf[0];
    response->returnCode = data[0];

    if (response->returnCode == RETAP_RETURN_OK) {
        appLen = (size_t)dataLen - 1;
        if (appLen > RETAP_MAX_PAYLOAD) {
            return false;
        }

        if (appLen > 0) {
            memcpy(response->data, &data[1], appLen);
        }
        response->dataLen = appLen;

        return true;
    }

    if (response->returnCode == RETAP_RETURN_FAIL) {
        if (dataLen < 2) {
            return false;
        }

        response->failureReason = data[1];

        appLen = (size_t)dataLen - 2;
        if (appLen > RETAP_MAX_PAYLOAD) {
            return false;
        }

        if (appLen > 0) {
            memcpy(response->data, &data[2], appLen);
        }
        response->dataLen = appLen;

        return true;
    }

    return false;
}

bool retap_response_is_ok(const RetapResponse *response)
{
    return response != NULL && response->returnCode == RETAP_RETURN_OK;
}

bool retap_response_is_fail(const RetapResponse *response)
{
    return response != NULL && response->returnCode == RETAP_RETURN_FAIL;
}

int retap_request_timeout_ms(const RetapRequest *request)
{
    if (request == NULL) {
        return RETAP_DEFAULT_TIMEOUT_MS;
    }

    switch (request->procedure) {
    case RETAP_PROC_SET_TILT:
        return RETAP_SET_TILT_TIMEOUT_MS;
    case RETAP_PROC_CALIBRATE:
        return RETAP_CALIBRATE_TIMEOUT_MS;
    case RETAP_PROC_SEND_CONFIG_DATA:
        return RETAP_CONFIG_TIMEOUT_MS;
    default:
        return RETAP_DEFAULT_TIMEOUT_MS;
    }
}

CtrlCode retap_failure_to_ctrl_code(uint8_t failureReason)
{
    switch (failureReason) {
    case RETAP_RC_BUSY:
    case RETAP_RC_DEVICE_DISABLED:
    case RETAP_RC_DOWNLOAD_IN_PROGRESS:
        return CtrlCodeBusy;

    case RETAP_RC_DATA_ERROR:
    case RETAP_RC_CHECKSUM_ERROR:
    case RETAP_RC_BLOCK_NUMBER_SEQUENCE:
    case RETAP_RC_SEGMENT_NUMBER_SEQUENCE:
    case RETAP_RC_PROCEDURE_SEQUENCE_ERROR:
    case RETAP_RC_UNKNOWN_PARAMETER:
    case RETAP_RC_TOO_MUCH_DATA:
    case RETAP_RC_READ_ONLY:
    case RETAP_RC_INVALID_SOFTWARE:
        return CtrlCodeFormatError;

    case RETAP_RC_NOT_SCALED:
        return CtrlCodeNotConfigured;

    case RETAP_RC_NOT_CALIBRATED:
    case RETAP_RC_POSITION_LOST:
        return CtrlCodeNotCalibrated;

    case RETAP_RC_OUT_OF_RANGE:
        return CtrlCodeOutOfRange;

    case RETAP_RC_WORKING_SOFTWARE_MISSING:
        return CtrlCodeWorkingSoftwareMissing;

    case RETAP_RC_UNKNOWN_PROCEDURE:
    case RETAP_RC_UNKNOWN_ANTENNA_NUMBER:
        return CtrlCodeUnsupportedProcedure;

    case RETAP_RC_ACTUATOR_DETECTION_FAIL:
    case RETAP_RC_ACTUATOR_JAM_PERMANENT:
    case RETAP_RC_ACTUATOR_JAM_TEMPORARY:
    case RETAP_RC_EEPROM_ERROR:
    case RETAP_RC_FLASH_ERASE_ERROR:
    case RETAP_RC_FLASH_ERROR:
    case RETAP_RC_OTHER_HARDWARE_ERROR:
    case RETAP_RC_OTHER_SOFTWARE_ERROR:
    case RETAP_RC_RAM_ERROR:
    case RETAP_RC_UART_ERROR:
    default:
        return CtrlCodeHardwareError;
    }
}


bool retap_decode_request(const uint8_t *buf,
                          size_t len,
                          RetapRequest *request)
{
    uint16_t dataLen;

    if (buf == NULL || request == NULL || len < RETAP_HEADER_LEN) {
        return false;
    }

    dataLen = (uint16_t)buf[1] | ((uint16_t)buf[2] << 8);
    if ((size_t)dataLen + RETAP_HEADER_LEN != len) {
        return false;
    }

    if (dataLen > RETAP_MAX_PAYLOAD) {
        return false;
    }

    retap_request_init(request, buf[0]);
    request->dataLen = dataLen;

    if (dataLen > 0) {
        memcpy(request->data, &buf[RETAP_HEADER_LEN], dataLen);
    }

    return true;
}

bool retap_encode_ok_response(uint8_t procedure,
                              const uint8_t *data,
                              size_t dataLen,
                              uint8_t *buf,
                              size_t size,
                              size_t *len)
{
    size_t appLen;

    if (buf == NULL || len == NULL) {
        return false;
    }

    if (data == NULL && dataLen != 0) {
        return false;
    }

    if (dataLen > RETAP_MAX_PAYLOAD) {
        return false;
    }

    appLen = dataLen + 1; /* OK + optional data. */
    if (size < RETAP_HEADER_LEN + appLen) {
        return false;
    }

    buf[0] = procedure;
    buf[1] = (uint8_t)(appLen & 0xFF);
    buf[2] = (uint8_t)((appLen >> 8) & 0xFF);
    buf[RETAP_HEADER_LEN] = RETAP_RETURN_OK;

    if (dataLen > 0) {
        memcpy(&buf[RETAP_HEADER_LEN + 1], data, dataLen);
    }

    *len = RETAP_HEADER_LEN + appLen;

    return true;
}

bool retap_encode_fail_response(uint8_t procedure,
                                uint8_t failureReason,
                                const uint8_t *extra,
                                size_t extraLen,
                                uint8_t *buf,
                                size_t size,
                                size_t *len)
{
    size_t appLen;

    if (buf == NULL || len == NULL) {
        return false;
    }

    if (extra == NULL && extraLen != 0) {
        return false;
    }

    if (extraLen > RETAP_MAX_PAYLOAD - 2) {
        return false;
    }

    appLen = extraLen + 2; /* FAIL + reason + optional extra. */
    if (size < RETAP_HEADER_LEN + appLen) {
        return false;
    }

    buf[0] = procedure;
    buf[1] = (uint8_t)(appLen & 0xFF);
    buf[2] = (uint8_t)((appLen >> 8) & 0xFF);
    buf[RETAP_HEADER_LEN] = RETAP_RETURN_FAIL;
    buf[RETAP_HEADER_LEN + 1] = failureReason;

    if (extraLen > 0) {
        memcpy(&buf[RETAP_HEADER_LEN + 2], extra, extraLen);
    }

    *len = RETAP_HEADER_LEN + appLen;

    return true;
}

const char *retap_return_code_str(uint8_t code)
{
    switch (code) {
    case RETAP_RC_OK:                       return "OK";
    case RETAP_RC_ACTUATOR_DETECTION_FAIL:  return "ActuatorDetectionFail";
    case RETAP_RC_ACTUATOR_JAM_PERMANENT:   return "ActuatorJamPermanent";
    case RETAP_RC_ACTUATOR_JAM_TEMPORARY:   return "ActuatorJamTemporary";
    case RETAP_RC_BLOCK_NUMBER_SEQUENCE:    return "BlockNumberSequenceError";
    case RETAP_RC_BUSY:                     return "Busy";
    case RETAP_RC_CHECKSUM_ERROR:           return "ChecksumError";
    case RETAP_RC_PROCEDURE_SEQUENCE_ERROR: return "ProcedureSequenceError";
    case RETAP_RC_DATA_ERROR:               return "DataError";
    case RETAP_RC_DEVICE_DISABLED:          return "DeviceDisabled";
    case RETAP_RC_EEPROM_ERROR:             return "EEPROMError";
    case RETAP_RC_FAIL:                     return "FAIL";
    case RETAP_RC_FLASH_ERASE_ERROR:        return "FlashEraseError";
    case RETAP_RC_FLASH_ERROR:              return "FlashError";
    case RETAP_RC_NOT_CALIBRATED:           return "NotCalibrated";
    case RETAP_RC_NOT_SCALED:               return "NotScaled";
    case RETAP_RC_OTHER_HARDWARE_ERROR:     return "OtherHardwareError";
    case RETAP_RC_OTHER_SOFTWARE_ERROR:     return "OtherSoftwareError";
    case RETAP_RC_OUT_OF_RANGE:             return "OutOfRange";
    case RETAP_RC_POSITION_LOST:            return "PositionLost";
    case RETAP_RC_RAM_ERROR:                return "RAMError";
    case RETAP_RC_SEGMENT_NUMBER_SEQUENCE:  return "SegmentNumberSequenceError";
    case RETAP_RC_UART_ERROR:               return "UARTError";
    case RETAP_RC_UNKNOWN_PROCEDURE:        return "UnknownProcedure";
    case RETAP_RC_READ_ONLY:                return "ReadOnly";
    case RETAP_RC_UNKNOWN_PARAMETER:        return "UnknownParameter";
    case RETAP_RC_UNKNOWN_ANTENNA_NUMBER:   return "UnknownAntennaNumber";
    case RETAP_RC_TOO_MUCH_DATA:            return "TooMuchData";
    case RETAP_RC_WORKING_SOFTWARE_MISSING: return "WorkingSoftwareMissing";
    case RETAP_RC_INVALID_SOFTWARE:         return "InvalidSoftware";
    case RETAP_RC_DOWNLOAD_IN_PROGRESS:     return "DownloadInProgress";
    default:                                return "UnknownReturnCode";
    }
}
