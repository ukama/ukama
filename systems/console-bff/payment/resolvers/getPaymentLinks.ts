/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import axios from "axios";
import { randomUUID } from "crypto";
import { Arg, Query, Resolver } from "type-graphql";

import { PAYMENT_ACCESS_TOKEN, PAYMENT_BASE_URL } from "../../common/configs";
import { logger } from "../../common/logger";
import { PaymentLinks } from "./types";

@Resolver()
export class GetPaymentLinks {
  @Query(() => PaymentLinks)
  async getPaymentLinks(
    @Arg("redirectUrl") redirectUrl: string
  ): Promise<PaymentLinks> {
    const redirectURLs: any = [];
    const data = JSON.stringify({
      depositId: randomUUID(),
      returnUrl: redirectUrl,
    });

    const config = {
      method: "post",
      maxBodyLength: Infinity,
      url: PAYMENT_BASE_URL,
      headers: {
        Authorization: `Bearer ${PAYMENT_ACCESS_TOKEN}`,
        "Content-Type": "application/json",
      },
      data: data,
    };

    await axios
      .request(config)
      .then(response => {
        redirectURLs.push({
          title: "Mobile money",
          type: "mobile_money",
          link: response.data.redirectUrl,
        });
      })
      .catch(error => {
        logger.error(error);
      });
    logger.info(redirectURLs);
    return { links: redirectURLs };
  }
}
