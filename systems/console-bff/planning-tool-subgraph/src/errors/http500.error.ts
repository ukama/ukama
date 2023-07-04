import { BaseError } from "./base.error";
import { HttpStatusCode } from "./codes";

export class HTTP500Error extends BaseError {
  constructor(description: string) {
    super("INTERNAL SERVER", HttpStatusCode.INTERNAL_SERVER, description, true);
  }
}
