/*
 * notification.h
 *
 *  Created on: Mar 10, 2021
 *      Author: vishal
 */

#ifndef SERVIFACE_NOTIFICATION_H_
#define SERVIFACE_NOTIFICATION_H_
#include <stdint.h>

typedef struct {
	uint32_t id;
	char uuid[ATTR_MAX_LEN];
	char uri[ATTR_MAX_LEN];
	uint32_t count;
	uint32_t format;
	uint32_t length;
	uint8_t * data;
}NotifyEvent;

typedef struct {
	uint32_t id;
	uint32_t status;
}NotifyEventResp;


int init_evt_list();
int notify_handler( char* uuid, uint8_t* uri, uint32_t count, uint32_t format, uint8_t* data, int length);

#endif /* SERVIFACE_NOTIFICATION_H_ */
