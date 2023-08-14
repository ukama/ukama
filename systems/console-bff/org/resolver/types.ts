import { Field, InputType, ObjectType } from "type-graphql";

import { UserResDto } from "../../user/resolver/types";

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
  role: string;

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
  certificate: string;

  @Field()
  is_deactivated: boolean;

  @Field()
  created_at: string;
}

@ObjectType()
export class OrgsAPIResDto {
  @Field()
  user: string;

  @Field(() => [OrgAPIDto])
  owner_of: OrgAPIDto[];

  @Field(() => [OrgAPIDto])
  member_of: OrgAPIDto[];
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
  certificate: string;

  @Field()
  isDeactivated: boolean;

  @Field()
  createdAt: string;
}

@ObjectType()
export class OrgsResDto {
  @Field()
  user: string;

  @Field(() => [OrgDto])
  ownerOf: OrgDto[];

  @Field(() => [OrgDto])
  memberOf: OrgDto[];
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
export class AddMemberInputDto {
  @Field()
  userId: string;

  @Field()
  role: string;
}

@InputType()
export class UpdateMemberInputDto {
  @Field()
  orgName: string;

  @Field()
  isDeactivated: boolean;

  @Field()
  role: string;
}

@InputType()
export class MemberInputDto {
  @Field()
  memberId: boolean;

  @Field()
  orgName: boolean;
}

@ObjectType()
export class UserAPIObj {
  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  id: string;

  @Field()
  phone: string;

  @Field()
  is_deactivated: boolean;

  @Field()
  auth_id: string;

  @Field()
  registered_since: string;
}

@ObjectType()
export class UserAPIResDto {
  @Field(() => [UserAPIObj])
  user: UserAPIObj;
}
