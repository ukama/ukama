import { IAlertMapper } from "./interface";
import { AlertDto } from "./types";

class AlertMapper implements IAlertMapper {
    dtoToDto = (data: AlertDto[]): AlertDto[] => {
        const alerts: AlertDto[] = [];

        for (let i = 0; i < data.length; i++) {
            if (data[i]) {
                const alert = {
                    id: data[i].id,
                    type: data[i].type,
                    title: data[i].title,
                    description: data[i].description,
                    alertDate: new Date(data[i].alertDate),
                };
                alerts.push(alert);
            }
        }

        return alerts;
    };
}
export default <IAlertMapper>new AlertMapper();
