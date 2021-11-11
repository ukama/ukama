import { Resolver, Arg, Query } from "type-graphql";
import { Service } from "typedi";
import { DataUsageDto } from "../types";
import { DataService } from "../service";
import { TIME_FILTER } from "../../../constants";

@Service()
@Resolver()
export class DataUsageResolver {
    constructor(private readonly dataService: DataService) {}

    @Query(() => DataUsageDto)
    async getDataUsage(
        @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER
    ): Promise<DataUsageDto | null> {
        return await this.dataService.getDataUsage(filter);
    }
}
