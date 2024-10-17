/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import dayjs from "dayjs";
import { GraphQLError } from "graphql";

import {
  SubscriberDto,
  SubscriberInputDto,
  SubscriberMetricsByNetworkDto,
  SubscribersResDto,
  UpdateSubscriberInputDto,
} from "../resolver/types";
import { CBooleanResponse } from "./../../common/types";
import {
  addSubscriberReqToSubscriberResDto,
  dtoToSubscriberResDto,
  dtoToSubscribersResDto,
} from "./mapper";

const VERSION = "v1";
const SUBSCRIBER = "subscriber";

class SubscriberApi extends RESTDataSource {
  addSubscriber = async (
    baseURL: string,
    req: SubscriberInputDto
  ): Promise<SubscriberDto> => {
    this.baseURL = baseURL;
    this.logger.info(
      `AddSubscriber [PUT]: ${this.baseURL}/${VERSION}/${SUBSCRIBER}`
    );
    return this.put(`/${VERSION}/${SUBSCRIBER}`, {
      body: {
        address: "none",
        email: req.email,
        phone: req.phone,
        id_serial: "none",
        gender: "undefined",
        name: req.name,
        network_id: req.network_id,
        proof_of_Identification: "default",
        dob: dayjs().subtract(10, "year").format(),
      },
    }).then(res => addSubscriberReqToSubscriberResDto(res));
  };

  updateSubscriber = async (
    baseURL: string,
    subscriberId: string,
    req: UpdateSubscriberInputDto
  ): Promise<CBooleanResponse> => {
    this.baseURL = baseURL;
    this.logger.info(
      `UpdateSubscriber [PATCH]: ${this.baseURL}/${VERSION}/${SUBSCRIBER}/${subscriberId}`
    );
    return this.patch(`/${VERSION}/${SUBSCRIBER}/${subscriberId}`, {
      body: { ...req },
    }).then(() => {
      return { success: true };
    });
  };

  deleteSubscriber = async (
    baseURL: string,
    subscriberId: string
  ): Promise<CBooleanResponse> => {
    this.baseURL = baseURL;
    this.logger.info(
      `DeleteSubscriber [DELETE]: ${this.baseURL}/${VERSION}/${SUBSCRIBER}/${subscriberId}`
    );
    return this.delete(`/${VERSION}/${SUBSCRIBER}/${subscriberId}`).then(() => {
      return {
        success: true,
      };
    });
  };

  getSubscriber = async (
    baseURL: string,
    subscriberId: string
  ): Promise<SubscriberDto> => {
    this.baseURL = baseURL;
    this.logger.info(
      `GetSubscriber [GET]: ${this.baseURL}/${VERSION}/${SUBSCRIBER}/${subscriberId}`
    );
    return this.get(`/${VERSION}/${SUBSCRIBER}/${subscriberId}`).then(res =>
      dtoToSubscriberResDto(res)
    );
  };

  getSubMetricsByNetwork = async (
    baseURL: string
  ): Promise<SubscriberMetricsByNetworkDto> => {
    this.baseURL = baseURL;
    return {
      total: 4,
      active: 1,
      inactive: 1,
      terminated: 1,
    };
  };

  getSubscribersByNetwork = async (
    baseURL: string,
    networkId: string
  ): Promise<SubscribersResDto> => {
    this.logger.info(
      `GetSubscribersByNetwork [GET]: ${baseURL}/${VERSION}/${SUBSCRIBER}s/networks/${networkId}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${SUBSCRIBER}s/networks/${networkId}`)
      .then(res => dtoToSubscribersResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default SubscriberApi;
