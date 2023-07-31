import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context,IdResponse } from "../../../common/types";

@Resolver()
export class DeletePackageResolver {

    @Mutation(() => IdResponse)
    @UseMiddleware(Authentication)
    async deletePackage(
        @Arg("packageId") packageId: string,
        @Ctx() ctx: Context
    ): Promise<IdResponse> {
        const { dataSources } = ctx;
        return dataSources.dataSource.deletePackage(packageId, parseHeaders(ctx));
    }
}
