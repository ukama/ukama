import { Service } from "typedi";
import { ResidentsResponse } from "./types";
import { IResidentService } from "./interface";
import { HTTP404Error, Messages } from "../../errors";
import { PaginationDto } from "../../common/types";
import ResidentMapper from "./mapper";
import { getPaginatedOutput } from "../../utils";
import ResidentIOMethods from "./io";

@Service()
export class ResidentService implements IResidentService {
    getResidents = async (req: PaginationDto): Promise<ResidentsResponse> => {
        const res = await ResidentIOMethods.getResidentsMethod(req);
        if (!res) throw new HTTP404Error(Messages.RESIDENTS_NOT_FOUND);
        const meta = getPaginatedOutput(
            req.pageNo,
            req.pageSize,
            res.data.length
        );
        const residents = ResidentMapper.dtoToDto(res.data.data);
        if (!residents) throw new HTTP404Error(Messages.RESIDENTS_NOT_FOUND);

        return {
            residents,
            meta,
        };
    };
}
