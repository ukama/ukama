import { IEsimMapper } from "./interface";
import { EsimDto, EsimResponse } from "./types";

class EsimMapper implements IEsimMapper {
    dtoToDto = (res: EsimResponse): EsimDto[] => {
        const esims = res.data;
        return esims;
    };
}
export default <IEsimMapper>new EsimMapper();
