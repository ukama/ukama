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
  AddPaymentInputDto,
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
  add = async (
    baseURL: string,
    req: AddPaymentInputDto
  ): Promise<PaymentDto> => {
    logger.info(`[POST] AddPayment: ${baseURL}/${VERSION}/payments`);
    logger.info(
      `[POST] AddPayment Payload: ${JSON.stringify({
        amount: req.amount,
        item_id: req.itemId,
        country: req.country,
        currency: req.currency,
        item_type: req.itemType,
        payer_email: req.payerEmail,
        payer_phone: req.payerPhone,
        description: req.description,
        payment_method: req.paymentMethod,
        metadata: {
          simId: req.simId,
          subscriberId: req.subscriberId,
        },
      })}`
    );
    this.baseURL = baseURL;

    return this.post(`/${VERSION}/payments`, {
      body: {
        amount: req.amount,
        item_id: req.itemId,
        country: req.country,
        currency: req.currency,
        item_type: req.itemType,
        payer_email: req.payerEmail,
        payer_phone: req.payerPhone,
        description: req.description,
        payment_method: req.paymentMethod,
        metadata: {
          simId: req.simId,
          targetId: req.subscriberId,
        },
      },
    }).then(res => dtoToPaymentDto(res));
  };

  update: any = async (
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
        deposited_amount: req.amount,
        description: req.description,
        paid_at: req.paidAt,
        payer_email: req.payerEmail,
        payer_name: req.payerName,
        payer_phone: req.payerPhone,
        payment_method: req.paymentMethod,
        status: req.status,
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

  processPayments = async (
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
