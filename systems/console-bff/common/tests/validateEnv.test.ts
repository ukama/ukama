/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * validateEnv reads config constants that are evaluated when the module is
 * first imported, so each case resets the module registry and sets env first.
 */
const REQUIRED = {
  AUTH_URL: "http://kratos:4433",
  NUCLEUS_API_GW: "http://nucleus:8080",
  INIT_API_GW: "http://init:8080",
  INVENTORY_API_GW: "http://inventory:8080",
};

describe("validateEnv", () => {
  const ORIGINAL = process.env;

  beforeEach(() => {
    jest.resetModules();
    process.env = { ...ORIGINAL, NODE_ENV: "development" };
  });

  afterAll(() => {
    process.env = ORIGINAL;
  });

  const load = () =>
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    require("../configs/validateEnv").validateEnv as () => void;

  it("passes when all required URLs are set", () => {
    Object.assign(process.env, REQUIRED);
    expect(() => load()()).not.toThrow();
  });

  it("throws when a required variable is missing", () => {
    Object.assign(process.env, REQUIRED, { AUTH_URL: "" });
    expect(() => load()()).toThrow(/AUTH_URL/);
  });

  it("throws when a required URL is malformed", () => {
    Object.assign(process.env, REQUIRED, { INIT_API_GW: "not-a-url" });
    expect(() => load()()).toThrow(/INIT_API_GW/);
  });

  it("rejects localhost URLs in production", () => {
    Object.assign(process.env, REQUIRED, {
      NODE_ENV: "production",
      JWT_SECRET: "a-strong-secret-value-for-testing-only",
      AUTH_URL: "http://localhost:4433",
    });
    expect(() => load()()).toThrow(/localhost/);
  });
});
