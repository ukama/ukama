import { RESTDataSource } from "@apollo/datasource-rest";
import { BoolResponse, THeaders } from "../../../common/types";
import { SERVER } from "../../../constants/endpoints";
import { getHeaders } from "../../../utils";
import SubscriberMapper from "./mapper";
import {
    SubscriberDto,
    SubscriberInputDto,
    SubscriberMetricsByNetworkDto,
    SubscribersResDto,
    UpdateSubscriberInputDto,
} from "../types";


export class SubscriberApi extends RESTDataSource{
    addSubscriber = async (
        req: SubscriberInputDto,
        headers: THeaders
    ): Promise<SubscriberDto> => {
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,{
            headers: getHeaders(headers),
            body: { ...req },
        }).then(res => 
            SubscriberMapper.dtoToSubscriberResDto(res));
    };

    updateSubscriber = async (
        subscriberId: string,
        req: UpdateSubscriberInputDto,
        headers: THeaders
    ): Promise<BoolResponse> => {
        return this.patch(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,{
            headers: getHeaders(headers),
            body: { ...req },
        }).then(res => 
            {
                return { success: true }
            });
    };
    
    deleteSubscriber = async (
        subscriberId: string,
        headers: THeaders
    ): Promise<BoolResponse> => {
        return this.delete(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,{headers: getHeaders(headers)}).then((res)=>{
            return {
                success: true,
            };
        });
    };

    getSubscriber = async (
        subscriberId: string,
        headers: THeaders
    ): Promise<SubscriberDto> => {
        return this.get(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,{headers: getHeaders(headers)}).then(res => 
            SubscriberMapper.dtoToSubscriberResDto(res));
    };

    getSubMetricsByNetwork =
        async (): Promise<SubscriberMetricsByNetworkDto> => {
            return {
                total: 4,
                active: 1,
                inactive: 1,
                terminated: 1,
            };
        };

    getSubscribersByNetwork = async (
        networkId: string,
        headers: THeaders
    ): Promise<SubscribersResDto> => {
        return this.get(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}s/networks/${networkId}`,{headers: getHeaders(headers)}).then(res => 
            SubscriberMapper.dtoToSubscribersResDto(res));
    };
}
