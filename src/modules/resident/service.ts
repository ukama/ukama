import { Service } from "typedi";
import { ResidentResponse, ResidentsResponse } from "./types";
import { IResidentService } from "./interface";
import { HTTP404Error, Messages } from "../../errors";
import { PaginationDto } from "../../common/types";
import ResidentMapper from "./mapper";
import { getPaginatedOutput } from "../../utils";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";

@Service()
export class ResidentService implements IResidentService {
    getResidents = async (req: PaginationDto): Promise<ResidentsResponse> => {
        const res = await catchAsyncIOMethod<ResidentResponse>({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_RESIDENTS,
            params: req,
        });
        const meta = getPaginatedOutput(req.pageNo, req.pageSize, res.length);
        const residents = ResidentMapper.dtoToDto(res);
        if (!residents) throw new HTTP404Error(Messages.RESIDENTS_NOT_FOUND);

        return {
            residents,
            meta,
        };
    };
}
