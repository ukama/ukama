import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseCookie } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { SimDataUsage } from "../types";

@Service()
@Resolver()
export class GetDataUsageResolver {
    constructor(private readonly simService: SimService) {}

    @Query(() => SimDataUsage)
    @UseMiddleware(Authentication)
    async getDataUsage(
        @Arg("simId") simId: string,
        @Ctx() ctx: Context
    ): Promise<SimDataUsage> {
        return await this.simService.getDataUsage(simId, parseCookie(ctx));
    }
}
