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
import { TIME_FILTER } from "../../common/enums";
import { parseHeaders } from "../../common/utils";
import { Context } from "../context";
import { ConnectedUserDto } from "./types";

@Resolver()
export class GetConnectedUsersResolver {
  @Query(() => ConnectedUserDto)
  @UseMiddleware(Authentication)
  async getConnectedUsers(
    @Arg("filter", () => TIME_FILTER) filter: TIME_FILTER,
    @PubSub() pubsub: PubSubEngine,
    @Ctx() ctx: Context
  ): Promise<ConnectedUserDto> {
    const { dataSources } = ctx;
    const user = dataSources.dataSource.getConnectedUsers(parseHeaders());
    pubsub.publish("getConnectedUsers", user);
    return user;
  }
}
