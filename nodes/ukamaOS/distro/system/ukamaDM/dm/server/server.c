/*******************************************************************************
 *
 * Copyright (c) 2013, 2014 Intel Corporation and others.
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
 *    David Navarro, Intel Corporation - initial API and implementation
 *    domedambrosio - Please refer to git log
 *    Simon Bernard - Please refer to git log
 *    Toby Jaffey - Please refer to git log
 *    Julien Vermillard - Please refer to git log
 *    Bosch Software Innovations GmbH - Please refer to git log
 *    Christian Renz - Please refer to git log
 *    Scott Bertin, AMETEK, Inc. - Please refer to git log
 *
 *******************************************************************************/

/*
 Copyright (c) 2013, 2014 Intel Corporation

 Redistribution and use in source and binary forms, with or without modification,
 are permitted provided that the following conditions are met:

     * Redistributions of source code must retain the above copyright notice,
       this list of conditions and the following disclaimer.
     * Redistributions in binary form must reproduce the above copyright notice,
       this list of conditions and the following disclaimer in the documentation
       and/or other materials provided with the distribution.
     * Neither the name of Intel Corporation nor the names of its contributors
       may be used to endorse or promote products derived from this software
       without specific prior written permission.

 THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
 IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
 INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
 BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
 OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 THE POSSIBILITY OF SUCH DAMAGE.

 David Navarro <david.navarro@intel.com>

*/


#include "liblwm2m.h"

#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdio.h>
#include <ctype.h>
#include <sys/select.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <sys/stat.h>
#include <errno.h>
#include <signal.h>
#include <inttypes.h>

#include "commandline.h"
#include "connection.h"
#include "liblwm2m.h"
#include "ifhandler.h"
#include "notification.h"
#include "request.h"

#define MAX_PACKET_SIZE 1024

#ifdef LWM2M_VERSION_1_0
#define URI_MAX_STRING_LEN    18      // /65535/65535/65535
#else
#define URI_MAX_STRING_LEN    24      // /65535/65535/65535/65535
#endif

static int g_quit = 0;

static void prv_print_error(uint8_t status)
{
    fprintf(stdout, "Error: ");
    print_status(stdout, status);
    fprintf(stdout, "\r\n");
}

static const char * prv_dump_version(lwm2m_version_t version)
{
    switch(version)
    {
    case VERSION_MISSING:
        return "Missing";
    case VERSION_UNRECOGNIZED:
        return "Unrecognized";
    case VERSION_1_0:
        return "1.0";
    case VERSION_1_1:
        return "1.1";
    default:
        return "";
    }
}

static void prv_dump_binding(lwm2m_binding_t binding)
{
    if(BINDING_UNKNOWN == binding)
    {
        fprintf(stdout, "\tbinding: \"Not specified\"\r\n");
    }
    else
    {
        const struct bindingTable
        {
            lwm2m_binding_t binding;
            const char *text;
        } bindingTable[] =
        {
            { BINDING_U, "UDP" },
            { BINDING_T, "TCP" },
            { BINDING_S, "SMS" },
            { BINDING_N, "Non-IP" },
            { BINDING_Q, "queue mode" },
        };
        size_t i;
        bool oneSeen = false;
        fprintf(stdout, "\tbinding: \"");
        for (i = 0; i < sizeof(bindingTable) / sizeof(bindingTable[0]); i++)
        {
            if ((binding & bindingTable[i].binding) != 0)
            {
                if (oneSeen)
                {
                    fprintf(stdout, ", %s", bindingTable[i].text);
                }
                else
                {
                    fprintf(stdout, "%s", bindingTable[i].text);
                    oneSeen = true;
                }
            }
        }
        fprintf(stdout, "\"\r\n");
    }
}

