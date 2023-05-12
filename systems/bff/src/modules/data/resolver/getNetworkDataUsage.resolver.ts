import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { TIME_FILTER } from "../../../constants";
import { DataService } from "../service";
import { DataUsageNetworkResponse } from "../types";

@Service()
@Resolver()
export class GetNetworkDataUsageResolver {
    constructor(private readonly dataService: DataService) {}

    @Query(() => DataUsageNetworkResponse)
    @UseMiddleware(Authentication)
    async getNetworkDataUsage(
        @Arg("networkId") networkId: string,
        @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER,
        @Ctx() ctx: Context
    ): Promise<DataUsageNetworkResponse> {
        return await this.dataService.getNetworkDataUsage(filter);
    }
}
