# analytics

Single read-only analytics service behind `api-gateway-analytics`.

It registers the business, customer and network gRPC APIs on one gRPC server.
The gateway remains a pure gateway and forwards to this service.
The collector remains the only writer/worker service.
