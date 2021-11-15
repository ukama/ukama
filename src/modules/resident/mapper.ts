import { IResidentMapper } from "./interface";
import { ResidentDto } from "./types";

class ResidentMapper implements IResidentMapper {
    dtoToDto = (data: ResidentDto[]): ResidentDto[] => {
        const residents: ResidentDto[] = [];

        for (let i = 0; i < data.length; i++) {
            if (data[i]) {
                const alert = {
                    id: data[i].id,
                    name: data[i].name,
                    usage: data[i].usage,
                };
                residents.push(alert);
            }
        }

        return residents;
    };
}
export default <IResidentMapper>new ResidentMapper();
