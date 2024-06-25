/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
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
