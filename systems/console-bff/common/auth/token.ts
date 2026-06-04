/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
import { createHmac, timingSafeEqual } from "crypto";

import { TOKEN_SECRET } from "../configs";

const SIGNATURE_SEPARATOR = ".";

const computeSignature = (payload: string): string =>
  createHmac("sha256", TOKEN_SECRET).update(payload).digest("base64url");

/**
 * Signs an opaque payload (base64-encoded session claims) with
 * HMAC-SHA256. Output format: `<payload>.<signature>`.
 */
export const signToken = (payload: string): string =>
  `${payload}${SIGNATURE_SEPARATOR}${computeSignature(payload)}`;

/**
 * Verifies a signed token and returns its payload, or null when the
 * token is missing, malformed, or its signature does not match.
 * Uses a timing-safe comparison to prevent signature oracle attacks.
 */
export const verifyToken = (token: string): string | null => {
  if (!token) return null;
  const separatorIndex = token.lastIndexOf(SIGNATURE_SEPARATOR);
  if (separatorIndex <= 0) return null;

  const payload = token.slice(0, separatorIndex);
  const signature = token.slice(separatorIndex + 1);
  const expected = computeSignature(payload);

  const signatureBuf = Buffer.from(signature);
  const expectedBuf = Buffer.from(expected);
  if (signatureBuf.length !== expectedBuf.length) return null;
  if (!timingSafeEqual(signatureBuf, expectedBuf)) return null;

  return payload;
};
