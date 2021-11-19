import { IESIMMapper } from "./interface";
import { EsimDto, EsimResponse } from "./types";

class ESIMMapper implements IESIMMapper {
    dtoToDto = (res: EsimResponse): EsimDto[] => {
        const esims = res.data;
        return esims;
    };
}
export default <IESIMMapper>new ESIMMapper();