static void prv_dump_client(lwm2m_client_t * targetP)
{
    lwm2m_client_object_t * objectP;

    fprintf(stdout, "Client #%d:\r\n", targetP->internalID);
    fprintf(stdout, "\tname: \"%s\"\r\n", targetP->name);
    fprintf(stdout, "\tversion: \"%s\"\r\n", prv_dump_version(targetP->version));
    prv_dump_binding(targetP->binding);
    if (targetP->msisdn) fprintf(stdout, "\tmsisdn: \"%s\"\r\n", targetP->msisdn);
    if (targetP->altPath) fprintf(stdout, "\talternative path: \"%s\"\r\n", targetP->altPath);
    fprintf(stdout, "\tlifetime: %d sec\r\n", targetP->lifetime);
    fprintf(stdout, "\tobjects: ");
    for (objectP = targetP->objectList; objectP != NULL ; objectP = objectP->next)
    {
        if (objectP->instanceList == NULL)
        {
            fprintf(stdout, "/%d, ", objectP->id);
        }
        else
        {
            lwm2m_list_t * instanceP;

            for (instanceP = objectP->instanceList; instanceP != NULL ; instanceP = instanceP->next)
            {
                fprintf(stdout, "/%d/%d, ", objectP->id, instanceP->id);
            }
        }
    }
    fprintf(stdout, "\r\n");
}

static void prv_output_clients(char * buffer,
                               void * user_data, void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    lwm2m_client_t * targetP;

    targetP = lwm2mH->clientList;

    if (targetP == NULL)
    {
        fprintf(stdout, "No client.\r\n");
        return;
    }

    for (targetP = lwm2mH->clientList ; targetP != NULL ; targetP = targetP->next)
    {
        prv_dump_client(targetP);
    }
}

static int prv_read_id(char * buffer,
                       uint16_t * idP)
{
    int nb;
    int value;

    nb = sscanf(buffer, "%d", &value);
    if (nb == 1)
    {
        if (value < 0 || value > LWM2M_MAX_ID)
        {
            nb = 0;
        }
        else
        {
            *idP = value;
        }
    }

    return nb;
}

static int prv_read_name(char * buffer,
		char * name)
{
	int length;
	int ret = 0;
	// find end of command name
	length = 0;
	while (buffer[length] != 0 && !isspace(buffer[length]&0xFF))
		length++;

	if (length) {
		strncpy(name, buffer, length);
		ret = 1;
	}
	return ret;
}

/* Convert client name to Id */
static int prv_name_to_id(char *name, void *user_data, uint16_t* idP) {
	int ret = 0;
	lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
	lwm2m_client_t * targetP;

	targetP = lwm2mH->clientList;

	if (targetP == NULL)
	{
		fprintf(stdout, "No client.\r\n");
		return ret;
	}

	// Loop through the list
	for (targetP = lwm2mH->clientList ; targetP != NULL ; targetP = targetP->next)
	{
		if ( (strcmp(name, targetP->name)) == 0) {
			*idP = targetP->internalID;
			ret = 1;
			break;
		}
	}

	return ret;
}

/* Convert client Id to Name */
static int prv_id_to_name(char *name, void *user_data, uint16_t* idP) {
	int ret = 0;
	lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
	lwm2m_client_t * targetP;

	targetP = lwm2mH->clientList;

	if (targetP == NULL)
	{
		fprintf(stdout, "No client.\r\n");
		return ret;
	}

	// Loop through the list
	for (targetP = lwm2mH->clientList ; targetP != NULL ; targetP = targetP->next)
	{
		if ( *idP == targetP->internalID) {
			strcpy(name, targetP->name);
			ret = 1;
			break;
		}
	}

	return ret;
}

static void prv_printUri(const lwm2m_uri_t * uriP)
{
    fprintf(stdout, "/%d", uriP->objectId);
    if (LWM2M_URI_IS_SET_INSTANCE(uriP))
        fprintf(stdout, "/%d", uriP->instanceId);
    else if (LWM2M_URI_IS_SET_RESOURCE(uriP))
        fprintf(stdout, "/");
    if (LWM2M_URI_IS_SET_RESOURCE(uriP))
            fprintf(stdout, "/%d", uriP->resourceId);
#ifndef LWM2M_VERSION_1_0
    else if (LWM2M_URI_IS_SET_RESOURCE_INSTANCE(uriP))
        fprintf(stdout, "/");
    if (LWM2M_URI_IS_SET_RESOURCE_INSTANCE(uriP))
            fprintf(stdout, "/%d", uriP->resourceInstanceId);
#endif
}

