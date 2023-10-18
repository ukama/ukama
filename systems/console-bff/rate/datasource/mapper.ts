import {
  DefaultMarkupAPIResDto,
  DefaultMarkupHistoryAPIResDto,
  DefaultMarkupHistoryResDto,
  DefaultMarkupResDto,
} from "../resolver/types";

export const dtoToDefaultMarkupHistoryDto = (
  res: DefaultMarkupHistoryAPIResDto
): DefaultMarkupHistoryResDto => {
  return {
    markupRates: res.markupRates.map(rate => ({
      createdAt: rate.createdAt,
      deletedAt: rate.deletedAt,
      Markup: rate.Markup,
    })),
  };
};

export const dtoToDefaultMarkupDto = (
  res: DefaultMarkupAPIResDto
): DefaultMarkupResDto => {
  return {
    markup: res.markup,
  };
};
