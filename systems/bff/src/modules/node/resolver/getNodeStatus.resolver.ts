import { Arg, Ctx, Query, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { GetNodeStatusInput, GetNodeStatusRes } from "../../node/types";
import { NodeService } from "../service";

@Service()
@Resolver()
export class GetNodeStatusResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => GetNodeStatusRes)
    @UseMiddleware(Authentication)
    async getNodeStatus(
        @Arg("data") data: GetNodeStatusInput,
        @Ctx() ctx: Context
    ): Promise<GetNodeStatusRes> {
        return this.nodeService.getNodeStatus(data, parseHeaders(ctx));
    }
}
