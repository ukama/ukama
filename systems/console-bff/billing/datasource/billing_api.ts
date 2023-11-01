/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

import { BILLING_API_GW } from "../../common/configs";
import { BillHistoryDto, BillResponse } from "../resolvers/types";
import { billHistoryDtoToDto, dtoToDto } from "./mapper";

const version = "/v1/invoices";
class BillingAPI extends RESTDataSource {
  baseURL = BILLING_API_GW + version;
  public getCurrentBill = async (): Promise<BillResponse> => {
    return this.get("/current")
      .then(res => dtoToDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  public getBillHistory = async (): Promise<BillHistoryDto[]> => {
    return this.get("/history")
      .then(res => billHistoryDtoToDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default BillingAPI;
