/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Query, Resolver } from "type-graphql";

import CURRENCIES from "../../common/data/currencies";
import { CurrencyRes } from "./types";

@Resolver()
export class GetCurrencySymbolResolver {
  @Query(() => CurrencyRes)
  async getCurrencySymbol(@Arg("code") code: string): Promise<CurrencyRes> {
    const currency = CURRENCIES.find(
      currency => currency.Code === code.toUpperCase()
    );
    return {
      code: currency?.Code ?? "",
      symbol: currency?.Symbol ?? "",
      image: currency?.SymbolImage ?? "",
    };
  }
}
