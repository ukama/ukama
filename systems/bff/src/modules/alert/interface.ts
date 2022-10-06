import { PaginationDto } from "../../common/types";
import { AlertDto, AlertResponse, AlertsResponse } from "./types";

export interface IAlertService {
    getAlerts(req: PaginationDto): Promise<AlertsResponse>;
}

export interface IAlertMapper {
    dtoToDto(res: AlertResponse): AlertDto[];
}