static void prv_result_callback(uint16_t clientID,
                                lwm2m_uri_t * uriP,
                                int status,
                                lwm2m_media_type_t format,
                                uint8_t * data,
                                int dataLength,
                                void * userData)
{
    fprintf(stdout, "\r\nClient #%d ", clientID);
    prv_printUri(uriP);
    fprintf(stdout, " : ");
    print_status(stdout, status);
    fprintf(stdout, "\r\n");

    output_data(stdout, format, data, dataLength, 1);

    fprintf(stdout, "\r\n> ");
    fflush(stdout);

    /* Sending response to LwM2M Gateway */
    if (userData) {
    	uint32_t reqid = *(uint32_t*)userData;
    	response_handler(reqid, status, format, data, dataLength);
    }

}

static void prv_notify_callback(uint16_t clientID,
                                lwm2m_uri_t * uriP,
                                int count,
                                lwm2m_media_type_t format,
                                uint8_t * data,
                                int dataLength,
                                void * userData)
{
    fprintf(stdout, "\r\n Observe response from client #%d ", clientID);
    prv_printUri(uriP);
    fprintf(stdout, " number %d\r\n", count);

    output_data(stdout, format, data, dataLength, 1);

    fprintf(stdout, "\r\n> ");
    fflush(stdout);

    /* Sending response to LwM2M Gateway */
    if (userData) {
    	uint32_t reqid = *(uint32_t*)userData;
    	response_handler(reqid, COAP_204_CHANGED, format, data, dataLength);
    }

}

static void prv_notify_cb(uint16_t clientID,
								char* name,
                                lwm2m_uri_t * uriP,
                                int count,
                                lwm2m_media_type_t format,
                                uint8_t * data,
                                int dataLength,
                                void * userData)
{
    fprintf(stdout, "\r\n Notify from client #%d ", clientID);
    prv_printUri(uriP);
    fprintf(stdout, " number %d\r\n", count);

    output_data(stdout, format, data, dataLength, 1);

    fprintf(stdout, "\r\n> ");
    fflush(stdout);

    /* URI to String */
    uint8_t baseUriStr[URI_MAX_STRING_LEN] = {'\0'};
    int baseUriLen = uri_toString(uriP, baseUriStr, URI_MAX_STRING_LEN, NULL);
    if (baseUriLen < 0){
        fprintf(stderr, "Failed to build URI for notification.\r\n> ");
        fflush(stderr);
    	return ;
    }

    fprintf(stdout, "Received notification from uuid %s for client id %d notification.\r\n> ", name, clientID);
    fflush(stdout);
    /* Sending notification to gateway */
    notify_handler(name, baseUriStr, count, format, data, dataLength);

}

static uint32_t prv_read_client(char * buffer,
                            void * user_data,
							void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char* end = NULL;
    int result;
    int ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    if (!check_end_of_args(end)) goto syntax_error;

    result = lwm2m_dm_read(lwm2mH, clientId, &uri, prv_result_callback, ctx);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);
    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}

static uint32_t prv_discover_client(char * buffer,
                                void * user_data,
								void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char* end = NULL;
    int result;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    if (!check_end_of_args(end)) goto syntax_error;

    result = lwm2m_dm_discover(lwm2mH, clientId, &uri, prv_result_callback, ctx);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);
    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}

static uint32_t prv_write_client(char * buffer,
                             void * user_data,
							 void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char * end = NULL;
    int result;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(end, &end);
    if (buffer[0] == 0) goto syntax_error;

    if (!check_end_of_args(end)) goto syntax_error;

    result = lwm2m_dm_write(lwm2mH, clientId, &uri, LWM2M_CONTENT_TEXT, (uint8_t *)buffer, end - buffer, prv_result_callback, ctx);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);
    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}


static uint32_t prv_time_client(char * buffer,
                            void * user_data,
							void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char * end = NULL;
    int result;
    lwm2m_attributes_t attr;
    int nb;
    int value;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    memset(&attr, 0, sizeof(lwm2m_attributes_t));
    attr.toSet = LWM2M_ATTR_FLAG_MIN_PERIOD | LWM2M_ATTR_FLAG_MAX_PERIOD;

    buffer = get_next_arg(end, &end);
    if (buffer[0] == 0) goto syntax_error;

    nb = sscanf(buffer, "%d", &value);
    if (nb != 1) goto syntax_error;
    if (value < 0) goto syntax_error;
    attr.minPeriod = value;

    buffer = get_next_arg(end, &end);
    if (buffer[0] == 0) goto syntax_error;

    nb = sscanf(buffer, "%d", &value);
    if (nb != 1) goto syntax_error;
    if (value < 0) goto syntax_error;
    attr.maxPeriod = value;

    if (!check_end_of_args(end)) goto syntax_error;

    result = lwm2m_dm_write_attributes(lwm2mH, clientId, &uri, &attr, prv_result_callback, ctx);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);
    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}


