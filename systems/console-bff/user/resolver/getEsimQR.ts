import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";

import { Authentication } from "../../common/auth";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { ESimQRCodeRes, GetESimQRCodeInput } from "./types";

@Resolver()
export class GetEsimQRResolver {
  @Query(() => ESimQRCodeRes)
  @UseMiddleware(Authentication)
  async getEsimQR(
    @Arg("data") data: GetESimQRCodeInput,
    @Ctx() ctx: Context
  ): Promise<ESimQRCodeRes | null> {
    const { dataSources } = ctx;
    return dataSources.dataSource.getEsimQRCode(data, parseHeaders());
  }
}
