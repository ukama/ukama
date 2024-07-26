/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import {
  NotificationResDto,
  UpdateNotificationResDto,
} from "../resolvers/types";
import { dtoToNotificationDto, dtoToUpdateNotificationDto } from "./mapper";

class NotificationApi extends RESTDataSource {
  getNotification = async (
    baseURL: string,
    id: string
  ): Promise<NotificationResDto> => {
    this.logger.info(
      `GetNotification [GET]: ${baseURL}/${VERSION}/event-notification/${id}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/event-notification/${id}`).then(res =>
      dtoToNotificationDto(res)
    );
  };
  updateNotification = async (
    baseURL: string,
    id: string,
    isRead: boolean
  ): Promise<UpdateNotificationResDto> => {
    this.logger.info(
      `UpdateNotification [POST]: ${baseURL}/${VERSION}/event-notification/${id}?is_read=${isRead}`
    );
    this.baseURL = baseURL;
    return this.post(
      `/${VERSION}/event-notification/${id}?is_read=${isRead}`
    ).then(res => dtoToUpdateNotificationDto(res));
  };
}

export default NotificationApi;
