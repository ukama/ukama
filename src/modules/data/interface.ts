import { DATA_BILL_FILTER, TIME_FILTER } from "../../constants";
import { DataBillDto, DataUsageDto } from "./types";

export interface IDataService {
    getDataUsage(filter: TIME_FILTER): Promise<DataUsageDto>;
    getDataBill(filter: DATA_BILL_FILTER): Promise<DataBillDto>;
}

export interface IDataMapper {
    dataUsageDtoToDto(res: DataUsageDto): DataUsageDto;
    dataBillDtoToDto(res: DataBillDto): DataBillDto;
}
