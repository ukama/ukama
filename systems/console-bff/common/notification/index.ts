/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NotificationsResDto } from "../../subscriptions/resolvers/types";
import { CONSOLE_APP_URL } from "./../configs/index";

type TEventKeyToAction = {
  title: string;
  action: string;
};

export const eventKeyToAction = (
  key: string,
  data: NotificationsResDto
): TEventKeyToAction => {
  switch (key) {
    case "EventNodeStateTransition":
      if (data.title && data.title.includes("Unknown")) {
        return {
          title: "Configure node",
          action: `${CONSOLE_APP_URL}/configure/check?step=1&flow=ins&nid=${data.resourceId}`,
        };
      }
      return { title: "Node state changed", action: "" };

    case "EventInvoiceGenerate":
      return {
        title: "Ukama bill ready. View now.",
        action: `${CONSOLE_APP_URL}/manage/billing`,
      };

    default:
      return { title: "Network Updated", action: "updated" };
  }
};