static uint32_t prv_attr_client(char * buffer,
                            void * user_data,
							void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char * end = NULL;
    int result;
    lwm2m_attributes_t attr;
    int nb;
    float value;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    memset(&attr, 0, sizeof(lwm2m_attributes_t));
    attr.toSet = LWM2M_ATTR_FLAG_LESS_THAN | LWM2M_ATTR_FLAG_GREATER_THAN;

    buffer = get_next_arg(end, &end);
    if (buffer[0] == 0) goto syntax_error;

    nb = sscanf(buffer, "%f", &value);
    if (nb != 1) goto syntax_error;
    attr.lessThan = value;

    buffer = get_next_arg(end, &end);
    if (buffer[0] == 0) goto syntax_error;

    nb = sscanf(buffer, "%f", &value);
    if (nb != 1) goto syntax_error;
    attr.greaterThan = value;

    buffer = get_next_arg(end, &end);
    if (buffer[0] != 0)
    {
        nb = sscanf(buffer, "%f", &value);
        if (nb != 1) goto syntax_error;
        attr.step = value;

        attr.toSet |= LWM2M_ATTR_FLAG_STEP;
    }

    if (!check_end_of_args(end)) goto syntax_error;

    result = lwm2m_dm_write_attributes(lwm2mH, clientId, &uri, &attr, prv_result_callback, ctx);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);

    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}


static uint32_t prv_clear_client(char * buffer,
                             void * user_data,
							 void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char * end = NULL;
    int result;
    lwm2m_attributes_t attr;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;

    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    memset(&attr, 0, sizeof(lwm2m_attributes_t));
    attr.toClear = LWM2M_ATTR_FLAG_LESS_THAN | LWM2M_ATTR_FLAG_GREATER_THAN | LWM2M_ATTR_FLAG_STEP | LWM2M_ATTR_FLAG_MIN_PERIOD | LWM2M_ATTR_FLAG_MAX_PERIOD ;

    buffer = get_next_arg(end, &end);
    if (!check_end_of_args(end)) goto syntax_error;

    result = lwm2m_dm_write_attributes(lwm2mH, clientId, &uri, &attr, prv_result_callback, NULL);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);
    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}


static uint32_t prv_exec_client(char * buffer,
                            void * user_data,
							void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char * end = NULL;
    int result;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(end, &end);


    if (buffer[0] == 0)
    {
        result = lwm2m_dm_execute(lwm2mH, clientId, &uri, 0, NULL, 0, prv_result_callback, ctx);
    }
    else
    {
        if (!check_end_of_args(end)) goto syntax_error;

        result = lwm2m_dm_execute(lwm2mH, clientId, &uri, LWM2M_CONTENT_TEXT, (uint8_t *)buffer, end - buffer, prv_result_callback, ctx);
    }

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);
    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}

static uint32_t  prv_create_client(char * buffer,
                              void * user_data,
							  void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char * end = NULL;
    int result;
    int64_t value;
    uint8_t * temp_buffer = NULL;
    int temp_length = 0;
    lwm2m_media_type_t format = LWM2M_CONTENT_TEXT;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    //Get Uri
    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    //Get Data to Post
    buffer = get_next_arg(end, &end);
    if (buffer[0] == 0) goto syntax_error;

    if (!check_end_of_args(end)) goto syntax_error;

   // TLV

   /* Client dependent part   */

    if (uri.objectId == 31024)
    {
        lwm2m_data_t * dataP;

        if (1 != sscanf(buffer, "%"PRId64, &value))
        {
            fprintf(stdout, "Invalid value !");
            return ret;
        }

        dataP = lwm2m_data_new(1);
        if (dataP == NULL)
        {
            fprintf(stdout, "Allocation error !");
            return ret;
        }
        lwm2m_data_encode_int(value, dataP);
        dataP->id = 1;

        format = LWM2M_CONTENT_TLV;
        //format = LWM2M_CONTENT_JSON;
        temp_length = lwm2m_data_serialize(NULL, 1, dataP, &format, &temp_buffer);
    }
   /* End Client dependent part*/

    //Create
    result = lwm2m_dm_create(lwm2mH, clientId, &uri, format, temp_buffer, temp_length, prv_result_callback, ctx);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);
    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}

