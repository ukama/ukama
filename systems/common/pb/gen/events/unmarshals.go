/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package events

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func UnmarshalNodeReleasedEvent(msg *anypb.Any, emsg string) (*NodeReleasedEvent, error) {
	p := &NodeReleasedEvent{}
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

func UnmarshalCappCreatedEvent(msg *anypb.Any, emsg string) (*CappCreatedEvent, error) {
	p := &CappCreatedEvent{}
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

func UnmarshalAddSiteEventRequest(msg *anypb.Any, emsg string) (*AddSiteEventRequest, error) {
	p := &AddSiteEventRequest{}
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

func UnmarshalEventSubscriberDeleted(msg *anypb.Any, emsg string) (*EventSubscriberDeleted, error) {
	p := &EventSubscriberDeleted{}
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

func UnmarshalNodeOfflineEvent(msg *anypb.Any, emsg string) (*NodeOfflineEvent, error) {
	p := &NodeOfflineEvent{}
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

func UnmarshalRawInvoice(msg *anypb.Any, emsg string) (*RawInvoice, error) {
	p := &RawInvoice{}
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

func UnmarshalFee(msg *anypb.Any, emsg string) (*Fee, error) {
	p := &Fee{}
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

func UnmarshalCustomer(msg *anypb.Any, emsg string) (*Customer, error) {
	p := &Customer{}
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

func UnmarshalSimAllocation(msg *anypb.Any, emsg string) (*SimAllocation, error) {
	p := &SimAllocation{}
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

func UnmarshalUpdateSiteEventRequest(msg *anypb.Any, emsg string) (*UpdateSiteEventRequest, error) {
	p := &UpdateSiteEventRequest{}
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

func UnmarshalCreatePackageEvent(msg *anypb.Any, emsg string) (*CreatePackageEvent, error) {
	p := &CreatePackageEvent{}
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

func UnmarshalUserAccounting(msg *anypb.Any, emsg string) (*UserAccounting, error) {
	p := &UserAccounting{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalInvoice(msg *anypb.Any, emsg string) (*Invoice, error) {
	p := &Invoice{}
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

func UnmarshalNodeStateUpdatedEvent(msg *anypb.Any, emsg string) (*NodeStateUpdatedEvent, error) {
	p := &NodeStateUpdatedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalAddMemberEventRequest(msg *anypb.Any, emsg string) (*AddMemberEventRequest, error) {
	p := &AddMemberEventRequest{}
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

func UnmarshalNodeCreatedEvent(msg *anypb.Any, emsg string) (*NodeCreatedEvent, error) {
	p := &NodeCreatedEvent{}
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

func UnmarshalFeeItem(msg *anypb.Any, emsg string) (*FeeItem, error) {
	p := &FeeItem{}
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

func UnmarshalEventBaserateUploaded(msg *anypb.Any, emsg string) (*EventBaserateUploaded, error) {
	p := &EventBaserateUploaded{}
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

func UnmarshalEventSubscriberAdded(msg *anypb.Any, emsg string) (*EventSubscriberAdded, error) {
	p := &EventSubscriberAdded{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalNetworkCreatedEvent(msg *anypb.Any, emsg string) (*NetworkCreatedEvent, error) {
	p := &NetworkCreatedEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalSimUsage(msg *anypb.Any, emsg string) (*SimUsage, error) {
	p := &SimUsage{}
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

func UnmarshalUpdatePackageEvent(msg *anypb.Any, emsg string) (*UpdatePackageEvent, error) {
	p := &UpdatePackageEvent{}
	err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
	if err != nil {
		log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
		return nil, err
	}
	return p, nil
}

func UnmarshalSimActivePackage(msg *anypb.Any, emsg string) (*SimActivePackage, error) {
	p := &SimActivePackage{}
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
