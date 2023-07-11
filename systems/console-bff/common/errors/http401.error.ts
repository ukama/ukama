import { BaseError } from "./base.error";
import { HttpStatusCode } from "./codes";

export class HTTP401Error extends BaseError {
  constructor(description: string) {
    super("UNAUTHORIZED ", HttpStatusCode.UNAUTHORIZED, description, true);
  }
}
