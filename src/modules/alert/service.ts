import { Service } from "typedi";
import { AlertsResponse } from "./types";
import { IAlertService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import { PaginationDto } from "../../common/types";
import AlertMapper from "./mapper";
import { getPaginatedOutput } from "../../utils";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";
import { API_METHOD_TYPE } from "../../constants";

@Service()
export class AlertService implements IAlertService {
    getAlerts = async (req: PaginationDto): Promise<AlertsResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_ALERTS,
            params: req,
        });
        if (checkError(res)) throw new Error(res.message);

        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const alerts = AlertMapper.dtoToDto(res);
        if (!alerts) throw new HTTP404Error(Messages.ALERTS_NOT_FOUND);
        return {
            alerts,
            meta,
        };
    };
}
