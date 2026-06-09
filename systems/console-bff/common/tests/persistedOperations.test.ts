/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import type { Request, Response } from "express";
import { mkdtempSync, writeFileSync } from "fs";
import { tmpdir } from "os";
import { join } from "path";

import { persistedOperations, sha256 } from "../middleware/persistedOperations";

const KNOWN_QUERY = "query Ping { ping }";

const writeManifest = (): string => {
  const dir = mkdtempSync(join(tmpdir(), "persisted-ops-"));
  const path = join(dir, "manifest.json");
  writeFileSync(path, JSON.stringify({ [sha256(KNOWN_QUERY)]: KNOWN_QUERY }));
  return path;
};

const makeRes = () => {
  let statusCode = 200;
  let body: unknown;
  const res = {
    status: (code: number) => {
      statusCode = code;
      return res;
    },
    json: (payload: unknown) => {
      body = payload;
      return res;
    },
  } as unknown as Response;
  return { res, getStatus: () => statusCode, getBody: () => body };
};

const makeReq = (query: string) =>
  ({ body: { query, operationName: "X" } }) as unknown as Request;

describe("persistedOperations middleware", () => {
  it("is a no-op when enforcement is off", () => {
    const middleware = persistedOperations({ enforced: false });
    const { res } = makeRes();
    const next = jest.fn();
    middleware(makeReq("query Evil { __typename }"), res, next);
    expect(next).toHaveBeenCalled();
  });

  it("allows operations present in the manifest", () => {
    const middleware = persistedOperations({
      enforced: true,
      manifestPath: writeManifest(),
    });
    const { res } = makeRes();
    const next = jest.fn();
    middleware(makeReq(KNOWN_QUERY), res, next);
    expect(next).toHaveBeenCalled();
  });

  it("rejects unknown operations with 403", () => {
    const middleware = persistedOperations({
      enforced: true,
      manifestPath: writeManifest(),
    });
    const { res, getStatus, getBody } = makeRes();
    const next = jest.fn();
    middleware(makeReq("query Evil { secrets }"), res, next);
    expect(next).not.toHaveBeenCalled();
    expect(getStatus()).toBe(403);
    expect(getBody()).toMatchObject({
      errors: [{ extensions: { code: "PERSISTED_OPERATION_NOT_FOUND" } }],
    });
  });

  it("optionally lets introspection through", () => {
    const middleware = persistedOperations({
      enforced: true,
      manifestPath: writeManifest(),
      allowIntrospection: true,
    });
    const { res } = makeRes();
    const next = jest.fn();
    middleware(makeReq("query I { __schema { types { name } } }"), res, next);
    expect(next).toHaveBeenCalled();
  });

  it("fails fast when the manifest is missing", () => {
    expect(() =>
      persistedOperations({
        enforced: true,
        manifestPath: "/does/not/exist.json",
      })
    ).toThrow("manifest not found");
  });
});
