/*******************************************************************************
 *
 * Copyright (c) 2015 Intel Corporation and others.
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v2.0
 * and Eclipse Distribution License v1.0 which accompany this distribution.
 *
 * The Eclipse Public License is available at
 *    http://www.eclipse.org/legal/epl-v20.html
 * The Eclipse Distribution License is available at
 *    http://www.eclipse.org/org/documents/edl-v10.php.
 *
 * Contributors:
 *    David Navarro, Intel Corporation - initial implementation
 *    Scott Bertin, AMETEK, Inc. - Please refer to git log
 *
 *******************************************************************************/

#include <stdlib.h>
#include <string.h>
#include <ctype.h> // isspace

#include "bootstrap_info.h"

typedef struct
{
    uint16_t    id;
    char *      uri;
    bool        isBootstrap;
    uint32_t    lifetime;
    uint8_t     securityMode;
    uint8_t *   publicKey;
    size_t      publicKeyLen;
    uint8_t *   secretKey;
    size_t      secretKeyLen;
    uint8_t *   serverKey;
    size_t      serverKeyLen;
#ifndef LWM2M_VERSION_1_0
    int         registrationPriorityOrder; // <0 when it doesn't exist
    int         initialRegistrationDelayTimer; // <0 when it doesn't exist
    int8_t      registrationFailureBlock; // <0 when it doesn't exist, 0 for false, > 0 for true
    int8_t      bootstrapOnRegistrationFailure; // <0 when it doesn't exist, 0 for false, > 0 for true
    int         communicationRetryCount; // <0 when it doesn't exist
    int         communicationRetryTimer; // <0 when it doesn't exist
    int         communicationSequenceDelayTimer; // <0 when it doesn't exist
    int         communicationSequenceRetryCount; // <0 when it doesn't exist
#endif
} read_server_t;

static int prv_find_next_section(FILE * fd,
                                 char * tag)
{
    char * line;
    size_t length;
    int found;

    line = NULL;
    length = 0;
    found = 0;
    while (found == 0
        && getline(&line, &length, fd) != -1)
    {
        if (line[0] == '[')
        {
            int i;

            length = strlen(line);
            i = 1;
            while (line[i] != ']') i++;

            if (i < length)
            {
                line[i] = 0;
                if (strcasecmp(line + 1, tag) == 0)
                {
                    found = 1;
                }
            }
        }
        lwm2m_free(line);
        line = NULL;
        length = 0;
    }

    return found;
}

// returns -1 for error, 0 if not found (end of section or file)
// and 1 if found
static int prv_read_key_value(FILE * fd,
                              char ** keyP,
                              char ** valueP)
{
    char * line;
    fpos_t prevPos;
    ssize_t res;
    size_t length;
    size_t start;
    size_t middle;
    size_t end;

    *keyP = NULL;
    *valueP = NULL;

    line = NULL;
    if (fgetpos(fd, &prevPos) != 0) return -1;
    while ((res = getline(&line, &length, fd)) != -1)
    {
        length = strlen(line);

        start = 0;
        while (start < length && isspace(line[start]&0xff)) start++;
        // ignore empty and commented lines
        if (start != length
         && line[start] != ';'
         && line[start] != '#')
        {
            break;
        }

        lwm2m_free(line);
        line = NULL;
        length = 0;
        if (fgetpos(fd, &prevPos) != 0) return -1;
    }

    if (res == -1) return -1;

    // end of section
    if (line[start] == '[')
    {
        lwm2m_free(line);
        fsetpos(fd, &prevPos);
        return 0;
    }

    middle = start;
    while (middle < length && line[middle] != '=') middle++;
    // invalid lines
    if (middle == start
     || middle == length)
    {
        lwm2m_free(line);
        return -1;
    }

    end = length - 2;
    while (end > middle && isspace(line[end]&0xff)) end--;
    // invalid lines
    if (end == middle)
    {
        lwm2m_free(line);
        return -1;
    }
    end += 1;

    line[middle] = 0;
    *keyP = strdup(line + start);
    if (*keyP == NULL)
    {
        lwm2m_free(line);
        return -1;
    }

    middle +=1 ;
    while (middle < end && isspace(line[middle]&0xff)) middle++;
    line[end] = 0;
    *valueP = strdup(line + middle);
    if (*valueP == NULL)
    {
        lwm2m_free(*keyP);
        *keyP = NULL;
        lwm2m_free(line);
        return -1;
    }

    lwm2m_free(line);

    return 1;
}

