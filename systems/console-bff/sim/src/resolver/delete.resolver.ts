import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { DeleteSimInputDto, DeleteSimResDto } from "../types";

@Resolver()
export class DeleteSimResolver {

    @Mutation(() => DeleteSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: DeleteSimInputDto,
        @Ctx() ctx: Context
    ): Promise<DeleteSimResDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.deleteSim(data, parseHeaders(ctx));
    }
}
