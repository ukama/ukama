import { Resolver, Query, Arg, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { ResidentService } from "../service";
import { PaginationDto } from "../../../common/types";
import { ResidentsResponse } from "../types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetResidentResolver {
    constructor(private readonly residentService: ResidentService) {}

    @Query(() => ResidentsResponse)
    @UseMiddleware(Authentication)
    async getResidents(
        @Arg("data") data: PaginationDto
    ): Promise<ResidentsResponse> {
        return this.residentService.getResidents(data);
    }
}
