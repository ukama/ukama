import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { BoolResponse, ParsedCookie } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError, HTTP404Error, Messages } from "../../errors";
import { IOrgService } from "./interface";
import OrgMapper from "./mapper";
import {
    AddOrgInputDto,
    MemberObj,
    OrgDto,
    OrgMembersResDto,
    OrgsResDto,
    UpdateMemberInputDto,
} from "./types";

@Service()
export class OrgService implements IOrgService {
    getOrgMembers = async (cookie: ParsedCookie): Promise<OrgMembersResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${cookie.orgName}/members`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToMembersResDto(res);
    };
    getOrgMember = async (cookie: ParsedCookie): Promise<MemberObj> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${cookie.orgName}/members/${cookie.userId}`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToMemberResDto(res);
    };
    removeMember = async (cookie: ParsedCookie): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.DELETE,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${cookie.orgName}/members/${cookie.userId}`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            success: true,
        };
    };
    getOrgs = async (cookie: ParsedCookie): Promise<OrgsResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${cookie.orgName}`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToOrgsResDto(res);
    };
    getOrg = async (orgName: string, cookie: ParsedCookie): Promise<OrgDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${orgName}`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToOrgResDto(res);
    };
    addOrg = async (
        req: AddOrgInputDto,
        cookie: ParsedCookie,
    ): Promise<OrgDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: SERVER.REGISTRY_ORGS_API_URL,
            body: req,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToOrgResDto(res);
    };
    addMember = async (
        userId: string,
        cookie: ParsedCookie,
    ): Promise<MemberObj> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${cookie.orgName}/members`,
            body: { user_uuid: userId },
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return OrgMapper.dtoToMemberResDto(res);
    };
    updateMember = async (
        memberId: string,
        req: UpdateMemberInputDto,
        cookie: ParsedCookie,
    ): Promise<BoolResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.REGISTRY_ORGS_API_URL}/${cookie.orgName}/members/${memberId}`,
            body: req,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        if (!res) throw new HTTP404Error(Messages.NODES_NOT_FOUND);
        return {
            success: true,
        };
    };
}
