import { Resolver, Arg, Mutation, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { AddNodeDto, AddNodeResponse } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class AddNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => AddNodeResponse)
    @UseMiddleware(Authentication)
    async addNode(
        @Arg("data")
        req: AddNodeDto
    ): Promise<AddNodeResponse | null> {
        return await this.nodeService.addNode(req);
    }
}
