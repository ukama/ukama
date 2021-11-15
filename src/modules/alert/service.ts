import { Service } from "typedi";
import { AlertsResponse } from "./types";
import { IAlertService } from "./interface";
import { HTTP404Error, Messages } from "../../errors";
import { PaginationDto } from "../../common/types";
import AlertMapper from "./mapper";
import { getPaginatedOutput } from "../../utils";
import { getAlertsMethod } from "./io";

@Service()
export class AlertService implements IAlertService {
    getAlerts = async (req: PaginationDto): Promise<AlertsResponse> => {
        const res = await getAlertsMethod(req);
        if (!res) throw new HTTP404Error(Messages.ALERTS_NOT_FOUND);
        const meta = getPaginatedOutput(
            req.pageNo,
            req.pageSize,
            res.data.length
        );
        const alerts = AlertMapper.dtoToDto(res.data.data);
        if (!alerts) throw new HTTP404Error(Messages.ALERTS_NOT_FOUND);
        return {
            alerts,
            meta,
        };
    };
}
