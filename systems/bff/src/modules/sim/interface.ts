import { SimAPIResDto, SimResDto, SimDetailsDto } from "./types";

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface ISimService {}

export interface ISimMapper {
    dtoToSimResDto(res: SimAPIResDto): SimResDto;
    dtoToSimDetailsDto(res: any): SimDetailsDto;
}