static int prv_readDigit(char digit)
{
    if (digit >= '0' && digit <= '9')
    {
        return (digit - '0');
    }
    if (digit >= 'a' && digit <= 'f')
    {
        return (10 + digit - 'a');
    }
    if (digit >= 'A' && digit <= 'F')
    {
        return (10 + digit - 'A');
    }

    return -1;
}

static size_t prv_readSecurityKey(char * value,
                                  uint8_t ** resultP)
{
    size_t length;
    size_t charIndex;
    size_t resIndex;

    if (strlen(value)%2 != 0) return 0;

    length = strlen(value) / 2;
    *resultP = (uint8_t *)lwm2m_malloc(length);
    if (*resultP == NULL) return 0;

    resIndex = 0;
    charIndex = 0;
    while (resIndex < length)
    {
        int tmp;

        tmp = prv_readDigit(value[charIndex]);
        if (tmp == -1) goto error;
        (*resultP)[resIndex] = (tmp & 0xFF);

        charIndex++;
        tmp = prv_readDigit(value[charIndex]);
        if (tmp == -1) goto error;
        (*resultP)[resIndex] = ((*resultP)[resIndex] << 4) + (tmp & 0xFF);

        charIndex++;
        resIndex++;
    }

    return length;

error:
    lwm2m_free(*resultP);
    *resultP = NULL;

    return 0;
}

static read_server_t * prv_read_next_server(FILE * fd)
{
    char * key;
    char * value;
    read_server_t * readSrvP;
    int res;

    if (prv_find_next_section(fd, "Server") == 0) return NULL;

    readSrvP = (read_server_t *)lwm2m_malloc(sizeof(read_server_t));
    if (readSrvP == NULL) return NULL;
    memset(readSrvP, 0, sizeof(read_server_t));
#ifndef LWM2M_VERSION_1_0
    readSrvP->registrationPriorityOrder = -1;
    readSrvP->initialRegistrationDelayTimer = -1;
    readSrvP->registrationFailureBlock = -1;
    readSrvP->bootstrapOnRegistrationFailure = -1;
    readSrvP->communicationRetryCount = -1;
    readSrvP->communicationRetryTimer = -1;
    readSrvP->communicationSequenceDelayTimer = -1;
    readSrvP->communicationSequenceRetryCount = -1;
#endif

    while((res = prv_read_key_value(fd, &key, &value)) == 1)
    {
        if (strcasecmp(key, "id") == 0)
        {
            int num;

            if (sscanf(value, "%d", &num) != 1) goto error;
            if (num <= 0 || num > LWM2M_MAX_ID) goto error;
            readSrvP->id = num;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "uri") == 0)
        {
            readSrvP->uri = value;
        }
        else if (strcasecmp(key, "bootstrap") == 0)
        {
            if (strcasecmp(value, "yes") == 0)
            {
                readSrvP->isBootstrap = true;
            }
            else if (strcasecmp(value, "no") == 0)
            {
                readSrvP->isBootstrap = false;
            }
            else goto error;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "lifetime") == 0)
        {
            int num;

            if (sscanf(value, "%d", &num) != 1) goto error;
            if (num <= 0) goto error;
            readSrvP->lifetime = num;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "security") == 0)
        {
            if (strcasecmp(value, "nosec") == 0)
            {
                readSrvP->securityMode = LWM2M_SECURITY_MODE_NONE;
            }
            else if (strcasecmp(value, "PSK") == 0)
            {
                readSrvP->securityMode = LWM2M_SECURITY_MODE_PRE_SHARED_KEY;
            }
            else if (strcasecmp(value, "RPK") == 0)
            {
                readSrvP->securityMode = LWM2M_SECURITY_MODE_RAW_PUBLIC_KEY;
            }
            else if (strcasecmp(value, "certificate") == 0)
            {
                readSrvP->securityMode = LWM2M_SECURITY_MODE_CERTIFICATE;
            }
            else goto error;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "public") == 0)
        {
            readSrvP->publicKeyLen = prv_readSecurityKey(value, &(readSrvP->publicKey));
            if (readSrvP->publicKeyLen == 0) goto error;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "server") == 0)
        {
            readSrvP->serverKeyLen = prv_readSecurityKey(value, &(readSrvP->serverKey));
            if (readSrvP->serverKeyLen == 0) goto error;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "secret") == 0)
        {
            readSrvP->secretKeyLen = prv_readSecurityKey(value, &(readSrvP->secretKey));
            if (readSrvP->secretKeyLen == 0) goto error;
            lwm2m_free(value);
        }
