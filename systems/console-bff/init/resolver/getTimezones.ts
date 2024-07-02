/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Query, Resolver } from "type-graphql";

import TimeZones from "../../common/data/timezones";
import { TimezoneDto, TimezoneRes } from "./types";

@Resolver()
export class GetTimezonesResolver {
  @Query(() => TimezoneRes)
  async getTimezones(): Promise<TimezoneRes> {
    const timezones: TimezoneDto[] = [];
    for (const timezone of TimeZones) {
      timezones.push({
        value: timezone.value,
        abbr: timezone.abbr,
        offset: timezone.offset,
        isdst: timezone.isdst,
        text: timezone.text,
        utc: timezone.utc,
      });
    }
    return {
      timezones: timezones,
    };
  }
}
