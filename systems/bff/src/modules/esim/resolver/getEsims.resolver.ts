import { Resolver, Query, UseMiddleware } from "type-graphql";
import { Service } from "typedi";
import { EsimDto } from "../types";
import { EsimService } from "../service";
import { Authentication } from "../../../common/Authentication";

@Service()
@Resolver()
export class GetEsimResolver {
    constructor(private readonly esimService: EsimService) {}

    @Query(() => [EsimDto])
    @UseMiddleware(Authentication)
    async getEsims(): Promise<EsimDto[]> {
        return this.esimService.getEsims();
    }
}
