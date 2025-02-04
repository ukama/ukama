/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

class ReportAPI extends RESTDataSource {
  async getGeneratedPdfReport(
    baseURL: string,
    id: string
  ): Promise<ArrayBuffer> {
    try {
      const response = await fetch(`${baseURL}/v1/pdf/${id}`, {
        headers: { Accept: "application/pdf" },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch PDF: ${response.status}`);
      }

      return await response.arrayBuffer();
    } catch (error) {
      throw new Error(
        `Failed to fetch PDF with ID ${id}: ${
          error instanceof Error ? error.message : "Unknown error"
        }`
      );
    }
  }
}

export default ReportAPI;
