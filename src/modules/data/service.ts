import { Service } from "typedi";
import { DataBillDto, DataUsageDto } from "./types";
import { IDataService } from "./interface";
import { HTTP404Error, Messages } from "../../errors";
import { DATA_BILL_FILTER, TIME_FILTER } from "../../constants";
import DataMapper from "./mapper";
import DataIOMethods from "./io";

@Service()
export class DataService implements IDataService {
    getDataUsage = async (filter: TIME_FILTER): Promise<DataUsageDto> => {
        const res = await DataIOMethods.getDataUsageMethod(filter);
        if (!res) throw new HTTP404Error(Messages.DATA_NOT_FOUND);

        const data = DataMapper.dataUsageDtoToDto(res.data.data);

        if (!data) throw new HTTP404Error(Messages.DATA_NOT_FOUND);

        return data;
    };
    getDataBill = async (filter: DATA_BILL_FILTER): Promise<DataBillDto> => {
        const res = await DataIOMethods.getDataBillMethod(filter);
        if (!res) throw new HTTP404Error(Messages.DATA_NOT_FOUND);
        const bill = DataMapper.dataBillDtoToDto(res.data.data);

        if (!bill) throw new HTTP404Error(Messages.DATA_NOT_FOUND);

        return bill;
    };
}
