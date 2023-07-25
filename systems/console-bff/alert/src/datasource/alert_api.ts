import { RESTDataSource } from "@apollo/datasource-rest";
import { AlertsResponse } from "../types";
import { PaginationDto } from "../../../common/types";
import AlertMapper from "./mapper";
import { getPaginatedOutput } from "../../../utils";
import { SERVER } from "../../../constants/endpoints";

export class AlertApi extends RESTDataSource {
    getAlerts = async (req: PaginationDto): Promise<AlertsResponse> => {
        return this.get(`${SERVER.GET_ALERTS}`,{
            params: {
                pageNo: `${req.pageNo}`,
                pageSize: `${req.pageSize}`,
            },}).then(res => {
                const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
                const alerts = AlertMapper.dtoToDto(res);
                return {
                    alerts,
                    meta,
                };
            });
    };
}
