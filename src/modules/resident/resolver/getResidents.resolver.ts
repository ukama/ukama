import { Resolver, Query, Arg } from "type-graphql";
import { Service } from "typedi";
import { ResidentService } from "../service";
import { PaginationDto } from "../../../common/types";
import { ResidentsResponse } from "../types";

@Service()
@Resolver()
export class GetDataBillResolver {
    constructor(private readonly residentService: ResidentService) {}

    @Query(() => ResidentsResponse)
    async getResidents(
        @Arg("data") data: PaginationDto
    ): Promise<ResidentsResponse> {
        return this.residentService.getResidents(data);
    }
}
