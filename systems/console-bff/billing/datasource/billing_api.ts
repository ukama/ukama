/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { BILLING_API_GW, VERSION } from "../../common/configs";
import { BillHistoryDto, BillResponse, InvoiceDto } from "../resolvers/types";
import { billHistoryDtoToDto, dtoToDto } from "./mapper";

const version = "/v1/invoices";
class BillingAPI extends RESTDataSource {
  baseURL = BILLING_API_GW + version;
  public getCurrentBill = async (): Promise<BillResponse> => {
    return this.get("/current").then(res => dtoToDto(res));
  };

  public getBillHistory = async (): Promise<BillHistoryDto[]> => {
    return this.get("/history").then(res => billHistoryDtoToDto(res));
  };

  getInvoice = async (invoiceId: string): Promise<InvoiceDto> => {
    return this.get(`/${VERSION}/invoice/${invoiceId}`).then(res => res);
  };
  GetInvoicesBySubscriber = async (
    subscriberId: string
  ): Promise<InvoiceDto[]> => {
    return this.get(`/${VERSION}?subscriber=${subscriberId}`).then(res => res);
  };
  GetInvoicesByNetwork = async (networkId: string): Promise<InvoiceDto[]> => {
    return this.get(`/${VERSION}?network=${networkId}`).then(res => res);
  };
  GetInvoicePDF = async (invoiceId: string): Promise<any> => {
    return this.get(`/${VERSION}/pdf/${invoiceId}`).then(res => res);
  };

  AddInvoice = async (rawInvoice: string): Promise<InvoiceDto> => {
    return this.post(`/${VERSION}`, {
      body: rawInvoice,
    }).then(res => res);
  };

  RemoveInvoice = async (invoiceId: string): Promise<any> => {
    return this.delete(`/${VERSION}/${invoiceId}`).then(res => res);
  };
}

export default BillingAPI;
