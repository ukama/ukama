/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "liblwm2m.h"
#include "inc/ereg.h"
#include "object_helper.h"
#include "objects/atten.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>
#include <netdb.h>
#include <arpa/inet.h>

uint8_t objh_set_bool_value(lwm2m_data_t * dataArray, bool * data) {
	int ret = 0;
	bool value;
	uint8_t result = COAP_400_BAD_REQUEST;
	if (1 == lwm2m_data_decode_bool(dataArray, &value))
	{
		if ((value == false) || (value == true) )
		{
			*data = value;
			result = COAP_204_CHANGED;
		}
		else
		{
			result = COAP_406_NOT_ACCEPTABLE;
		}
	}
	return result;
}

uint8_t objh_set_int_value(lwm2m_data_t * dataArray, uint32_t * data) {
	int ret = 0;
	int64_t value;
	uint8_t result = COAP_400_BAD_REQUEST;
	if (1 == lwm2m_data_decode_int(dataArray, &value))
	{
		if (value >= 0 && value <= 0xFFFFFFFF)
		{
			*data = value;
			result = COAP_204_CHANGED;
		}
		else
		{
			result = COAP_406_NOT_ACCEPTABLE;
		}
	}
	return result;
}

uint8_t objh_set_str_value(lwm2m_data_t * dataArray, char* data) {
	int64_t value;
	uint8_t result = COAP_400_BAD_REQUEST;
	if ( dataArray->type == LWM2M_TYPE_STRING
			&& dataArray->value.asBuffer.length > 0 ) {
		if (data) {
			lwm2m_free(data);
		}
		size_t szstr = dataArray->value.asBuffer.length + 1;
		data = (char *)lwm2m_malloc(szstr);
		if (data) {
			memset(data, 0, szstr);
			strncpy(data, (char*)dataArray->value.asBuffer.buffer, szstr);
			result = COAP_204_CHANGED;
		} else {
			result =  COAP_500_INTERNAL_SERVER_ERROR;
		}
	}
	if (data) {
		lwm2m_free(data);
	}
	return result;
}

uint8_t objh_set_double_value(lwm2m_data_t * dataArray, double * data) {
	int ret = 0;
	double value;
	uint8_t result = COAP_400_BAD_REQUEST;
	if (1 == lwm2m_data_decode_float(dataArray, &value))
	{
		if (value >= 0 && value <= 0x7FFFFFFFFFFF)
		{
			*data = value;
			result = COAP_204_CHANGED;
		}
		else
		{
			result = COAP_406_NOT_ACCEPTABLE;
		}
	}
	return result;
}

uint8_t objh_send_data_ukama_edr(uint16_t instanceId, uint16_t rid, int objectType, void* data, size_t *size) {
	int ret = 0;
	ret = ereg_write_inst(instanceId, objectType, rid, data, size);
	if (ret) {
		return COAP_500_INTERNAL_SERVER_ERROR;
	} else {
		return COAP_204_CHANGED;
	}
}

int objh_store_data(char* filename, char* data, int size) {
	int ret = 0;
	FILE * fd;
	int blocks = 1;

	/* Open file to write */
	fd = fopen(filename, "w+");
	if (fd == NULL) {
		fprintf(stderr, "Opening file %s failed.\r\n", filename);
		return -1;
	}

	/* Write data */
	ret = fwrite(data,size, blocks, fd);
	if (ret != blocks ) {
		fprintf(stderr, "Failed to write all data to file %s. Expected %d bytes Written %d bytes.\r\n", filename, size, ret);
		ret = -1;
		goto close;
	} else {
		fprintf(stdout, "Wrote %d blocks to file %s.\r\n", blocks, filename);
		ret = 0;
	}

	close:
	/* Close file */
	fclose(fd);
	return ret;
}

int objh_parse_addr(char* url, int size, char** addr) {

	char *local_url = (char *) malloc(sizeof(char) * (strlen(url) + 1));
	char *token;
	char *token_host;
	char *token_port;

	char *token_ptr;
	char *host_token_ptr;

	char *path = NULL;
	char *protocol = NULL;
	char *host_addr = NULL;
	char *host_port = NULL;

	int host_exists;
	bool verify_host = true;
	int port;

	int ret = 0;
	// Copy our string
	strcpy(local_url, url);

	/* Protocol */
	token = strtok_r(local_url, ":", &token_ptr);
	if (token) {
		protocol = (char *) malloc(sizeof(char) * strlen(token) + 1);
		if (protocol) {
			strcpy(protocol, token);
			fprintf(stdout, "protocol: %s\n", protocol);
		}
	} else {
		ret = -1;
		fprintf(stderr, "Error parsing protocol.\n");
		goto error;
	}

	/* Host:Port */
	token = strtok_r(NULL, "/", &token_ptr);
	if (token) {

		host_port = (char *) malloc(sizeof(char) * (strlen(token) + 1));
		if (host_port) {
			strcpy(host_port, token);
		} else {
			ret = -1;
			goto error;
		}

	} else {
		ret = -1;
		fprintf(stderr, "Error parsing host_port.\n");
		goto error;
	}
	fprintf(stdout, "host_port: %s\n", host_port);

	/* Host */
	token_host = strtok_r(host_port, ":", &host_token_ptr);
	if (token_host) {
		host_addr = (char *) malloc(sizeof(char) * strlen(token_host) + 1);
		if (host_addr){
			strcpy(host_addr, token_host);
			fprintf(stdout, "host_addr: %s\n", host_addr);
		} else {
			ret = -1;
			fprintf(stderr, "Error in memory allocation for host_addr.\n");
			goto error;
		}
	} else {
		ret = -1;
		fprintf(stderr, "Error parsing host.\n");
		goto error;
	}

	/* Port */
	token_port = strtok_r(NULL, ":", &host_token_ptr);
	if (token_port) {
		port = atoi(token_port);
		fprintf(stdout, "token_port: %d\n", port);
	} else {
		port = 0;
		ret = -1;
		fprintf(stderr, "Error parsing token_port.\n");
		/* No need to report error */
	}

	*addr = host_addr;

	/* Only error case */
	error:
	if (protocol) free(protocol);
	if (host_port) free(host_port);
	if (ret) {
		if (host_addr) {
			free(host_addr);
			host_addr = NULL;
		}
	}

	return ret;
}
