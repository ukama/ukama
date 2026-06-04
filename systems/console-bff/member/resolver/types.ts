/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { IsBoolean, IsNotEmpty, IsUUID } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";

// User types are owned by the user module. Re-exported here so member's
// datasource/mapper keep their existing imports while these types are defined
// exactly once (avoids a duplicate-type clash when the schema is consolidated).
export {
  UserAPIObj,
  UserAPIResDto,
  UserResDto,
} from "../../user/resolver/types";

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
  @IsUUID()
  userId: string;

  @Field()
  @IsNotEmpty()
  role: string;
}

@InputType()
export class UpdateMemberInputDto {
  @Field()
  @IsBoolean()
  isDeactivated: boolean;

  @Field()
  @IsNotEmpty()
  role: string;
}

// (Removed dead MemberInputDto — unused, and its memberId field was
// mistyped as boolean.)

// UserResDto / UserAPIObj / UserAPIResDto now come from the user module
// (re-exported above) — definitions removed to avoid duplicate GraphQL types.
