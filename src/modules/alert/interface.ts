import { PaginationDto } from "../../common/types";
import { AlertDto, AlertsResponse } from "./types";

export interface IAlertService {
    getAlerts(req: PaginationDto): Promise<AlertsResponse>;
}

export interface IAlertMapper {
    dtoToDto(data: AlertDto[]): AlertDto[];
}
