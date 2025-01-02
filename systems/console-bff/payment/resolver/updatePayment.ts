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
import { PaymentDto, UpdatePaymentInputDto } from "./types";

@Resolver()
export class UpdatePaymentResolver {
  @Mutation(() => PaymentDto)
  async updatePayment(
    @Arg("data") data: UpdatePaymentInputDto,
    @Ctx() ctx: Context
  ): Promise<PaymentDto> {
    const { dataSources, baseURL } = ctx;
    let payment: PaymentDto | undefined = undefined;

    logger.info(`Updating payment for bill: ${data.id}`);

    try {
      const updatedPaymentRes = await dataSources.dataSource.updatePayment(
        baseURL,
        data
      );

      if (updatedPaymentRes && updatedPaymentRes.id) {
        logger.info("Payment updated successfully");

        try {
          const tokenRes = await dataSources.dataSource.getToken(
            baseURL,
            updatedPaymentRes.id
          );

          if (tokenRes && tokenRes.token) {
            logger.info("Token retrieved successfully");

            try {
              const processPaymentRes =
                await dataSources.dataSource.processPayment(baseURL, {
                  id: updatedPaymentRes.id,
                  token: tokenRes.token,
                  correspondent: "test",
                });

              if (processPaymentRes && processPaymentRes.payment) {
                payment = processPaymentRes.payment;
                logger.info("Payment processed successfully");
              } else {
                logger.error("Failed to process payment: No payment returned");
                throw new Error("Failed to process payment");
              }
            } catch (processError) {
              logger.error("Error processing payment:", processError);
              throw new Error("Failed to process payment");
            }
          } else {
            logger.error("Failed to retrieve token");
            throw new Error("Failed to retrieve token");
          }
        } catch (tokenError) {
          logger.error("Error getting token:", tokenError);
          throw new Error("Failed to get token");
        }
      } else {
        logger.error("Failed to update payment");
        throw new Error("Failed to update payment");
      }
    } catch (updateError) {
      logger.error("Error updating payment:", updateError);
      throw new Error("Failed to update payment");
    }

    return payment;
  }
}
