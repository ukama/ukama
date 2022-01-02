import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { NodeMetaDataDto } from "../../node/types";

@Service()
@Resolver()
export class GetNodeMetaDataSubscriptionResolver {
    @Subscription(() => NodeMetaDataDto, {
        topics: "getNodeMetaData",
    })
    async getNodeMetaData(
        @Root() nodeMeta: NodeMetaDataDto
    ): Promise<NodeMetaDataDto> {
        return nodeMeta;
    }
}
