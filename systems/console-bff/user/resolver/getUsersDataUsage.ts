import {
  Arg,
  Ctx,
  PubSub,
  PubSubEngine,
  Query,
  Resolver,
  UseMiddleware,
} from "type-graphql";

import { Authentication } from "../../common/auth";
import { Context } from "../context";
import { DataUsageInputDto, GetUserDto } from "./types";

@Resolver()
export class GetUsersDataUsageResolver {
  @Query(() => [GetUserDto])
  @UseMiddleware(Authentication)
  async getUsersDataUsage(
    @Arg("data") data: DataUsageInputDto,
    @PubSub() pubsub: PubSubEngine,
    @Ctx() ctx: Context
  ): Promise<GetUserDto[]> {
    const { dataSources } = ctx;

    const users: GetUserDto[] = [];
    if (data.ids.length > 0) {
      for (let i = 0; i < data.ids.length; i++) {
        const user = await dataSources.dataSource.getUser(data.ids[i]);
        pubsub.publish("getUsersSub", user);
        // users.push(user);
      }
    }
    return users;
  }
}
