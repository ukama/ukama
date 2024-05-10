# Site Registry

Site Registry manages sites within the Ukama network.

## Description
The Site Registry is responsible for handling various operations related to sites within the Ukama network. Each site belongs to a network, and networks are associated with organizations. Thus, the relationships are organized as follows: every site belongs to a network, and every network belongs to an organization.

## Service definition
The Site Registry provides the following RPC definitions:

``` proto
service SiteService {
    rpc Add(AddRequest) returns (AddResponse);
    rpc Get(GetRequest) returns (GetResponse);
    rpc GetSites(GetSitesRequest) returns (GetSitesResponse);
}
