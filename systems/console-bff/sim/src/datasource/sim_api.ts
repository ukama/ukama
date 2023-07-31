import { RESTDataSource } from "@apollo/datasource-rest";
import { THeaders } from "../../../common/types";
import { SERVER } from "../../../constants/endpoints";
import { getHeaders } from "../../../utils";
import generateTokenFromIccid from "../../../utils/generateSimToken";
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
} from "../types";

export class SimApi extends RESTDataSource {
    uploadSims = async (
        req: UploadSimsInputDto,
        headers: THeaders
    ): Promise<UploadSimsResDto> => {
        return this.put(`${SERVER.SUBSCRIBER_SIMPOOL_API_URL}/upload`, {
            headers: getHeaders(headers),
            body: {
                data: req.data,
                sim_type: req.simType,
            },
          }).then(res => res);
    };

    allocateSim = async (
        req: AllocateSimInputDto,
        headers: THeaders
    ): Promise<SimDto> => {
        const token = generateTokenFromIccid(
            req.iccid,
            process.env.ENCRYPTION_KEY || ""
        );
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
            headers: getHeaders(headers),
            body: {
                ...req,
                sim_token: token,
            },
          }).then(res => SimMapper.dtoToSimResDto(res));
    };

    toggleSimStatus = async (
        req: ToggleSimStatusInputDto,
        headers: THeaders
    ): Promise<SimStatusResDto> => {
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`, {
            headers: getHeaders(headers),
            body: {
                simId: req.simId,
                status: req.status,
            },
          }).then(res => res);
    };

    getSim = async (
        req: GetSimInputDto,
        headers: THeaders
    ): Promise<SimDto> => {
        //TODO: body is in get request

        // const res = await catchAsyncIOMethod({
        //     type: API_METHOD_TYPE.GET,
        //     path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
        //     body: {
        //         simId: req.simId,
        //     },
        //     headers: getHeaders(headers),
        // });
        // if (checkError(res)) throw new Error(res.message);
        // return SimMapper.dtoToSimResDto(res);

        return this.get(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,{
            headers: getHeaders(headers),
            params:{
                simId: req.simId,
            }
        }).then(res => SimMapper.dtoToSimResDto(res));
    };

    getSims = async (type: string, headers: THeaders): Promise<SimsResDto> => {
        return this.get(`${SERVER.SUBSCRIBER_SIMPOOL_API_URL}/sims/${type}`,{
            headers: getHeaders(headers)
        }).then(res => SimMapper.dtoToSimsDto(res));
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
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,{
            headers: getHeaders(headers),
            body: {
                subscriberId: req.subscriberId,
            },
        }).then(res => SimMapper.dtoToSimDetailsDto(res));
    };

    getSimByNetworkId = async (
        req: GetSimByNetworkInputDto,
        headers: THeaders
    ): Promise<SimDetailsDto> => {
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,{
            headers: getHeaders(headers),
            body: {
                networkId: req.networkId,
            },
        }).then(res => SimMapper.dtoToSimDetailsDto(res));
    };

    deleteSim = async (
        req: DeleteSimInputDto,
        headers: THeaders
    ): Promise<DeleteSimResDto> => {
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,{
            headers: getHeaders(headers),
            body: {
                simId: req.simId,
            },
        }).then(res => res);
    };

    addPackegeToSim = async (
        req: AddPackageToSimInputDto,
        headers: THeaders
    ): Promise<AddPackageSimResDto> => {
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,{
            headers: getHeaders(headers),
            body: {
                ...req
            },
        }).then(res => res);
    };

    removePackageFromSim = async (
        req: RemovePackageFormSimInputDto,
        headers: THeaders
    ): Promise<RemovePackageFromSimResDto> => {
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,{
            headers: getHeaders(headers),
            body: {
                ...req
            },
        }).then(res => res);
    };

    getPackagesForSim = async (
        req: GetPackagesForSimInputDto,
        headers: THeaders
    ): Promise<GetPackagesForSimResDto> => {
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,{
            headers: getHeaders(headers),
            body: {
                simId: req.simId,
            },
        }).then(res => res);
    };

    getSimPoolStats = async (
        type: string,
        headers: THeaders
    ): Promise<SimPoolStatsDto> => {
        return this.get(`${SERVER.SUBSCRIBER_SIMPOOL_API_URL}/stats/${type}`,{
            headers: getHeaders(headers),
        }).then(res => res);
    };

    setActivePackageForSim = async (
        req: SetActivePackageForSimInputDto,
        headers: THeaders
    ): Promise<SetActivePackageForSimResDto> => {
        return this.put(`${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,{
            headers: getHeaders(headers),
            body: {
                ...req
            },
        }).then(res => res);
    };
}
