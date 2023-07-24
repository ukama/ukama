import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { AddOrgInputDto, OrgDto } from "../types";

@Resolver()
export class AddOrgResolver {

    @Mutation(() => OrgDto)
    @UseMiddleware(Authentication)
    async addOrg(
        @Arg("data") data: AddOrgInputDto,
        @Ctx() ctx: Context
    ): Promise<OrgDto> {
        const { dataSources } = ctx;
        return dataSources.dataSource.addOrg(data, parseHeaders(ctx));
    }
}
