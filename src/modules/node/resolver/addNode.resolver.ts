import { Resolver, Arg, Mutation } from "type-graphql";
import { Service } from "typedi";
import { AddNodeDto, AddNodeResponse } from "../types";
import { NodeService } from "../service";

@Service()
@Resolver()
export class AddNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Mutation(() => AddNodeResponse)
    async addNode(
        @Arg("data")
        req: AddNodeDto
    ): Promise<AddNodeResponse | null> {
        return await this.nodeService.addNode(req);
    }
}
