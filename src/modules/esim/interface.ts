import { EsimDto, EsimResponse } from "./types";

export interface IEsimService {
    getEsims(): Promise<EsimDto[]>;
}

export interface IEsimMapper {
    dtoToDto(data: EsimResponse): EsimDto[];
}
