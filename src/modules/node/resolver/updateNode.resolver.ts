import { Resolver, Arg, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { UpdateNodeDto, UpdateNodeResponse } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class UpdateNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => UpdateNodeResponse)
    @UseMiddleware(Authentication)
    async updateNode(
        @Arg("data")
        req: UpdateNodeDto
    ): Promise<UpdateNodeResponse | null> {
        return await this.nodeService.updateNode(req);
    }
}
