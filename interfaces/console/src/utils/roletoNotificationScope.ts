import { Notification_Scope, Role_Type } from '@/client/graphql/generated';

// mapping for role to scope
export const RoleToNotificationScopes: {
  [key in Role_Type]: Notification_Scope[];
} = {
  [Role_Type.RoleOwner]: [
    Notification_Scope.ScopeOrg,
    Notification_Scope.ScopeNetworks,
    Notification_Scope.ScopeNetwork,
    Notification_Scope.ScopeSites,
    Notification_Scope.ScopeSite,
    Notification_Scope.ScopeSubscribers,
    Notification_Scope.ScopeSubscriber,
    Notification_Scope.ScopeUsers,
    Notification_Scope.ScopeUser,
    Notification_Scope.ScopeNode,
  ],
  [Role_Type.RoleAdmin]: [
    Notification_Scope.ScopeOrg,
    Notification_Scope.ScopeNetworks,
    Notification_Scope.ScopeNetwork,
    Notification_Scope.ScopeSites,
    Notification_Scope.ScopeSite,
    Notification_Scope.ScopeSubscribers,
    Notification_Scope.ScopeSubscriber,
    Notification_Scope.ScopeUsers,
    Notification_Scope.ScopeUser,
    Notification_Scope.ScopeNode,
  ],
  [Role_Type.RoleNetworkOwner]: [
    Notification_Scope.ScopeNetwork,
    Notification_Scope.ScopeSite,
    Notification_Scope.ScopeSites,
    Notification_Scope.ScopeSubscribers,
    Notification_Scope.ScopeSubscriber,
    Notification_Scope.ScopeUsers,
    Notification_Scope.ScopeUser,
    Notification_Scope.ScopeNode,
  ],
  [Role_Type.RoleVendor]: [Notification_Scope.ScopeNetwork],
  [Role_Type.RoleUser]: [Notification_Scope.ScopeUser],
  [Role_Type.RoleInvalid]: [],
};
