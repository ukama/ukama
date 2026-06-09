/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Arg, Ctx, Query, Resolver } from "type-graphql";

import type { AppContext } from "../../server/context";
import { GetPdfReportUrlDto } from "./types";

@Resolver()
export class GetPdfGeneratedReportResolver {
  @Query(() => GetPdfReportUrlDto)
  async getGeneratedPdfReport(
    @Arg("id") id: string,
    @Ctx() ctx: AppContext
  ): Promise<GetPdfReportUrlDto> {
    const { dataSources } = ctx;
    const baseURL = await ctx.urls.url("report");

    try {
      const buffer = await dataSources.report.getGeneratedPdfReport(
        baseURL,
        id
      );

      if (!buffer || buffer.byteLength < 100) {
        throw new Error("Generated PDF appears to be empty or corrupted");
      }

      const base64 = Buffer.from(buffer).toString("base64");

      return {
        id,
        filename: `${id}.pdf`,
        contentType: "application/pdf",
        downloadUrl: `data:application/pdf;base64,${base64}`,
      };
    } catch (error) {
      throw new Error(
        `Failed to process PDF with ID ${id}: ${
          error instanceof Error ? error.message : "Unknown error"
        }`
      );
    }
  }
}
