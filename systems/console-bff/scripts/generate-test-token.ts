/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

/**
 * One-time migration helper: converts a legacy base64 session token to a
 * signed JWT suitable for use as TOKEN in .env for integration tests.
 *
 * Usage:
 *   TOKEN=<old-base64-token> JWT_SECRET=<secret> npx ts-node scripts/generate-test-token.ts
 *
 * Output: prints the new JWT — paste it as TOKEN=<jwt> in your .env
 */
import "dotenv/config";
import jwt from "jsonwebtoken";

const legacyToken = process.env.TOKEN ?? "";
const secret =
  process.env.JWT_SECRET ?? "change-me-in-production-use-32-bytes!";

if (!legacyToken) {
  console.error("Error: TOKEN env var is required (the old base64 token)");
  process.exit(1);
}

let orgId = "";
let orgName = "";
let userId = "";
let name = "";
let email = "";
let role = "";
let isEmailVerified = false;
let isShowWelcome = false;
let country = "";
let currency = "";

try {
  const decoded = Buffer.from(legacyToken, "base64").toString("utf-8");
  const parts = decoded.split(";");
  [orgId, orgName, userId, name, email, role] = parts;
  isEmailVerified = parts[6] === "true";
  isShowWelcome = parts[7] === "true";
  country = parts[8] ?? "";
  currency = parts[9] ?? "";
  console.error("Decoded legacy token fields:", {
    orgId,
    orgName,
    userId,
    name,
    email,
    role,
  });
} catch {
  console.error("Error: TOKEN is not a valid base64-encoded legacy token");
  process.exit(1);
}

const token = jwt.sign(
  {
    orgId,
    orgName,
    userId,
    name,
    email,
    role,
    isEmailVerified,
    isShowWelcome,
    country,
    currency,
  },
  secret,
  { expiresIn: "365d" }
);

console.log(token);
