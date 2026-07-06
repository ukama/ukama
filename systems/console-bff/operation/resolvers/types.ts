/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Field, Float, InputType, ObjectType } from "type-graphql";

import { NODE_TYPE } from "../../common/enums";

/**
 * An operation row from the operation system. `status` is the operation
 * lifecycle state as returned by the gateway (PENDING/RUNNING/SUCCESS/
 * FAILED/TIMEOUT/CANCELLED) and `fencingToken` is the monotonic guard token
 * for the operation's lease.
 */
@ObjectType()
export class OperationDto {
  @Field()
  id: string;

  @Field()
  type: string;

  @Field()
  system: string;

  @Field()
  status: string;

  @Field(() => Float)
  fencingToken: number;

  @Field({ nullable: true })
  requestedBy?: string;

  @Field({ nullable: true })
  idempotencyKey?: string;

  @Field()
  resourceKey: string;

  @Field({ nullable: true })
  leaseExpiresAt?: string;

  @Field({ nullable: true })
  error?: string;

  @Field({ nullable: true })
  startedAt?: string;

  @Field({ nullable: true })
  terminalAt?: string;

  @Field({ nullable: true })
  createdAt?: string;
}

@ObjectType()
export class ResourceLockDto {
  @Field()
  locked: boolean;

  @Field(() => OperationDto, { nullable: true })
  operation?: OperationDto;
}

@InputType()
export class GetOperationInputDto {
  @Field()
  id: string;
}

@InputType()
export class GetResourceLockInputDto {
  @Field()
  resourceKey: string;
}

/**
 * Aggregated operation status for the console's action UI. The operation
 * system locks are always keyed `node:<NODE_ID>`; these view types let the
 * console ask a single question per page ("is this node / site busy, and
 * which actions are available?") instead of reasoning about per-node locks
 * itself. Computed statelessly per request (see ../logic.ts).
 */
@ObjectType()
export class ActionAvailabilityDto {
  @Field()
  available: boolean;

  /** Human-readable reason when `available` is false (for tooltips). */
  @Field({ nullable: true })
  reason?: string;
}

@ObjectType()
export class NodeOperationStatusDto {
  @Field()
  nodeId: string;

  @Field(() => NODE_TYPE, { nullable: true })
  type?: NODE_TYPE;

  /** True while a non-terminal operation holds an unexpired lease. */
  @Field()
  busy: boolean;

  /** The active operation, if any (null when idle). */
  @Field(() => OperationDto, { nullable: true })
  operation?: OperationDto;
}

/**
 * Per-action availability for the site page. Each action depends on a
 * different physical node (RF → amplifier/anode, service → tower/tnode,
 * restartSite → every node), so they clear independently.
 */
@ObjectType()
export class SiteActionsDto {
  @Field(() => ActionAvailabilityDto)
  restartSite: ActionAvailabilityDto;

  @Field(() => ActionAvailabilityDto)
  rf: ActionAvailabilityDto;

  @Field(() => ActionAvailabilityDto)
  service: ActionAvailabilityDto;
}

@ObjectType()
export class SiteOperationStatusDto {
  @Field()
  siteId: string;

  /** True if any of the site's nodes is busy. */
  @Field()
  busy: boolean;

  /** True if a per-node lock read failed (fail-open — actions still allowed). */
  @Field()
  degraded: boolean;

  @Field(() => [NodeOperationStatusDto])
  nodes: NodeOperationStatusDto[];

  @Field(() => SiteActionsDto)
  actions: SiteActionsDto;
}

@InputType()
export class NodeOperationStatusInputDto {
  @Field()
  nodeId: string;
}

@InputType()
export class SiteOperationStatusInputDto {
  @Field()
  siteId: string;
}
