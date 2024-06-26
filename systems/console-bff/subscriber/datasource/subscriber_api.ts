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

import { SUBSCRIBER_API_GW } from "../../common/configs";
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
  baseURL = SUBSCRIBER_API_GW;

  addSubscriber = async (req: SubscriberInputDto): Promise<SubscriberDto> => {
    this.logger.info(`Request Url: ${this.baseURL}/${VERSION}/${SUBSCRIBER}`);
    return this.put(`/${VERSION}/${SUBSCRIBER}`, {
      body: {
        address: "none",
        email: req.email,
        phone: req.phone,
        id_serial: "none",
        gender: "undefined",
        last_name: req.last_name,
        first_name: req.first_name,
        network_id: req.network_id,
        proof_of_Identification: "default",
        dob: dayjs().subtract(10, "year").format(),
      },
    }).then(res => addSubscriberReqToSubscriberResDto(res));
  };

  updateSubscriber = async (
    subscriberId: string,
    req: UpdateSubscriberInputDto
  ): Promise<CBooleanResponse> => {
    this.logger.info(
      `Request Url: ${this.baseURL}/${VERSION}/${SUBSCRIBER}/${subscriberId}`
    );
    return this.patch(`/${VERSION}/${SUBSCRIBER}/${subscriberId}`, {
      body: { ...req },
    }).then(() => {
      return { success: true };
    });
  };

  deleteSubscriber = async (
    subscriberId: string
  ): Promise<CBooleanResponse> => {
    this.logger.info(
      `Request Url: ${this.baseURL}/${VERSION}/${SUBSCRIBER}/${subscriberId}`
    );
    return this.delete(`/${VERSION}/${SUBSCRIBER}/${subscriberId}`).then(() => {
      return {
        success: true,
      };
    });
  };

  getSubscriber = async (subscriberId: string): Promise<SubscriberDto> => {
    this.logger.info(
      `Request Url: ${this.baseURL}/${VERSION}/${SUBSCRIBER}/${subscriberId}`
    );
    return this.get(`/${VERSION}/${SUBSCRIBER}/${subscriberId}`).then(res =>
      dtoToSubscriberResDto(res)
    );
  };

  getSubMetricsByNetwork = async (): Promise<SubscriberMetricsByNetworkDto> => {
    return {
      total: 4,
      active: 1,
      inactive: 1,
      terminated: 1,
    };
  };

  getSubscribersByNetwork = async (
    networkId: string
  ): Promise<SubscribersResDto> => {
    return this.get(`/${VERSION}/${SUBSCRIBER}s/networks/${networkId}`)
      .then(res => dtoToSubscribersResDto(res))

      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default SubscriberApi;
