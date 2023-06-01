import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { THeaders } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError } from "../../errors";
import { getHeaders } from "../../utils";
import generateTokenFromIccid from "../../utils/generateSimToken";
import { ISimService } from "./interface";
import SimMapper from "./mapper";
import {
    AddPackageSimResDto,
    AddPackageToSimInputDto,
    AllocateSimInputDto,
    DeleteSimInputDto,
    DeleteSimResDto,
    GetPackagesForSimInputDto,
    GetPackagesForSimResDto,
    GetSimByNetworkInputDto,
    GetSimBySubscriberIdInputDto,
    GetSimInputDto,
    RemovePackageFormSimInputDto,
    RemovePackageFromSimResDto,
    SetActivePackageForSimInputDto,
    SetActivePackageForSimResDto,
    SimDataUsage,
    SimDetailsDto,
    SimDto,
    SimPoolStatsDto,
    SimStatusResDto,
    SimsResDto,
    ToggleSimStatusInputDto,
    UploadSimsInputDto,
    UploadSimsResDto,
} from "./types";

@Service()
export class SimService implements ISimService {
    uploadSims = async (
        req: UploadSimsInputDto,
        headers: THeaders
    ): Promise<UploadSimsResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_SIMPOOL_API_URL}/upload`,
            body: {
                data: req.data,
                sim_type: req.simType,
            },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    allocateSim = async (
        req: AllocateSimInputDto,
        headers: THeaders
    ): Promise<SimDto> => {
        const token = generateTokenFromIccid(
            req.iccid,
            process.env.ENCRYPTION_KEY || ""
        );
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                ...req,
                sim_token: token,
            },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimResDto(res);
    };
    toggleSimStatus = async (
        req: ToggleSimStatusInputDto,
        headers: THeaders
    ): Promise<SimStatusResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                simId: req.simId,
                status: req.status,
            },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    getSim = async (
        req: GetSimInputDto,
        headers: THeaders
    ): Promise<SimDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                simId: req.simId,
            },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimResDto(res);
    };
    getSims = async (type: string, headers: THeaders): Promise<SimsResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.SUBSCRIBER_SIMPOOL_API_URL}/sims/${type}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimsDto(res);
    };
    getDataUsage = async (
        simId: string,
        headers: THeaders
    ): Promise<SimDataUsage> => {
        // const res = await catchAsyncIOMethod({
        //     type: API_METHOD_TYPE.GET,
        //     path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}/${simId}`,
        //     getHeaders(headers)
        // });
        // if (checkError(res)) throw new Error(res.message);
        return {
            usage: "1240",
        };
    };
    getSimBySubscriberId = async (
        req: GetSimBySubscriberIdInputDto,
        headers: THeaders
    ): Promise<SimDetailsDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                subscriberId: req.subscriberId,
            },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimDetailsDto(res);
    };
    getSimByNetworkId = async (
        req: GetSimByNetworkInputDto,
        headers: THeaders
    ): Promise<SimDetailsDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                networkId: req.networkId,
            },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimDetailsDto(res);
    };
    deleteSim = async (
        req: DeleteSimInputDto,
        headers: THeaders
    ): Promise<DeleteSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                simId: req.simId,
            },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    addPackegeToSim = async (
        req: AddPackageToSimInputDto,
        headers: THeaders
    ): Promise<AddPackageSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    removePackageFromSim = async (
        req: RemovePackageFormSimInputDto,
        headers: THeaders
    ): Promise<RemovePackageFromSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    getPackagesForSim = async (
        req: GetPackagesForSimInputDto,
        headers: THeaders
    ): Promise<GetPackagesForSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: {
                simId: req.simId,
            },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    getSimPoolStats = async (
        type: string,
        headers: THeaders
    ): Promise<SimPoolStatsDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.SUBSCRIBER_SIMPOOL_API_URL}/stats/${type}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
    setActivePackageForSim = async (
        req: SetActivePackageForSimInputDto,
        headers: THeaders
    ): Promise<SetActivePackageForSimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        return res;
    };
}
