import { IAlertMapper } from "./interface";
import { AlertDto, AlertResponse } from "./types";

class AlertMapper implements IAlertMapper {
    dtoToDto = (res: AlertResponse): AlertDto[] => {
        const alerts: AlertDto[] = [];

        for (let i = 0; i < res.data.length; i++) {
            if (res.data[i]) {
                const alert = {
                    id: res.data[i].id,
                    type: res.data[i].type,
                    title: res.data[i].title,
                    description: res.data[i].description,
                    alertDate: new Date(res.data[i].alertDate),
                };
                alerts.push(alert);
            }
        }

        return alerts;
    };
}
export default <IAlertMapper>new AlertMapper();
