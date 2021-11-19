import { EsimDto, EsimResponse } from "./types";

export interface IESIMService {
    getEsims(): Promise<EsimDto[]>;
}

export interface IESIMMapper {
    dtoToDto(data: EsimResponse): EsimDto[];
}
