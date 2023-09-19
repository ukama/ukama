import { ApiMethodDataDto } from "../types";
import { axiosErrorHandler } from "./../errors/index";
import ApiMethods from "./client";

export const asyncRestCall = async (req: ApiMethodDataDto): Promise<any> => {
  try {
    return await ApiMethods.fetch(req);
  } catch (error) {
    return axiosErrorHandler(error);
  }
};
