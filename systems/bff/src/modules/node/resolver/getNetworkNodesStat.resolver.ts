import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NodeService } from "../service";
import { NodeStatsResponse } from "../types";

@Service()
@Resolver()
export class GetNetworkNodesStatResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodeStatsResponse)
    @UseMiddleware(Authentication)
    async getNetworkNodesStat(
        @Arg("networkId")
        networkId: string,
        @Ctx() ctx: Context
    ): Promise<NodeStatsResponse> {
        return this.nodeService.getNodesStats();
    }
}
