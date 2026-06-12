/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Shared home KPI query. Both the Business and Network home screens render a
 * KPI strip parameterised by `lens`; the BFF routes each lens to its analytics
 * overview endpoint (business/home | network/overview) and returns its `kpis`.
 * Home sites come live from the registry (`sitesView`), not analytics, so
 * they're not modelled here.
 */
import { Field, InputType, ObjectType, registerEnumType } from "type-graphql";

import { KpiDto } from "./shared";

export enum HomeLens {
  BUSINESS = "business",
  NETWORK = "network",
}
registerEnumType(HomeLens, {
  name: "HomeLens",
  description: "Which lens's home data to fetch.",
});

@InputType()
export class HomeViewInput {
  @Field(() => HomeLens)
  lens: HomeLens;

  @Field({ nullable: true })
  networkId?: string;

  @Field({ nullable: true })
  period?: string;

  @Field({ nullable: true })
  from?: string;

  @Field({ nullable: true })
  to?: string;

  @Field({ nullable: true })
  timezone?: string;
}

@ObjectType()
export class HomeKpis {
  @Field(() => [KpiDto])
  kpis: KpiDto[];
}
