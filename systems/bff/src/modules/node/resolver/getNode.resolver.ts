import { Service } from "typedi";
import { NodeResponse } from "../types";
import { NodeService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, Query, UseMiddleware, Arg, Ctx } from "type-graphql";

@Service()
@Resolver()
export class GetNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodeResponse)
    @UseMiddleware(Authentication)
    async getNode(
        @Arg("nodeId")
        nodeId: string,
        @Ctx() ctx: Context
    ): Promise<NodeResponse> {
        return this.nodeService.getNode(nodeId, parseCookie(ctx));
    }
}
