/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { AlertDto, AlertResponse } from "../resolver/types";

export const dtoToDto = (res: AlertResponse): AlertDto[] => {
  const alerts: AlertDto[] = [];

  for (const alert of res.data) {
    const alertObj = {
      id: alert.id,
      type: alert.type,
      title: alert.title,
      description: alert.description,
      alertDate: new Date(alert.alertDate),
    };
    alerts.push(alertObj);
  }

  return alerts;
};
