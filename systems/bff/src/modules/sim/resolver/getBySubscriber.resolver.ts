import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SimService } from "../service";
import { GetSimBySubscriberIdInputDto, SimDetailsDto } from "../types";

@Service()
@Resolver()
export class GetSimBySubscriberResolver {
    constructor(private readonly simService: SimService) {}

    @Query(() => SimDetailsDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: GetSimBySubscriberIdInputDto,
        @Ctx() ctx: Context
    ): Promise<SimDetailsDto> {
        return await this.simService.getSimBySubscriberId(
            data,
            parseHeaders(ctx)
        );
    }
}
