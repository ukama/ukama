/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import type { Request, Response } from "express";
import { gunzipSync } from "zlib";

import { compression } from "../middleware/compression";

type SentBody = string | Buffer | undefined;

const makeRes = () => {
  const headers = new Map<string, unknown>();
  let sent: SentBody;
  const res = {
    send: (body?: SentBody) => {
      sent = body;
      return res;
    },
    setHeader: (name: string, value: unknown) => {
      headers.set(name.toLowerCase(), value);
    },
    getHeader: (name: string) => headers.get(name.toLowerCase()),
  } as unknown as Response;
  return { res, headers, getSent: () => sent };
};

const makeReq = (acceptEncoding?: string) =>
  ({
    headers: acceptEncoding ? { "accept-encoding": acceptEncoding } : {},
  }) as unknown as Request;

const run = (req: Request, res: Response) =>
  new Promise<void>(resolve => compression()(req, res, () => resolve()));

describe("compression middleware", () => {
  const bigBody = JSON.stringify({ data: "x".repeat(5000) });

  it("gzips large responses when the client accepts gzip", async () => {
    const { res, headers, getSent } = makeRes();
    await run(makeReq("gzip, deflate"), res);
    res.send(bigBody);
    expect(headers.get("content-encoding")).toBe("gzip");
    expect(headers.get("vary")).toBe("Accept-Encoding");
    const sent = getSent();
    expect(Buffer.isBuffer(sent)).toBe(true);
    expect(gunzipSync(sent as Buffer).toString()).toBe(bigBody);
  });

  it("skips small payloads", async () => {
    const { res, headers, getSent } = makeRes();
    await run(makeReq("gzip"), res);
    res.send("tiny");
    expect(headers.get("content-encoding")).toBeUndefined();
    expect(getSent()).toBe("tiny");
  });

  it("skips clients that do not accept gzip", async () => {
    const { res, headers, getSent } = makeRes();
    await run(makeReq(), res);
    res.send(bigBody);
    expect(headers.get("content-encoding")).toBeUndefined();
    expect(getSent()).toBe(bigBody);
  });

  it("does not double-encode already-encoded responses", async () => {
    const { res, getSent } = makeRes();
    await run(makeReq("gzip"), res);
    res.setHeader("content-encoding", "br");
    res.send(bigBody);
    expect(getSent()).toBe(bigBody);
  });
});
