import { Service } from "typedi";
import { EsimResponse, EsimDto } from "./types";
import { IESIMService } from "./interface";
import NodeMapper from "./mapper";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";

@Service()
export class ESIMService implements IESIMService {
    public getEsims = async (): Promise<EsimDto[]> => {
        const res = await catchAsyncIOMethod<EsimResponse>({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_ESIMS,
        });

        const esims = NodeMapper.dtoToDto(res);

        return esims;
    };
}
