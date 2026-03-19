/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Field, InputType, ObjectType } from "type-graphql";

import { SOFTWARE_STATUS } from "../../common/enums";

@ObjectType()
export class App {
  @Field()
  name: string;

  @Field()
  space: string;

  @Field()
  notes: string;

  @Field(() => [String])
  metricsKeys: string[];
}

@ObjectType()
export class Apps {
  @Field(() => [App])
  apps: App[];
}

@InputType()
export class GetSoftwaresInput {
  @Field()
  name: string;

  @Field()
  nodeId: string;

  @Field(() => SOFTWARE_STATUS)
  status: SOFTWARE_STATUS;
}

@ObjectType()
export class Softwares {
  @Field(() => [Software])
  software: Software[];
}

@ObjectType()
export class Software {
  @Field()
  id: string;

  @Field()
  releaseDate: string;

  @Field()
  nodeId: string;

  @Field(() => SOFTWARE_STATUS)
  status: SOFTWARE_STATUS;

  @Field(() => [String])
  changeLog: string[];

  @Field()
  currentVersion: string;

  @Field()
  desiredVersion: string;

  @Field()
  name: string;

  @Field()
  space: string;

  @Field()
  notes: string;

  @Field(() => [String])
  metricsKeys: string[];

  @Field()
  createdAt: string;

  @Field()
  updatedAt: string;
}

@InputType()
export class UpdateSoftwareInputDto {
  @Field()
  name: string;

  @Field()
  nodeId: string;

  @Field()
  tag: string;
}

@ObjectType()
export class StringResponse {
  @Field()
  message: string;
}
