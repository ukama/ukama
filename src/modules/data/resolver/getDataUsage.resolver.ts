import { Resolver, Arg, Query, UseMiddleware } from "type-graphql";
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
        @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER
    ): Promise<DataUsageDto | null> {
        return await this.dataService.getDataUsage(filter);
    }
}
