import { Arg, Mutation, Resolver } from "type-graphql";

import { updateAttributes } from "./../../common/auth/authCalls";
import { UserFistVisitInputDto, UserFistVisitResDto } from "./types";

@Resolver()
export class updateFirstVisitResolver {
  @Mutation(() => UserFistVisitResDto)
  async updateFirstVisit(
    @Arg("data") data: UserFistVisitInputDto
  ): Promise<UserFistVisitResDto> {
    const user = await updateAttributes(
      data.userId,
      data.email,
      data.name,
      "",
      data.firstVisit
    );
    return user;
  }
}
