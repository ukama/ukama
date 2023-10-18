import { BaseError } from "./base.error";
import { HttpStatusCode } from "./codes";

export class HTTP400Error extends BaseError {
  constructor(description: string) {
    super("BAD REQUEST", HttpStatusCode.BAD_REQUEST, description, true);
  }
}
