import { Service } from "typedi";
import { NodeService } from "../service";
import { parseCookie } from "../../../common";
import { Context } from "../../../common/types";
import { DeleteNodeRes } from "../../user/types";
import { Authentication } from "../../../common/Authentication";
import { Resolver, Arg, Mutation, UseMiddleware, Ctx } from "type-graphql";

@Service()
@Resolver()
export class DeleteNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => DeleteNodeRes)
    @UseMiddleware(Authentication)
    async deleteNode(
        @Arg("id")
        id: string,
        @Ctx() ctx: Context,
    ): Promise<DeleteNodeRes | null> {
        return this.nodeService.deleteNode(id, parseCookie(ctx));
    }
}
