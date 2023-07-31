import { RESTDataSource } from "@apollo/datasource-rest";
import { BoolResponse, THeaders } from "../../../common/types";
import { SERVER } from "../../../constants/endpoints";
import { getHeaders } from "../../../utils";
import OrgMapper from "./mapper";
import {
    AddMemberInputDto,
    AddOrgInputDto,
    MemberObj,
    OrgDto,
    OrgMembersResDto,
    OrgsResDto,
    UpdateMemberInputDto,
} from "../types";


export class OrgApi extends RESTDataSource {
    getOrgMembers = async (headers: THeaders): Promise<OrgMembersResDto> => {
        return this.get(`${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members`,{headers: getHeaders(headers)}).then(res => 
            OrgMapper.dtoToMembersResDto(res));
    };

    getOrgMember = async (headers: THeaders): Promise<MemberObj> => {
        return this.get(`${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members/${headers.userId}`,{headers: getHeaders(headers)}).then(res => 
            OrgMapper.dtoToMemberResDto(res));
    };

    removeMember = async (headers: THeaders): Promise<BoolResponse> => {
        return this.delete(`${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members/${headers.userId}`,{headers: getHeaders(headers)}).then((res)=>{
            return {
                success: true,
            };
        });
    };

    getOrgs = async (headers: THeaders): Promise<OrgsResDto> => {
        return this.get(`${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}`,{headers: getHeaders(headers)}).then(res => 
            OrgMapper.dtoToOrgsResDto(res));
    };

    getOrg = async (orgName: string, headers: THeaders): Promise<OrgDto> => {
        return this.get(`${SERVER.REGISTRY_ORGS_API_URL}/${orgName}`,{headers: getHeaders(headers)}).then(res => 
            OrgMapper.dtoToOrgResDto(res));
    };

    addOrg = async (
        req: AddOrgInputDto,
        headers: THeaders
    ): Promise<OrgDto> => {
        return this.post(SERVER.REGISTRY_ORGS_API_URL, {
            headers: getHeaders(headers),
            body: req,
          }).then(res => OrgMapper.dtoToOrgResDto(res));
    };

    addMember = async (
        data: AddMemberInputDto,
        headers: THeaders
    ): Promise<MemberObj> => {
        return this.post(`${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members`, {
            headers: getHeaders(headers),
            body: { user_uuid: headers.userId, role: data.role },
          }).then(res => OrgMapper.dtoToMemberResDto(res));
    };
    
    updateMember = async (
        memberId: string,
        req: UpdateMemberInputDto,
        headers: THeaders
    ): Promise<BoolResponse> => {
        return this.post(`${SERVER.REGISTRY_ORGS_API_URL}/${headers.orgName}/members/${memberId}`, {
            headers: getHeaders(headers),
            body: req,
          }).then(res => {
            return {
                success: true,
          }});
    };
}
