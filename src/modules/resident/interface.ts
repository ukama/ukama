import { PaginationDto } from "../../common/types";
import { ResidentDto, ResidentsResponse } from "./types";

export interface IResidentService {
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
}

export interface IResidentMapper {
    dtoToDto(data: ResidentDto[]): ResidentDto[];
}
