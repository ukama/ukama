import { RESTDataSource } from "@apollo/datasource-rest";

import {
  SubscriberDto,
  SubscriberInputDto,
  SubscriberMetricsByNetworkDto,
  SubscribersResDto,
  UpdateSubscriberInputDto,
} from "../resolver/types";
import { BoolResponse } from "./../../common/types";
import { SERVER } from "./../../constants/endpoints";
import { dtoToSubscriberResDto, dtoToSubscribersResDto } from "./mapper";

class SubscriberApi extends RESTDataSource {
  addSubscriber = async (req: SubscriberInputDto): Promise<SubscriberDto> => {
    return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
      body: { ...req },
    }).then(res => dtoToSubscriberResDto(res));
  };

  updateSubscriber = async (
    subscriberId: string,
    req: UpdateSubscriberInputDto
  ): Promise<BoolResponse> => {
    return this.patch(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`, {
      body: { ...req },
    }).then(res => {
      return { success: true };
    });
  };

  deleteSubscriber = async (subscriberId: string): Promise<BoolResponse> => {
    return this.delete(
      `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`
    ).then(res => {
      return {
        success: true,
      };
    });
  };

  getSubscriber = async (subscriberId: string): Promise<SubscriberDto> => {
    return this.get(
      `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,
      {}
    ).then(res => dtoToSubscriberResDto(res));
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
    return this.get(
      `${SERVER.SUBSCRIBER_REGISTRY_API_URL}s/networks/${networkId}`
    ).then(res => dtoToSubscribersResDto(res));
  };
}

export default SubscriberApi;
