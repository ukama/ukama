import { AlertDto, AlertResponse } from "../resolver/types";

export const dtoToDto = (res: AlertResponse): AlertDto[] => {
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
