import { IDataMapper } from "./interface";
import { DataUsageDto } from "./types";

class DataMapper implements IDataMapper {
    dataUsageDtoToDto = (res: DataUsageDto): DataUsageDto => {
        return {
            id: res.id,
            dataConsumed: res.dataConsumed,
            dataPackage: res.dataPackage,
        };
    };
}
export default <IDataMapper>new DataMapper();
