import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { AlertsResponse } from "../../alert/types";

@Service()
@Resolver()
export class GetAlertsSubscriptionResolver {
    @Subscription(() => AlertsResponse, {
        topics: "GETALERTS",
    })
    async GETALERTS(@Root() alerts: AlertsResponse): Promise<AlertsResponse> {
        return alerts;
    }
}
