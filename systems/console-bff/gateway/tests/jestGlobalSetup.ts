/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

/**
 * Jest globalSetup: ensures process.env.TOKEN is a valid signed JWT before
 * any test file is loaded. If TOKEN is already a JWT it is used as-is;
 * if it looks like a legacy base64 token it is migrated on-the-fly.
 */
import "dotenv/config";
import jwt from "jsonwebtoken";

export default async function globalSetup() {
  const rawToken = process.env.TOKEN ?? "";
  const secret =
    process.env.JWT_SECRET ?? "change-me-in-production-use-32-bytes!";

  if (!rawToken) return;

  // Already a JWT — verify it is still valid (may throw if expired or wrong secret)
  if (rawToken.split(".").length === 3) {
    try {
      jwt.verify(rawToken, secret);
      return; // valid JWT, nothing to do
    } catch {
      // fall through to re-sign with current secret
    }
  }

  // Legacy base64 token — decode and re-sign as JWT
  try {
    const decoded = Buffer.from(rawToken, "base64").toString("utf-8");
    const parts = decoded.split(";");
    const [orgId, orgName, userId, name, email, role] = parts;
    const isEmailVerified = parts[6] === "true";
    const isShowWelcome = parts[7] === "true";
    const country = parts[8] ?? "";
    const currency = parts[9] ?? "";

    const jwtToken = jwt.sign(
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
      { expiresIn: "24h" }
    );

    process.env.TOKEN = jwtToken;
  } catch {
    console.warn(
      "[jestGlobalSetup] TOKEN is neither a valid JWT nor a base64 legacy token — tests may fail"
    );
  }
}
