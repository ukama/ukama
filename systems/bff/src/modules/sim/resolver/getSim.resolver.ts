import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { GetSimInputDto, SimDto } from "../types";

@Service()
@Resolver()
export class GetSimResolver {
    constructor(private readonly simService: SimService) {}

    @Query(() => SimDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: GetSimInputDto,
        @Ctx() ctx: Context
    ): Promise<SimDto> {
        return await this.simService.getSim(data, parseHeaders(ctx));
    }
}
