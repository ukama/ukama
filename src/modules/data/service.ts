import { Service } from "typedi";
import { DataUsageDto } from "./types";
import { IDataService } from "./interface";
import { HTTP404Error, Messages } from "../../errors";
import { TIME_FILTER } from "../../constants";
import DataMapper from "./mapper";
import { SERVER } from "../../constants/endpoints";
import { getDataUsageMethod } from "./io";

@Service()
export class DataService implements IDataService {
    getDataUsage = async (filter: TIME_FILTER): Promise<DataUsageDto> => {
        const res = await getDataUsageMethod(
            SERVER.GET_DATA_USAGE,
            `${filter}`,
            null
        );
        if (!res) throw new HTTP404Error(Messages.DATA_NOT_FOUND);

        const data = DataMapper.dataUsageDtoToDto(res.data.data);

        if (!data) throw new HTTP404Error(Messages.DATA_NOT_FOUND);

        return data;
    };
}
