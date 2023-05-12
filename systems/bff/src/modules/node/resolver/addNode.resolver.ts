import { Service } from "typedi";
import { NodeService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";
import { AddNodeDto, AddNodeResponse } from "../types";

@Service()
@Resolver()
export class AddNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => AddNodeResponse)
    @UseMiddleware(Authentication)
    async addNode(
        @Arg("data")
        req: AddNodeDto,
        @Ctx() ctx: Context
    ): Promise<AddNodeResponse> {
        return await this.nodeService.addNode(req, parseCookie(ctx));
    }
}
