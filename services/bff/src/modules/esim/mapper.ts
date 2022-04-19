import { IEsimMapper } from "./interface";
import { EsimDto, EsimResponse } from "./types";

class EsimMapper implements IEsimMapper {
    dtoToDto = (res: EsimResponse): EsimDto[] => {
        return res.data;
    };
}
export default <IEsimMapper>new EsimMapper();
