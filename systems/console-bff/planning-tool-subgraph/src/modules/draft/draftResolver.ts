import "reflect-metadata";
import { Arg, Ctx, Mutation, Query, Resolver } from "type-graphql";
import { Context } from "../../common/context";
import {
  AddDraftInput,
  DeleteDraftRes,
  DeleteLinkRes,
  DeleteSiteRes,
  Draft,
  LinkInput,
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
      include: {
        sites: { include: { location: true, links: true } },
        events: true,
      },
    });
    return dr;
  }
  @Query(() => [Draft])
  async getDrafts(@Arg("userId") userId: string, @Ctx() ctx: Context) {
    const dr = await ctx.prisma.draft.findMany({
      where: { userId: userId },
      include: {
        sites: { include: { location: true, links: true } },
        events: true,
      },
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

  @Mutation(() => Draft)
  async updateSite(
    @Arg("data") data: SiteInput,
    @Arg("siteId") siteId: string,
    @Arg("draftId") draftId: string,
    @Ctx() ctx: Context
  ) {
    const dr = await ctx.prisma.draft.update({
      where: { id: draftId },
      data: {
        sites: {
          update: {
            where: {
              id: siteId,
            },
            data: {
              name: data.siteName,
              height: data.height,
              apOption: data.apOption,
              solarUptime: data.solarUptime,
              isSetlite: data.isSetlite,
              location: {
                update: {
                  lat: data.lat,
                  lng: data.lng,
                  address: data.address,
                },
              },
            },
          },
        },
      },
      include: {
        sites: { include: { location: true, links: true } },
        events: true,
      },
    });
    return dr;
  }

  @Mutation(() => Draft)
  async addSite(
    @Arg("draftId") draftId: string,
    @Arg("data") data: SiteInput,
    @Ctx() ctx: Context
  ) {
    const dr = await ctx.prisma.draft.update({
      where: {
        id: draftId,
      },
      data: {
        lastSaved: data.lastSaved,
        sites: {
          create: [
            {
              status: "up",
              name: data.siteName,
              height: data.height,
              apOption: data.apOption,
              isSetlite: data.isSetlite,
              solarUptime: data.solarUptime,
              location: {
                create: {
                  id: data.locationId,
                  lat: data.lat,
                  lng: data.lng,
                  address: data.address,
                },
              },
            },
          ],
        },
      },
      include: {
        sites: { include: { location: true, links: true } },
        events: true,
      },
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
      include: {
        sites: { include: { location: true, links: true } },
        events: true,
      },
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
      },
      include: {
        sites: { include: { location: true, links: true } },
        events: true,
      },
    });
    return dr;
  }

  @Mutation(() => TLocation)
  async updateLocation(
    @Arg("draftId") draftId: string,
    @Arg("data") data: LocationInput,
    @Arg("locationId") locationId: string,
    @Ctx() ctx: Context
  ) {
    const l = await ctx.prisma.location.update({
      where: { id: locationId },
      data: {
        id: locationId,
        lat: data.lat,
        lng: data.lng,
        address: data.address,
      },
    });
    await ctx.prisma.draft.update({
      where: {
        id: draftId,
      },
      data: {
        lastSaved: data.lastSaved,
      },
    });
    return l;
  }

  @Mutation(() => DeleteDraftRes)
  async deleteDraft(@Arg("id") id: string, @Ctx() ctx: Context) {
    await ctx.prisma.draft.delete({
      where: { id: id },
    });
    return { id: id };
  }

  @Mutation(() => Site)
  async addLink(
    @Arg("draftId") draftId: string,
    @Arg("siteId") siteId: string,
    @Arg("data") data: LinkInput,
    @Ctx() ctx: Context
  ) {
    const l = await ctx.prisma.site.update({
      where: { id: siteId },
      data: {
        links: {
          create: {
            data: data.data,
            linkWith: data.linkWith,
          },
        },
      },
      include: { location: true, links: true },
    });
    await ctx.prisma.draft.update({
      where: {
        id: draftId,
      },
      data: {
        lastSaved: data.lastSaved,
      },
    });
    return l;
  }

  @Mutation(() => DeleteLinkRes)
  async deleteLink(
    @Arg("linkId") linkId: string,
    @Arg("draftId") draftId: string,
    @Arg("lastSaved") lastSaved: number,
    @Ctx() ctx: Context
  ) {
    await ctx.prisma.link.delete({
      where: { id: linkId },
    });
    await ctx.prisma.draft.update({
      where: {
        id: draftId,
      },
      data: {
        lastSaved: lastSaved,
      },
    });
    return { id: linkId };
  }

  @Mutation(() => DeleteSiteRes)
  async deleteSite(@Arg("id") id: string, @Ctx() ctx: Context) {
    await ctx.prisma.site.delete({
      where: { id: id },
    });
    return { id: id };
  }
}
