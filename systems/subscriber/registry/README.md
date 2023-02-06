# Subscriber Registry

The Subscriber Registry sub-system allows you to manage subscribers. It has the following features:

      ├── README.md
      ├── bin
      │      ├── integration
      │      ├── subscriber-registry
      ├── cmd
      │      ├── server
      │      │      ├── main.go
      │      ├── version
      │      │      ├── version.go
      ├── coverage.out
      ├── dockerfile
      ├── generate-dir-tree.sh
      ├── go.mod
      ├── go.sum
      ├── makefile
      ├── mocks
      │      ├── SubscriberRepo.go
      ├── pb
      │      ├── gen
      │      │      ├── mocks
      │      │      │      ├── SubscriberRegistryServiceClient.go
      │      │      │      ├── SubscriberRegistryServiceServer.go
      │      │      │      ├── UnsafeSubscriberRegistryServiceServer.go
      │      │      ├── subscriber.pb.go
      │      │      ├── subscriber.validator.pb.go
      │      │      ├── subscriber_grpc.pb.go
      │      ├── subscriber.proto
      ├── pkg
      │      ├── config.go
      │      ├── db
      │      │      ├── model.go
      │      │      ├── subscriber_repo.go
      │      │      ├── subscriber_repo_test.go
      │      ├── global.go
      │      ├── server
      │      │      ├── event.go
      │      │      ├── subscriber.go
      │      │      ├── subscriber_test.go
      ├── template.tmpl
      ├── test
      │      ├── integration
      │      │      ├── susbcriber_test.go

## Features

- Add: Adds a new subscriber to the database.

- Delete: Deletes a subscriber from the database and sends an event message with the subscriber's ID to a message bus.

- Get: Retrieves a subscriber from the database by ID.

- Update: Updates a subscriber in the database.

- GetByNetwork: Retrieve subscriber by NetworkID.

## Prerequisites

To use this sub-registry system, you will need the following:

- A gRPC client library, such as [grpcurl](https://github.com/fullstorydev/grpcurl) or [Evans](https://github.com/ktr0731/evans).

- The Protocol Buffers (Protobuf) compiler, to generate the necessary code for your client and server applications.

## Installation

- Clone this repository to your local machine:

  `git clone https://github.com/ukama/ukama.git`

* Navigate to the root directory of the repository:

```
cd/systems/subscriber/subsciber
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

Make sure that the subscriber registry and message bus broker are running and listening for connections. You can start the server by running `make server`.

### Add

The Add method allows you to add a new subscriber to the database. It takes the following parameters:

```
orgID : A UUID representing the org ID of the subscriber.
networkID`: A string representing the network ID of the subscriber.
firstName: A string representing the first name of the subscriber.
lastName: A string representing the last name of the subscriber.
email: A string representing the email address of the subscriber.
phoneNumber: A string representing the phone number of the subscriber.
gender: A string representing the gender of the subscriber.

```

#### Example usage:

```

grpcurl -d '{"networkID": "123456","orgID:" "137849","firstName": "John", "lastName": "Doe", "email": "john.doe@example.com", "phoneNumber": "123-456-7890", "gender": "M"}' -plaintext localhost:9090 SubscriberService.Add

```

### Get

The Get method allows you to retrieve a subscriber from the database by ID. It takes the following parameter:

```
subscriberID: A string representing the ID of the subscriber to retrieve.
```

#### Example usage:

```
grpcurl -d '{"subscriberID": "your_subscriber_id"}' -plaintext localhost:9090 SubscriberService.Get

```

Update

The Update method allows you to update an existing subscriber in the database. It takes the following parameters:

```
subscriberID: A string representing the ID of the subscriber to be updated.
firstName: A string representing the first name of the subscriber.
lastName: A string representing the last name of the subscriber.
email: A string representing the email address of the subscriber.
phoneNumber: A string representing the phone number of the subscriber.
gender: A string representing the gender of the subscriber.
```

#### Example usage:

```
grpcurl -d '{"subscriberID": "123456", "firstName": "John", "lastName": "Doe", "email": "john.doe@example.com", "phoneNumber": "123-456-7890", "gender": "M"}' -plaintext localhost:9090 SubscriberService.Update

```

### GetByNetwork

The GetByNetwork method allows you to retrieve all subscribers from the database by network ID. It takes the following parameter:

networkID: A string representing the network ID of the subscriber to retrieve.
This will return all subscribers with the matching network ID, or an error if no such subscriber is found.

```
grpcurl -d '{"networkID": "123456"}' -plaintext localhost:9090 SubscriberService.GetByNetwork

```
