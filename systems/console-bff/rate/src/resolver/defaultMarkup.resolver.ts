import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { BoolResponse, Context } from "../../../common/types";
import { DefaultMarkupInputDto } from "../types";

@Resolver()
export class DefaultMarkupResolver {

    @Mutation(() => BoolResponse)
    @UseMiddleware(Authentication)
    async defaultMarkup(
        @Arg("data") data: DefaultMarkupInputDto,
        @Ctx() ctx: Context
    ): Promise<BoolResponse> {
        const { dataSources } = ctx;
        return dataSources.dataSource.defaultMarkup(data, parseHeaders(ctx));
    }
}
