import {
    MemberObj,
    OrgAPIResDto,
    OrgDto,
    OrgMemberAPIResDto,
    OrgMembersAPIResDto,
    OrgMembersResDto,
    OrgsAPIResDto,
    OrgsResDto,
} from "./types";

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface IOrgService {}

export interface IOrgMapper {
    dtoToMembersResDto(res: OrgMembersAPIResDto): OrgMembersResDto;
    dtoToMemberResDto(res: OrgMemberAPIResDto): MemberObj;
    dtoToOrgsResDto(res: OrgsAPIResDto): OrgsResDto;
    dtoToOrgResDto(res: OrgAPIResDto): OrgDto;
}
