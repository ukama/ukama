import "reflect-metadata";
import { Arg, Ctx, Mutation, Query, Resolver } from "type-graphql";
import { Context } from "../../common/context";
import {
  AddDraftInput,
  DeleteDraftRes,
  Draft,
  LocationInput,
  Site,
  SiteInput,
  Event as TEvent,
  Location as TLocation,
  UpdateEventInput,
} from "./types";

@Resolver(Draft)
export class DraftResolver {
  @Query(() => Draft)
  async getDraft(@Arg("id") id: string, @Ctx() ctx: Context) {
    const dr = await ctx.prisma.draft.findUnique({
      where: { id: id },
      include: { sites: { include: { location: true } }, events: true },
    });
    return dr;
  }
  @Query(() => [Draft])
  async getDrafts(@Arg("userId") userId: string, @Ctx() ctx: Context) {
    const dr = await ctx.prisma.draft.findMany({
      where: { userId: userId },
      include: { sites: { include: { location: true } }, events: true },
    });
    return dr;
  }
  @Mutation(() => TEvent)
  async updateEvent(
    @Arg("data") data: UpdateEventInput,
    @Arg("eventId") eventId: string,
    @Ctx() ctx: Context
  ) {
    const dr = await ctx.prisma.event.update({
      where: { id: eventId },
      data: {
        value: data.value,
        operation: data.operation,
        createdAt: new Date().toISOString(),
      },
    });
    return dr;
  }

  @Mutation(() => Site)
  async updateSite(
    @Arg("data") data: SiteInput,
    @Arg("siteId") siteId: string,
    @Ctx() ctx: Context
  ) {
    const dr = await ctx.prisma.site.update({
      where: { id: siteId },
      data: {
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
            lastSaved: data.lastSaved,
          },
        },
      },
      include: { Draft: true, location: true },
    });
    return dr;
  }

  @Mutation(() => Site)
  async addSite(@Arg("data") data: SiteInput, @Ctx() ctx: Context) {
    const dr = await ctx.prisma.site.create({
      data: {
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
            lastSaved: data.lastSaved,
          },
        },
      },
      include: { Draft: true, location: true },
    });
    return dr;
  }

  @Mutation(() => Draft)
  async updateDraftName(
    @Arg("name") name: string,
    @Arg("id") id: string,
    @Ctx() ctx: Context
  ) {
    const dr = await ctx.prisma.draft.update({
      where: { id: id },
      data: {
        name: name,
      },
      include: { sites: { include: { location: true } }, events: true },
    });
    return dr;
  }

  @Mutation(() => Draft)
  async addDraft(@Arg("data") data: AddDraftInput, @Ctx() ctx: Context) {
    const dr = await ctx.prisma.draft.create({
      data: {
        name: data.name,
        userId: data.userId,
        createdAt: new Date().toISOString(),
      },
      include: { sites: { include: { location: true } }, events: true },
    });
    return dr;
  }

  @Mutation(() => TLocation)
  async updateLocation(
    @Arg("locationId") locationId: string,
    @Arg("data") data: LocationInput,
    @Ctx() ctx: Context
  ) {
    const dr = await ctx.prisma.location.update({
      where: { id: locationId },
      data: {
        lat: data.lat,
        lng: data.lng,
        address: data.address,
      },
    });
    return dr;
  }

  @Mutation(() => DeleteDraftRes)
  async deleteDraft(@Arg("id") id: string, @Ctx() ctx: Context) {
    await ctx.prisma.draft.delete({
      where: { id: id },
    });
    return { id: id };
  }
}
