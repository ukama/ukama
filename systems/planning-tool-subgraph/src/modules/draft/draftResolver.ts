import "reflect-metadata";
import { Arg, Ctx, Mutation, Query, Resolver } from "type-graphql";
import { Context } from "../../common/context";
import {
  AddDraftInput,
  Draft,
  UpdateEventInput,
  UpdateSiteInput,
} from "./types";

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
  @Query(() => [Draft])
  async getDrafts(@Arg("userId") userId: string, @Ctx() ctx: Context) {
    const dr = await ctx.prisma.draft.findMany({
      where: { userId: userId },
      include: { site: { include: { location: true } }, events: true },
    });
    return dr;
  }
  @Mutation(() => Draft)
  async updateEvent(
    @Arg("data") data: UpdateEventInput,
    @Arg("draftId") draftId: number,
    @Ctx() ctx: Context
  ) {
    const dr = await ctx.prisma.draft.update({
      where: { id: draftId },
      data: {
        events: {
          create: {
            value: data.value,
            operation: data.operation,
            createdAt: new Date().toISOString(),
          },
        },
      },
      include: { site: { include: { location: true } }, events: true },
    });
    return dr;
  }

  @Mutation(() => Draft)
  async updateSite(
    @Arg("data") data: UpdateSiteInput,
    @Arg("id") id: number,
    @Ctx() ctx: Context
  ) {
    const dr = await ctx.prisma.draft.update({
      where: { id: id },
      data: {
        lastSaved: data.lastSaved,
        updatedAt: new Date().toISOString(),
        site: {
          create: {
            name: data.siteName,
            height: data.height,
            apOption: data.apOption,
            solarUptime: data.solarUptime,
            isSetlite: data.isSetlite,
            location: {
              create: {
                lat: data.lat,
                lng: data.lng,
                address: data.address,
              },
            },
          },
        },
      },
      include: { site: { include: { location: true } }, events: true },
    });
    return dr;
  }

  @Mutation(() => Draft)
  async updateDraftName(
    @Arg("name") name: string,
    @Arg("id") id: number,
    @Ctx() ctx: Context
  ) {
    const dr = await ctx.prisma.draft.update({
      where: { id: id },
      data: {
        name: name,
      },
      include: { site: { include: { location: true } }, events: true },
    });
    return dr;
  }

  @Mutation(() => Draft)
  async addDraft(@Arg("data") data: AddDraftInput, @Ctx() ctx: Context) {
    const dr = await ctx.prisma.draft.create({
      data: {
        name: data.name,
        userId: data.userId,
        lastSaved: data.lastSaved,
        createdAt: new Date().toISOString(),
        site: {
          create: {
            name: "",
            height: 0,
            apOption: "",
            solarUptime: 0,
            isSetlite: false,
            location: {
              create: {
                lat: "",
                lng: "",
                address: "",
              },
            },
          },
        },
      },
      include: { site: { include: { location: true } }, events: true },
    });
    return dr;
  }
}
