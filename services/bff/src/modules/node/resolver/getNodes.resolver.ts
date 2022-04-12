import { Resolver, Query, Arg, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { NodesResponse } from "../types";
import { NodeService } from "../service";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetNodeResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => NodesResponse)
    @UseMiddleware(Authentication)
    async getNodes(@Arg("data") data: PaginationDto): Promise<NodesResponse> {
        return this.nodeService.getNodes(data);
    }
}
