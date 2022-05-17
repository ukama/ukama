import { Service } from "typedi";
import { UserService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { ESimQRCodeRes, GetESimQRCodeInput, UserResDto } from "../types";
import { Resolver, Arg, Ctx, Mutation, UseMiddleware } from "type-graphql";

@Service()
@Resolver()
export class GetEsimQRResolver {
    constructor(private readonly userService: UserService) {}

    @Mutation(() => UserResDto)
    @UseMiddleware(Authentication)
    async getEsimQR(
        @Arg("data") data: GetESimQRCodeInput,
        @Ctx() ctx: Context
    ): Promise<ESimQRCodeRes | null> {
        return this.userService.getEsimQRCode(data, parseCookie(ctx));
    }
}
