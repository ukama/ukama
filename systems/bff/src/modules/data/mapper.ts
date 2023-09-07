import { IDataMapper } from "./interface";
import {
    DataBillDto,
    DataBillResponse,
    DataUsageDto,
    DataUsageResponse,
} from "./types";

class DataMapper implements IDataMapper {
    dataUsageDtoToDto = (res: DataUsageResponse): DataUsageDto => {
        return res.data;
    };
    dataBillDtoToDto = (res: DataBillResponse): DataBillDto => {
        return res.data;
    };
}
export default <IDataMapper>new DataMapper();
