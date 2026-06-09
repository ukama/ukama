/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Assigns a correlation id to every request: reuses an inbound x-request-id
 * (e.g. from an upstream proxy) or generates one. The id is echoed in the
 * response header, written onto req.headers so the Apollo context can read it,
 * and surfaced in access logs — letting a single request be traced across the
 * gateway and (via propagation) the subgraphs.
 */
import { randomUUID } from "crypto";
import type { NextFunction, Request, Response } from "express";

export const REQUEST_ID_HEADER = "x-request-id";

export const requestId = () => {
  return (req: Request, res: Response, next: NextFunction): void => {
    const inbound = req.header(REQUEST_ID_HEADER);
    const id = inbound && inbound.trim() !== "" ? inbound : randomUUID();
    req.headers[REQUEST_ID_HEADER] = id;
    res.setHeader(REQUEST_ID_HEADER, id);
    next();
  };
};
