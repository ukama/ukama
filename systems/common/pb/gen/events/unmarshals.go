package events
import (
"google.golang.org/protobuf/types/known/anypb"
"google.golang.org/protobuf/proto"
log "github.com/sirupsen/logrus"
)
func unmarshalUpdatePackageEvent(msg *anypb.Any, emsg string) (*UpdatePackageEvent, error) {
  p := &UpdatePackageEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalSimAllocation(msg *anypb.Any, emsg string) (*SimAllocation, error) {
  p := &SimAllocation{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalUpdateSiteEventRequest(msg *anypb.Any, emsg string) (*UpdateSiteEventRequest, error) {
  p := &UpdateSiteEventRequest{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalPublishServiceStatusUp(msg *anypb.Any, emsg string) (*PublishServiceStatusUp, error) {
  p := &PublishServiceStatusUp{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNodeOfflineEvent(msg *anypb.Any, emsg string) (*NodeOfflineEvent, error) {
  p := &NodeOfflineEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalSimActivePackage(msg *anypb.Any, emsg string) (*SimActivePackage, error) {
  p := &SimActivePackage{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalEventSubscriberAdded(msg *anypb.Any, emsg string) (*EventSubscriberAdded, error) {
  p := &EventSubscriberAdded{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalEventSubscriberDeleted(msg *anypb.Any, emsg string) (*EventSubscriberDeleted, error) {
  p := &EventSubscriberDeleted{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalUserAccountingEvent(msg *anypb.Any, emsg string) (*UserAccountingEvent, error) {
  p := &UserAccountingEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalEventBaserateUploaded(msg *anypb.Any, emsg string) (*EventBaserateUploaded, error) {
  p := &EventBaserateUploaded{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalSubscription(msg *anypb.Any, emsg string) (*Subscription, error) {
  p := &Subscription{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalInvoice(msg *anypb.Any, emsg string) (*Invoice, error) {
  p := &Invoice{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalCreatePackageEvent(msg *anypb.Any, emsg string) (*CreatePackageEvent, error) {
  p := &CreatePackageEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalCustomer(msg *anypb.Any, emsg string) (*Customer, error) {
  p := &Customer{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNetworkCreatedEvent(msg *anypb.Any, emsg string) (*NetworkCreatedEvent, error) {
  p := &NetworkCreatedEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalAddSiteEventRequest(msg *anypb.Any, emsg string) (*AddSiteEventRequest, error) {
  p := &AddSiteEventRequest{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalFeeItem(msg *anypb.Any, emsg string) (*FeeItem, error) {
  p := &FeeItem{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalAddMemberEventRequest(msg *anypb.Any, emsg string) (*AddMemberEventRequest, error) {
  p := &AddMemberEventRequest{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalDeletePackageEvent(msg *anypb.Any, emsg string) (*DeletePackageEvent, error) {
  p := &DeletePackageEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalEventInvitationCreated(msg *anypb.Any, emsg string) (*EventInvitationCreated, error) {
  p := &EventInvitationCreated{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNodeUpdatedEvent(msg *anypb.Any, emsg string) (*NodeUpdatedEvent, error) {
  p := &NodeUpdatedEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalCappCreatedEvent(msg *anypb.Any, emsg string) (*CappCreatedEvent, error) {
  p := &CappCreatedEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalUserAccounting(msg *anypb.Any, emsg string) (*UserAccounting, error) {
  p := &UserAccounting{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalDefaultMarkupUpdate(msg *anypb.Any, emsg string) (*DefaultMarkupUpdate, error) {
  p := &DefaultMarkupUpdate{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNotificationDeletedEvent(msg *anypb.Any, emsg string) (*NotificationDeletedEvent, error) {
  p := &NotificationDeletedEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalEventSimsUploaded(msg *anypb.Any, emsg string) (*EventSimsUploaded, error) {
  p := &EventSimsUploaded{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalOrgIPUpdateEvent(msg *anypb.Any, emsg string) (*OrgIPUpdateEvent, error) {
  p := &OrgIPUpdateEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNodeDeletedEvent(msg *anypb.Any, emsg string) (*NodeDeletedEvent, error) {
  p := &NodeDeletedEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNodeOnlineEvent(msg *anypb.Any, emsg string) (*NodeOnlineEvent, error) {
  p := &NodeOnlineEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalSimUsage(msg *anypb.Any, emsg string) (*SimUsage, error) {
  p := &SimUsage{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalRawInvoice(msg *anypb.Any, emsg string) (*RawInvoice, error) {
  p := &RawInvoice{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNodeAssignedEvent(msg *anypb.Any, emsg string) (*NodeAssignedEvent, error) {
  p := &NodeAssignedEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalEvent(msg *anypb.Any, emsg string) (*Event, error) {
  p := &Event{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalMarkupUpdate(msg *anypb.Any, emsg string) (*MarkupUpdate, error) {
  p := &MarkupUpdate{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNodeCreatedEvent(msg *anypb.Any, emsg string) (*NodeCreatedEvent, error) {
  p := &NodeCreatedEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNodeStateUpdatedEvent(msg *anypb.Any, emsg string) (*NodeStateUpdatedEvent, error) {
  p := &NodeStateUpdatedEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalFee(msg *anypb.Any, emsg string) (*Fee, error) {
  p := &Fee{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNotification(msg *anypb.Any, emsg string) (*Notification, error) {
  p := &Notification{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalPayment(msg *anypb.Any, emsg string) (*Payment, error) {
  p := &Payment{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalNodeReleasedEvent(msg *anypb.Any, emsg string) (*NodeReleasedEvent, error) {
  p := &NodeReleasedEvent{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

func unmarshalEventResponse(msg *anypb.Any, emsg string) (*EventResponse, error) {
  p := &EventResponse{}
  err := anypb.UnmarshalTo(msg, p, proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
  if err != nil {
    log.Errorf("%s : %+v. Error %s.", emsg, msg, err.Error())
    return nil, err
  }
  return p, nil
}

