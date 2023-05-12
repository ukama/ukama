import {
    DefaultMarkupAPIResDto,
    DefaultMarkupHistoryAPIResDto,
    DefaultMarkupHistoryResDto,
    DefaultMarkupResDto,
} from "./types";

// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface IRateService {}

export interface IRateMapper {
    dtoToDefaultMarkupDto(res: DefaultMarkupAPIResDto): DefaultMarkupResDto;
    dtoToDefaultMarkupHistoryDto(
        res: DefaultMarkupHistoryAPIResDto
    ): DefaultMarkupHistoryResDto;
}
