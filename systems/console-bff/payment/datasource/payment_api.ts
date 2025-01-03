/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import { logger } from "../../common/logger";
import {
  CorrespondentsResDto,
  GetPaymentsInputDto,
  PaymentDto,
  PaymentsDto,
  ProcessPaymentDto,
  ProcessPaymentInputDto,
  TokenResDto,
  UpdatePaymentInputDto,
} from "../resolver/types";
import {
  dtoToCorspondantsDto,
  dtoToPaymentDto,
  dtoToPaymentsDto,
  dtoToProcessPaymentsDto,
} from "./mapper";

class PaymentAPI extends RESTDataSource {
  updatePayment: any = async (
    baseURL: string,
    req: UpdatePaymentInputDto
  ): Promise<PaymentDto> => {
    logger.info(
      `[PUT] UpdatePayment: ${baseURL}/${VERSION}/payments/${req.id}`
    );
    this.baseURL = baseURL;
    return this.put(`/${VERSION}/payments/${req.id}`, {
      body: {
        country: req.country,
        currency: req.currency,
        payer_email: req.payerEmail,
        payer_name: req.payerName,
        payment_method: req.paymentMethod,
      },
    }).then(res => dtoToPaymentDto(res));
  };

  getPayment = async (
    baseURL: string,
    paymentId: string
  ): Promise<PaymentDto> => {
    logger.info(
      `[GET] GetPayment: ${baseURL}/${VERSION}/payments/${paymentId}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/payments/${paymentId}`, {}).then(res =>
      dtoToPaymentDto(res)
    );
  };

  getToken = async (
    baseURL: string,
    paymentId: string
  ): Promise<TokenResDto> => {
    logger.info(`[GET] GetToken: ${baseURL}/${VERSION}/tokens/${paymentId}`);
    this.baseURL = baseURL;

    return this.get(`/${VERSION}/tokens/${paymentId}`, {}).then(res => res);
  };

  getPayments = async (
    baseURL: string,
    data: GetPaymentsInputDto
  ): Promise<PaymentsDto> => {
    this.baseURL = baseURL;
    let params = "sort=true";
    if (data.paymentMethod) {
      params = params + `&payment_method=${data.paymentMethod}`;
    }
    if (data.type) {
      params = params + `&item_type=${data.type}`;
    }
    if (data.status) {
      params = params + `&status=${data.status}`;
    }
    logger.info(`[GET] GetPayments: ${baseURL}/${VERSION}/payments?${params}`);
    return this.get(`/${VERSION}/payments?${params}`).then(res =>
      dtoToPaymentsDto(res)
    );
  };

  processPayment = async (
    baseURL: string,
    req: ProcessPaymentInputDto
  ): Promise<ProcessPaymentDto> => {
    logger.info(
      `[PATCH] ProcessPayments: ${this.baseURL}/${VERSION}/payments/${req.id}`
    );
    this.baseURL = baseURL;

    return this.patch(`/${VERSION}/payments/${req.id}`, {
      body: {
        correspondent: req.correspondent,
        token: req.token,
      },
    }).then(res => dtoToProcessPaymentsDto(res));
  };

  getCorrespondents = async (
    baseURL: string,
    phoneNumber: string,
    paymentMethod: string
  ): Promise<CorrespondentsResDto> => {
    logger.info(
      `[GET] GetCorrespondents: ${baseURL}/${VERSION}/correspondents/${phoneNumber}?payment_method=${paymentMethod}`
    );
    this.baseURL = baseURL;

    return this.get(
      `/${VERSION}/correspondents/${phoneNumber}?payment_method=${paymentMethod}`
    ).then(res => dtoToCorspondantsDto(res));
  };
}

export default PaymentAPI;
