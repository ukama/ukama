import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { IdResponse, ParsedCookie } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError, HTTP404Error, Messages } from "../../errors";
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
        cookie: ParsedCookie,
    ): Promise<PackageDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.DATA_PLAN_PACKAGES_API_URL,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return PackageMapper.dtoToPackageDto(res);
    };
    getPackages = async (cookie: ParsedCookie): Promise<PackagesResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.DATA_PLAN_PACKAGES_API_URL}/org/${cookie.orgId}`,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return PackageMapper.dtoToPackagesDto(res);
    };
    addPackage = async (
        req: AddPackageInputDto,
        cookie: ParsedCookie,
    ): Promise<PackageDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.DATA_PLAN_PACKAGES_API_URL,
            body: req,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return PackageMapper.dtoToPackageDto(res);
    };
    deletePackage = async (
        packageId: string,
        cookie: ParsedCookie,
    ): Promise<IdResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.DATA_PLAN_PACKAGES_API_URL}/${packageId}`,
            headers: cookie.header,
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
        cookie: ParsedCookie,
    ): Promise<PackageDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PATCH,
            path: `${SERVER.DATA_PLAN_PACKAGES_API_URL}/${packageId}`,
            body: req,
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return PackageMapper.dtoToPackageDto(res);
    };
}
