import { Field, InputType, ObjectType } from "type-graphql";
import { MEMBER_ROLES } from "../../constants";
import { UserResDto } from "../user/types";

@ObjectType()
export class MemberAPIObj {
    @Field()
    uuid: string;

    @Field()
    user_id: string;

    @Field()
    org_id: string;

    @Field()
    member_since: string;

    @Field()
    is_deactivated: boolean;

    @Field()
    role: string;
}

@ObjectType()
export class OrgMembersAPIResDto {
    @Field()
    org: string;

    @Field(() => [MemberAPIObj])
    members: MemberAPIObj[];
}

@ObjectType()
export class OrgMemberAPIResDto {
    @Field(() => MemberAPIObj)
    member: MemberAPIObj;
}

@ObjectType()
export class MemberObj {
    @Field()
    uuid: string;

    @Field()
    userId: string;

    @Field()
    orgId: string;

    @Field()
    isDeactivated: boolean;

    @Field()
    role: MEMBER_ROLES;

    @Field({ nullable: true })
    memberSince: string;

    @Field(() => UserResDto)
    user?: UserResDto;
}

@ObjectType()
export class OrgMembersResDto {
    @Field()
    org: string;

    @Field(() => [MemberObj])
    members: MemberObj[];
}

@ObjectType()
export class OrgAPIDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    owner: string;

    @Field()
    certificate: boolean;

    @Field()
    is_deactivated: boolean;

    @Field()
    created_at: string;
}

@ObjectType()
export class OrgsAPIResDto {
    @Field()
    owner: string;

    @Field(() => [OrgAPIDto])
    orgs: OrgAPIDto[];
}

@ObjectType()
export class OrgAPIResDto {
    @Field(() => OrgAPIDto)
    org: OrgAPIDto;
}

@ObjectType()
export class OrgDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    owner: string;

    @Field()
    certificate: boolean;

    @Field()
    isDeactivated: boolean;

    @Field()
    createdAt: string;
}

@ObjectType()
export class OrgsResDto {
    @Field()
    owner: string;

    @Field(() => [OrgDto])
    orgs: OrgDto[];
}

@ObjectType()
export class OrgResDto {
    @Field(() => OrgDto)
    org: OrgDto;
}

@InputType()
export class AddOrgInputDto {
    @Field()
    certificate: string;

    @Field()
    name: string;

    @Field()
    owner_uuid: string;
}

@InputType()
export class UpdateMemberInputDto {
    @Field()
    isDeactivated: boolean;
}