#ifndef LWM2M_VERSION_1_0
        else if (strcasecmp(key, "registrationPriorityOrder") == 0)
        {
            int num;

            if (sscanf(value, "%d", &num) != 1) goto error;
            if (num < 0) goto error;
            readSrvP->registrationPriorityOrder = num;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "initialRegistrationDelay") == 0)
        {
            int num;

            if (sscanf(value, "%d", &num) != 1) goto error;
            if (num < 0) goto error;
            readSrvP->initialRegistrationDelayTimer = num;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "registrationFailureBlock") == 0)
        {
            if (strcasecmp(value, "yes") == 0)
            {
                readSrvP->registrationFailureBlock = true;
            }
            else if (strcasecmp(value, "no") == 0)
            {
                readSrvP->registrationFailureBlock = false;
            }
            else goto error;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "bootstrapOnRegistrationFailure") == 0)
        {
            if (strcasecmp(value, "yes") == 0)
            {
                readSrvP->bootstrapOnRegistrationFailure = true;
            }
            else if (strcasecmp(value, "no") == 0)
            {
                readSrvP->bootstrapOnRegistrationFailure = false;
            }
            else goto error;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "communicationRetryCount") == 0)
        {
            int num;

            if (sscanf(value, "%d", &num) != 1) goto error;
            if (num < 0) goto error;
            readSrvP->communicationRetryCount = num;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "communicationRetryTimer") == 0)
        {
            int num;

            if (sscanf(value, "%d", &num) != 1) goto error;
            if (num < 0) goto error;
            readSrvP->communicationRetryTimer = num;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "communicationSequenceDelayTimer") == 0)
        {
            int num;

            if (sscanf(value, "%d", &num) != 1) goto error;
            if (num < 0) goto error;
            readSrvP->communicationSequenceDelayTimer = num;
            lwm2m_free(value);
        }
        else if (strcasecmp(key, "communicationSequenceRetryCount") == 0)
        {
            int num;

            if (sscanf(value, "%d", &num) != 1) goto error;
            if (num < 0) goto error;
            readSrvP->communicationSequenceRetryCount = num;
            lwm2m_free(value);
        }
#endif
        else
        {
            // ignore key for now
            lwm2m_free(value);
        }
        lwm2m_free(key);
    }

    if (res == -1) goto error;
    if (readSrvP->id == 0
     || readSrvP->uri == 0
     || (readSrvP->securityMode != LWM2M_SECURITY_MODE_NONE
          && (readSrvP->publicKey == NULL || readSrvP->secretKey == NULL))
     || (readSrvP->serverKey == NULL
          && (readSrvP->securityMode == LWM2M_SECURITY_MODE_RAW_PUBLIC_KEY
           || readSrvP->securityMode == LWM2M_SECURITY_MODE_CERTIFICATE)))
    {
        goto error;
    }

    return readSrvP;

error:
    if (readSrvP != NULL)
    {
        if (readSrvP->uri != NULL) lwm2m_free(readSrvP->uri);
        if (readSrvP->publicKey != NULL) lwm2m_free(readSrvP->publicKey);
        if (readSrvP->secretKey != NULL) lwm2m_free(readSrvP->secretKey);
        if (readSrvP->serverKey != NULL) lwm2m_free(readSrvP->serverKey);
        lwm2m_free(readSrvP);
    }
    if (key != NULL) lwm2m_free(key);
    if (value != NULL) lwm2m_free(value);

    return NULL;
}

