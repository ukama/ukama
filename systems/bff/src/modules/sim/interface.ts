import {
    GetSimAPIResDto,
    SimAPIResDto,
    SimDetailsDto,
    SimDto,
    SimsAPIResDto,
    SimsResDto,
} from "./types";

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface ISimService {}

export interface ISimMapper {
    dtoToSimResDto(res: SimAPIResDto): SimDto;
    dtoToSimDetailsDto(res: GetSimAPIResDto): SimDetailsDto;
    dtoToSimsDto(res: SimsAPIResDto): SimsResDto;
}
