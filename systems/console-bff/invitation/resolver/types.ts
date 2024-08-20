/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { IsEmail } from "class-validator";
import { Field, InputType, ObjectType } from "type-graphql";

import { INVITATION_STATUS } from "../../common/enums";

@InputType()
export class CreateInvitationInputDto {
  @Field()
  name: string;

  @Field()
  @IsEmail()
  email: string;

  @Field()
  role: string;
}

@ObjectType()
export class InvitationDto {
  @Field()
  email: string;

  @Field()
  expireAt: string;

  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  role: string;

  @Field()
  link: string;

  @Field()
  userId: string;

  @Field(() => INVITATION_STATUS)
  status: INVITATION_STATUS;
}

@ObjectType()
export class InvitationsDto {
  @Field(() => [InvitationDto])
  invitations: InvitationDto[];
}

@ObjectType()
export class InvitationAPIDto {
  @Field()
  email: string;

  @Field()
  name: string;

  @Field()
  expireAt: string;

  @Field()
  userId: string;

  @Field()
  role: string;

  @Field()
  id: string;

  @Field()
  link: string;

  @Field(() => INVITATION_STATUS)
  status: INVITATION_STATUS;
}

@ObjectType()
export class InvitationAPIResDto {
  @Field(() => InvitationAPIDto)
  invitation: InvitationAPIDto;
}
@ObjectType()
export class InvitationsAPIResDto {
  @Field(() => [InvitationAPIDto])
  invitations: InvitationAPIDto[];
}

@ObjectType()
export class InvitationsResDto {
  @Field(() => [InvitationDto])
  invitations: InvitationDto[];
}

@InputType()
export class GetInvitationInputDto {
  @Field()
  id: string;
}

@InputType()
export class GetInvitationByOrgInputDto {}

@ObjectType()
export class UpdateInvitationResDto {
  @Field()
  id: string;
}

@ObjectType()
export class DeleteInvitationResDto {
  @Field()
  id: string;
}

@InputType()
export class UpateInvitationInputDto {
  @Field()
  id: string;

  @Field()
  email: string;

  @Field(() => INVITATION_STATUS)
  status: INVITATION_STATUS;
}