static int prv_add_server(bs_info_t * infoP,
                          read_server_t * dataP)
{
    lwm2m_data_t * tlvP;
    int size;
    bs_server_tlv_t * serverP;
    lwm2m_media_type_t format;
    int res;

    switch (dataP->securityMode)
    {
    case LWM2M_SECURITY_MODE_NONE:
        size = 4;
        break;
    case LWM2M_SECURITY_MODE_PRE_SHARED_KEY:
        size = 6;
        break;
    case LWM2M_SECURITY_MODE_RAW_PUBLIC_KEY:
    case LWM2M_SECURITY_MODE_CERTIFICATE:
        size = 7;
        break;
    default:
        return -1;
    }

    serverP = (bs_server_tlv_t *)lwm2m_malloc(sizeof(bs_server_tlv_t));
    if (serverP == NULL) return -1;
    memset(serverP, 0, sizeof(bs_server_tlv_t));

    serverP->id = dataP->id;

    tlvP = lwm2m_data_new(size);
    if (tlvP == NULL) goto error;

    // LWM2M Server URI
    tlvP[0].id = LWM2M_SECURITY_URI_ID;
    lwm2m_data_encode_string(dataP->uri, tlvP);

    // Bootstrap Server
    tlvP[1].id = LWM2M_SECURITY_BOOTSTRAP_ID;
    lwm2m_data_encode_bool(dataP->isBootstrap, tlvP + 1);

    // Short Server ID
    tlvP[2].id = LWM2M_SECURITY_SHORT_SERVER_ID;
    lwm2m_data_encode_int(dataP->id, tlvP + 2);

    // Security Mode
    tlvP[3].id = LWM2M_SECURITY_SECURITY_ID;
    lwm2m_data_encode_int(dataP->securityMode, tlvP + 3);

    if (size > 4)
    {
        tlvP[4].id = LWM2M_SECURITY_PUBLIC_KEY_ID;
        lwm2m_data_encode_opaque(dataP->publicKey, dataP->publicKeyLen, tlvP + 4);

        tlvP[5].id = LWM2M_SECURITY_SECRET_KEY_ID;
        lwm2m_data_encode_opaque(dataP->secretKey, dataP->secretKeyLen, tlvP + 5);

        if (size == 7)
        {
            tlvP[6].id = LWM2M_SECURITY_SERVER_PUBLIC_KEY_ID;
            lwm2m_data_encode_opaque(dataP->serverKey, dataP->serverKeyLen, tlvP + 5);
        }
    }

    format = LWM2M_CONTENT_TLV;
    res = lwm2m_data_serialize(NULL, size, tlvP, &format, &(serverP->securityData));
    if (res <= 0) goto error;
    serverP->securityLen = (size_t)res;
    lwm2m_data_free(size, tlvP);

    if (dataP->isBootstrap == false)
    {
#ifndef LWM2M_VERSION_1_0
        int i;
#endif
        size = 4;
#ifndef LWM2M_VERSION_1_0
        if (dataP->registrationPriorityOrder >= 0) size++;
        if (dataP->initialRegistrationDelayTimer >= 0) size++;
        if (dataP->registrationFailureBlock >= 0) size++;
        if (dataP->bootstrapOnRegistrationFailure >= 0) size++;
        if (dataP->communicationRetryCount >= 0) size++;
        if (dataP->communicationRetryTimer >= 0) size++;
        if (dataP->communicationSequenceDelayTimer >= 0) size++;
        if (dataP->communicationSequenceRetryCount >= 0) size++;
#endif

        tlvP = lwm2m_data_new(size);
        if (tlvP == NULL) goto error;

        // Short Server ID
        tlvP[0].id = LWM2M_SERVER_SHORT_ID_ID;
        lwm2m_data_encode_int(dataP->id, tlvP);

        // Lifetime
        tlvP[1].id = LWM2M_SERVER_LIFETIME_ID;
        lwm2m_data_encode_int(dataP->lifetime, tlvP + 1);

        // Notification Storing
        tlvP[2].id = LWM2M_SERVER_STORING_ID;
        lwm2m_data_encode_bool(false, tlvP + 2);

        // Binding
        tlvP[3].id = LWM2M_SERVER_BINDING_ID;
        lwm2m_data_encode_string("U", tlvP + 3);

#ifndef LWM2M_VERSION_1_0
        i = 3;
        if (dataP->registrationPriorityOrder >= 0)
        {
            tlvP[++i].id = LWM2M_SERVER_REG_ORDER_ID;
            lwm2m_data_encode_uint(dataP->registrationPriorityOrder, tlvP + i);
        }
        if (dataP->initialRegistrationDelayTimer >= 0)
        {
            tlvP[++i].id = LWM2M_SERVER_INITIAL_REG_DELAY_ID;
            lwm2m_data_encode_uint(dataP->initialRegistrationDelayTimer, tlvP + i);
        }
        if (dataP->registrationFailureBlock >= 0)
        {
            tlvP[++i].id = LWM2M_SERVER_REG_FAIL_BLOCK_ID;
            lwm2m_data_encode_bool(dataP->registrationFailureBlock > 0, tlvP + i);
        }
        if (dataP->bootstrapOnRegistrationFailure >= 0)
        {
            tlvP[++i].id = LWM2M_SERVER_REG_FAIL_BOOTSTRAP_ID;
            lwm2m_data_encode_bool(dataP->bootstrapOnRegistrationFailure > 0, tlvP + i);
        }
        if (dataP->communicationRetryCount >= 0)
        {
            tlvP[++i].id = LWM2M_SERVER_COMM_RETRY_COUNT_ID;
            lwm2m_data_encode_uint(dataP->communicationRetryCount, tlvP + i);
        }
        if (dataP->communicationRetryTimer >= 0)
        {
            tlvP[++i].id = LWM2M_SERVER_COMM_RETRY_TIMER_ID;
            lwm2m_data_encode_uint(dataP->communicationRetryTimer, tlvP + i);
        }
        if (dataP->communicationSequenceDelayTimer >= 0)
        {
            tlvP[++i].id = LWM2M_SERVER_SEQ_DELAY_TIMER_ID;
            lwm2m_data_encode_uint(dataP->communicationSequenceDelayTimer, tlvP + i);
        }
        if (dataP->communicationSequenceRetryCount >= 0)
        {
            tlvP[++i].id = LWM2M_SERVER_SEQ_RETRY_COUNT_ID;
            lwm2m_data_encode_uint(dataP->communicationSequenceRetryCount, tlvP + i);
        }
#endif
        res = lwm2m_data_serialize(NULL, size, tlvP, &format, &(serverP->serverData));
        if (res <= 0) goto error;
        serverP->serverLen = res;
        lwm2m_data_free(size, tlvP);
    }

    infoP->serverList = (bs_server_tlv_t *)LWM2M_LIST_ADD(infoP->serverList, serverP);

    return 0;

error:
    if (tlvP != NULL) lwm2m_data_free(size, tlvP);
    if (serverP->securityData != NULL) lwm2m_free(serverP->securityData);
    if (serverP->serverData != NULL) lwm2m_free(serverP->serverData);
    lwm2m_free(serverP);

    return -1;
}

