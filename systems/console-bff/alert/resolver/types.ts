/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, ObjectType } from "type-graphql";

import { ALERT_TYPE } from "../../common/enums";
import { PaginationResponse } from "../../common/types";

@ObjectType()
export class AlertDto {
  @Field({ nullable: true })
  id: string;

  @Field(() => ALERT_TYPE)
  type: ALERT_TYPE;

  @Field({ nullable: true })
  title: string;

  @Field({ nullable: true })
  description: string;

  @Field({ nullable: true })
  alertDate: Date;
}

@ObjectType()
export class AlertsResponse extends PaginationResponse {
  @Field(() => [AlertDto])
  alerts: AlertDto[];
}

@ObjectType()
export class AlertResponse {
  @Field()
  status: string;

  @Field(() => [AlertDto])
  data: AlertDto[];

  @Field()
  length: number;
}
