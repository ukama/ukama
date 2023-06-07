import "reflect-metadata";
import { Arg, Ctx, Query, Resolver } from "type-graphql";
import { Context } from "../../common/context";
import { APOptions } from "../../common/enums";
import { Draft } from "./types";

@Resolver(Draft)
export class DraftResolver {
  @Query(() => Draft)
  async getDraft(@Arg("id") id: string, @Ctx() ctx: Context): Promise<Draft> {
    return {
      id: id,
      name: "Draft 1",
      lastSaved: 1686084992,
      site: {
        name: "My Site",
        height: 0,
        apOption: APOptions.ONE_TO_ONE,
        solarUptime: 90,
        isSetlite: false,
        location: {
          lat: "0.123",
          lng: "0.456",
          address: "123 Main Street",
        },
      },
      events: [
        {
          id: "event-id",
          operation: "attributeName",
          value: "attributeValue",
        },
      ],
    };
  }
}
