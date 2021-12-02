import {
    Resolver,
    Arg,
    Query,
    UseMiddleware,
    PubSub,
    PubSubEngine,
} from "type-graphql";
import { Service } from "typedi";
import { DataUsageDto } from "../types";
import { DataService } from "../service";
import { TIME_FILTER } from "../../../constants";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class DataUsageResolver {
    constructor(private readonly dataService: DataService) {}

    @Query(() => DataUsageDto)
    @UseMiddleware(Authentication)
    async getDataUsage(
        @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER,
        @PubSub() pubsub: PubSubEngine
    ): Promise<DataUsageDto | null> {
        const data = await this.dataService.getDataUsage(filter);
        await pubsub.publish("DATAUSAGE", data);
        return data;
    }
}
