import { IResidentMapper } from "./interface";
import { ResidentDto, ResidentResponse } from "./types";

class ResidentMapper implements IResidentMapper {
    dtoToDto = (res: ResidentResponse): ResidentDto[] => {
        const residents = res.data;
        return residents;
    };
}
export default <IResidentMapper>new ResidentMapper();