static uint32_t prv_delete_client(char * buffer,
                              void * user_data,
							  void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char* end = NULL;
    int result;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    if (!check_end_of_args(end)) goto syntax_error;

    result = lwm2m_dm_delete(lwm2mH, clientId, &uri, prv_result_callback, ctx);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);
    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}

static uint32_t prv_observe_client(char * buffer,
                               void * user_data,
							   void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char* end = NULL;
    int result;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    if (!check_end_of_args(end)) goto syntax_error;

    result = lwm2m_observe(lwm2mH, clientId, &uri, prv_notify_callback, prv_notify_cb, ctx);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);
    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}

static uint32_t prv_cancel_client(char * buffer,
                              void * user_data,
							  void *ctx)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) user_data;
    uint16_t clientId;
    lwm2m_uri_t uri;
    char* end = NULL;
    int result;
    uint32_t ret = COAP_500_INTERNAL_SERVER_ERROR;
    char name[64]= {"\0"};
    result = prv_read_name(buffer, name);
    if (result == 0) goto syntax_error;

    result = prv_name_to_id(name, user_data, &clientId);
    if (result == 0) goto syntax_error;

    buffer = get_next_arg(buffer, &end);
    if (buffer[0] == 0) goto syntax_error;

    result = lwm2m_stringToUri(buffer, end - buffer, &uri);
    if (result == 0) goto syntax_error;

    if (!check_end_of_args(end)) goto syntax_error;

    result = lwm2m_observe_cancel(lwm2mH, clientId, &uri, prv_result_callback, ctx);

    if (result == 0)
    {
        fprintf(stdout, "OK");
        ret = 0;
    }
    else
    {
        prv_print_error(result);

    }
    return ret;

syntax_error:
    fprintf(stdout, "Syntax error !");
    return ret;
}

static void prv_monitor_callback(uint16_t clientID,
                                 lwm2m_uri_t * uriP,
                                 int status,
                                 lwm2m_media_type_t format,
                                 uint8_t * data,
                                 int dataLength,
                                 void * userData)
{
    lwm2m_context_t * lwm2mH = (lwm2m_context_t *) userData;
    lwm2m_client_t * targetP;

    switch (status)
    {
    case COAP_201_CREATED:
        fprintf(stdout, "\r\nNew client #%d registered.\r\n", clientID);

        targetP = (lwm2m_client_t *)lwm2m_list_find((lwm2m_list_t *)lwm2mH->clientList, clientID);

        prv_dump_client(targetP);
        break;

    case COAP_202_DELETED:
        fprintf(stdout, "\r\nClient #%d unregistered.\r\n", clientID);
        break;

    case COAP_204_CHANGED:
        fprintf(stdout, "\r\nClient #%d updated.\r\n", clientID);

        targetP = (lwm2m_client_t *)lwm2m_list_find((lwm2m_list_t *)lwm2mH->clientList, clientID);

        prv_dump_client(targetP);
        break;

    default:
        fprintf(stdout, "\r\nMonitor callback called with an unknown status: %d.\r\n", status);
        break;
    }

    fprintf(stdout, "\r\n> ");
    fflush(stdout);
}


static void prv_quit(char * buffer,
                     void * user_data)
{
    g_quit = 1;
}

void handle_sigint(int signum)
{
    g_quit = 2;
}

