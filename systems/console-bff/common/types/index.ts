import { API_METHOD_TYPE } from "../enums";

export type ApiMethodDataDto = {
  url: String;
  method: API_METHOD_TYPE;
  params?: any;
  headers?: any;
  body?: any;
};

export type ErrorType = {
  message: string;
  code: number;
  description?: string;
};
