import "reflect-metadata";
import { Arg, Ctx, Mutation, Query, Resolver } from "type-graphql";

import { asyncRestCall } from "../../../common/axiosClient";
import { PLANNING_API_URL, PLANNING_BUCKET } from "../../../common/configs";
import { API_METHOD_TYPE } from "../../../common/enums";
import { PrismaContext } from "../../../common/prisma";
import {
  AddDraftInput,
  CoverageInput,
  CoverageRes,
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
  async getDraft(@Arg("id") id: string, @Ctx() ctx: PrismaContext) {
    const dr = await ctx.prisma.draft.findUnique({
      where: { id: id },
      include: {
        links: true,
        sites: { include: { location: true } },
        events: true,
      },
    });
    return dr;
  }
  @Query(() => [Draft])
  async getDrafts(@Arg("userId") userId: string, @Ctx() ctx: PrismaContext) {
    const dr = await ctx.prisma.draft.findMany({
      where: { userId: userId },
      include: {
        links: true,
        sites: { include: { location: true } },
        events: true,
      },
    });
    return dr;
  }
  @Mutation(() => TEvent)
  async updateEvent(
    @Arg("data") data: UpdateEventInput,
    @Arg("eventId") eventId: string,
    @Ctx() ctx: PrismaContext
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
    @Ctx() ctx: PrismaContext
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
        links: true,
        sites: { include: { location: true } },
        events: true,
      },
    });
    return dr;
  }

  @Mutation(() => Draft)
  async addSite(
    @Arg("draftId") draftId: string,
    @Arg("data") data: SiteInput,
    @Ctx() ctx: PrismaContext
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
              url: "",
              east: 0,
              west: 0,
              north: 0,
              south: 0,
              populationUrl: "",
              populationCovered: 0,
              totalBoxesCovered: 0,
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
        links: true,
        sites: { include: { location: true } },
        events: true,
      },
    });
    return dr;
  }

  @Mutation(() => Draft)
  async updateDraftName(
    @Arg("name") name: string,
    @Arg("id") id: string,
    @Ctx() ctx: PrismaContext
  ) {
    const dr = await ctx.prisma.draft.update({
      where: { id: id },
      data: {
        name: name,
      },
      include: {
        links: true,
        sites: { include: { location: true } },
        events: true,
      },
    });
    return dr;
  }

  @Mutation(() => Draft)
  async addDraft(@Arg("data") data: AddDraftInput, @Ctx() ctx: PrismaContext) {
    const dr = await ctx.prisma.draft.create({
      data: {
        name: data.name,
        userId: data.userId,
        lastSaved: data.lastSaved,
        createdAt: new Date().toISOString(),
      },
      include: {
        links: true,
        sites: { include: { location: true } },
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
    @Ctx() ctx: PrismaContext
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
  async deleteDraft(@Arg("id") id: string, @Ctx() ctx: PrismaContext) {
    await ctx.prisma.draft.delete({
      where: { id: id },
    });
    return { id: id };
  }

  @Mutation(() => Draft)
  async addLink(
    @Arg("draftId") draftId: string,
    @Arg("data") data: LinkInput,
    @Ctx() ctx: PrismaContext
  ) {
    const l = await ctx.prisma.draft.update({
      where: { id: draftId },
      data: {
        lastSaved: data.lastSaved,
        links: {
          create: {
            siteA: data.siteA,
            siteB: data.siteB,
          },
        },
      },
      include: {
        links: true,
        sites: { include: { location: true } },
        events: true,
      },
    });
    return l;
  }

  @Mutation(() => DeleteLinkRes)
  async deleteLink(
    @Arg("linkId") linkId: string,
    @Arg("draftId") draftId: string,
    @Arg("lastSaved") lastSaved: number,
    @Ctx() ctx: PrismaContext
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
  async deleteSite(@Arg("id") id: string, @Ctx() ctx: PrismaContext) {
    await ctx.prisma.site.delete({
      where: { id: id },
    });
    return { id: id };
  }

  @Mutation(() => Site)
  async coverage(
    @Arg("data") data: CoverageInput,
    @Arg("siteId") siteId: string,
    @Ctx() ctx: PrismaContext
  ) {
    const config = {
      method: API_METHOD_TYPE.POST,
      url: `${PLANNING_API_URL}/coverage`,
      headers: {
        "Content-Type": "application/json",
      },
      data: JSON.stringify({
        mode: data.mode,
        sites: [
          {
            latitude: data.lat,
            longitude: data.lng,
            transmitter_height: data.height,
          },
        ],
      }),
    };

    const res = await asyncRestCall(config);
    const index = `lat${data.lat.toString().replace(".", "_")}lon${data.lng
      .toString()
      .replace(".", "_")}`;
    const c: CoverageRes = {
      east: res.data.east,
      west: res.data.west,
      north: res.data.north,
      south: res.data.south,
      url: `https://${PLANNING_BUCKET}.s3.amazonaws.com/${res.data.url}`,
      populationData: {
        populationCovered: res.data.population_data[index].population_covered,
        totalBoxesCovered: res.data.population_data[index].total_boxes_covered,
        url: `https://${PLANNING_BUCKET}.s3.amazonaws.com/${res.data.population_data[index].url}`,
      },
    };
    const m = await ctx.prisma.site.update({
      where: { id: siteId },
      data: {
        url: c.url,
        east: c.east,
        west: c.west,
        north: c.north,
        south: c.south,
        populationUrl: c.populationData.url,
        populationCovered: c.populationData.populationCovered,
        totalBoxesCovered: c.populationData.totalBoxesCovered,
      },
      include: {
        location: true,
      },
    });

    return m;
  }
}