void print_usage(void)
{
    fprintf(stderr, "Usage: lwm2mserver [OPTION]\r\n");
    fprintf(stderr, "Launch a LWM2M server on localhost.\r\n\n");
    fprintf(stdout, "Options:\r\n");
    fprintf(stdout, "  -4\t\tUse IPv4 connection. Default: IPv6 connection\r\n");
    fprintf(stdout, "  -l PORT\tSet the local UDP port of the Server. Default: "LWM2M_STANDARD_PORT_STR"\r\n");
    fprintf(stdout, "\r\n");
}

 command_desc_t commands[] =
    {
            {"list", "List registered clients.", NULL, prv_output_clients, NULL},
            {"read", "Read from a client.", " read CLIENT# URI\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri to read such as /3, /3/0/2, /1024/11, /1024/0/1\r\n"
                                            "Result will be displayed asynchronously.", prv_read_client, NULL},
            {"disc", "Discover resources of a client.", " disc CLIENT# URI\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri to discover such as /3, /3/0/2, /1024/11, /1024/0/1\r\n"
                                            "Result will be displayed asynchronously.", prv_discover_client, NULL},
            {"write", "Write to a client.", " write CLIENT# URI DATA\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri to write to such as /3, /3/0/2, /1024/11, /1024/0/1\r\n"
                                            "   DATA: data to write\r\n"
                                            "Result will be displayed asynchronously.", prv_write_client, NULL},
            {"time", "Write time-related attributes to a client.", " time CLIENT# URI PMIN PMAX\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri to write attributes to such as /3, /3/0/2, /1024/11, /1024/0/1\r\n"
                                            "   PMIN: Minimum period\r\n"
                                            "   PMAX: Maximum period\r\n"
                                            "Result will be displayed asynchronously.", prv_time_client, NULL},
            {"attr", "Write value-related attributes to a client.", " attr CLIENT# URI LT GT [STEP]\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri to write attributes to such as /3/0/2, /1024/0/1\r\n"
                                            "   LT: \"Less than\" value\r\n"
                                            "   GT: \"Greater than\" value\r\n"
                                            "   STEP: \"Step\" value\r\n"
                                            "Result will be displayed asynchronously.", prv_attr_client, NULL},
            {"clear", "Clear attributes of a client.", " clear CLIENT# URI\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri to clear attributes of such as /3, /3/0/2, /1024/11, /1024/0/1\r\n"
                                            "Result will be displayed asynchronously.", prv_clear_client, NULL},
            {"exec", "Execute a client resource.", " exec CLIENT# URI\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri of the resource to execute such as /3/0/2\r\n"
                                            "Result will be displayed asynchronously.", prv_exec_client, NULL},
            {"del", "Delete a client Object instance.", " del CLIENT# URI\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri of the instance to delete such as /1024/11\r\n"
                                            "Result will be displayed asynchronously.", prv_delete_client, NULL},
            {"create", "Create an Object instance.", " create CLIENT# URI DATA\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri to which create the Object Instance such as /1024, /1024/45 \r\n"
                                            "   DATA: data to initialize the new Object Instance (0-255 for object 31024) \r\n"
                                            "Result will be displayed asynchronously.", prv_create_client, NULL},
            {"observe", "Observe from a client.", " observe CLIENT# URI\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri to observe such as /3, /3/0/2, /1024/11\r\n"
                                            "Result will be displayed asynchronously.", prv_observe_client, NULL},
            {"cancel", "Cancel an observe.", " cancel CLIENT# URI\r\n"
                                            "   CLIENT#: client number as returned by command 'list'\r\n"
                                            "   URI: uri on which to cancel an observe such as /3, /3/0/2, /1024/11\r\n"
                                            "Result will be displayed asynchronously.", prv_cancel_client, NULL},

            {"q", "Quit the server.", NULL, prv_quit, NULL},

            COMMAND_END_LIST
    };


