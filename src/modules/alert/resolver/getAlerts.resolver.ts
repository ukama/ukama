import { Resolver, Query, Arg } from "type-graphql";
import { Service } from "typedi";
import { AlertsResponse } from "../types";
import { AlertService } from "../service";
import { PaginationDto } from "../../../common/types";

@Service()
@Resolver()
export class GetDataBillResolver {
    constructor(private readonly alertService: AlertService) {}

    @Query(() => AlertsResponse)
    async getAlerts(@Arg("data") data: PaginationDto): Promise<AlertsResponse> {
        return this.alertService.getAlerts(data);
    }
}
