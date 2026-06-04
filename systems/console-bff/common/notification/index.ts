/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { CONSOLE_APP_URL } from "./../configs/index";

type TEventKeyToAction = {
  title: string;
  action: string;
};

/** Minimal structural input so any notification shape can derive a redirect
 *  (decouples this shared helper from the subscriptions module's types). */
interface NotificationLike {
  title?: string;
  resourceId?: string;
}

export const eventKeyToAction = (
  key: string,
  data: NotificationLike
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
