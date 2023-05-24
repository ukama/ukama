import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UserService } from "../service";
import { ESimQRCodeRes, GetESimQRCodeInput } from "../types";

@Service()
@Resolver()
export class GetEsimQRResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => ESimQRCodeRes)
    @UseMiddleware(Authentication)
    async getEsimQR(
        @Arg("data") data: GetESimQRCodeInput,
        @Ctx() ctx: Context
    ): Promise<ESimQRCodeRes | null> {
        return this.userService.getEsimQRCode(data, parseHeaders(ctx));
    }
}
