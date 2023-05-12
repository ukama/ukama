import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { AllocateSimInputDto, SimResDto } from "../types";

@Service()
@Resolver()
export class AllocateSimResolver {
    constructor(private readonly simService: SimService) {}

    @Mutation(() => SimResDto)
    @UseMiddleware(Authentication)
    async allocateSim(
        @Arg("data") data: AllocateSimInputDto,
        @Ctx() ctx: Context
    ): Promise<SimResDto> {
        return await this.simService.allocateSim(data, parseCookie(ctx));
    }
}
