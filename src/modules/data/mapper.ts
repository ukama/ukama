import { IDataMapper } from "./interface";
import { DataBillDto, DataUsageDto } from "./types";

class DataMapper implements IDataMapper {
    dataUsageDtoToDto = (res: DataUsageDto): DataUsageDto => {
        return {
            id: res.id,
            dataConsumed: res.dataConsumed,
            dataPackage: res.dataPackage,
        };
    };
    dataBillDtoToDto = (res: DataBillDto): DataBillDto => {
        return {
            id: res.id,
            dataBill: res.dataBill,
            billDue: res.billDue,
        };
    };
}
export default <IDataMapper>new DataMapper();