static bs_endpoint_info_t * prv_read_next_endpoint(FILE * fd)
{
    char * key;
    char * value;
    bs_endpoint_info_t * endptP;
    int res;
    bs_command_t * cmdP;

    if (prv_find_next_section(fd, "Endpoint") == 0) return NULL;

    endptP = (bs_endpoint_info_t *)lwm2m_malloc(sizeof(bs_endpoint_info_t));
    if (endptP == NULL) return NULL;
    memset(endptP, 0, sizeof(bs_endpoint_info_t));

    cmdP = NULL;

    while((res = prv_read_key_value(fd, &key, &value)) == 1)
    {
        if (strcasecmp(key, "Name") == 0)
        {
            endptP->name = value;
        }
        else if (strcasecmp(key, "Delete") == 0)
        {
            lwm2m_uri_t uri;


            if (lwm2m_stringToUri(value, strlen(value), &uri) == 0) goto error;

            cmdP = (bs_command_t *)lwm2m_malloc(sizeof(bs_command_t));
            if (cmdP == NULL) goto error;
            memset(cmdP, 0, sizeof(bs_command_t));

            cmdP->operation = BS_DELETE;
            if (LWM2M_URI_IS_SET_OBJECT(&uri))
            {
                cmdP->uri = (lwm2m_uri_t *)lwm2m_malloc(sizeof(lwm2m_uri_t));
                if (cmdP->uri == NULL) goto error;
                memcpy(cmdP->uri, &uri, sizeof(lwm2m_uri_t));
            }

            lwm2m_free(value);
        }
        else if (strcasecmp(key, "Server") == 0)
        {
            int num;

            if (sscanf(value, "%d", &num) != 1) goto error;
            if (num <= 0 || num > LWM2M_MAX_ID) goto error;

            cmdP = (bs_command_t *)lwm2m_malloc(sizeof(bs_command_t));
            if (cmdP == NULL) goto error;
            memset(cmdP, 0, sizeof(bs_command_t));
            cmdP->next = (bs_command_t *)lwm2m_malloc(sizeof(bs_command_t));
            if (cmdP->next == NULL) goto error;
            memset(cmdP->next, 0, sizeof(bs_command_t));

            cmdP->operation = BS_WRITE_SECURITY;
            cmdP->serverId = num;
            cmdP->next->operation = BS_WRITE_SERVER;
            cmdP->next->serverId = num;

            lwm2m_free(value);
        }
        else
        {
            // ignore key for now
            lwm2m_free(value);
        }
        lwm2m_free(key);

        if (cmdP != NULL)
        {
            if (endptP->commandList == NULL)
            {
                endptP->commandList = cmdP;
            }
            else
            {
                bs_command_t * parentP;

                parentP = endptP->commandList;
                while (parentP->next != NULL)
                {
                    parentP = parentP->next;
                }
                parentP->next = cmdP;
            }
            cmdP = NULL;
        }
    }
    if (endptP->commandList != NULL)
    {
        bs_command_t * parentP;

        cmdP = (bs_command_t *)lwm2m_malloc(sizeof(bs_command_t));
        if (cmdP == NULL) goto error;
        memset(cmdP, 0, sizeof(bs_command_t));

        cmdP->operation = BS_FINISH;

        parentP = endptP->commandList;
        while (parentP->next != NULL)
        {
            parentP = parentP->next;
        }
        parentP->next = cmdP;
    }

    return endptP;

error:
    if (key != NULL) lwm2m_free(key);
    if (value != NULL) lwm2m_free(value);
    while (cmdP != NULL)
    {
        bs_command_t * tempP;

        if (cmdP->uri != NULL) lwm2m_free(cmdP->uri);
        tempP = cmdP;
        cmdP = cmdP->next;
        lwm2m_free(tempP);
    }
    if (endptP != NULL)
    {
        if (endptP->name != NULL) lwm2m_free(endptP->name);
        while (endptP->commandList != NULL)
        {
            cmdP = endptP->commandList;
            endptP->commandList =endptP->commandList->next;

            if (cmdP->uri != NULL) lwm2m_free(cmdP->uri);
            lwm2m_free(cmdP);
        }
        lwm2m_free(endptP);
    }

    return NULL;
}

