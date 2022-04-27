/*
 * request.h
 *
 *  Created on: Mar 10, 2021
 *      Author: vishal
 */

#ifndef SERVIFACE_REQUEST_H_
#define SERVIFACE_REQUEST_H_

#include <stdint.h>

typedef struct {
	uint32_t reqid;
	int sock;
	uint32_t length;
	char* msg;
} RequestMsg;


typedef struct {
	uint32_t reqid;
	uint32_t status;
	uint32_t format;
	uint32_t length;
	uint8_t * msg; // DATA_OFFSET depends on this structure
}ResponseMsg;

int init_req_list();
int response_handler(uint32_t reqid, uint32_t status, uint32_t format, uint8_t* data, int length);
void serve_request(char* msg, int sock);

#endif /* SERVIFACE_REQUEST_H_ */
