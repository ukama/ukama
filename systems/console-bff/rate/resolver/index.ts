/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NonEmptyArray } from "type-graphql";

import { DefaultMarkupResolver } from "./defaultMarkup";
import { GetDefaultMarkupResolver } from "./getDefaultMarkup";
import { GetDefaultMarkupHistoryResolver } from "./getDefaultMarkupHistory";

const resolvers: NonEmptyArray<any> = [
  DefaultMarkupResolver,
  GetDefaultMarkupResolver,
  GetDefaultMarkupHistoryResolver,
];

export default resolvers;
