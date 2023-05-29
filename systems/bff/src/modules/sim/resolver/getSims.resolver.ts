import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { SimsResDto } from "../types";

@Service()
@Resolver()
export class GetSimsResolver {
    constructor(private readonly simService: SimService) {}

    @Query(() => SimsResDto)
    @UseMiddleware(Authentication)
    async getSims(
        @Arg("type") type: string,
        @Ctx() ctx: Context
    ): Promise<SimsResDto> {
        return await this.simService.getSims(type, parseHeaders(ctx));
    }
}
