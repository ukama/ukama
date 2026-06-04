/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { signToken, verifyToken } from "../auth/token";
import { parseToken } from "../utils";

const b64 = (s: string) => Buffer.from(s).toString("base64");
const nowSec = () => Math.floor(Date.now() / 1000);

// orgId;orgName;userId;name;email;role;verified;welcome;country;currency[;exp]
const claims = (exp?: number) =>
  `org-1;ukama;user-1;Admin;a@x.com;ROLE_OWNER;true;false;COD;cdf${
    exp !== undefined ? `;${exp}` : ""
  }`;

describe("token signing", () => {
  it("round-trips a signed payload", () => {
    const payload = b64(claims());
    expect(verifyToken(signToken(payload))).toBe(payload);
  });

  it("rejects a tampered signature", () => {
    const t = signToken(b64(claims()));
    const tampered = t.slice(0, -1) + (t.endsWith("A") ? "B" : "A");
    expect(verifyToken(tampered)).toBeNull();
  });

  it("rejects an unsigned payload", () => {
    expect(verifyToken(b64(claims()))).toBeNull();
  });
});

describe("parseToken", () => {
  it("returns claims for a valid, unexpired token", () => {
    const t = signToken(b64(claims(nowSec() + 3600)));
    expect(parseToken(t, "orgId")).toBe("org-1");
    expect(parseToken(t, "orgName")).toBe("ukama");
    expect(parseToken(t, "userId")).toBe("user-1");
  });

  it("rejects an expired token", () => {
    const t = signToken(b64(claims(nowSec() - 1)));
    expect(() => parseToken(t, "orgId")).toThrow();
  });

  it("rejects a tampered token", () => {
    const t = signToken(b64(claims(nowSec() + 3600)));
    const tampered = t.slice(0, -1) + (t.endsWith("A") ? "B" : "A");
    expect(() => parseToken(tampered, "orgId")).toThrow();
  });

  it("accepts a legacy token without an exp claim", () => {
    const t = signToken(b64(claims()));
    expect(parseToken(t, "orgId")).toBe("org-1");
  });

  it("returns undefined for an empty token", () => {
    expect(parseToken("", "orgId")).toBeUndefined();
  });
});
