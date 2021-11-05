import { Resolver, Query } from "type-graphql";
import { Service } from "typedi";

@Service()
@Resolver()
export class TestQueryResolver {
    @Query(() => String)
    async testQuery(): Promise<string> {
        return "Hello World";
    }
}
