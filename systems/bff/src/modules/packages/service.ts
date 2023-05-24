import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { IdResponse, THeaders } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { HTTP404Error, Messages, checkError } from "../../errors";
import { getHeaders } from "../../utils";
import { IPackageService } from "./interface";
import PackageMapper from "./mapper";
import {
    AddPackageInputDto,
    PackageDto,
    PackagesResDto,
    UpdatePackageInputDto,
} from "./types";
@Service()
export class PackageService implements IPackageService {
    getPackage = async (
        packageId: string,
        headers: THeaders
    ): Promise<PackageDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.DATA_PLAN_PACKAGES_API_URL}/${packageId}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return PackageMapper.dtoToPackageDto(res);
    };
    getPackages = async (headers: THeaders): Promise<PackagesResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.DATA_PLAN_PACKAGES_API_URL}/org/${headers.orgId}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return PackageMapper.dtoToPackagesDto(res);
    };
    addPackage = async (
        req: AddPackageInputDto,
        headers: THeaders
    ): Promise<PackageDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.DATA_PLAN_PACKAGES_API_URL,
            body: req,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return PackageMapper.dtoToPackageDto(res);
    };
    deletePackage = async (
        packageId: string,
        headers: THeaders
    ): Promise<IdResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.DATA_PLAN_PACKAGES_API_URL}/${packageId}`,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            uuid: packageId,
        };
    };
    updatePackage = async (
        packageId: string,
        req: UpdatePackageInputDto,
        headers: THeaders
    ): Promise<PackageDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.DATA_PLAN_PACKAGES_API_URL}/${packageId}`,
            body: req,
            headers: getHeaders(headers),
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return PackageMapper.dtoToPackageDto(res);
    };
}
