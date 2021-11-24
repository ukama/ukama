import { Resolver, Query, Arg, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { AlertsResponse } from "../types";
import { AlertService } from "../service";
import { PaginationDto } from "../../../common/types";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetAlertsResolver {
    constructor(private readonly alertService: AlertService) {}

    @Query(() => AlertsResponse)
    @UseMiddleware(Authentication)
    async getAlerts(@Arg("data") data: PaginationDto): Promise<AlertsResponse> {
        return this.alertService.getAlerts(data);
    }
}
