/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Field, Float, InputType, Int, ObjectType } from "type-graphql";

/**
 * An operation row from the operation system. `status` is the operation
 * lifecycle state as returned by the gateway (PENDING/RUNNING/SUCCESS/
 * FAILED/TIMEOUT/CANCELLED) and `fencingToken` is the monotonic guard token
 * the caller must echo back on markRunning.
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
export class StartOperationResponseDto {
  @Field(() => OperationDto, { nullable: true })
  operation?: OperationDto;

  @Field(() => OperationDto, { nullable: true })
  conflictingOperation?: OperationDto;
}

@ObjectType()
export class ResourceLockDto {
  @Field()
  locked: boolean;

  @Field(() => OperationDto, { nullable: true })
  operation?: OperationDto;
}

@InputType()
export class StartOperationInputDto {
  @Field()
  type: string;

  @Field()
  system: string;

  @Field()
  resourceKey: string;

  @Field({ nullable: true })
  requestedBy?: string;

  @Field({ nullable: true })
  idempotencyKey?: string;

  @Field(() => Int, { nullable: true })
  leaseSeconds?: number;
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

@InputType()
export class MarkOperationRunningInputDto {
  @Field()
  id: string;

  @Field(() => Float)
  fencingToken: number;
}

@InputType()
export class ForceUnlockInputDto {
  @Field()
  id: string;

  @Field()
  reason: string;
}
