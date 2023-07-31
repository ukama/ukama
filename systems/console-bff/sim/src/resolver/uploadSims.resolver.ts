import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { UploadSimsInputDto, UploadSimsResDto } from "../types";

@Resolver()
export class UploadSimsResolver {

    @Mutation(() => UploadSimsResDto)
    @UseMiddleware(Authentication)
    async uploadSims(
        @Arg("data") data: UploadSimsInputDto,
        @Ctx() ctx: Context
    ): Promise<UploadSimsResDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.uploadSims(data, parseHeaders(ctx));
    }
}
