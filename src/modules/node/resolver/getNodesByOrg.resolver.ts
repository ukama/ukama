import { Resolver, Query, UseMiddleware, Arg } from "type-graphql";
import { Service } from "typedi";
import { OrgNodeResponse } from "../types";
import { NodeService } from "../service";
import { Authentication } from "../../../common/Authentication";
import { LoginDto } from "../../auth/types";

@Service()
@Resolver()
export class GetNodesByOrgResolver {
    constructor(private readonly nodeService: NodeService) {}

    @Query(() => OrgNodeResponse)
    @UseMiddleware(Authentication)
    async getNodesByOrg(@Arg("data") data: LoginDto): Promise<OrgNodeResponse> {
        return this.nodeService.getNodesByOrg(data);
    }
}
