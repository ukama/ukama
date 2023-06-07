import "reflect-metadata";
import { Arg, Ctx, Query, Resolver } from "type-graphql";
import { Context } from "../../common/context";
import { Draft } from "./types";

@Resolver(Draft)
export class DraftResolver {
  @Query(() => Draft)
  async getDraft(@Arg("id") id: number, @Ctx() ctx: Context) {
    const dr = await ctx.prisma.draft.findUnique({
      where: { id: id },
      include: { site: { include: { location: true } }, events: true },
    });
    return dr;
  }
}
