/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { Server, createServer } from "http";
import { AddressInfo } from "net";

import { createTimedFetch } from "../datasource";

describe("createTimedFetch", () => {
  let server: Server;
  let url: string;

  beforeAll(done => {
    server = createServer((req, res) => {
      const delay = Number(
        new URL(req.url ?? "/", "http://x").searchParams.get("delay")
      );
      setTimeout(() => res.end("ok"), delay);
    });
    server.listen(0, () => {
      url = `http://127.0.0.1:${(server.address() as AddressInfo).port}/`;
      done();
    });
  });

  afterAll(done => {
    server.closeAllConnections();
    server.close(done);
  });

  it("returns a response that arrives inside the limit", async () => {
    const res = await createTimedFetch(500)(`${url}?delay=50`);
    expect(await res.text()).toBe("ok");
  });

  it("aborts a response that arrives after the limit", async () => {
    await expect(createTimedFetch(50)(`${url}?delay=500`)).rejects.toThrow(
      /abort/i
    );
  });
});
