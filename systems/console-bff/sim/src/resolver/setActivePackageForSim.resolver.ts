import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import {
    SetActivePackageForSimInputDto,
    SetActivePackageForSimResDto,
} from "../types";


@Resolver()
export class SetActivePackageForSimResolver {

    @Mutation(() => SetActivePackageForSimResDto)
    @UseMiddleware(Authentication)
    async getSim(
        @Arg("data") data: SetActivePackageForSimInputDto,
        @Ctx() ctx: Context
    ): Promise<SetActivePackageForSimResDto> {
        const { dataSources } = ctx;
        return await dataSources.dataSource.setActivePackageForSim(
            data,
            parseHeaders(ctx)
        );
    }
}
