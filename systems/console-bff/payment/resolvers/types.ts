/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class LinksDTO {
  @Field()
  link: string;

  @Field()
  type: string;

  @Field()
  title: string;
}

@ObjectType()
export class PaymentLinks {
  @Field(() => [LinksDTO])
  links: LinksDTO[];
}

@InputType()
export class PaymentLinksInput {
  @Field()
  amount: number;

  @Field()
  msisdn: string;

  @Field()
  country: string;

  @Field()
  currency: string;

  @Field()
  reason: string;

  @Field()
  redirectUrl: string;
}

@ObjectType()
export class TokenParserDto {
  @Field()
  id: string;

  @Field()
  orgName: string;

  @Field()
  for: string;

  @Field()
  countryCode: string;

  @Field()
  phoneNumber: string;

  @Field()
  amount: number;

  @Field()
  currency: string;
}
