import { BaseError } from "./base.error";
import { HttpStatusCode } from "./codes";

export class HTTP404Error extends BaseError {
  constructor(description: string) {
    super("NOT_FOUND ", HttpStatusCode.NOT_FOUND, description, true);
  }
}
