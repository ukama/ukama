import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { NodeService } from "../service";
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
        return await this.nodeService.addNode(req, parseHeaders(ctx));
    }
}
