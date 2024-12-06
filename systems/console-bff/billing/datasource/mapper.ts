/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { GetReportResDto, GetReportsResDto } from "../resolvers/types";

export const dtoToReportsDto = (res: GetReportsResDto): GetReportsResDto => {
  return res;
};

export const dtoToReportDto = (res: GetReportResDto): GetReportResDto => {
  return res;
};
