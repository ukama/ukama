/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import "reflect-metadata";

import { mapNodeInterfaces } from "../../health/datasource/mapper";

const NODE_ID = "uk-sa2637-tnode-v0-6ea4";

describe("mapNodeInterfaces", () => {
  it("carries cellular service and radio state from the health response", () => {
    const res = {
      interfaces: {
        cellular: { available: true, error: "", service: "on" },
        radio: { available: true, state: "off" },
        gps: { available: true, lock: true },
        fem: null,
      },
    };
    expect(mapNodeInterfaces(NODE_ID, res)).toEqual({
      nodeId: NODE_ID,
      cellular: { available: true, error: "", service: "on" },
      radio: { available: true, state: "off" },
    });
  });

  it("defaults omitted fields to empty values", () => {
    const res = { interfaces: { cellular: {}, radio: { available: true } } };
    expect(mapNodeInterfaces(NODE_ID, res)).toEqual({
      nodeId: NODE_ID,
      cellular: { available: false, error: "", service: "" },
      radio: { available: true, state: "" },
    });
  });

  it("leaves interfaces the node did not report undefined", () => {
    expect(mapNodeInterfaces(NODE_ID, { interfaces: {} })).toEqual({
      nodeId: NODE_ID,
      cellular: undefined,
      radio: undefined,
    });
    expect(mapNodeInterfaces(NODE_ID, null)).toEqual({
      nodeId: NODE_ID,
      cellular: undefined,
      radio: undefined,
    });
  });
});
