import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { BoolResponse, THeaders } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { HTTP404Error, Messages, checkError } from "../../errors";
import { getHeaders } from "../../utils";
import { IOrgService } from "./interface";
import OrgMapper from "./mapper";
import {
    AddMemberInputDto,
    AddOrgInputDto,
    MemberObj,
    OrgDto,
    OrgMembersResDto,
    OrgsResDto,
    UpdateMemberInputDto,
} from "./types";

@Service()
export class OrgService implements IOrgService {
    getOrgMembers = async (headers: THeaders): Promise<OrgMembersResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members`,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToMembersResDto(res);
    };
    getOrgMember = async (headers: THeaders): Promise<MemberObj> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members/${headers.userId}`,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToMemberResDto(res);
    };
    removeMember = async (headers: THeaders): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members/${headers.userId}`,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            success: true,
        };
    };
    getOrgs = async (headers: THeaders): Promise<OrgsResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}`,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToOrgsResDto(res);
    };
    getOrg = async (orgName: string, headers: THeaders): Promise<OrgDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${orgName}`,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToOrgResDto(res);
    };
    addOrg = async (
        req: AddOrgInputDto,
        headers: THeaders
    ): Promise<OrgDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.REGISTRY_ORGS_API_URL,
            body: req,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToOrgResDto(res);
    };
    addMember = async (
        data: AddMemberInputDto,
        headers: THeaders
    ): Promise<MemberObj> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members`,
            body: { user_uuid: headers.userId, role: data.role },
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToMemberResDto(res);
    };
    updateMember = async (
        memberId: string,
        req: UpdateMemberInputDto,
        headers: THeaders
    ): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members/${memberId}`,
            body: req,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            success: true,
        };
    };
}
