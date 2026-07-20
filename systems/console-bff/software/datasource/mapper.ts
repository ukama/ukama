/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import {
  PromoteReleaseResponse,
  ReleaseCatalog,
  Softwares,
  StringResponse,
} from "../resolvers/types";

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

// The api-gateway returns protobuf JSON, so some keys are snake_case
// (json_name overrides); accept both spellings defensively.
interface RawPromoteRelease {
  message?: string;
  name?: string;
  desired_version?: string;
  desiredVersion?: string;
}

interface RawRelease {
  name?: string;
  type?: string;
  version?: string;
  available?: boolean;
  chunked?: boolean;
  desired?: boolean;
  uploaded_at?: string;
  uploadedAt?: string;
}

interface RawReleaseCatalog {
  releases?: RawRelease[];
}

export const mapPromoteRelease = (
  res: RawPromoteRelease
): PromoteReleaseResponse => {
  return {
    message: res?.message ?? "",
    name: res?.name ?? "",
    desiredVersion: res?.desired_version ?? res?.desiredVersion ?? "",
  };
};

export const mapReleaseCatalog = (res: RawReleaseCatalog): ReleaseCatalog => {
  return {
    releases: (res?.releases ?? []).map(r => ({
      name: r?.name ?? "",
      type: r?.type ?? "",
      version: r?.version ?? "",
      available: Boolean(r?.available),
      chunked: Boolean(r?.chunked),
      desired: Boolean(r?.desired),
      uploadedAt: r?.uploaded_at ?? r?.uploadedAt ?? "",
    })),
  };
};
