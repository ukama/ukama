import { TIME_FILTER } from "../../constants";
import { DataUsageDto } from "./types";

export interface IDataService {
    getDataUsage(filter: TIME_FILTER): Promise<DataUsageDto>;
}

export interface IDataMapper {
    dataUsageDtoToDto(res: DataUsageDto): DataUsageDto;
}
