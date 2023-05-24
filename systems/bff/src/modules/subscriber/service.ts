import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { BoolResponse, THeaders } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError } from "../../errors";
import { getHeaders } from "../../utils";
import { ISubscriberService } from "./interface";
import SubscriberMapper from "./mapper";
import {
    SubscriberDto,
    SubscriberInputDto,
    SubscriberMetricsByNetworkDto,
    SubscribersResDto,
    UpdateSubscriberInputDto,
} from "./types";

@Service()
export class SubscriberService implements ISubscriberService {
    addSubscriber = async (
        req: SubscriberInputDto,
        headers: THeaders
    ): Promise<SubscriberDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return SubscriberMapper.dtoToSubscriberResDto(res);
    };
    updateSubscriber = async (
        subscriberId: string,
        req: UpdateSubscriberInputDto,
        headers: THeaders
    ): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,
            body: { ...req },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return { success: true };
    };
    deleteSubscriber = async (
        subscriberId: string,
        headers: THeaders
    ): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return { success: true };
    };
    getSubscriber = async (
        subscriberId: string,
        headers: THeaders
    ): Promise<SubscriberDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return SubscriberMapper.dtoToSubscriberResDto(res);
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
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}s/networks/${networkId}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return SubscriberMapper.dtoToSubscribersResDto(res);
    };
}
