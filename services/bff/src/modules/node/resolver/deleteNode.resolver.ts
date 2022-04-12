import { Resolver, Arg, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { DeactivateResponse } from "../../user/types";

@Service()
@Resolver()
export class DeleteNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => DeactivateResponse)
    @UseMiddleware(Authentication)
    async deleteNode(
        @Arg("id")
        id: string
    ): Promise<DeactivateResponse | null> {
        return this.nodeService.deleteNode(id);
    }
}
