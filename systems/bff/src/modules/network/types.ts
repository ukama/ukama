import { Field, InputType, ObjectType } from "type-graphql";
import { NETWORK_STATUS } from "../../constants";

@ObjectType()
export class SiteAPIDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    network_id: string;

    @Field()
    is_deactivated: string;

    @Field()
    created_at: string;
}

@ObjectType()
export class NetworkStatusDto {
    @Field()
    liveNode: number;

    @Field()
    totalNodes: number;

    @Field(() => NETWORK_STATUS)
    status: NETWORK_STATUS;
}

@ObjectType()
export class NetworkStatusResponse {
    @Field()
    status: string;

    @Field(() => NetworkStatusDto)
    data: NetworkStatusDto;
}

@ObjectType()
export class NetworkAPIDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    org_id: string;

    @Field()
    is_deactivated: string;

    @Field()
    created_at: string;
}
@ObjectType()
export class NetworkAPIResDto {
    @Field(() => NetworkAPIDto)
    network: NetworkAPIDto;
}

@ObjectType()
export class NetworksAPIResDto {
    @Field()
    org_id: string;

    @Field(() => [NetworkAPIDto])
    networks: NetworkAPIDto[];
}

@ObjectType()
export class NetworkDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    orgId: string;

    @Field()
    isDeactivated: string;

    @Field()
    createdAt: string;
}

@ObjectType()
export class NetworksResDto {
    @Field()
    orgId: string;

    @Field(() => [NetworkDto])
    networks: NetworkDto[];
}

@ObjectType()
export class SiteDto {
    @Field()
    id: string;

    @Field()
    name: string;

    @Field()
    networkId: string;

    @Field()
    isDeactivated: string;

    @Field()
    createdAt: string;
}

@ObjectType()
export class SitesResDto {
    @Field()
    networkId: string;

    @Field(() => [SiteDto])
    sites: SiteDto[];
}

@ObjectType()
export class SiteAPIResDto {
    @Field(() => SiteAPIDto)
    site: SiteAPIDto;
}

@ObjectType()
export class SitesAPIResDto {
    @Field()
    network_id: string;

    @Field(() => [SiteAPIDto])
    sites: SiteAPIDto[];
}

@InputType()
export class AddNetworkInputDto {
    @Field()
    network_name: string;

    @Field()
    org: string;
}

@InputType()
export class AddSiteInputDto {
    @Field()
    site: string;
}
