/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import {
  BillHistoryDto,
  BillHistoryResponse,
  BillResponse,
  CurrentBillResponse,
} from "../resolvers/types";

export const dtoToDto = (res: CurrentBillResponse): BillResponse => {
  const bill = res.data;
  let total = 0;
  for (const sub of bill) {
    const subTotal = sub.subtotal;
    total = total + subTotal;
  }
  return {
    bill,
    total,
    dueDate: "10-10-2021",
    billMonth: "11-10-2021",
  };
};

export const billHistoryDtoToDto = (
  res: BillHistoryResponse
): BillHistoryDto[] => {
  return res.data;
};
