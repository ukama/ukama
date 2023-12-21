/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { randomUUID } from "crypto";
import fs from "fs";
import jwt from "jsonwebtoken";
import { Arg, Query, Resolver } from "type-graphql";

import { PaymentLinks, PaymentLinksInput } from "./types";

@Resolver()
export class GetPaymentLinks {
  @Query(() => PaymentLinks)
  async getPaymentLinks(
    @Arg("data") data: PaymentLinksInput
  ): Promise<PaymentLinks> {
    const privateKey = fs.readFileSync("./private_key.pem");
    const token = jwt.sign(
      {
        data: {
          id: randomUUID(),
          amount: data.amount,
          msisdn: data.msisdn,
          country: data.country,
          currency: "GHS",
          reason: data.reason,
        },
        sub: "payment",
        iat: Math.floor(new Date().getTime() / 1000),
        nbf: Math.floor(new Date().getTime() / 1000),
        exp: 90000,
        token: "session_token",
      },
      privateKey,
      { algorithm: "RS256" }
    );

    const redirectURLs: any = [];
    redirectURLs.push({
      title: "Mobile money",
      type: "mobile_money",
      link: `http://localhost:3000/mobile_money?token=${token}`,
    });
    redirectURLs.push({
      title: "Stripe",
      type: "stripe_payment",
      link: `http://localhost:3000/stripe_payment?token=${token}`,
    });
    redirectURLs.push({
      title: "Cash",
      type: "cash",
      link: ``,
    });

    return { links: redirectURLs };
  }
}
