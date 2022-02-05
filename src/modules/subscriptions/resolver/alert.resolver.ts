import { Resolver, Root, Subscription } from "type-graphql";
import { Service } from "typedi";
import { AlertsResponse, AlertDto } from "../../alert/types";

@Service()
@Resolver()
export class GetAlertsSubscriptionResolver {
    @Subscription(() => AlertDto, {
        topics: "getAlerts",
    })
    async getAlerts(@Root() alerts: AlertsResponse): Promise<AlertDto> {
        return alerts.alerts[alerts.alerts.length - 1];
    }
}
