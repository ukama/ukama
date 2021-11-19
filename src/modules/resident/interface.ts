import { PaginationDto } from "../../common/types";
import { ResidentDto, ResidentResponse, ResidentsResponse } from "./types";

export interface IResidentService {
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
}

export interface IResidentMapper {
    dtoToDto(res: ResidentResponse): ResidentDto[];
}
