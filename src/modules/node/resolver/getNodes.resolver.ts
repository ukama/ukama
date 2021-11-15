import { Resolver, Query, Arg } from "type-graphql";
import { Service } from "typedi";
import { NodesResponse } from "../types";
import { NodeService } from "../service";
import { PaginationDto } from "../../../common/types";

@Service()
@Resolver()
export class GetNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodesResponse)
    async getNodes(@Arg("data") data: PaginationDto): Promise<NodesResponse> {
        return this.nodeService.getNodes(data);
    }
}
