import { Service } from "typedi";
import { UserService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { ESimQRCodeRes, GetESimQRCodeInput } from "../types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, Arg, Ctx, Query, UseMiddleware } from "type-graphql";

@Service()
@Resolver()
export class GetEsimQRResolver {
    constructor(private readonly userService: UserService) {}

    @Query(() => ESimQRCodeRes)
    @UseMiddleware(Authentication)
    async getEsimQR(
        @Arg("data") data: GetESimQRCodeInput,
        @Ctx() ctx: Context,
    ): Promise<ESimQRCodeRes | null> {
        return this.userService.getEsimQRCode(data, parseCookie(ctx));
    }
}
