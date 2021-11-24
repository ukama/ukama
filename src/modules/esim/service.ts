import { Service } from "typedi";
import { EsimDto } from "./types";
import { IEsimService } from "./interface";
import EsimMapper from "./mapper";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError } from "../../errors";

@Service()
export class EsimService implements IEsimService {
    public getEsims = async (): Promise<EsimDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_ESIMS,
        });

        if (checkError(res)) throw new Error(res.message);

        const esims = EsimMapper.dtoToDto(res);

        return esims;
    };
}
