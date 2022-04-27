import { Service } from "typedi";
import { DataBillDto, DataUsageDto } from "./types";
import { IDataService } from "./interface";
import { checkError, HTTP404Error, Messages } from "../../errors";
import {
    API_METHOD_TYPE,
    DATA_BILL_FILTER,
    TIME_FILTER,
} from "../../constants";
import DataMapper from "./mapper";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";

@Service()
export class DataService implements IDataService {
    getDataUsage = async (filter: TIME_FILTER): Promise<DataUsageDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_DATA_USAGE,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);

        const data = DataMapper.dataUsageDtoToDto(res);

        if (!data) throw new HTTP404Error(Messages.DATA_NOT_FOUND);

        return data;
    };

    getDataBill = async (filter: DATA_BILL_FILTER): Promise<DataBillDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_DATA_BILL,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);

        const bill = DataMapper.dataBillDtoToDto(res);

        if (!bill) throw new HTTP404Error(Messages.DATA_NOT_FOUND);

        return bill;
    };
}
