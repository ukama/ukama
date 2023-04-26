import { IRateMapper } from "./interface";
import {
    DefaultMarkupAPIResDto,
    DefaultMarkupHistoryAPIResDto,
    DefaultMarkupHistoryResDto,
    DefaultMarkupResDto,
} from "./types";

class RateMapper implements IRateMapper {
    dtoToDefaultMarkupHistoryDto(
        res: DefaultMarkupHistoryAPIResDto,
    ): DefaultMarkupHistoryResDto {
        return {
            markupRates: res.markupRates.map(rate => ({
                createdAt: rate.createdAt,
                deletedAt: rate.deletedAt,
                Markup: rate.Markup,
            })),
        };
    }
    dtoToDefaultMarkupDto = (
        res: DefaultMarkupAPIResDto,
    ): DefaultMarkupResDto => {
        return {
            markup: res.markup,
        };
    };
}
export default <IRateMapper>new RateMapper();
