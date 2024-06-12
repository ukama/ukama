/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Query, Resolver } from "type-graphql";

import COUNTRIES from "../../common/data/countries";
import { CountriesRes, CountryDto } from "./types";

@Resolver()
export class GetCountriesResolver {
  @Query(() => CountriesRes)
  async getCountries(): Promise<CountriesRes> {
    const countries: CountryDto[] = [];
    for (const country of COUNTRIES) {
      countries.push({
        name: country.name,
        code: country.code,
      });
    }
    return {
      countries: countries,
    };
  }
}