bs_info_t *  bs_get_info(FILE * fd)
{
    bs_info_t * infoP;
    read_server_t * readSrvP;
    bs_endpoint_info_t * cltInfoP;

    infoP = (bs_info_t *)lwm2m_malloc(sizeof(bs_info_t));
    if (infoP == NULL) return NULL;
    memset(infoP, 0, sizeof(bs_info_t));

    do
    {
        readSrvP = prv_read_next_server(fd);
        if (readSrvP != NULL)
        {
            if (prv_add_server(infoP, readSrvP) != 0) goto error;
        }
    } while (readSrvP != NULL);

    rewind(fd);
    do
    {
        cltInfoP = prv_read_next_endpoint(fd);
        if (cltInfoP != NULL)
        {
            cltInfoP->next = infoP->endpointList;
            infoP->endpointList = cltInfoP;
        }
    } while (cltInfoP != NULL);

    // check validity
    if (infoP->endpointList == NULL) goto error;

    cltInfoP = infoP->endpointList;
    while (cltInfoP != NULL)
    {
        bs_endpoint_info_t * otherP;
        bs_command_t * cmdP;
        bs_command_t * parentP;

        // check names are unique
        otherP = cltInfoP->next;
        while (otherP != NULL)
        {
            if (cltInfoP->name == NULL)
            {
                if (otherP->name == NULL) goto error;
            }
            else
            {
                if (otherP->name != NULL
                 && strcmp(cltInfoP->name, otherP->name) == 0)
                {
                    goto error;
                }
            }
            otherP = otherP->next;
        }

        // check servers exist
        cmdP = cltInfoP->commandList;
        parentP = NULL;
        // be careful: iterator changes are inside the switch/case
        while (cmdP != NULL)
        {
            switch (cmdP->operation)
            {
            case BS_WRITE_SECURITY:
                if (LWM2M_LIST_FIND(infoP->serverList, cmdP->serverId) == NULL) goto error;
                parentP = cmdP;
                cmdP = cmdP->next;
                break;

            case BS_WRITE_SERVER:
            {
                bs_server_tlv_t * serverP;

                serverP = (bs_server_tlv_t *)LWM2M_LIST_FIND(infoP->serverList, cmdP->serverId);
                if (serverP == NULL) goto error;
                if (serverP->serverData == NULL)
                {
                    // this is a Bootstrap server, remove this command
                    if (parentP == NULL)
                    {
                        cltInfoP->commandList = cmdP->next;
                        lwm2m_free(cmdP);
                        cmdP = cltInfoP->commandList;
                    }
                    else
                    {
                        parentP->next = cmdP->next;
                        lwm2m_free(cmdP);
                        cmdP = parentP->next;
                    }
                }
                else
                {
                    cmdP = cmdP->next;
                }
            }
            break;

            case BS_DELETE:
            default:
                parentP = cmdP;
                cmdP = cmdP->next;
                break;
            }
        }

        cltInfoP = cltInfoP->next;
    }

    return infoP;

error:
    bs_free_info(infoP);
    return NULL;
}

