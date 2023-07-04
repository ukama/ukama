import ApiMethods from "../api";
import { axiosErrorHandler } from "../errors";
import { ApiMethodDataDto } from "./types";

export const catchAsyncIOMethod = async (
  req: ApiMethodDataDto
): Promise<any> => {
  try {
    return await ApiMethods.fetch(req);
  } catch (error) {
    return axiosErrorHandler(error);
  }
};
