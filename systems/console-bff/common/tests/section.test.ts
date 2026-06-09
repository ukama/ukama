/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import "reflect-metadata";

import { SectionErrorCollector, withSection } from "../../dashboard/section";
import { SectionErrorCode } from "../../dashboard/types";

describe("withSection", () => {
  it("returns the value and records nothing on success", async () => {
    const collector = new SectionErrorCollector();
    const result = await withSection(collector, "kpis", async () => 42);
    expect(result).toBe(42);
    expect(collector.list()).toEqual([]);
  });

  it("returns null and records a typed error on failure", async () => {
    const collector = new SectionErrorCollector();
    const result = await withSection(collector, "kpis", async () => {
      throw new Error("upstream exploded");
    });
    expect(result).toBeNull();
    expect(collector.list()).toEqual([
      {
        section: "kpis",
        code: SectionErrorCode.INTERNAL,
        message: "Failed to load kpis",
      },
    ]);
  });

  it("maps HTTP-like status codes", async () => {
    const collector = new SectionErrorCollector();
    await withSection(collector, "a", async () => {
      throw { code: 404, message: "not found" };
    });
    await withSection(collector, "b", async () => {
      throw { extensions: { response: { status: 403 } } };
    });
    await withSection(collector, "c", async () => {
      throw { code: 500, message: "internal upstream" };
    });
    expect(collector.list().map(e => e.code)).toEqual([
      SectionErrorCode.NOT_FOUND,
      SectionErrorCode.FORBIDDEN,
      SectionErrorCode.UPSTREAM_ERROR,
    ]);
  });

  it("times out slow sections with UPSTREAM_TIMEOUT", async () => {
    const collector = new SectionErrorCollector();
    const result = await withSection(
      collector,
      "slow",
      () => new Promise(resolve => setTimeout(resolve, 1000)),
      20
    );
    expect(result).toBeNull();
    expect(collector.list()[0].code).toBe(SectionErrorCode.UPSTREAM_TIMEOUT);
  });

  it("keeps errors isolated per collector instance (request-scoped)", async () => {
    const a = new SectionErrorCollector();
    const b = new SectionErrorCollector();
    await withSection(a, "x", async () => {
      throw new Error("fail");
    });
    expect(a.list()).toHaveLength(1);
    expect(b.list()).toHaveLength(0);
  });
});
