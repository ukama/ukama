import { Resolver, Arg, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { DeleteNodeResponse } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class DeleteNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => DeleteNodeResponse)
    @UseMiddleware(Authentication)
    async deleteNode(
        @Arg("id")
        id: string
    ): Promise<DeleteNodeResponse | null> {
        return this.nodeService.deleteNode(id);
    }
}
