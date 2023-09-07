## Registry Invitation gRPC

This repository contains the gRPC definition for the registry Invitation service. The service allows users to manage invitations for Ukama organizations.
Installation
To use this gRPC service, you will need to have the following dependencies installed:
Protocol Buffers (protoc)
Go

## Usage

To generate the Go code from the protobuf definition, run the following command:
```
make gen
```
This will generate the necessary Go files for the service.

## Service Definition

The service is defined in the invitation.proto file. It includes the following RPC methods:

Invitations
`Add`: Adds a new invitation.
`Get`: Retrieves an invitation by ID.
`UpdateStatus`: Updates the status of an invitation.
`Delete`: Deletes an invitation.
`GetByOrg`: Retrieves all invitations for a specific organization.
`GetInvitationByEmail`: Retrieves an invitation by email.

### Message Definitions

The service uses the following message definitions:
`GetInvitationByEmailRequest`: Request message for retrieving an invitation by email.
`GetInvitationByEmailResponse`: Response message for retrieving an invitation by email.
`AddInvitationRequest`: Request message for adding a new invitation.
`GetInvitationByOrgRequest`: Request message for retrieving invitations by organization.
`GetInvitationByOrgResponse`: Response message for retrieving invitations by organization.
`AddInvitationResponse`: Response message for adding a new invitation.
`GetInvitationRequest`: Request message for retrieving an invitation by ID.
`GetInvitationResponse`: Response message for retrieving an invitation by ID.
`DeleteInvitationRequest`: Request message for deleting an invitation.
`DeleteInvitationResponse`: Response message for deleting an invitation.
`UpdateInvitationStatusRequest`: Request message for updating the status of an invitation.
`UpdateInvitationStatusResponse`: Response message for updating the status of an invitation.
`Invitation`: Represents an invitation with its properties.
`StatusType`: Enum for the status of an invitation.
`RoleType`: Enum for the role of an invitation.

