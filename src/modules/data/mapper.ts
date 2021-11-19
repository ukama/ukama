import { IDataMapper } from "./interface";
import {
    DataBillDto,
    DataBillResponse,
    DataUsageDto,
    DataUsageResponse,
} from "./types";

class DataMapper implements IDataMapper {
    dataUsageDtoToDto = (res: DataUsageResponse): DataUsageDto => {
        const dataUsage = res.data;
        return dataUsage;
    };
    dataBillDtoToDto = (res: DataBillResponse): DataBillDto => {
        const dataBill = res.data;
        return dataBill;
    };
}
export default <IDataMapper>new DataMapper();