int main(int argc, char *argv[])
{
    int sock;
    fd_set readfds;
    struct timeval tv;
    int result;
    lwm2m_context_t * lwm2mH = NULL;
    int i;
    connection_t * connList = NULL;
    int addressFamily = AF_INET6;
    int opt;
    const char * localPort = LWM2M_STANDARD_PORT_STR;

    opt = 1;
    while (opt < argc)
    {
        if (argv[opt] == NULL
            || argv[opt][0] != '-'
            || argv[opt][2] != 0)
        {
            print_usage();
            return 0;
        }
        switch (argv[opt][1])
        {
        case '4':
            addressFamily = AF_INET;
            break;
        case 'l':
            opt++;
            if (opt >= argc)
            {
                print_usage();
                return 0;
            }
            localPort = argv[opt];
            break;
        default:
            print_usage();
            return 0;
        }
        opt += 1;
    }

    sock = create_socket(localPort, addressFamily);
    if (sock < 0)
    {
        fprintf(stderr, "Error opening socket: %d\r\n", errno);
        return -1;
    }

    lwm2mH = lwm2m_init(NULL);
    if (NULL == lwm2mH)
    {
        fprintf(stderr, "lwm2m_init() failed\r\n");
        return -1;
    }

    signal(SIGINT, handle_sigint);

    /*
     * Server started Aysnc messages from Server
     */
    pthread_t conn_id = connection_handler_start((void*) lwm2mH);
    fprintf(stdout, "Connection Handler server started with thread id %ld.\r\n", conn_id);

    for (i = 0 ; commands[i].name != NULL ; i++)
    {
        commands[i].userData = (void *)lwm2mH;
    }
    fprintf(stdout, "> "); fflush(stdout);

    lwm2m_set_monitoring_callback(lwm2mH, prv_monitor_callback, lwm2mH);

    while (0 == g_quit)
    {
        FD_ZERO(&readfds);
        FD_SET(sock, &readfds);
        FD_SET(STDIN_FILENO, &readfds);

        tv.tv_sec = 60;
        tv.tv_usec = 0;

        result = lwm2m_step(lwm2mH, &(tv.tv_sec));
        if (result != 0)
        {
            fprintf(stderr, "lwm2m_step() failed: 0x%X\r\n", result);
            return -1;
        }

        result = select(FD_SETSIZE, &readfds, 0, 0, &tv);

        if ( result < 0 )
        {
            if (errno != EINTR)
            {
              fprintf(stderr, "Error in select(): %d\r\n", errno);
            }
        }
        else if (result > 0)
        {
            uint8_t buffer[MAX_PACKET_SIZE];
            int numBytes;

            if (FD_ISSET(sock, &readfds))
            {
                struct sockaddr_storage addr;
                socklen_t addrLen;

                addrLen = sizeof(addr);
                numBytes = recvfrom(sock, buffer, MAX_PACKET_SIZE, 0, (struct sockaddr *)&addr, &addrLen);

                if (numBytes == -1)
                {
                    fprintf(stderr, "Error in recvfrom(): %d\r\n", errno);
                }
                else
                {
                    char s[INET6_ADDRSTRLEN];
                    in_port_t port;
                    connection_t * connP;

					s[0] = 0;
                    if (AF_INET == addr.ss_family)
                    {
                        struct sockaddr_in *saddr = (struct sockaddr_in *)&addr;
                        inet_ntop(saddr->sin_family, &saddr->sin_addr, s, INET6_ADDRSTRLEN);
                        port = saddr->sin_port;
                    }
                    else if (AF_INET6 == addr.ss_family)
                    {
                        struct sockaddr_in6 *saddr = (struct sockaddr_in6 *)&addr;
                        inet_ntop(saddr->sin6_family, &saddr->sin6_addr, s, INET6_ADDRSTRLEN);
                        port = saddr->sin6_port;
                    }

                    fprintf(stderr, "%d bytes received from [%s]:%hu\r\n", numBytes, s, ntohs(port));
                    output_buffer(stderr, buffer, numBytes, 0);

                    connP = connection_find(connList, &addr, addrLen);
                    if (connP == NULL)
                    {
                        connP = connection_new_incoming(connList, sock, (struct sockaddr *)&addr, addrLen);
                        if (connP != NULL)
                        {
                            connList = connP;
                        }
                    }
                    if (connP != NULL)
                    {
                        lwm2m_handle_packet(lwm2mH, buffer, numBytes, connP);
                    }
                }
            }
            else if (FD_ISSET(STDIN_FILENO, &readfds))
            {
                numBytes = read(STDIN_FILENO, buffer, MAX_PACKET_SIZE - 1);

                if (numBytes > 1)
                {
                    buffer[numBytes] = 0;
                    handle_command(commands, (char*)buffer);
                    fprintf(stdout, "\r\n");
                }
                if (g_quit == 0)
                {
                    fprintf(stdout, "> ");
                    fflush(stdout);
                }
                else
                {
                    fprintf(stdout, "\r\n");
                }
            }
        }
    }

    lwm2m_close(lwm2mH);
    close(sock);
    connection_free(connList);

    connection_handler_stop(conn_id);
#ifdef MEMORY_TRACE
    if (g_quit == 1)
    {
        trace_print(0, 1);
    }
#endif

    return 0;
}

uint32_t handle_server_req(RequestMsg* req) {
	  int ret = COAP_404_NOT_FOUND;
#ifdef DEBUG_GATEWAYIF
	  fprintf(stdout," Received command for id %llu sock %d length %d msg: %s\r\n", req->reqid, req->sock, req->length, req->msg);
	  fflush(stdout);
#endif
	  ret = handle_request_command(commands, req->msg, &req->reqid);
	  return ret;
}
