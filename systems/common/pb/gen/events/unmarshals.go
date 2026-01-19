/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package events

import (
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func UnmarshalAddMemberEventRequest(msg *anypb.Any, emsg string) (*AddMemberEventRequest, error) {
	p := &AddMemberEventRequest{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalAsrActivated(msg *anypb.Any, emsg string) (*AsrActivated, error) {
	p := &AsrActivated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalAsrInactivated(msg *anypb.Any, emsg string) (*AsrInactivated, error) {
	p := &AsrInactivated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalAsrUpdated(msg *anypb.Any, emsg string) (*AsrUpdated, error) {
	p := &AsrUpdated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalCDRReported(msg *anypb.Any, emsg string) (*CDRReported, error) {
	p := &CDRReported{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalCappCreatedEvent(msg *anypb.Any, emsg string) (*CappCreatedEvent, error) {
	p := &CappCreatedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalCapps(msg *anypb.Any, emsg string) (*Capps, error) {
	p := &Capps{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalCreatePackageEvent(msg *anypb.Any, emsg string) (*CreatePackageEvent, error) {
	p := &CreatePackageEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalCustomer(msg *anypb.Any, emsg string) (*Customer, error) {
	p := &Customer{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalDefaultMarkupUpdate(msg *anypb.Any, emsg string) (*DefaultMarkupUpdate, error) {
	p := &DefaultMarkupUpdate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalDeleteMemberEventRequest(msg *anypb.Any, emsg string) (*DeleteMemberEventRequest, error) {
	p := &DeleteMemberEventRequest{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalDeletePackageEvent(msg *anypb.Any, emsg string) (*DeletePackageEvent, error) {
	p := &DeletePackageEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEnforceNodeStateEvent(msg *anypb.Any, emsg string) (*EnforceNodeStateEvent, error) {
	p := &EnforceNodeStateEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEvent(msg *anypb.Any, emsg string) (*Event, error) {
	p := &Event{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventAddSite(msg *anypb.Any, emsg string) (*EventAddSite, error) {
	p := &EventAddSite{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventArtifactChunkReady(msg *anypb.Any, emsg string) (*EventArtifactChunkReady, error) {
	p := &EventArtifactChunkReady{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventArtifactUploaded(msg *anypb.Any, emsg string) (*EventArtifactUploaded, error) {
	p := &EventArtifactUploaded{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventBaserateUploaded(msg *anypb.Any, emsg string) (*EventBaserateUploaded, error) {
	p := &EventBaserateUploaded{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventInventoryNodeComponentAdd(msg *anypb.Any, emsg string) (*EventInventoryNodeComponentAdd, error) {
	p := &EventInventoryNodeComponentAdd{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventInvitationCreated(msg *anypb.Any, emsg string) (*EventInvitationCreated, error) {
	p := &EventInvitationCreated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventInvitationDeleted(msg *anypb.Any, emsg string) (*EventInvitationDeleted, error) {
	p := &EventInvitationDeleted{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventInvitationUpdated(msg *anypb.Any, emsg string) (*EventInvitationUpdated, error) {
	p := &EventInvitationUpdated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventNetworkCreate(msg *anypb.Any, emsg string) (*EventNetworkCreate, error) {
	p := &EventNetworkCreate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventNetworkDelete(msg *anypb.Any, emsg string) (*EventNetworkDelete, error) {
	p := &EventNetworkDelete{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventOperatorCdrReport(msg *anypb.Any, emsg string) (*EventOperatorCdrReport, error) {
	p := &EventOperatorCdrReport{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventOrgCreate(msg *anypb.Any, emsg string) (*EventOrgCreate, error) {
	p := &EventOrgCreate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventOrgRegisterUser(msg *anypb.Any, emsg string) (*EventOrgRegisterUser, error) {
	p := &EventOrgRegisterUser{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventRegistryNodeAssign(msg *anypb.Any, emsg string) (*EventRegistryNodeAssign, error) {
	p := &EventRegistryNodeAssign{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventRegistryNodeAttach(msg *anypb.Any, emsg string) (*EventRegistryNodeAttach, error) {
	p := &EventRegistryNodeAttach{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventRegistryNodeCreate(msg *anypb.Any, emsg string) (*EventRegistryNodeCreate, error) {
	p := &EventRegistryNodeCreate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventRegistryNodeDelete(msg *anypb.Any, emsg string) (*EventRegistryNodeDelete, error) {
	p := &EventRegistryNodeDelete{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventRegistryNodeDettach(msg *anypb.Any, emsg string) (*EventRegistryNodeDettach, error) {
	p := &EventRegistryNodeDettach{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventRegistryNodeRelease(msg *anypb.Any, emsg string) (*EventRegistryNodeRelease, error) {
	p := &EventRegistryNodeRelease{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventRegistryNodeStatusUpdate(msg *anypb.Any, emsg string) (*EventRegistryNodeStatusUpdate, error) {
	p := &EventRegistryNodeStatusUpdate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventRegistryNodeUpdate(msg *anypb.Any, emsg string) (*EventRegistryNodeUpdate, error) {
	p := &EventRegistryNodeUpdate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventResponse(msg *anypb.Any, emsg string) (*EventResponse, error) {
	p := &EventResponse{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimActivation(msg *anypb.Any, emsg string) (*EventSimActivation, error) {
	p := &EventSimActivation{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimActivePackage(msg *anypb.Any, emsg string) (*EventSimActivePackage, error) {
	p := &EventSimActivePackage{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimAddPackage(msg *anypb.Any, emsg string) (*EventSimAddPackage, error) {
	p := &EventSimAddPackage{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimAllocation(msg *anypb.Any, emsg string) (*EventSimAllocation, error) {
	p := &EventSimAllocation{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimDeactivation(msg *anypb.Any, emsg string) (*EventSimDeactivation, error) {
	p := &EventSimDeactivation{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimPackageExpire(msg *anypb.Any, emsg string) (*EventSimPackageExpire, error) {
	p := &EventSimPackageExpire{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimRemovePackage(msg *anypb.Any, emsg string) (*EventSimRemovePackage, error) {
	p := &EventSimRemovePackage{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimTermination(msg *anypb.Any, emsg string) (*EventSimTermination, error) {
	p := &EventSimTermination{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimUsage(msg *anypb.Any, emsg string) (*EventSimUsage, error) {
	p := &EventSimUsage{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSimsUploaded(msg *anypb.Any, emsg string) (*EventSimsUploaded, error) {
	p := &EventSimsUploaded{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSubscriberAdded(msg *anypb.Any, emsg string) (*EventSubscriberAdded, error) {
	p := &EventSubscriberAdded{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSubscriberDeleted(msg *anypb.Any, emsg string) (*EventSubscriberDeleted, error) {
	p := &EventSubscriberDeleted{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventSubscriberUpdate(msg *anypb.Any, emsg string) (*EventSubscriberUpdate, error) {
	p := &EventSubscriberUpdate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventUpdateSite(msg *anypb.Any, emsg string) (*EventUpdateSite, error) {
	p := &EventUpdateSite{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventUserCreate(msg *anypb.Any, emsg string) (*EventUserCreate, error) {
	p := &EventUserCreate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventUserDeactivate(msg *anypb.Any, emsg string) (*EventUserDeactivate, error) {
	p := &EventUserDeactivate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventUserDelete(msg *anypb.Any, emsg string) (*EventUserDelete, error) {
	p := &EventUserDelete{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalEventUserUpdate(msg *anypb.Any, emsg string) (*EventUserUpdate, error) {
	p := &EventUserUpdate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalFee(msg *anypb.Any, emsg string) (*Fee, error) {
	p := &Fee{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalFeeItem(msg *anypb.Any, emsg string) (*FeeItem, error) {
	p := &FeeItem{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalMarkupUpdate(msg *anypb.Any, emsg string) (*MarkupUpdate, error) {
	p := &MarkupUpdate{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalMeshRegisterEvent(msg *anypb.Any, emsg string) (*MeshRegisterEvent, error) {
	p := &MeshRegisterEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeAssignedEvent(msg *anypb.Any, emsg string) (*NodeAssignedEvent, error) {
	p := &NodeAssignedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeChanged(msg *anypb.Any, emsg string) (*NodeChanged, error) {
	p := &NodeChanged{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeCreatedEvent(msg *anypb.Any, emsg string) (*NodeCreatedEvent, error) {
	p := &NodeCreatedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeDeletedEvent(msg *anypb.Any, emsg string) (*NodeDeletedEvent, error) {
	p := &NodeDeletedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeOfflineEvent(msg *anypb.Any, emsg string) (*NodeOfflineEvent, error) {
	p := &NodeOfflineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeOnlineEvent(msg *anypb.Any, emsg string) (*NodeOnlineEvent, error) {
	p := &NodeOnlineEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeReleasedEvent(msg *anypb.Any, emsg string) (*NodeReleasedEvent, error) {
	p := &NodeReleasedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeStateChangeEvent(msg *anypb.Any, emsg string) (*NodeStateChangeEvent, error) {
	p := &NodeStateChangeEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeStateUpdatedEvent(msg *anypb.Any, emsg string) (*NodeStateUpdatedEvent, error) {
	p := &NodeStateUpdatedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNodeUpdatedEvent(msg *anypb.Any, emsg string) (*NodeUpdatedEvent, error) {
	p := &NodeUpdatedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNotification(msg *anypb.Any, emsg string) (*Notification, error) {
	p := &Notification{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNotificationDeletedEvent(msg *anypb.Any, emsg string) (*NotificationDeletedEvent, error) {
	p := &NotificationDeletedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalOrgIPUpdateEvent(msg *anypb.Any, emsg string) (*OrgIPUpdateEvent, error) {
	p := &OrgIPUpdateEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalPayment(msg *anypb.Any, emsg string) (*Payment, error) {
	p := &Payment{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalProfile(msg *anypb.Any, emsg string) (*Profile, error) {
	p := &Profile{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalProfileAdded(msg *anypb.Any, emsg string) (*ProfileAdded, error) {
	p := &ProfileAdded{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalProfileRemoved(msg *anypb.Any, emsg string) (*ProfileRemoved, error) {
	p := &ProfileRemoved{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalProfileUpdated(msg *anypb.Any, emsg string) (*ProfileUpdated, error) {
	p := &ProfileUpdated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalPublishServiceStatusUp(msg *anypb.Any, emsg string) (*PublishServiceStatusUp, error) {
	p := &PublishServiceStatusUp{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalRawReport(msg *anypb.Any, emsg string) (*RawReport, error) {
	p := &RawReport{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalReport(msg *anypb.Any, emsg string) (*Report, error) {
	p := &Report{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalResource(msg *anypb.Any, emsg string) (*Resource, error) {
	p := &Resource{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalSessionCreated(msg *anypb.Any, emsg string) (*SessionCreated, error) {
	p := &SessionCreated{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalSessionDestroyed(msg *anypb.Any, emsg string) (*SessionDestroyed, error) {
	p := &SessionDestroyed{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalSimRemoved(msg *anypb.Any, emsg string) (*SimRemoved, error) {
	p := &SimRemoved{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalSimUploaded(msg *anypb.Any, emsg string) (*SimUploaded, error) {
	p := &SimUploaded{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalStoreRunningAppsInfoEvent(msg *anypb.Any, emsg string) (*StoreRunningAppsInfoEvent, error) {
	p := &StoreRunningAppsInfoEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalSubscriber(msg *anypb.Any, emsg string) (*Subscriber, error) {
	p := &Subscriber{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalSubscription(msg *anypb.Any, emsg string) (*Subscription, error) {
	p := &Subscription{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalSystem(msg *anypb.Any, emsg string) (*System, error) {
	p := &System{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalUpdateMemberEventRequest(msg *anypb.Any, emsg string) (*UpdateMemberEventRequest, error) {
	p := &UpdateMemberEventRequest{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalUpdatePackageEvent(msg *anypb.Any, emsg string) (*UpdatePackageEvent, error) {
	p := &UpdatePackageEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalUserAccounting(msg *anypb.Any, emsg string) (*UserAccounting, error) {
	p := &UserAccounting{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalUserAccountingEvent(msg *anypb.Any, emsg string) (*UserAccountingEvent, error) {
	p := &UserAccountingEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalWebhook(msg *anypb.Any, emsg string) (*Webhook, error) {
	p := &Webhook{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalWebhookDeletedEvent(msg *anypb.Any, emsg string) (*WebhookDeletedEvent, error) {
	p := &WebhookDeletedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

