import { Service } from "typedi";
import { NodeService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, Query, UseMiddleware, Arg, Ctx } from "type-graphql";
import { GetNodeStatusInput, GetNodeStatusRes } from "../../node/types";

@Service()
@Resolver()
export class GetNodeStatusResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => GetNodeStatusRes)
    @UseMiddleware(Authentication)
    async getNodeStatus(
        @Arg("data") data: GetNodeStatusInput,
        @Ctx() ctx: Context,
    ): Promise<GetNodeStatusRes> {
        return this.nodeService.getNodeStatus(data, parseCookie(ctx));
    }
}
