import { Arg, Query, Resolver } from "type-graphql";

import { getIdentity } from "../../common/auth/authCalls";
import { WhoamiDto } from "./types";

@Resolver()
export class WhoamiResolver {
  @Query(() => WhoamiDto)
  async whoami(@Arg("userId") userId: string): Promise<WhoamiDto> {
    return await getIdentity(userId, "");
  }
}
