import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { DeleteNodeRes } from "../../user/types";
import { NodeService } from "../service";

@Service()
@Resolver()
export class DeleteNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => DeleteNodeRes)
    @UseMiddleware(Authentication)
    async deleteNode(
        @Arg("id")
        id: string,
        @Ctx() ctx: Context
    ): Promise<DeleteNodeRes | null> {
        return this.nodeService.deleteNode(id, parseHeaders(ctx));
    }
}
