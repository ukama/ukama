import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { AllocateSimInputDto, SimDto } from "../types";

@Service()
@Resolver()
export class AllocateSimResolver {
    constructor(private readonly simService: SimService) {}

    @Mutation(() => SimDto)
    @UseMiddleware(Authentication)
    async allocateSim(
        @Arg("data") data: AllocateSimInputDto,
        @Ctx() ctx: Context
    ): Promise<SimDto> {
        return await this.simService.allocateSim(data, parseHeaders(ctx));
    }
}