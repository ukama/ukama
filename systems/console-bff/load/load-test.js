/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * k6 load test for the consolidated console-bff (CONSOLIDATION-DESIGN P3).
 * Establishes baseline SLOs and a regression guardrail.
 *
 * Run:
 *   BASE_URL=http://localhost:8080 \
 *   UKAMA_SESSION=<cookie> TOKEN=<cookie> \
 *   k6 run load/load-test.js
 *
 * Health probes run unauthenticated; the GraphQL query runs only when
 * UKAMA_SESSION + TOKEN are provided (otherwise it's skipped, so the smoke
 * profile still works without credentials).
 */
import http from "k6/http";
import { check, group, sleep } from "k6";

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";
const SESSION = __ENV.UKAMA_SESSION || "";
const TOKEN = __ENV.TOKEN || "";
const authed = SESSION !== "" && TOKEN !== "";

export const options = {
  scenarios: {
    // Warm baseline → ramp → hold, then ramp down.
    api: {
      executor: "ramping-vus",
      startVUs: 1,
      stages: [
        { duration: "30s", target: 10 },
        { duration: "1m", target: 25 },
        { duration: "30s", target: 0 },
      ],
    },
  },
  // SLOs — tune against observed baselines.
  thresholds: {
    http_req_failed: ["rate<0.01"], // <1% errors
    "http_req_duration{kind:health}": ["p(95)<150"],
    "http_req_duration{kind:graphql}": ["p(95)<800"],
  },
};

const cookieHeader = `ukama_session=${SESSION}; token=${TOKEN}`;

export default function () {
  group("health", () => {
    const h = http.get(`${BASE_URL}/healthz`, { tags: { kind: "health" } });
    check(h, { "healthz 200": r => r.status === 200 });
    const ready = http.get(`${BASE_URL}/readyz`, { tags: { kind: "health" } });
    check(ready, { "readyz 200": r => r.status === 200 });
  });

  if (authed) {
    group("graphql", () => {
      const res = http.post(
        `${BASE_URL}/graphql`,
        JSON.stringify({ query: "query { getOrgs { user ownerOf { id name } } }" }),
        {
          headers: { "content-type": "application/json", cookie: cookieHeader },
          tags: { kind: "graphql" },
        }
      );
      check(res, {
        "graphql 200": r => r.status === 200,
        "graphql has data": r => r.body && r.body.includes('"getOrgs"'),
      });
    });
  }

  sleep(1);
}
