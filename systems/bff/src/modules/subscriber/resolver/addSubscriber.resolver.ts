import { Arg, Ctx, Mutation, Resolver, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { parseHeaders } from "../../../common";
import { Authentication } from "../../../common/Authentication";
import { Context } from "../../../common/types";
import { SubscriberService } from "../service";
import { SubscriberDto, SubscriberInputDto } from "../types";

@Service()
@Resolver()
export class AddSubscriberResolver {
    constructor(private readonly userService: SubscriberService) {}

    @Mutation(() => SubscriberDto)
    @UseMiddleware(Authentication)
    async addSubscriber(
        @Arg("data") data: SubscriberInputDto,
        @Ctx() ctx: Context
    ): Promise<SubscriberDto> {
        return await this.userService.addSubscriber(data, parseHeaders(ctx));
    }
}
