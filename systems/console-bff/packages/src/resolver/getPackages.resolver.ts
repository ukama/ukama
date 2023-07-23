import { Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { PackagesResDto } from "../types";

@Service()
@Resolver()
export class GetPackagesResolver {
    @Query(() => PackagesResDto)
    @UseMiddleware(Authentication)
    async getPackages(@Ctx() ctx: Context): Promise<PackagesResDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.getPackages(parseHeaders(ctx));
    }
}
