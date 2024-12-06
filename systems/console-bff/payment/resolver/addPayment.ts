/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Mutation, Resolver } from "type-graphql";

import { logger } from "../../common/logger";
import { Context } from "../context";
import { AddPaymentInputDto, PaymentDto } from "./types";

@Resolver()
export class AddPaymentResolver {
  @Mutation(() => PaymentDto)
  async addPayment(
    @Arg("data") data: AddPaymentInputDto,
    @Ctx() ctx: Context
  ): Promise<PaymentDto> {
    const { dataSources, baseURL } = ctx;
    let payment: PaymentDto | undefined = undefined;
    logger.info(`Adding payment: ${baseURL}`);
    try {
      const paymentRes = await dataSources.dataSource.add(baseURL, data);
      if (paymentRes && paymentRes.id) {
        try {
          logger.info("Payment added successfully");
          const tokenRes = await dataSources.dataSource.getToken(
            baseURL,
            paymentRes.id
          );
          if (tokenRes && tokenRes.token) {
            logger.info("Token Res");
            try {
              const processPayment =
                await dataSources.dataSource.processPayments(baseURL, {
                  id: paymentRes.id,
                  token: tokenRes.token,
                  correspondent: data.correspondent || "",
                });
              payment = processPayment.payment;
              logger.info("Process Payment");
            } catch (error) {
              logger.error("Error processing payment:", error);
              throw new Error("Failed to process payment");
            }
          } else {
            throw new Error("Failed to retrieve token");
          }
        } catch (error) {
          logger.error("Error getting token:", error);
          throw new Error("Failed to get token");
        }
      } else {
        logger.error("Failed to add payment");
        throw new Error("Failed to add payment");
      }
    } catch (error) {
      logger.error("Error adding payment:", error);
      throw new Error("Failed to add payment");
    }

    return payment;
  }
}