void bs_free_info(bs_info_t * infoP)
{
    if (infoP == NULL) return;

    while (infoP->serverList != NULL)
    {
        bs_server_tlv_t * targetP;

        targetP = infoP->serverList;
        infoP->serverList = infoP->serverList->next;

        if (targetP->securityData != NULL) lwm2m_free(targetP->securityData);
        if (targetP->serverData != NULL) lwm2m_free(targetP->serverData);

        lwm2m_free(targetP);
    }

    while (infoP->endpointList != NULL)
    {
        bs_endpoint_info_t * targetP;

        targetP = infoP->endpointList;
        infoP->endpointList = infoP->endpointList->next;

        if (targetP->name != NULL) lwm2m_free(targetP->name);
        while (targetP->commandList != NULL)
        {
            bs_command_t * cmdP;

            cmdP = targetP->commandList;
            targetP->commandList =targetP->commandList->next;

            if (cmdP->uri != NULL) lwm2m_free(cmdP->uri);
            lwm2m_free(cmdP);
        }

        lwm2m_free(targetP);
    }

    lwm2m_free(infoP);
}


void print_bs_info(bs_info_t *data)
{
    if (data == NULL)
    {
        return;
    }

    bs_server_tlv_t *parent_server = data->serverList;
    bs_endpoint_info_t *parent_ep = data->endpointList;

    printf ("\nPrinting Server List\n");

    /* Print server list. */
    while (parent_server != NULL)
    {
        printf("id = %d", parent_server->id);
        printf ("\n");

        printf ("\nSecurity Data\n");
        output_tlv(stdout, parent_server->securityData, parent_server->securityLen, 0);

        printf ("\n");
        printf ("\nServer Data\n");

        output_tlv(stdout, parent_server->serverData, parent_server->serverLen, 0);

        printf ("\n");

        parent_server = parent_server->next;
    }

    printf ("\nPrinting Endpoint List\n");

    while (parent_ep)
    {
        printf ("\nEnd Point Name = %s\n", parent_ep->name);

        printf ("\nCommand List Name = %s\n", parent_ep->name);

        bs_command_t * parent_commands = parent_ep->commandList;

        while (parent_commands)
        {
            if (parent_commands->uri)
            {
                printf("object id : %d", parent_commands->uri->objectId);
                printf ("\n");
                printf("instance id : %d", parent_commands->uri->instanceId);
                printf ("\n");
                printf("resource id : %d", parent_commands->uri->resourceId);
                printf ("\n");
                printf("resource instance id : %d", parent_commands->uri->resourceInstanceId);
                printf ("\n");
            }

            parent_commands = parent_commands->next;
        }

        printf ("\n");

        parent_ep = parent_ep->next;
    }
}

