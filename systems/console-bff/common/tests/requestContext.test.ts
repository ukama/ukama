/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import {
  getRequestId,
  runWithRequestId,
  setRequestId,
} from "../logger/requestContext";

describe("request context (AsyncLocalStorage)", () => {
  it("exposes the id inside runWithRequestId", () => {
    runWithRequestId("req-1", () => {
      expect(getRequestId()).toBe("req-1");
    });
  });

  it("has no id outside any request scope", () => {
    expect(getRequestId()).toBeUndefined();
  });

  it("keeps concurrent requests isolated", async () => {
    const seen: Record<string, string | undefined> = {};
    await Promise.all([
      new Promise<void>(resolve =>
        runWithRequestId("a", async () => {
          await new Promise(r => setTimeout(r, 20));
          seen.a = getRequestId();
          resolve();
        })
      ),
      new Promise<void>(resolve =>
        runWithRequestId("b", async () => {
          await new Promise(r => setTimeout(r, 5));
          seen.b = getRequestId();
          resolve();
        })
      ),
    ]);
    expect(seen).toEqual({ a: "a", b: "b" });
  });

  it("binds an id with setRequestId (callback-free)", async () => {
    await runWithRequestId("outer", async () => {
      setRequestId("bound");
      expect(getRequestId()).toBe("bound");
    });
  });

  it("ignores an empty id in setRequestId", () => {
    runWithRequestId("keep", () => {
      setRequestId("");
      expect(getRequestId()).toBe("keep");
    });
  });
});
