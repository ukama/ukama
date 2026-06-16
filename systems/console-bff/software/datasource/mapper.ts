/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Softwares, StringResponse } from "../resolvers/types";

export const mapSoftwares = (softwares: Softwares): Softwares => {
  return {
    software: softwares.software.map(software => ({
      id: software.id,
      releaseDate: software.releaseDate,
      nodeId: software.nodeId,
      status: software.status,
      changeLog: software.changeLog,
      currentVersion: software.currentVersion,
      desiredVersion: software.desiredVersion,
      name: software.name,
      space: software.space,
      notes: software.notes,
      metricsKeys: software.metricsKeys,
      createdAt: software.createdAt,
      updatedAt: software.updatedAt,
    })),
  };
};

export const mapUpdateSoftware = (
  updateSoftware: StringResponse
): StringResponse => {
  return {
    message: updateSoftware.message,
  };
};
