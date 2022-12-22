# Subscriber Registry

The Subscriber Registry sub-system contains a gRPC server that allows you to manage subscribers. The following methods are exposed:

## Features

- Add: Adds a new subscriber to the database.

- Delete: Deletes a subscriber from the database.

- Get: Retrieves a subscriber from the database by ID.

- Update: Updates a subscriber in the database.

## Prerequisites

To use this sub-registry system, you will need the following:

- A gRPC client library, such as [grpcurl](https://github.com/fullstorydev/grpcurl) or [Evans](https://github.com/ktr0731/evans).

- The Protocol Buffers (Protobuf) compiler, to generate the necessary code for your client and server applications.

## Installation

- Clone this repository to your local machine:

  `git clone https://github.com/your_username/subscriber-grpc-server.git`

* Navigate to the root directory of the repository:

```
cd/subscriber
```

- Compile the Protobuf files using the Protobuf compiler:

```
 make gen
```

- To start the subscriber reg server

```
 make server
```

## Demo

Make sure that the subscriber registry is running and listening for connections. You can start the server by running `make server`.

### Add subscriber

```
grpcurl -d '{"networkID": "123456", "firstName": "John", "lastName": "Doe", "email": "john.doe@example.com", "phoneNumber": "123-456-7890", "gender": "M"}' -plaintext localhost:9090 SubscriberService.Add

```

### Get subscriber

```
grpcurl -d '{"subscriberID": "your_subscriber_id"}' -plaintext localhost:9090 SubscriberService.Get

```

### Update subscriber

```
grpcurl -d '{"subscriberID": "123456", "firstName": "John", "lastName": "Doe", "email": "john.doe@example.com", "phoneNumber": "123-456-7890", "gender": "M"}' -plaintext localhost:9090 SubscriberService.Update

```

### Get subscriber by network

```
grpcurl -d '{"networkID": "123456"}' -plaintext localhost:9090 SubscriberService.GetByNetwork

```
