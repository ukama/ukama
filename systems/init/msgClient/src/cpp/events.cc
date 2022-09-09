/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <iostream>
#include <memory>
#include <string>

#include <grpcpp/ext/proto_server_reflection_plugin.h>
#include <grpcpp/grpcpp.h>
#include <grpcpp/health_check_service_interface.h>

#include "events.grpc.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;

using events::EventListener;
using events::RegistrationRequest;
using events::Response;

class EventsServiceImpl final : public EventListener::Service {
  Status SystemRegistrationEvent(ServerContext *context,
								 const RegistrationRequest *request,
								 Response *reply) {
	
	// request->name();
	// request->ip();
	// request->port();

	  return Status::OK;
  }
};


#ifdef __cplusplus
extern "C" {
#endif

/*
 * run_grpc_server -- run GRPC server on given address (host:port)
 *
 */
void run_grpc_server(char *address) {
  
	std::string server_address(address);
	EventsServiceImpl service;

	grpc::EnableDefaultHealthCheckService(true);
	grpc::reflection::InitProtoReflectionServerBuilderPlugin();
	ServerBuilder builder;
	// Listen on the given address without any authentication mechanism.
	builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
	// Register "service" as the instance through which we'll communicate with
	// clients. In this case it corresponds to an *synchronous* service.
	builder.RegisterService(&service);
	// Finally assemble the server.
	std::unique_ptr<Server> server(builder.BuildAndStart());
	std::cout << "GRPC server listening on " << server_address << std::endl;

	// Wait for the server to shutdown. Note that some other thread must be
	// responsible for shutting down the server for this call to ever return.
	server->Wait();
}

#ifdef __cplusplus
}
#endif
