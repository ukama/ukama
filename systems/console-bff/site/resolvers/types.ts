import { Field, InputType, ObjectType } from "type-graphql";

@ObjectType()
export class SiteDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  networkId: string;

  @Field()
  backhaulId: string;

  @Field()
  powerId: string;

  @Field()
  accessId: string;

  @Field()
  spectrumId: string;

  @Field()
  switchId: string;

  @Field()
  isDeactivated: boolean;

  @Field()
  latitude: string;

  @Field()
  longitude: string;

  @Field()
  installDate: string;

  @Field()
  createdAt: string;

  @Field()
  location: string;
}

@ObjectType()
export class SitesResDto {
  @Field(() => [SiteDto])
  sites: SiteDto[];
}

@ObjectType()
export class SiteAPIDto {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field()
  network_id: string;

  @Field()
  backhaul_id: string;

  @Field()
  category: string;

  @Field()
  spectrum_id: string;

  @Field()
  power_id: string;

  @Field()
  access_id: string;

  @Field()
  switch_id: string;

  @Field()
  location: string;

  @Field()
  is_deactivated: boolean;

  @Field()
  latitude: string;

  @Field()
  longitude: string;

  @Field()
  install_date: string;

  @Field()
  created_at: string;
}

@ObjectType()
export class SitesAPIResDto {
  @Field(() => [SiteAPIDto])
  sites: SiteAPIDto[];
}

@ObjectType()
export class SiteAPIResDto {
  @Field(() => SiteAPIDto)
  site: SiteAPIDto;
}

@InputType()
export class AddSiteInputDto {
  @Field()
  name: string;

  @Field()
  network_id: string;

  @Field()
  backhaul_id: string;

  @Field()
  power_id: string;

  @Field()
  access_id: string;

  @Field()
  spectrum_id: string;

  @Field()
  switch_id: string;

  @Field()
  latitude: string;

  @Field()
  longitude: string;

  @Field()
  install_date: string;

  @Field()
  location: string;
}

@InputType()
export class UpdateSiteInputDto {
  @Field()
  name: string;
}

@InputType()
export class SitesInputDto {
  @Field({ nullable: true })
  networkId?: string;
}
