/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import jwt from "jsonwebtoken";
import { Arg, Query, Resolver } from "type-graphql";

import { formatCountryCode, formatPhoneNumber } from "../../common/utils";
import { TokenParserDto } from "./types";

@Resolver()
export class TokenParserResolver {
  @Query(() => TokenParserDto)
  async tokenParser(@Arg("data") data: string): Promise<TokenParserDto> {
    const decoded: any = jwt.decode(data, { complete: true });
    const { id, amount, msisdn, currency, reason } = decoded?.payload
      .data as any;
    return {
      id: id,
      orgName: "Ukama",
      for: reason,
      countryCode: formatCountryCode(msisdn),
      phoneNumber: formatPhoneNumber(msisdn),
      amount: parseInt(amount),
      currency: currency,
    };
  }
}
