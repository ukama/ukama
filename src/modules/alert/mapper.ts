import { IAlertMapper } from "./interface";
import { AlertDto, AlertResponse } from "./types";

class AlertMapper implements IAlertMapper {
    dtoToDto = (res: AlertResponse): AlertDto[] => {
        const alerts: AlertDto[] = [];

        for (const alert of res.data) {
            const alertObj = {
                id: alert.id,
                type: alert.type,
                title: alert.title,
                description: alert.description,
                alertDate: new Date(alert.alertDate),
            };
            alerts.push(alertObj);
        }

        return alerts;
    };
}
export default <IAlertMapper>new AlertMapper();
