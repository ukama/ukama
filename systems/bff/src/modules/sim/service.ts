import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { ParsedCookie } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError } from "../../errors";
import { ISimService } from "./interface";
import SimMapper from "./mapper";
import {
    SetActivePackageForSimInputDto,
    DeleteSimInputDto,
    SetActivePackageForSimResDto,
    RemovePackageFromSimResDto,
    GetPackagesForSimResDto,
    GetPackagesForSimInputDto,
    RemovePackageFormSimInputDto,
    AddPackageSimResDto,
    AddPackageToSimInputDto,
    DeleteSimResDto,
    GetSimByNetworkInputDto,
    AllocateSimInputDto,
    GetSimInputDto,
    GetSimBySubscriberIdInputDto,
    SimStatusResDto,
    SimDetailsDto,
    SimResDto,
    ToggleSimStatusInputDto,
} from "./types";

@Service()
export class SimService implements ISimService {
    allocateSim = async (
        req: AllocateSimInputDto,
        cookie: ParsedCookie
    ): Promise<SimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimResDto(res);
    };
    toggleSimStatus = async (
        req: ToggleSimStatusInputDto,
        cookie: ParsedCookie
    ): Promise<SimStatusResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                simId: req.simId,
                status: req.status,
            },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    getSim = async (
        req: GetSimInputDto,
        cookie: ParsedCookie
    ): Promise<SimDetailsDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                simId: req.simId,
            },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimDetailsDto(res);
    };
    getSimBySubscriberId = async (
        req: GetSimBySubscriberIdInputDto,
        cookie: ParsedCookie
    ): Promise<SimDetailsDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                subscriberId: req.subscriberId,
            },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimDetailsDto(res);
    };
    getSimByNetworkId = async (
        req: GetSimByNetworkInputDto,
        cookie: ParsedCookie
    ): Promise<SimDetailsDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                networkId: req.networkId,
            },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimDetailsDto(res);
    };
    deleteSim = async (
        req: DeleteSimInputDto,
        cookie: ParsedCookie
    ): Promise<DeleteSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                simId: req.simId,
            },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    addPackegeToSim = async (
        req: AddPackageToSimInputDto,
        cookie: ParsedCookie
    ): Promise<AddPackageSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    removePackageFromSim = async (
        req: RemovePackageFormSimInputDto,
        cookie: ParsedCookie
    ): Promise<RemovePackageFromSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    getPackagesForSim = async (
        req: GetPackagesForSimInputDto,
        cookie: ParsedCookie
    ): Promise<GetPackagesForSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                simId: req.simId,
            },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    setActivePackageForSim = async (
        req: SetActivePackageForSimInputDto,
        cookie: ParsedCookie
    ): Promise<SetActivePackageForSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
}
