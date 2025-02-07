/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class MemberAPIDto {
  @Field()
  user_id: string;

  @Field()
  org_id: string;

  @Field()
  member_id: string;

  @Field()
  member_since: string;

  @Field()
  is_deactivated: boolean;

  @Field()
  role: string;
}

@ObjectType()
export class MembersAPIResDto {
  @Field(() => [MemberAPIDto])
  members: MemberAPIDto[];
}

@ObjectType()
export class MemberAPIResDto {
  @Field(() => MemberAPIDto)
  member: MemberAPIDto;
}

@ObjectType()
export class MemberDto {
  @Field()
  userId: string;

  @Field()
  name?: string;

  @Field()
  email?: string;

  @Field()
  memberId: string;

  @Field()
  isDeactivated: boolean;

  @Field()
  role: string;

  @Field({ nullable: true })
  memberSince: string;
}

@ObjectType()
export class MembersResDto {
  @Field(() => [MemberDto])
  members: MemberDto[];
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
  isDeactivated: boolean;

  @Field()
  role: string;
}

@InputType()
export class MemberInputDto {
  @Field()
  memberId: boolean;
}

@ObjectType()
export class UserResDto {
  @Field()
  name: string;

  @Field()
  email: string;

  @Field()
  uuid: string;

  @Field()
  phone: string;

  @Field()
  isDeactivated: boolean;

  @Field()
  authId: string;

  @Field()
  registeredSince: string;
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
