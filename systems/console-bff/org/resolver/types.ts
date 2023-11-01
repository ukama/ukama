import { Field, ObjectType } from "type-graphql";

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
