import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { GetSimByNetworkInputDto, SimDetailsDto } from "../types";

@Service()
@Resolver()
export class GetSimByNetworkResolver {
    constructor(private readonly simService: SimService) {}

    @Mutation(() => SimDetailsDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: GetSimByNetworkInputDto,
        @Ctx() ctx: Context
    ): Promise<SimDetailsDto> {
        return await this.simService.getSimByNetworkId(data, parseCookie(ctx));
    }
}
