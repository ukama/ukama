import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { BoolResponse, ParsedCookie } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError } from "../../errors";
import { ISubscriberService } from "./interface";
import SubscriberMapper from "./mapper";
import {
    SubscriberAPIResDto,
    SubscriberDto,
    SubscriberInputDto,
    SubscriberMetricsByNetworkDto,
    UpdateSubscriberInputDto,
} from "./types";

@Service()
export class SubscriberService implements ISubscriberService {
    dtoToSubscriberResDto(res: SubscriberAPIResDto): SubscriberDto {
        throw new Error("Method not implemented.");
    }
    addSubscriber = async (
        req: SubscriberInputDto,
        cookie: ParsedCookie
    ): Promise<SubscriberDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return SubscriberMapper.dtoToSubscriberResDto(res);
    };
    updateSubscriber = async (
        subscriberId: string,
        req: UpdateSubscriberInputDto,
        cookie: ParsedCookie
    ): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,
            body: { ...req },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return { success: true };
    };
    deleteSubscriber = async (
        subscriberId: string,
        cookie: ParsedCookie
    ): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return { success: true };
    };
    getSubscriber = async (
        subscriberId: string,
        cookie: ParsedCookie
    ): Promise<SubscriberDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${subscriberId}`,
            headers: cookie.header,
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
}
