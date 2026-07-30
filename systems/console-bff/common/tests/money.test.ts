/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import {
  majorUnit,
  normalizeKpiValue,
  normalizeReportCell,
} from "../../analytics/money";

describe("majorUnit", () => {
  it("uses the org currency code, lowercased", () => {
    expect(majorUnit("USD")).toBe("usd");
    expect(majorUnit("eur")).toBe("eur");
  });

  it("falls back to 'major' when the token carries no currency claim", () => {
    expect(majorUnit(undefined)).toBe("major");
    expect(majorUnit("")).toBe("major");
    expect(majorUnit("   ")).toBe("major");
  });
});

describe("normalizeKpiValue", () => {
  it("converts cents to major units and rewrites the unit", () => {
    const out = normalizeKpiValue(
      { value: 700, unit: "cents", op: "SUM" },
      "usd"
    );
    expect(out.value).toBe(7);
    expect(out.unit).toBe("usd");
  });

  it("converts trend changeAbs and prevValue, but not changePct", () => {
    const out = normalizeKpiValue(
      {
        value: 1000,
        unit: "cents",
        op: "SUM",
        trend: {
          direction: "up",
          changePct: 25,
          changeAbs: 200,
          prevValue: 800,
          hasPrevious: true,
        },
      },
      "usd"
    );
    expect(out.value).toBe(10);
    expect(out.trend?.changeAbs).toBe(2);
    expect(out.trend?.prevValue).toBe(8);
    expect(out.trend?.changePct).toBe(25);
    expect(out.trend?.direction).toBe("up");
  });

  it("leaves null trend members null", () => {
    const out = normalizeKpiValue(
      {
        value: 500,
        unit: "cents",
        op: "SUM",
        trend: { direction: "new", changeAbs: null, prevValue: null },
      },
      "usd"
    );
    expect(out.trend?.changeAbs).toBeNull();
    expect(out.trend?.prevValue).toBeNull();
  });

  it("does not convert a COUNT op, and corrects its mislabelled unit", () => {
    // The gateway tags `unit` per KPI, not per op, so REVENUE/COUNT arrives
    // labelled "cents" while actually being a purchase count.
    const out = normalizeKpiValue(
      { value: 42, unit: "cents", op: "COUNT" },
      "usd"
    );
    expect(out.value).toBe(42);
    expect(out.unit).toBe("count");
  });

  it("passes non-monetary units through untouched", () => {
    for (const unit of ["count", "bytes", "percent", "", undefined]) {
      const out = normalizeKpiValue({ value: 1024, unit, op: "SUM" }, "usd");
      expect(out.value).toBe(1024);
      expect(out.unit).toBe(unit);
    }
  });

  it("scales fractional cents from AVG rollups without rounding", () => {
    const out = normalizeKpiValue(
      { value: 333.3333, unit: "cents", op: "AVG" },
      "usd"
    );
    expect(out.value).toBeCloseTo(3.333333, 6);
  });

  it("preserves the trend-consistency invariant value - prev === changeAbs", () => {
    // ukama-lab's kpi_contract (require_trend_consistency) recomputes this,
    // so the three members must stay mutually consistent after conversion.
    for (const [value, prev, abs] of [
      [500, 400, 100],
      [333.3333, 300.7, 32.6333],
      [1, 3, -2],
    ]) {
      const out = normalizeKpiValue(
        {
          value,
          unit: "cents",
          op: "SUM",
          trend: { changeAbs: abs, prevValue: prev, changePct: 0 },
        },
        "usd"
      );
      expect(out.value - (out.trend?.prevValue ?? 0)).toBeCloseTo(
        out.trend?.changeAbs ?? 0,
        9
      );
    }
  });
});

describe("normalizeReportCell", () => {
  it("converts a money cell tagged cents", () => {
    const out = normalizeReportCell(
      { value: 700, unit: "cents", format: "money" },
      "usd"
    );
    expect(out.value).toBe(7);
    expect(out.unit).toBe("usd");
  });

  it("converts on format alone when the composer emitted an empty unit", () => {
    // Entities with no rollup rows in the window hit the zero-cell early
    // return in composer.go: value 0, unit "", format still "money".
    const out = normalizeReportCell(
      { value: 0, unit: "", format: "money" },
      "usd"
    );
    expect(out.value).toBe(0);
    expect(out.unit).toBe("usd");
  });

  it("leaves non-money columns alone", () => {
    const sold = normalizeReportCell({ value: 3, unit: "count" }, "usd");
    expect(sold.value).toBe(3);
    expect(sold.unit).toBe("count");

    const data = normalizeReportCell(
      { value: 2097152, unit: "bytes", format: "bytes" },
      "usd"
    );
    expect(data.value).toBe(2097152);
    expect(data.unit).toBe("bytes");
  });
});