void free_bs_info(bs_info_t *data)
{
    bs_endpoint_info_t *parent_ep = data->endpointList;
    bs_server_tlv_t    *parent_server = data->serverList;

    while (parent_server)
    {
        /* free the security data memory. */
        lwm2m_free(parent_server->securityData);
        /* free the server data memory. */
        lwm2m_free(parent_server->serverData);
        parent_server = parent_server->next;
    }

    /* free the server list. */
    LWM2M_LIST_FREE(data->serverList);
    data->serverList = NULL;

    while(parent_ep)
    {
        bs_command_t * parent_commands = parent_ep->commandList;

        while (parent_commands)
        {
            /* free the command's uri memory. */
            lwm2m_free(parent_commands->uri);

            parent_commands = parent_commands->next;
        }

        /* free the command list. */
        LWM2M_LIST_FREE(data->endpointList->commandList);
        data->endpointList->commandList = NULL;
        parent_ep = parent_ep->next;
    }

    /* free the endpoint list. */
    LWM2M_LIST_FREE(data->endpointList);
    data->endpointList = NULL;

    lwm2m_free(data);
    data = NULL;
}

void bs_add_server(server_info_t *obj, FILE *fd)
{
    if (obj!= NULL && obj->uri != NULL && fd != NULL)
    {
            char *yes = "yes";
            char *no = "no";
            char *isbootstrap = "empty";
            if (obj->isBootstrap == true)
            {
                isbootstrap = yes;
            }
            else if (obj->isBootstrap == false)
            {
                isbootstrap = no;
            }

            char *nosec = "nosec";
            char *secmode = "empty";
            if (obj->securityMode == LWM2M_SECURITY_MODE_NONE)
            {
                secmode = nosec;
            }
            // todo: add other security modes
            // else if (strcasecmp(value, "PSK") == 0)
            // {
            //     readSrvP->securityMode = LWM2M_SECURITY_MODE_PRE_SHARED_KEY;
            // }
            // else if (strcasecmp(value, "RPK") == 0)
            // {
            //     readSrvP->securityMode = LWM2M_SECURITY_MODE_RAW_PUBLIC_KEY;
            // }
            // else if (strcasecmp(value, "certificate") == 0)
            // {
            //     readSrvP->securityMode = LWM2M_SECURITY_MODE_CERTIFICATE;
            // }

        fprintf(fd, "[server]\nid=%d\nuri=%s\nbootstrap=%s\nlifetime=%d\nsecurity=%s\n", obj->id, obj->uri, isbootstrap, obj->lifetime, secmode);
    }
}

void bs_add_endpoint(char *ep_name,int server_id,  FILE *fd)
{
    if (ep_name!= NULL && fd != NULL)
    {
        fprintf(fd, "[Endpoint]\nname=%s\nDelete=%s\nServer=%d\n",ep_name,"/0",server_id);
    }
}


void bs_write_file()
{
    FILE * fd;

    fd = fopen("test_file.ini", "w");

    fprintf(fd, "[server]\nid=%d\nuri=%s\nbootstrap=%s\nlifetime=%d\nsecurity=%s\n", 1, "coap://localhost:5683","no",300,"NoSec");
    fprintf(fd, "[server]\nid=%d\nuri=%s\nbootstrap=%s\nlifetime=%d\nsecurity=%s\n", 1, "coap://localhost:5683","no",300,"NoSec");

    server_info_t obj;
    memset(&obj, 0, sizeof(obj));
    obj.id = 1;
    obj.isBootstrap = 0;
    obj.lifetime = 300;
    obj.securityMode = LWM2M_SECURITY_MODE_NONE;
    char *uri = "coap://someurl.com";
    obj.uri = uri;

    bs_add_server(&obj, fd);

    fprintf(fd, "[Endpoint]\nname=%s\nDelete=%s\nServer=%d\n","testendpoint","/0",1);
    fprintf(fd, "[Endpoint]\nname=%s\nDelete=%s\nServer=%d\n","testendpoint","/0",1);

    bs_add_endpoint("test123*",10, fd);

    fclose(fd);
}

