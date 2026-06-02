/** Internal type. DO NOT USE DIRECTLY. */
type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
/** Internal type. DO NOT USE DIRECTLY. */
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
import type * as Types from '../../../../schema-types';

import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
const defaultOptions = {} as const;
export type NodeFragment = { id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } }>, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } };

export type GetNodeQueryVariables = Exact<{
  data: Types.NodeInput;
}>;


export type GetNodeQuery = { getNode: { id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } }>, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } } };

export type GetNodesQueryVariables = Exact<{
  data: Types.NodesFilterInput;
}>;


export type GetNodesQuery = { getNodes: { nodes: Array<{ id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } }>, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } }> } };

export type DeleteNodeMutationVariables = Exact<{
  data: Types.NodeInput;
}>;


export type DeleteNodeMutation = { deleteNodeFromOrg: { id: string } };

export type AttachNodeMutationVariables = Exact<{
  data: Types.AttachNodeInput;
}>;


export type AttachNodeMutation = { attachNode: { success: boolean } };

export type DetachhNodeMutationVariables = Exact<{
  data: Types.NodeInput;
}>;


export type DetachhNodeMutation = { detachhNode: { success: boolean } };

export type AddNodeMutationVariables = Exact<{
  data: Types.AddNodeInput;
}>;


export type AddNodeMutation = { addNode: { id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } }>, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } } };

export type ReleaseNodeFromSiteMutationVariables = Exact<{
  data: Types.NodeInput;
}>;


export type ReleaseNodeFromSiteMutation = { releaseNodeFromSite: { success: boolean } };

export type AddNodeToSiteMutationVariables = Exact<{
  data: Types.AddNodeToSiteInput;
}>;


export type AddNodeToSiteMutation = { addNodeToSite: { success: boolean } };

export type UpdateNodeStateMutationVariables = Exact<{
  data: Types.UpdateNodeStateInput;
}>;


export type UpdateNodeStateMutation = { updateNodeState: { id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } }>, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } } };

export type GetNodesForSiteQueryVariables = Exact<{
  siteId: string;
}>;


export type GetNodesForSiteQuery = { getNodesForSite: { nodes: Array<{ id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } }>, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } }> } };

export type UpdateNodeMutationVariables = Exact<{
  data: Types.UpdateNodeInput;
}>;


export type UpdateNodeMutation = { updateNode: { id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, attached: Array<{ id: string, name: string, latitude: string, longitude: string, type: Types.NodeTypeEnum, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } }>, site: { nodeId: string | null, siteId: string | null, networkId: string | null, addedAt: string | null }, status: { connectivity: string, state: string } } };

export type GetNodeAppsQueryVariables = Exact<{
  data: Types.NodeAppsChangeLogInput;
}>;


export type GetNodeAppsQuery = { getNodeApps: { type: Types.NodeTypeEnum, apps: Array<{ name: string, date: number, version: string, cpu: string, memory: string, notes: string }> } };

export type GetNodeStateQueryVariables = Exact<{
  getNodeStateId: string;
}>;


export type GetNodeStateQuery = { getNodeState: { id: string, nodeId: string, previousStateId: string | null, previousState: Types.NodeStateEnum | null, currentState: Types.NodeStateEnum, createdAt: string } };

export type RestartNodeMutationVariables = Exact<{
  data: Types.RestartNodeInputDto;
}>;


export type RestartNodeMutation = { restartNode: { success: boolean } };

export type ToggleInternetSwitchMutationVariables = Exact<{
  data: Types.ToggleInternetSwitchInputDto;
}>;


export type ToggleInternetSwitchMutation = { toggleInternetSwitch: { success: boolean } };

export type ToggleRfStatusMutationVariables = Exact<{
  data: Types.ToggleRfStatusInputDto;
}>;


export type ToggleRfStatusMutation = { toggleRFStatus: { success: boolean } };

export type ToggleServiceMutationVariables = Exact<{
  data: Types.ToggleRfStatusInputDto;
}>;


export type ToggleServiceMutation = { toggleService: { success: boolean } };

export type GetHealthReportQueryVariables = Exact<{
  data: Types.GetHealthReportInputDto;
}>;


export type GetHealthReportQuery = { getHealthReport: { id: string, nodeId: string, timestamp: string, system: Array<{ id: string, healthId: string, name: string, value: string }>, capps: Array<{ id: string, space: string, name: string, tag: string, status: string, resources: Array<{ id: string, cappId: string, name: string, value: string }> }> } };

export type MemberFragment = { role: string, userId: string, isDeactivated: boolean, memberSince: string | null, id: string };

export type GetMembersQueryVariables = Exact<{ [key: string]: never; }>;


export type GetMembersQuery = { getMembers: { members: Array<{ name: string, email: string, role: string, userId: string, isDeactivated: boolean, memberSince: string | null, id: string }> } };

export type GetMemberQueryVariables = Exact<{
  memberId: string;
}>;


export type GetMemberQuery = { getMember: { role: string, userId: string, isDeactivated: boolean, memberSince: string | null, id: string } };

export type AddMemberMutationVariables = Exact<{
  data: Types.AddMemberInputDto;
}>;


export type AddMemberMutation = { addMember: { role: string, userId: string, isDeactivated: boolean, memberSince: string | null, id: string } };

export type RemoveMemberMutationVariables = Exact<{
  memberId: string;
}>;


export type RemoveMemberMutation = { removeMember: { success: boolean } };

export type UpdateMemberMutationVariables = Exact<{
  memberId: string;
  data: Types.UpdateMemberInputDto;
}>;


export type UpdateMemberMutation = { updateMember: { success: boolean } };

export type GetMemberByUserIdQueryVariables = Exact<{
  userId: string;
}>;


export type GetMemberByUserIdQuery = { getMemberByUserId: { userId: string, name: string, email: string, memberId: string, isDeactivated: boolean, role: string, memberSince: string | null } };

export type OrgFragment = { id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean };

export type GetOrgsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgsQuery = { getOrgs: { user: string, ownerOf: Array<{ id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }>, memberOf: Array<{ id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }> } };

export type GetOrgQueryVariables = Exact<{ [key: string]: never; }>;


export type GetOrgQuery = { getOrg: { id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean } };

export type PackageRateFragment = { rate: { sms_mo: string, sms_mt: number, data: number, amount: number } };

export type PackageMarkupFragment = { markup: { baserate: string, markup: number } };

export type SimPackagesFragment = { id: string, package_id: string, start_date: string, end_date: string, is_active: boolean };

export type SubscriberSimsFragment = { subscriberId: string, sims: Array<{ id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, isPhysical: boolean }> };

export type PackageFragment = { uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { baserate: string, markup: number } };

export type GetPackagesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetPackagesQuery = { getPackages: { packages: Array<{ uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { baserate: string, markup: number } }> } };

export type GetPackageQueryVariables = Exact<{
  packageId: string;
}>;


export type GetPackageQuery = { getPackage: { uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { baserate: string, markup: number } } };

export type GetSimsBySubscriberQueryVariables = Exact<{
  data: Types.GetSimBySubscriberInputDto;
}>;


export type GetSimsBySubscriberQuery = { getSimsBySubscriber: { subscriberId: string, sims: Array<{ id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, isPhysical: boolean }> } };

export type AddPackageMutationVariables = Exact<{
  data: Types.AddPackageInputDto;
}>;


export type AddPackageMutation = { addPackage: { uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { baserate: string, markup: number } } };

export type RemovePackageForSimMutationVariables = Exact<{
  data: Types.RemovePackageFormSimInputDto;
}>;


export type RemovePackageForSimMutation = { removePackageForSim: { packageId: string | null } };

export type DeletePackageMutationVariables = Exact<{
  packageId: string;
}>;


export type DeletePackageMutation = { deletePackage: { uuid: string } };

export type GetPackagesForSimQueryVariables = Exact<{
  data: Types.GetPackagesForSimInputDto;
}>;


export type GetPackagesForSimQuery = { getPackagesForSim: { sim_id: string, packages: Array<{ id: string, package_id: string, start_date: string, end_date: string, is_active: boolean }> } };

export type AddPackagesToSimMutationVariables = Exact<{
  data: Types.AddPackagesToSimInputDto;
}>;


export type AddPackagesToSimMutation = { addPackagesToSim: { packages: Array<{ packageId: string | null }> } };

export type DeleteSimMutationVariables = Exact<{
  data: Types.DeleteSimInputDto;
}>;


export type DeleteSimMutation = { deleteSim: { simId: string | null } };

export type UpdatePacakgeMutationVariables = Exact<{
  packageId: string;
  data: Types.UpdatePackageInputDto;
}>;


export type UpdatePacakgeMutation = { updatePackage: { uuid: string, name: string, active: boolean, duration: number, simType: string, createdAt: string, deletedAt: string, updatedAt: string, smsVolume: number, dataVolume: number, voiceVolume: number, ulbr: string, dlbr: string, type: string, dataUnit: string, voiceUnit: string, messageUnit: string, flatrate: boolean, currency: string, from: string, to: string, country: string, provider: string, apn: string, ownerId: string, amount: number, rate: { sms_mo: string, sms_mt: number, data: number, amount: number }, markup: { baserate: string, markup: number } } };

export type PaymentFragment = { id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, extra: string, createdAt: string };

export type UpdatePaymentMutationVariables = Exact<{
  data: Types.UpdatePaymentInputDto;
}>;


export type UpdatePaymentMutation = { updatePayment: { id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, createdAt: string } };

export type ProcessPaymentMutationVariables = Exact<{
  data: Types.ProcessPaymentInputDto;
}>;


export type ProcessPaymentMutation = { processPayment: { payment: { id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, createdAt: string } } };

export type GetPaymentQueryVariables = Exact<{
  paymentId: string;
}>;


export type GetPaymentQuery = { getPayment: { id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, extra: string, createdAt: string } };

export type GetPaymentsQueryVariables = Exact<{
  data: Types.GetPaymentsInputDto;
}>;


export type GetPaymentsQuery = { getPayments: { payments: Array<{ id: string, itemId: string, itemType: string, amount: string, currency: string, paymentMethod: string, depositedAmount: string, paidAt: string, payerName: string, payerEmail: string, payerPhone: string, correspondent: string, country: string, description: string, status: string, failureReason: string, extra: string, createdAt: string }> } };

export type CustomerFragment = { externalId: string, name: string, email: string | null, addressLine1: string | null, legalName: string | null, legalNumber: string | null, phone: string | null, currency: string, timezone: string | null, vatRate: number, createdAt: string };

export type SubscriptionFragment = { externalCustomerId: string, externalId: string, planCode: string, name: string | null, status: string, createdAt: string, startedAt: string, canceledAt: string | null, terminatedAt: string | null };

export type FeeFragment = { taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { type: string, code: string, name: string } };

export type RawReportFragment = { issuingDate: string, paymentDueDate: string, paymentOverdue: boolean, invoiceType: string, status: string, paymentStatus: string, feesAmountCents: string, taxesAmountCents: string, subTotalExcludingTaxesAmountCents: string, subTotalIncludingTaxesAmountCents: string, vatAmountCents: string, vatAmountCurrency: string | null, totalAmountCents: string, currency: string, fileUrl: string, customer: { externalId: string, name: string, email: string | null, addressLine1: string | null, legalName: string | null, legalNumber: string | null, phone: string | null, currency: string, timezone: string | null, vatRate: number, createdAt: string }, subscriptions: Array<{ externalCustomerId: string, externalId: string, planCode: string, name: string | null, status: string, createdAt: string, startedAt: string, canceledAt: string | null, terminatedAt: string | null }>, fees: Array<{ taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { type: string, code: string, name: string } }> };

export type GetReportsQueryVariables = Exact<{
  data: Types.GetReportsInputDto;
}>;


export type GetReportsQuery = { getReports: { reports: Array<{ id: string, ownerId: string, ownerType: string, networkId: string, period: string, type: string, isPaid: boolean, createdAt: string, rawReport: { issuingDate: string, paymentDueDate: string, paymentOverdue: boolean, invoiceType: string, status: string, paymentStatus: string, feesAmountCents: string, taxesAmountCents: string, subTotalExcludingTaxesAmountCents: string, subTotalIncludingTaxesAmountCents: string, vatAmountCents: string, vatAmountCurrency: string | null, totalAmountCents: string, currency: string, fileUrl: string, customer: { externalId: string, name: string, email: string | null, addressLine1: string | null, legalName: string | null, legalNumber: string | null, phone: string | null, currency: string, timezone: string | null, vatRate: number, createdAt: string }, subscriptions: Array<{ externalCustomerId: string, externalId: string, planCode: string, name: string | null, status: string, createdAt: string, startedAt: string, canceledAt: string | null, terminatedAt: string | null }>, fees: Array<{ taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { type: string, code: string, name: string } }> } }> } };

export type GetReportQueryVariables = Exact<{
  id: string;
}>;


export type GetReportQuery = { getReport: { report: { id: string, ownerId: string, ownerType: string, networkId: string, period: string, type: string, isPaid: boolean, createdAt: string, rawReport: { issuingDate: string, paymentDueDate: string, paymentOverdue: boolean, invoiceType: string, status: string, paymentStatus: string, feesAmountCents: string, taxesAmountCents: string, subTotalExcludingTaxesAmountCents: string, subTotalIncludingTaxesAmountCents: string, vatAmountCents: string, vatAmountCurrency: string | null, totalAmountCents: string, currency: string, fileUrl: string, customer: { externalId: string, name: string, email: string | null, addressLine1: string | null, legalName: string | null, legalNumber: string | null, phone: string | null, currency: string, timezone: string | null, vatRate: number, createdAt: string }, subscriptions: Array<{ externalCustomerId: string, externalId: string, planCode: string, name: string | null, status: string, createdAt: string, startedAt: string, canceledAt: string | null, terminatedAt: string | null }>, fees: Array<{ taxesAmountCents: string, taxesPreciseAmount: string, totalAmountCents: string, totalAmountCurrency: string, eventsCount: string, units: number, item: { type: string, code: string, name: string } }> } } } };

export type GetSimPoolStatsQueryVariables = Exact<{
  data: Types.GetSimsInput;
}>;


export type GetSimPoolStatsQuery = { getSimPoolStats: { total: number, available: number, consumed: number, failed: number, esim: number, physical: number } };

export type GetSimsFromPoolQueryVariables = Exact<{
  data: Types.GetSimsInput;
}>;


export type GetSimsFromPoolQuery = { getSimsFromPool: { sims: Array<{ id: string, qrCode: string, iccid: string, msisdn: string, isAllocated: boolean, isFailed: boolean, simType: string, smApAddress: string, activationCode: string, createdAt: string, deletedAt: string, updatedAt: string, isPhysical: boolean }> } };

export type UploadSimsMutationVariables = Exact<{
  data: Types.UploadSimsInputDto;
}>;


export type UploadSimsMutation = { uploadSims: { iccid: Array<string> } };

export type SimPackageFragment = { id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean };

export type SimFragment = { id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, isPhysical: boolean, trafficPolicy: number, firstActivatedOn: string, lastActivatedOn: string, activationsCount: string, deactivationsCount: string, allocatedAt: string, syncStatus: string, package: { id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean } | null };

export type SimAllocationPackageFragment = { id: string | null, packageId: string | null, startDate: string | null, endDate: string | null, isActive: boolean | null };

export type SimAllocationFragment = { id: string, subscriber_id: string, network_id: string, iccid: string, msisdn: string, imsi: string | null, type: string, status: string, is_physical: boolean, traffic_policy: number, allocated_at: string, sync_status: string, package: { id: string | null, packageId: string | null, startDate: string | null, endDate: string | null, isActive: boolean | null } };

export type AllocateSimMutationVariables = Exact<{
  data: Types.AllocateSimInputDto;
}>;


export type AllocateSimMutation = { allocateSim: { id: string, subscriber_id: string, network_id: string, iccid: string, msisdn: string, imsi: string | null, type: string, status: string, is_physical: boolean, traffic_policy: number, allocated_at: string, sync_status: string, package: { id: string | null, packageId: string | null, startDate: string | null, endDate: string | null, isActive: boolean | null } } };

export type ToggleSimStatusMutationVariables = Exact<{
  data: Types.ToggleSimStatusInputDto;
}>;


export type ToggleSimStatusMutation = { toggleSimStatus: { simId: string | null } };

export type GetSimQueryVariables = Exact<{
  data: Types.GetSimInputDto;
}>;


export type GetSimQuery = { getSim: { id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, isPhysical: boolean, trafficPolicy: number, firstActivatedOn: string, lastActivatedOn: string, activationsCount: string, deactivationsCount: string, allocatedAt: string, syncStatus: string, package: { id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean } | null } };

export type GetSimsQueryVariables = Exact<{
  data: Types.ListSimsInput;
}>;


export type GetSimsQuery = { getSims: { sims: Array<{ id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, isPhysical: boolean, trafficPolicy: number, firstActivatedOn: string, lastActivatedOn: string, activationsCount: string, deactivationsCount: string, allocatedAt: string, syncStatus: string, package: { id: string, packageId: string, startDate: string, endDate: string, defaultDuration: string, isActive: boolean, asExpired: boolean } | null }> } };

export type SubscriberSimFragment = { sim: Array<{ id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status: string | null, isPhysical: boolean | null, package: { id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null };

export type SubscriberFragment = { uuid: string, address: string, dob: string, email: string, name: string, gender: string, idSerial: string, networkId: string, phone: string, proofOfIdentification: string, sim: Array<{ id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status: string | null, isPhysical: boolean | null, package: { id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null };

export type AddSubscriberMutationVariables = Exact<{
  data: Types.SubscriberInputDto;
}>;


export type AddSubscriberMutation = { addSubscriber: { uuid: string, address: string, dob: string, email: string, name: string, gender: string, idSerial: string, networkId: string, phone: string, proofOfIdentification: string, sim: Array<{ id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status: string | null, isPhysical: boolean | null, package: { id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null } };

export type GetSubscriberQueryVariables = Exact<{
  subscriberId: string;
}>;


export type GetSubscriberQuery = { getSubscriber: { uuid: string, address: string, dob: string, email: string, name: string, gender: string, idSerial: string, networkId: string, phone: string, proofOfIdentification: string, sim: Array<{ id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status: string | null, isPhysical: boolean | null, package: { id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null } };

export type UpdateSubscriberMutationVariables = Exact<{
  subscriberId: string;
  data: Types.UpdateSubscriberInputDto;
}>;


export type UpdateSubscriberMutation = { updateSubscriber: { success: boolean } };

export type DeleteSubscriberMutationVariables = Exact<{
  subscriberId: string;
}>;


export type DeleteSubscriberMutation = { deleteSubscriber: { success: boolean } };

export type GetSubscribersByNetworkQueryVariables = Exact<{
  networkId: string;
}>;


export type GetSubscribersByNetworkQuery = { getSubscribersByNetwork: { subscribers: Array<{ uuid: string, address: string, dob: string, email: string, name: string, gender: string, idSerial: string, networkId: string, phone: string, proofOfIdentification: string, sim: Array<{ id: string, subscriberId: string, networkId: string, iccid: string, msisdn: string, imsi: string, type: string, status: string, allocatedAt: string, sync_status: string | null, isPhysical: boolean | null, package: { id: string, package_id: string, start_date: string, end_date: string, is_active: boolean, created_at: string, updated_at: string } | null }> | null }> } };

export type GetSubscriberMetricsByNetworkQueryVariables = Exact<{
  networkId: string;
}>;


export type GetSubscriberMetricsByNetworkQuery = { getSubscriberMetricsByNetwork: { total: number, active: number, inactive: number, terminated: number } };

export type GetGeneratedPdfReportQueryVariables = Exact<{
  Id: string;
}>;


export type GetGeneratedPdfReportQuery = { getGeneratedPdfReport: { contentType: string, filename: string, downloadUrl: string } };

export type UserFragment = { name: string, uuid: string, email: string, phone: string, authId: string, isDeactivated: boolean, registeredSince: string };

export type WhoamiQueryVariables = Exact<{ [key: string]: never; }>;


export type WhoamiQuery = { whoami: { user: { name: string, uuid: string, email: string, phone: string, authId: string, isDeactivated: boolean, registeredSince: string }, ownerOf: Array<{ id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }>, memberOf: Array<{ id: string, name: string, owner: string, country: string, currency: string, createdAt: string, certificate: string, isDeactivated: boolean }> } };

export type GetUserQueryVariables = Exact<{
  userId: string;
}>;


export type GetUserQuery = { getUser: { name: string, uuid: string, email: string, phone: string, authId: string, isDeactivated: boolean, registeredSince: string } };

export type UNetworkFragment = { id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> };

export type GetNetworksQueryVariables = Exact<{ [key: string]: never; }>;


export type GetNetworksQuery = { getNetworks: { networks: Array<{ id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> }> } };

export type GetNetworkQueryVariables = Exact<{
  networkId: string;
}>;


export type GetNetworkQuery = { getNetwork: { id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> } };

export type AddNetworkMutationVariables = Exact<{
  data: Types.AddNetworkInputDto;
}>;


export type AddNetworkMutation = { addNetwork: { id: string, name: string, isDefault: boolean, budget: number, overdraft: number, trafficPolicy: number, isDeactivated: boolean, paymentLinks: boolean, createdAt: string, countries: Array<string>, networks: Array<string> } };

export type SetDefaultNetworkMutationVariables = Exact<{
  data: Types.SetDefaultNetworkInputDto;
}>;


export type SetDefaultNetworkMutation = { setDefaultNetwork: { success: boolean } };

export type USiteFragment = { id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string };

export type GetSiteQueryVariables = Exact<{
  siteId: string;
}>;


export type GetSiteQuery = { getSite: { id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string } };

export type AddSiteMutationVariables = Exact<{
  data: Types.AddSiteInputDto;
}>;


export type AddSiteMutation = { addSite: { id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string } };

export type GetSitesQueryVariables = Exact<{
  data: Types.SitesInputDto;
}>;


export type GetSitesQuery = { getSites: { sites: Array<{ id: string, name: string, networkId: string, backhaulId: string, powerId: string, accessId: string, spectrumId: string, switchId: string, isDeactivated: boolean, latitude: string, longitude: string, installDate: string, createdAt: string, location: string }> } };

export type UpdateSiteMutationVariables = Exact<{
  siteId: string;
  data: Types.UpdateSiteInputDto;
}>;


export type UpdateSiteMutation = { updateSite: { name: string } };

export type UComponentFragment = { id: string, inventoryId: string, type: string, userId: string, description: string, category: string, datasheetUrl: string, imageUrl: string, partNumber: string, manufacturer: string, managed: string, warranty: number, specification: string };

export type GetComponentByIdQueryVariables = Exact<{
  componentId: string;
}>;


export type GetComponentByIdQuery = { getComponentById: { id: string, inventoryId: string, type: string, userId: string, description: string, category: string, datasheetUrl: string, imageUrl: string, partNumber: string, manufacturer: string, managed: string, warranty: number, specification: string } };

export type GetComponentsByUserIdQueryVariables = Exact<{
  data: Types.ComponentTypeInputDto;
}>;


export type GetComponentsByUserIdQuery = { getComponentsByUserId: { components: Array<{ id: string, inventoryId: string, type: string, userId: string, description: string, category: string, datasheetUrl: string, imageUrl: string, partNumber: string, manufacturer: string, managed: string, warranty: number, specification: string }> } };

export type InvitationFragment = { email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Types.Invitation_Status };

export type CreateInvitationMutationVariables = Exact<{
  data: Types.CreateInvitationInputDto;
}>;


export type CreateInvitationMutation = { createInvitation: { email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Types.Invitation_Status } };

export type GetInvitationsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetInvitationsQuery = { getInvitations: { invitations: Array<{ email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Types.Invitation_Status }> } };

export type DeleteInvitationMutationVariables = Exact<{
  deleteInvitationId: string;
}>;


export type DeleteInvitationMutation = { deleteInvitation: { id: string } };

export type UpdateInvitationMutationVariables = Exact<{
  data: Types.UpateInvitationInputDto;
}>;


export type UpdateInvitationMutation = { updateInvitation: { id: string } };

export type GetInvitationsByEmailQueryVariables = Exact<{
  email: string;
}>;


export type GetInvitationsByEmailQuery = { getInvitationsByEmail: { invitations: Array<{ email: string, expireAt: string, id: string, name: string, role: string, link: string, userId: string, status: Types.Invitation_Status }> } };

export type GetCountriesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetCountriesQuery = { getCountries: { countries: Array<{ name: string, code: string }> } };

export type GetCurrencySymbolQueryVariables = Exact<{
  code: string;
}>;


export type GetCurrencySymbolQuery = { getCurrencySymbol: { code: string, symbol: string, image: string } };

export type GetTimezonesQueryVariables = Exact<{ [key: string]: never; }>;


export type GetTimezonesQuery = { getTimezones: { timezones: Array<{ value: string, abbr: string, offset: number, isdst: boolean, text: string, utc: Array<string> }> } };

export type UpdateNotificationMutationVariables = Exact<{
  isRead: boolean;
  updateNotificationId: string;
}>;


export type UpdateNotificationMutation = { updateNotification: { id: string } };

export type GetDataUsagesQueryVariables = Exact<{
  data: Types.SimUsagesInputDto;
}>;


export type GetDataUsagesQuery = { getDataUsages: { usages: Array<{ usage: string, simId: string }> } };

export type GetAppsQueryVariables = Exact<{ [key: string]: never; }>;


export type GetAppsQuery = { getApps: { apps: Array<{ name: string, space: string, notes: string, metricsKeys: Array<string> }> } | null };

export type SoftwareQueryVariables = Exact<{
  data: Types.GetSoftwaresInput;
}>;


export type SoftwareQuery = { getSoftwares: { software: Array<{ id: string, releaseDate: string, nodeId: string, status: Types.SoftwareStatusEnum, changeLog: Array<string>, currentVersion: string, desiredVersion: string, name: string, space: string, notes: string, metricsKeys: Array<string>, createdAt: string, updatedAt: string }> } };

export type UpdateSoftwareMutationVariables = Exact<{
  data: Types.UpdateSoftwareInputDto;
}>;


export type UpdateSoftwareMutation = { updateSoftware: { message: string } };

export const NodeFragmentDoc = gql`
    fragment node on Node {
  id
  name
  latitude
  longitude
  type
  attached {
    id
    name
    latitude
    longitude
    type
    site {
      nodeId
      siteId
      networkId
      addedAt
    }
    status {
      connectivity
      state
    }
  }
  site {
    nodeId
    siteId
    networkId
    addedAt
  }
  status {
    connectivity
    state
  }
}
    `;
export const MemberFragmentDoc = gql`
    fragment member on MemberDto {
  role
  userId
  id: memberId
  isDeactivated
  memberSince
}
    `;
export const OrgFragmentDoc = gql`
    fragment Org on OrgDto {
  id
  name
  owner
  country
  currency
  createdAt
  certificate
  isDeactivated
}
    `;
export const SimPackagesFragmentDoc = gql`
    fragment SimPackages on SimToPackagesDto {
  id
  package_id
  start_date
  end_date
  is_active
}
    `;
export const SubscriberSimsFragmentDoc = gql`
    fragment SubscriberSims on SubscriberToSimsDto {
  subscriberId
  sims {
    id
    subscriberId
    networkId
    iccid
    msisdn
    imsi
    type
    status
    allocatedAt
    isPhysical
  }
}
    `;
export const PackageRateFragmentDoc = gql`
    fragment PackageRate on PackageDto {
  rate {
    sms_mo
    sms_mt
    data
    amount
  }
}
    `;
export const PackageMarkupFragmentDoc = gql`
    fragment PackageMarkup on PackageDto {
  markup {
    baserate
    markup
  }
}
    `;
export const PackageFragmentDoc = gql`
    fragment Package on PackageDto {
  uuid
  name
  active
  duration
  simType
  createdAt
  deletedAt
  updatedAt
  smsVolume
  dataVolume
  voiceVolume
  ulbr
  dlbr
  type
  dataUnit
  voiceUnit
  messageUnit
  flatrate
  currency
  from
  to
  country
  provider
  apn
  ownerId
  amount
  ...PackageRate
  ...PackageMarkup
}
    ${PackageRateFragmentDoc}
${PackageMarkupFragmentDoc}`;
export const PaymentFragmentDoc = gql`
    fragment payment on PaymentDto {
  id
  itemId
  itemType
  amount
  currency
  paymentMethod
  depositedAmount
  paidAt
  payerName
  payerEmail
  payerPhone
  correspondent
  country
  description
  status
  failureReason
  extra
  createdAt
}
    `;
export const CustomerFragmentDoc = gql`
    fragment customer on CustomerDto {
  externalId
  name
  email
  addressLine1
  legalName
  legalNumber
  phone
  currency
  timezone
  vatRate
  createdAt
}
    `;
export const SubscriptionFragmentDoc = gql`
    fragment subscription on SubscriptionDto {
  externalCustomerId
  externalId
  planCode
  name
  status
  createdAt
  startedAt
  canceledAt
  terminatedAt
}
    `;
export const FeeFragmentDoc = gql`
    fragment fee on FeeDto {
  taxesAmountCents
  taxesPreciseAmount
  totalAmountCents
  totalAmountCurrency
  eventsCount
  units
  item {
    type
    code
    name
  }
}
    `;
export const RawReportFragmentDoc = gql`
    fragment rawReport on RawReportDto {
  issuingDate
  paymentDueDate
  paymentOverdue
  invoiceType
  status
  paymentStatus
  feesAmountCents
  taxesAmountCents
  subTotalExcludingTaxesAmountCents
  subTotalIncludingTaxesAmountCents
  vatAmountCents
  vatAmountCurrency
  totalAmountCents
  currency
  fileUrl
  customer {
    ...customer
  }
  subscriptions {
    ...subscription
  }
  fees {
    ...fee
  }
}
    ${CustomerFragmentDoc}
${SubscriptionFragmentDoc}
${FeeFragmentDoc}`;
export const SimPackageFragmentDoc = gql`
    fragment SimPackage on SimPackage {
  id
  packageId
  startDate
  endDate
  defaultDuration
  isActive
  asExpired
}
    `;
export const SimFragmentDoc = gql`
    fragment Sim on SimDto {
  id
  subscriberId
  networkId
  iccid
  msisdn
  imsi
  type
  status
  isPhysical
  trafficPolicy
  firstActivatedOn
  lastActivatedOn
  activationsCount
  deactivationsCount
  allocatedAt
  syncStatus
  package {
    ...SimPackage
  }
}
    ${SimPackageFragmentDoc}`;
export const SimAllocationPackageFragmentDoc = gql`
    fragment SimAllocationPackage on SimAllocatePackageDto {
  id
  packageId
  startDate
  endDate
  isActive
}
    `;
export const SimAllocationFragmentDoc = gql`
    fragment SimAllocation on AllocateSimAPIDto {
  id
  subscriber_id
  network_id
  package {
    ...SimAllocationPackage
  }
  iccid
  msisdn
  imsi
  type
  status
  is_physical
  traffic_policy
  allocated_at
  sync_status
}
    ${SimAllocationPackageFragmentDoc}`;
export const SubscriberSimFragmentDoc = gql`
    fragment SubscriberSim on SubscriberDto {
  sim {
    id
    subscriberId
    networkId
    iccid
    msisdn
    imsi
    type
    status
    allocatedAt
    sync_status
    isPhysical
    package {
      id
      package_id
      start_date
      end_date
      is_active
      created_at
      updated_at
    }
  }
}
    `;
export const SubscriberFragmentDoc = gql`
    fragment Subscriber on SubscriberDto {
  uuid
  address
  dob
  email
  name
  gender
  idSerial
  networkId
  phone
  proofOfIdentification
  ...SubscriberSim
}
    ${SubscriberSimFragmentDoc}`;
export const UserFragmentDoc = gql`
    fragment User on UserResDto {
  name
  uuid
  email
  phone
  authId
  isDeactivated
  registeredSince
}
    `;
export const UNetworkFragmentDoc = gql`
    fragment UNetwork on NetworkDto {
  id
  name
  isDefault
  budget
  overdraft
  trafficPolicy
  isDeactivated
  paymentLinks
  createdAt
  countries
  networks
}
    `;
export const USiteFragmentDoc = gql`
    fragment USite on SiteDto {
  id
  name
  networkId
  backhaulId
  powerId
  accessId
  spectrumId
  switchId
  isDeactivated
  latitude
  longitude
  installDate
  createdAt
  location
}
    `;
export const UComponentFragmentDoc = gql`
    fragment UComponent on ComponentDto {
  id
  inventoryId
  type
  userId
  description
  category
  datasheetUrl
  imageUrl
  partNumber
  manufacturer
  managed
  warranty
  specification
}
    `;
export const InvitationFragmentDoc = gql`
    fragment Invitation on InvitationDto {
  email
  expireAt
  id
  name
  role
  link
  userId
  status
}
    `;
export const GetNodeDocument = gql`
    query GetNode($data: NodeInput!) {
  getNode(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodeQuery__
 *
 * To run a query within a React component, call `useGetNodeQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodeQuery(baseOptions: Apollo.QueryHookOptions<GetNodeQuery, GetNodeQueryVariables> & ({ variables: GetNodeQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
      }
export function useGetNodeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
        }
// @ts-ignore
export function useGetNodeSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeQuery, GetNodeQueryVariables>;
export function useGetNodeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeQuery | undefined, GetNodeQueryVariables>;
export function useGetNodeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeQuery, GetNodeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeQuery, GetNodeQueryVariables>(GetNodeDocument, options);
        }
export type GetNodeQueryHookResult = ReturnType<typeof useGetNodeQuery>;
export type GetNodeLazyQueryHookResult = ReturnType<typeof useGetNodeLazyQuery>;
export type GetNodeSuspenseQueryHookResult = ReturnType<typeof useGetNodeSuspenseQuery>;
export type GetNodeQueryResult = Apollo.QueryResult<GetNodeQuery, GetNodeQueryVariables>;
export const GetNodesDocument = gql`
    query GetNodes($data: NodesFilterInput!) {
  getNodes(data: $data) {
    nodes {
      ...node
    }
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodesQuery__
 *
 * To run a query within a React component, call `useGetNodesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodesQuery(baseOptions: Apollo.QueryHookOptions<GetNodesQuery, GetNodesQueryVariables> & ({ variables: GetNodesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
      }
export function useGetNodesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
        }
// @ts-ignore
export function useGetNodesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodesQuery, GetNodesQueryVariables>;
export function useGetNodesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodesQuery | undefined, GetNodesQueryVariables>;
export function useGetNodesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesQuery, GetNodesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodesQuery, GetNodesQueryVariables>(GetNodesDocument, options);
        }
export type GetNodesQueryHookResult = ReturnType<typeof useGetNodesQuery>;
export type GetNodesLazyQueryHookResult = ReturnType<typeof useGetNodesLazyQuery>;
export type GetNodesSuspenseQueryHookResult = ReturnType<typeof useGetNodesSuspenseQuery>;
export type GetNodesQueryResult = Apollo.QueryResult<GetNodesQuery, GetNodesQueryVariables>;
export const DeleteNodeDocument = gql`
    mutation deleteNode($data: NodeInput!) {
  deleteNodeFromOrg(data: $data) {
    id
  }
}
    `;
export type DeleteNodeMutationFn = Apollo.MutationFunction<DeleteNodeMutation, DeleteNodeMutationVariables>;

/**
 * __useDeleteNodeMutation__
 *
 * To run a mutation, you first call `useDeleteNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteNodeMutation, { data, loading, error }] = useDeleteNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useDeleteNodeMutation(baseOptions?: Apollo.MutationHookOptions<DeleteNodeMutation, DeleteNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteNodeMutation, DeleteNodeMutationVariables>(DeleteNodeDocument, options);
      }
export type DeleteNodeMutationHookResult = ReturnType<typeof useDeleteNodeMutation>;
export type DeleteNodeMutationResult = Apollo.MutationResult<DeleteNodeMutation>;
export type DeleteNodeMutationOptions = Apollo.BaseMutationOptions<DeleteNodeMutation, DeleteNodeMutationVariables>;
export const AttachNodeDocument = gql`
    mutation attachNode($data: AttachNodeInput!) {
  attachNode(data: $data) {
    success
  }
}
    `;
export type AttachNodeMutationFn = Apollo.MutationFunction<AttachNodeMutation, AttachNodeMutationVariables>;

/**
 * __useAttachNodeMutation__
 *
 * To run a mutation, you first call `useAttachNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAttachNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [attachNodeMutation, { data, loading, error }] = useAttachNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAttachNodeMutation(baseOptions?: Apollo.MutationHookOptions<AttachNodeMutation, AttachNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AttachNodeMutation, AttachNodeMutationVariables>(AttachNodeDocument, options);
      }
export type AttachNodeMutationHookResult = ReturnType<typeof useAttachNodeMutation>;
export type AttachNodeMutationResult = Apollo.MutationResult<AttachNodeMutation>;
export type AttachNodeMutationOptions = Apollo.BaseMutationOptions<AttachNodeMutation, AttachNodeMutationVariables>;
export const DetachhNodeDocument = gql`
    mutation detachhNode($data: NodeInput!) {
  detachhNode(data: $data) {
    success
  }
}
    `;
export type DetachhNodeMutationFn = Apollo.MutationFunction<DetachhNodeMutation, DetachhNodeMutationVariables>;

/**
 * __useDetachhNodeMutation__
 *
 * To run a mutation, you first call `useDetachhNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDetachhNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [detachhNodeMutation, { data, loading, error }] = useDetachhNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useDetachhNodeMutation(baseOptions?: Apollo.MutationHookOptions<DetachhNodeMutation, DetachhNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DetachhNodeMutation, DetachhNodeMutationVariables>(DetachhNodeDocument, options);
      }
export type DetachhNodeMutationHookResult = ReturnType<typeof useDetachhNodeMutation>;
export type DetachhNodeMutationResult = Apollo.MutationResult<DetachhNodeMutation>;
export type DetachhNodeMutationOptions = Apollo.BaseMutationOptions<DetachhNodeMutation, DetachhNodeMutationVariables>;
export const AddNodeDocument = gql`
    mutation addNode($data: AddNodeInput!) {
  addNode(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;
export type AddNodeMutationFn = Apollo.MutationFunction<AddNodeMutation, AddNodeMutationVariables>;

/**
 * __useAddNodeMutation__
 *
 * To run a mutation, you first call `useAddNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNodeMutation, { data, loading, error }] = useAddNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNodeMutation(baseOptions?: Apollo.MutationHookOptions<AddNodeMutation, AddNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNodeMutation, AddNodeMutationVariables>(AddNodeDocument, options);
      }
export type AddNodeMutationHookResult = ReturnType<typeof useAddNodeMutation>;
export type AddNodeMutationResult = Apollo.MutationResult<AddNodeMutation>;
export type AddNodeMutationOptions = Apollo.BaseMutationOptions<AddNodeMutation, AddNodeMutationVariables>;
export const ReleaseNodeFromSiteDocument = gql`
    mutation releaseNodeFromSite($data: NodeInput!) {
  releaseNodeFromSite(data: $data) {
    success
  }
}
    `;
export type ReleaseNodeFromSiteMutationFn = Apollo.MutationFunction<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>;

/**
 * __useReleaseNodeFromSiteMutation__
 *
 * To run a mutation, you first call `useReleaseNodeFromSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useReleaseNodeFromSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [releaseNodeFromSiteMutation, { data, loading, error }] = useReleaseNodeFromSiteMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useReleaseNodeFromSiteMutation(baseOptions?: Apollo.MutationHookOptions<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>(ReleaseNodeFromSiteDocument, options);
      }
export type ReleaseNodeFromSiteMutationHookResult = ReturnType<typeof useReleaseNodeFromSiteMutation>;
export type ReleaseNodeFromSiteMutationResult = Apollo.MutationResult<ReleaseNodeFromSiteMutation>;
export type ReleaseNodeFromSiteMutationOptions = Apollo.BaseMutationOptions<ReleaseNodeFromSiteMutation, ReleaseNodeFromSiteMutationVariables>;
export const AddNodeToSiteDocument = gql`
    mutation addNodeToSite($data: AddNodeToSiteInput!) {
  addNodeToSite(data: $data) {
    success
  }
}
    `;
export type AddNodeToSiteMutationFn = Apollo.MutationFunction<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>;

/**
 * __useAddNodeToSiteMutation__
 *
 * To run a mutation, you first call `useAddNodeToSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNodeToSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNodeToSiteMutation, { data, loading, error }] = useAddNodeToSiteMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNodeToSiteMutation(baseOptions?: Apollo.MutationHookOptions<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>(AddNodeToSiteDocument, options);
      }
export type AddNodeToSiteMutationHookResult = ReturnType<typeof useAddNodeToSiteMutation>;
export type AddNodeToSiteMutationResult = Apollo.MutationResult<AddNodeToSiteMutation>;
export type AddNodeToSiteMutationOptions = Apollo.BaseMutationOptions<AddNodeToSiteMutation, AddNodeToSiteMutationVariables>;
export const UpdateNodeStateDocument = gql`
    mutation updateNodeState($data: UpdateNodeStateInput!) {
  updateNodeState(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;
export type UpdateNodeStateMutationFn = Apollo.MutationFunction<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>;

/**
 * __useUpdateNodeStateMutation__
 *
 * To run a mutation, you first call `useUpdateNodeStateMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNodeStateMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNodeStateMutation, { data, loading, error }] = useUpdateNodeStateMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateNodeStateMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>(UpdateNodeStateDocument, options);
      }
export type UpdateNodeStateMutationHookResult = ReturnType<typeof useUpdateNodeStateMutation>;
export type UpdateNodeStateMutationResult = Apollo.MutationResult<UpdateNodeStateMutation>;
export type UpdateNodeStateMutationOptions = Apollo.BaseMutationOptions<UpdateNodeStateMutation, UpdateNodeStateMutationVariables>;
export const GetNodesForSiteDocument = gql`
    query getNodesForSite($siteId: String!) {
  getNodesForSite(siteId: $siteId) {
    nodes {
      ...node
    }
  }
}
    ${NodeFragmentDoc}`;

/**
 * __useGetNodesForSiteQuery__
 *
 * To run a query within a React component, call `useGetNodesForSiteQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodesForSiteQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodesForSiteQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useGetNodesForSiteQuery(baseOptions: Apollo.QueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables> & ({ variables: GetNodesForSiteQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
      }
export function useGetNodesForSiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
        }
// @ts-ignore
export function useGetNodesForSiteSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>;
export function useGetNodesForSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodesForSiteQuery | undefined, GetNodesForSiteQueryVariables>;
export function useGetNodesForSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>(GetNodesForSiteDocument, options);
        }
export type GetNodesForSiteQueryHookResult = ReturnType<typeof useGetNodesForSiteQuery>;
export type GetNodesForSiteLazyQueryHookResult = ReturnType<typeof useGetNodesForSiteLazyQuery>;
export type GetNodesForSiteSuspenseQueryHookResult = ReturnType<typeof useGetNodesForSiteSuspenseQuery>;
export type GetNodesForSiteQueryResult = Apollo.QueryResult<GetNodesForSiteQuery, GetNodesForSiteQueryVariables>;
export const UpdateNodeDocument = gql`
    mutation UpdateNode($data: UpdateNodeInput!) {
  updateNode(data: $data) {
    ...node
  }
}
    ${NodeFragmentDoc}`;
export type UpdateNodeMutationFn = Apollo.MutationFunction<UpdateNodeMutation, UpdateNodeMutationVariables>;

/**
 * __useUpdateNodeMutation__
 *
 * To run a mutation, you first call `useUpdateNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNodeMutation, { data, loading, error }] = useUpdateNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateNodeMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNodeMutation, UpdateNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNodeMutation, UpdateNodeMutationVariables>(UpdateNodeDocument, options);
      }
export type UpdateNodeMutationHookResult = ReturnType<typeof useUpdateNodeMutation>;
export type UpdateNodeMutationResult = Apollo.MutationResult<UpdateNodeMutation>;
export type UpdateNodeMutationOptions = Apollo.BaseMutationOptions<UpdateNodeMutation, UpdateNodeMutationVariables>;
export const GetNodeAppsDocument = gql`
    query getNodeApps($data: NodeAppsChangeLogInput!) {
  getNodeApps(data: $data) {
    apps {
      name
      date
      version
      cpu
      memory
      notes
    }
    type
  }
}
    `;

/**
 * __useGetNodeAppsQuery__
 *
 * To run a query within a React component, call `useGetNodeAppsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeAppsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeAppsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetNodeAppsQuery(baseOptions: Apollo.QueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables> & ({ variables: GetNodeAppsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
      }
export function useGetNodeAppsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
        }
// @ts-ignore
export function useGetNodeAppsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeAppsQuery, GetNodeAppsQueryVariables>;
export function useGetNodeAppsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeAppsQuery | undefined, GetNodeAppsQueryVariables>;
export function useGetNodeAppsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeAppsQuery, GetNodeAppsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeAppsQuery, GetNodeAppsQueryVariables>(GetNodeAppsDocument, options);
        }
export type GetNodeAppsQueryHookResult = ReturnType<typeof useGetNodeAppsQuery>;
export type GetNodeAppsLazyQueryHookResult = ReturnType<typeof useGetNodeAppsLazyQuery>;
export type GetNodeAppsSuspenseQueryHookResult = ReturnType<typeof useGetNodeAppsSuspenseQuery>;
export type GetNodeAppsQueryResult = Apollo.QueryResult<GetNodeAppsQuery, GetNodeAppsQueryVariables>;
export const GetNodeStateDocument = gql`
    query GetNodeState($getNodeStateId: String!) {
  getNodeState(id: $getNodeStateId) {
    id
    nodeId
    previousStateId
    previousState
    currentState
    createdAt
  }
}
    `;

/**
 * __useGetNodeStateQuery__
 *
 * To run a query within a React component, call `useGetNodeStateQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNodeStateQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNodeStateQuery({
 *   variables: {
 *      getNodeStateId: // value for 'getNodeStateId'
 *   },
 * });
 */
export function useGetNodeStateQuery(baseOptions: Apollo.QueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables> & ({ variables: GetNodeStateQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNodeStateQuery, GetNodeStateQueryVariables>(GetNodeStateDocument, options);
      }
export function useGetNodeStateLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNodeStateQuery, GetNodeStateQueryVariables>(GetNodeStateDocument, options);
        }
// @ts-ignore
export function useGetNodeStateSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeStateQuery, GetNodeStateQueryVariables>;
export function useGetNodeStateSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>): Apollo.UseSuspenseQueryResult<GetNodeStateQuery | undefined, GetNodeStateQueryVariables>;
export function useGetNodeStateSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNodeStateQuery, GetNodeStateQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNodeStateQuery, GetNodeStateQueryVariables>(GetNodeStateDocument, options);
        }
export type GetNodeStateQueryHookResult = ReturnType<typeof useGetNodeStateQuery>;
export type GetNodeStateLazyQueryHookResult = ReturnType<typeof useGetNodeStateLazyQuery>;
export type GetNodeStateSuspenseQueryHookResult = ReturnType<typeof useGetNodeStateSuspenseQuery>;
export type GetNodeStateQueryResult = Apollo.QueryResult<GetNodeStateQuery, GetNodeStateQueryVariables>;
export const RestartNodeDocument = gql`
    mutation RestartNode($data: RestartNodeInputDto!) {
  restartNode(data: $data) {
    success
  }
}
    `;
export type RestartNodeMutationFn = Apollo.MutationFunction<RestartNodeMutation, RestartNodeMutationVariables>;

/**
 * __useRestartNodeMutation__
 *
 * To run a mutation, you first call `useRestartNodeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRestartNodeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [restartNodeMutation, { data, loading, error }] = useRestartNodeMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useRestartNodeMutation(baseOptions?: Apollo.MutationHookOptions<RestartNodeMutation, RestartNodeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RestartNodeMutation, RestartNodeMutationVariables>(RestartNodeDocument, options);
      }
export type RestartNodeMutationHookResult = ReturnType<typeof useRestartNodeMutation>;
export type RestartNodeMutationResult = Apollo.MutationResult<RestartNodeMutation>;
export type RestartNodeMutationOptions = Apollo.BaseMutationOptions<RestartNodeMutation, RestartNodeMutationVariables>;
export const ToggleInternetSwitchDocument = gql`
    mutation ToggleInternetSwitch($data: ToggleInternetSwitchInputDto!) {
  toggleInternetSwitch(data: $data) {
    success
  }
}
    `;
export type ToggleInternetSwitchMutationFn = Apollo.MutationFunction<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>;

/**
 * __useToggleInternetSwitchMutation__
 *
 * To run a mutation, you first call `useToggleInternetSwitchMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleInternetSwitchMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleInternetSwitchMutation, { data, loading, error }] = useToggleInternetSwitchMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleInternetSwitchMutation(baseOptions?: Apollo.MutationHookOptions<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>(ToggleInternetSwitchDocument, options);
      }
export type ToggleInternetSwitchMutationHookResult = ReturnType<typeof useToggleInternetSwitchMutation>;
export type ToggleInternetSwitchMutationResult = Apollo.MutationResult<ToggleInternetSwitchMutation>;
export type ToggleInternetSwitchMutationOptions = Apollo.BaseMutationOptions<ToggleInternetSwitchMutation, ToggleInternetSwitchMutationVariables>;
export const ToggleRfStatusDocument = gql`
    mutation ToggleRFStatus($data: ToggleRFStatusInputDto!) {
  toggleRFStatus(data: $data) {
    success
  }
}
    `;
export type ToggleRfStatusMutationFn = Apollo.MutationFunction<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>;

/**
 * __useToggleRfStatusMutation__
 *
 * To run a mutation, you first call `useToggleRfStatusMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleRfStatusMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleRfStatusMutation, { data, loading, error }] = useToggleRfStatusMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleRfStatusMutation(baseOptions?: Apollo.MutationHookOptions<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>(ToggleRfStatusDocument, options);
      }
export type ToggleRfStatusMutationHookResult = ReturnType<typeof useToggleRfStatusMutation>;
export type ToggleRfStatusMutationResult = Apollo.MutationResult<ToggleRfStatusMutation>;
export type ToggleRfStatusMutationOptions = Apollo.BaseMutationOptions<ToggleRfStatusMutation, ToggleRfStatusMutationVariables>;
export const ToggleServiceDocument = gql`
    mutation ToggleService($data: ToggleRFStatusInputDto!) {
  toggleService(data: $data) {
    success
  }
}
    `;
export type ToggleServiceMutationFn = Apollo.MutationFunction<ToggleServiceMutation, ToggleServiceMutationVariables>;

/**
 * __useToggleServiceMutation__
 *
 * To run a mutation, you first call `useToggleServiceMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleServiceMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleServiceMutation, { data, loading, error }] = useToggleServiceMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleServiceMutation(baseOptions?: Apollo.MutationHookOptions<ToggleServiceMutation, ToggleServiceMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleServiceMutation, ToggleServiceMutationVariables>(ToggleServiceDocument, options);
      }
export type ToggleServiceMutationHookResult = ReturnType<typeof useToggleServiceMutation>;
export type ToggleServiceMutationResult = Apollo.MutationResult<ToggleServiceMutation>;
export type ToggleServiceMutationOptions = Apollo.BaseMutationOptions<ToggleServiceMutation, ToggleServiceMutationVariables>;
export const GetHealthReportDocument = gql`
    query GetHealthReport($data: GetHealthReportInputDto!) {
  getHealthReport(data: $data) {
    id
    nodeId
    timestamp
    system {
      id
      healthId
      name
      value
    }
    capps {
      id
      space
      name
      tag
      status
      resources {
        id
        cappId
        name
        value
      }
    }
  }
}
    `;

/**
 * __useGetHealthReportQuery__
 *
 * To run a query within a React component, call `useGetHealthReportQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetHealthReportQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetHealthReportQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetHealthReportQuery(baseOptions: Apollo.QueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables> & ({ variables: GetHealthReportQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetHealthReportQuery, GetHealthReportQueryVariables>(GetHealthReportDocument, options);
      }
export function useGetHealthReportLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetHealthReportQuery, GetHealthReportQueryVariables>(GetHealthReportDocument, options);
        }
// @ts-ignore
export function useGetHealthReportSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetHealthReportQuery, GetHealthReportQueryVariables>;
export function useGetHealthReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetHealthReportQuery | undefined, GetHealthReportQueryVariables>;
export function useGetHealthReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetHealthReportQuery, GetHealthReportQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetHealthReportQuery, GetHealthReportQueryVariables>(GetHealthReportDocument, options);
        }
export type GetHealthReportQueryHookResult = ReturnType<typeof useGetHealthReportQuery>;
export type GetHealthReportLazyQueryHookResult = ReturnType<typeof useGetHealthReportLazyQuery>;
export type GetHealthReportSuspenseQueryHookResult = ReturnType<typeof useGetHealthReportSuspenseQuery>;
export type GetHealthReportQueryResult = Apollo.QueryResult<GetHealthReportQuery, GetHealthReportQueryVariables>;
export const GetMembersDocument = gql`
    query GetMembers {
  getMembers {
    members {
      ...member
      name
      email
    }
  }
}
    ${MemberFragmentDoc}`;

/**
 * __useGetMembersQuery__
 *
 * To run a query within a React component, call `useGetMembersQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMembersQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMembersQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetMembersQuery(baseOptions?: Apollo.QueryHookOptions<GetMembersQuery, GetMembersQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMembersQuery, GetMembersQueryVariables>(GetMembersDocument, options);
      }
export function useGetMembersLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMembersQuery, GetMembersQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMembersQuery, GetMembersQueryVariables>(GetMembersDocument, options);
        }
// @ts-ignore
export function useGetMembersSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetMembersQuery, GetMembersQueryVariables>): Apollo.UseSuspenseQueryResult<GetMembersQuery, GetMembersQueryVariables>;
export function useGetMembersSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMembersQuery, GetMembersQueryVariables>): Apollo.UseSuspenseQueryResult<GetMembersQuery | undefined, GetMembersQueryVariables>;
export function useGetMembersSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMembersQuery, GetMembersQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetMembersQuery, GetMembersQueryVariables>(GetMembersDocument, options);
        }
export type GetMembersQueryHookResult = ReturnType<typeof useGetMembersQuery>;
export type GetMembersLazyQueryHookResult = ReturnType<typeof useGetMembersLazyQuery>;
export type GetMembersSuspenseQueryHookResult = ReturnType<typeof useGetMembersSuspenseQuery>;
export type GetMembersQueryResult = Apollo.QueryResult<GetMembersQuery, GetMembersQueryVariables>;
export const GetMemberDocument = gql`
    query GetMember($memberId: String!) {
  getMember(id: $memberId) {
    ...member
  }
}
    ${MemberFragmentDoc}`;

/**
 * __useGetMemberQuery__
 *
 * To run a query within a React component, call `useGetMemberQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMemberQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMemberQuery({
 *   variables: {
 *      memberId: // value for 'memberId'
 *   },
 * });
 */
export function useGetMemberQuery(baseOptions: Apollo.QueryHookOptions<GetMemberQuery, GetMemberQueryVariables> & ({ variables: GetMemberQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMemberQuery, GetMemberQueryVariables>(GetMemberDocument, options);
      }
export function useGetMemberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMemberQuery, GetMemberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMemberQuery, GetMemberQueryVariables>(GetMemberDocument, options);
        }
// @ts-ignore
export function useGetMemberSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetMemberQuery, GetMemberQueryVariables>): Apollo.UseSuspenseQueryResult<GetMemberQuery, GetMemberQueryVariables>;
export function useGetMemberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMemberQuery, GetMemberQueryVariables>): Apollo.UseSuspenseQueryResult<GetMemberQuery | undefined, GetMemberQueryVariables>;
export function useGetMemberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMemberQuery, GetMemberQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetMemberQuery, GetMemberQueryVariables>(GetMemberDocument, options);
        }
export type GetMemberQueryHookResult = ReturnType<typeof useGetMemberQuery>;
export type GetMemberLazyQueryHookResult = ReturnType<typeof useGetMemberLazyQuery>;
export type GetMemberSuspenseQueryHookResult = ReturnType<typeof useGetMemberSuspenseQuery>;
export type GetMemberQueryResult = Apollo.QueryResult<GetMemberQuery, GetMemberQueryVariables>;
export const AddMemberDocument = gql`
    mutation addMember($data: AddMemberInputDto!) {
  addMember(data: $data) {
    ...member
  }
}
    ${MemberFragmentDoc}`;
export type AddMemberMutationFn = Apollo.MutationFunction<AddMemberMutation, AddMemberMutationVariables>;

/**
 * __useAddMemberMutation__
 *
 * To run a mutation, you first call `useAddMemberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddMemberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addMemberMutation, { data, loading, error }] = useAddMemberMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddMemberMutation(baseOptions?: Apollo.MutationHookOptions<AddMemberMutation, AddMemberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddMemberMutation, AddMemberMutationVariables>(AddMemberDocument, options);
      }
export type AddMemberMutationHookResult = ReturnType<typeof useAddMemberMutation>;
export type AddMemberMutationResult = Apollo.MutationResult<AddMemberMutation>;
export type AddMemberMutationOptions = Apollo.BaseMutationOptions<AddMemberMutation, AddMemberMutationVariables>;
export const RemoveMemberDocument = gql`
    mutation removeMember($memberId: String!) {
  removeMember(id: $memberId) {
    success
  }
}
    `;
export type RemoveMemberMutationFn = Apollo.MutationFunction<RemoveMemberMutation, RemoveMemberMutationVariables>;

/**
 * __useRemoveMemberMutation__
 *
 * To run a mutation, you first call `useRemoveMemberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemoveMemberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removeMemberMutation, { data, loading, error }] = useRemoveMemberMutation({
 *   variables: {
 *      memberId: // value for 'memberId'
 *   },
 * });
 */
export function useRemoveMemberMutation(baseOptions?: Apollo.MutationHookOptions<RemoveMemberMutation, RemoveMemberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemoveMemberMutation, RemoveMemberMutationVariables>(RemoveMemberDocument, options);
      }
export type RemoveMemberMutationHookResult = ReturnType<typeof useRemoveMemberMutation>;
export type RemoveMemberMutationResult = Apollo.MutationResult<RemoveMemberMutation>;
export type RemoveMemberMutationOptions = Apollo.BaseMutationOptions<RemoveMemberMutation, RemoveMemberMutationVariables>;
export const UpdateMemberDocument = gql`
    mutation updateMember($memberId: String!, $data: UpdateMemberInputDto!) {
  updateMember(memberId: $memberId, data: $data) {
    success
  }
}
    `;
export type UpdateMemberMutationFn = Apollo.MutationFunction<UpdateMemberMutation, UpdateMemberMutationVariables>;

/**
 * __useUpdateMemberMutation__
 *
 * To run a mutation, you first call `useUpdateMemberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateMemberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateMemberMutation, { data, loading, error }] = useUpdateMemberMutation({
 *   variables: {
 *      memberId: // value for 'memberId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateMemberMutation(baseOptions?: Apollo.MutationHookOptions<UpdateMemberMutation, UpdateMemberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateMemberMutation, UpdateMemberMutationVariables>(UpdateMemberDocument, options);
      }
export type UpdateMemberMutationHookResult = ReturnType<typeof useUpdateMemberMutation>;
export type UpdateMemberMutationResult = Apollo.MutationResult<UpdateMemberMutation>;
export type UpdateMemberMutationOptions = Apollo.BaseMutationOptions<UpdateMemberMutation, UpdateMemberMutationVariables>;
export const GetMemberByUserIdDocument = gql`
    query GetMemberByUserId($userId: String!) {
  getMemberByUserId(userId: $userId) {
    userId
    name
    email
    memberId
    isDeactivated
    role
    memberSince
  }
}
    `;

/**
 * __useGetMemberByUserIdQuery__
 *
 * To run a query within a React component, call `useGetMemberByUserIdQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetMemberByUserIdQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetMemberByUserIdQuery({
 *   variables: {
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useGetMemberByUserIdQuery(baseOptions: Apollo.QueryHookOptions<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables> & ({ variables: GetMemberByUserIdQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>(GetMemberByUserIdDocument, options);
      }
export function useGetMemberByUserIdLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>(GetMemberByUserIdDocument, options);
        }
// @ts-ignore
export function useGetMemberByUserIdSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>;
export function useGetMemberByUserIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetMemberByUserIdQuery | undefined, GetMemberByUserIdQueryVariables>;
export function useGetMemberByUserIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>(GetMemberByUserIdDocument, options);
        }
export type GetMemberByUserIdQueryHookResult = ReturnType<typeof useGetMemberByUserIdQuery>;
export type GetMemberByUserIdLazyQueryHookResult = ReturnType<typeof useGetMemberByUserIdLazyQuery>;
export type GetMemberByUserIdSuspenseQueryHookResult = ReturnType<typeof useGetMemberByUserIdSuspenseQuery>;
export type GetMemberByUserIdQueryResult = Apollo.QueryResult<GetMemberByUserIdQuery, GetMemberByUserIdQueryVariables>;
export const GetOrgsDocument = gql`
    query getOrgs {
  getOrgs {
    user
    ownerOf {
      ...Org
    }
    memberOf {
      ...Org
    }
  }
}
    ${OrgFragmentDoc}`;

/**
 * __useGetOrgsQuery__
 *
 * To run a query within a React component, call `useGetOrgsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrgsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrgsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetOrgsQuery(baseOptions?: Apollo.QueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrgsQuery, GetOrgsQueryVariables>(GetOrgsDocument, options);
      }
export function useGetOrgsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrgsQuery, GetOrgsQueryVariables>(GetOrgsDocument, options);
        }
// @ts-ignore
export function useGetOrgsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>): Apollo.UseSuspenseQueryResult<GetOrgsQuery, GetOrgsQueryVariables>;
export function useGetOrgsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>): Apollo.UseSuspenseQueryResult<GetOrgsQuery | undefined, GetOrgsQueryVariables>;
export function useGetOrgsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgsQuery, GetOrgsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetOrgsQuery, GetOrgsQueryVariables>(GetOrgsDocument, options);
        }
export type GetOrgsQueryHookResult = ReturnType<typeof useGetOrgsQuery>;
export type GetOrgsLazyQueryHookResult = ReturnType<typeof useGetOrgsLazyQuery>;
export type GetOrgsSuspenseQueryHookResult = ReturnType<typeof useGetOrgsSuspenseQuery>;
export type GetOrgsQueryResult = Apollo.QueryResult<GetOrgsQuery, GetOrgsQueryVariables>;
export const GetOrgDocument = gql`
    query getOrg {
  getOrg {
    ...Org
  }
}
    ${OrgFragmentDoc}`;

/**
 * __useGetOrgQuery__
 *
 * To run a query within a React component, call `useGetOrgQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetOrgQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetOrgQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetOrgQuery(baseOptions?: Apollo.QueryHookOptions<GetOrgQuery, GetOrgQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetOrgQuery, GetOrgQueryVariables>(GetOrgDocument, options);
      }
export function useGetOrgLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetOrgQuery, GetOrgQueryVariables>(GetOrgDocument, options);
        }
// @ts-ignore
export function useGetOrgSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>): Apollo.UseSuspenseQueryResult<GetOrgQuery, GetOrgQueryVariables>;
export function useGetOrgSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>): Apollo.UseSuspenseQueryResult<GetOrgQuery | undefined, GetOrgQueryVariables>;
export function useGetOrgSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetOrgQuery, GetOrgQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetOrgQuery, GetOrgQueryVariables>(GetOrgDocument, options);
        }
export type GetOrgQueryHookResult = ReturnType<typeof useGetOrgQuery>;
export type GetOrgLazyQueryHookResult = ReturnType<typeof useGetOrgLazyQuery>;
export type GetOrgSuspenseQueryHookResult = ReturnType<typeof useGetOrgSuspenseQuery>;
export type GetOrgQueryResult = Apollo.QueryResult<GetOrgQuery, GetOrgQueryVariables>;
export const GetPackagesDocument = gql`
    query getPackages {
  getPackages {
    packages {
      ...Package
    }
  }
}
    ${PackageFragmentDoc}`;

/**
 * __useGetPackagesQuery__
 *
 * To run a query within a React component, call `useGetPackagesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackagesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackagesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetPackagesQuery(baseOptions?: Apollo.QueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackagesQuery, GetPackagesQueryVariables>(GetPackagesDocument, options);
      }
export function useGetPackagesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackagesQuery, GetPackagesQueryVariables>(GetPackagesDocument, options);
        }
// @ts-ignore
export function useGetPackagesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagesQuery, GetPackagesQueryVariables>;
export function useGetPackagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagesQuery | undefined, GetPackagesQueryVariables>;
export function useGetPackagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesQuery, GetPackagesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackagesQuery, GetPackagesQueryVariables>(GetPackagesDocument, options);
        }
export type GetPackagesQueryHookResult = ReturnType<typeof useGetPackagesQuery>;
export type GetPackagesLazyQueryHookResult = ReturnType<typeof useGetPackagesLazyQuery>;
export type GetPackagesSuspenseQueryHookResult = ReturnType<typeof useGetPackagesSuspenseQuery>;
export type GetPackagesQueryResult = Apollo.QueryResult<GetPackagesQuery, GetPackagesQueryVariables>;
export const GetPackageDocument = gql`
    query getPackage($packageId: String!) {
  getPackage(packageId: $packageId) {
    ...Package
  }
}
    ${PackageFragmentDoc}`;

/**
 * __useGetPackageQuery__
 *
 * To run a query within a React component, call `useGetPackageQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackageQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackageQuery({
 *   variables: {
 *      packageId: // value for 'packageId'
 *   },
 * });
 */
export function useGetPackageQuery(baseOptions: Apollo.QueryHookOptions<GetPackageQuery, GetPackageQueryVariables> & ({ variables: GetPackageQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
      }
export function useGetPackageLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
        }
// @ts-ignore
export function useGetPackageSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackageQuery, GetPackageQueryVariables>;
export function useGetPackageSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackageQuery | undefined, GetPackageQueryVariables>;
export function useGetPackageSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackageQuery, GetPackageQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackageQuery, GetPackageQueryVariables>(GetPackageDocument, options);
        }
export type GetPackageQueryHookResult = ReturnType<typeof useGetPackageQuery>;
export type GetPackageLazyQueryHookResult = ReturnType<typeof useGetPackageLazyQuery>;
export type GetPackageSuspenseQueryHookResult = ReturnType<typeof useGetPackageSuspenseQuery>;
export type GetPackageQueryResult = Apollo.QueryResult<GetPackageQuery, GetPackageQueryVariables>;
export const GetSimsBySubscriberDocument = gql`
    query getSimsBySubscriber($data: GetSimBySubscriberInputDto!) {
  getSimsBySubscriber(data: $data) {
    ...SubscriberSims
  }
}
    ${SubscriberSimsFragmentDoc}`;

/**
 * __useGetSimsBySubscriberQuery__
 *
 * To run a query within a React component, call `useGetSimsBySubscriberQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimsBySubscriberQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimsBySubscriberQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimsBySubscriberQuery(baseOptions: Apollo.QueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables> & ({ variables: GetSimsBySubscriberQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
      }
export function useGetSimsBySubscriberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
        }
// @ts-ignore
export function useGetSimsBySubscriberSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>;
export function useGetSimsBySubscriberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsBySubscriberQuery | undefined, GetSimsBySubscriberQueryVariables>;
export function useGetSimsBySubscriberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>(GetSimsBySubscriberDocument, options);
        }
export type GetSimsBySubscriberQueryHookResult = ReturnType<typeof useGetSimsBySubscriberQuery>;
export type GetSimsBySubscriberLazyQueryHookResult = ReturnType<typeof useGetSimsBySubscriberLazyQuery>;
export type GetSimsBySubscriberSuspenseQueryHookResult = ReturnType<typeof useGetSimsBySubscriberSuspenseQuery>;
export type GetSimsBySubscriberQueryResult = Apollo.QueryResult<GetSimsBySubscriberQuery, GetSimsBySubscriberQueryVariables>;
export const AddPackageDocument = gql`
    mutation addPackage($data: AddPackageInputDto!) {
  addPackage(data: $data) {
    ...Package
  }
}
    ${PackageFragmentDoc}`;
export type AddPackageMutationFn = Apollo.MutationFunction<AddPackageMutation, AddPackageMutationVariables>;

/**
 * __useAddPackageMutation__
 *
 * To run a mutation, you first call `useAddPackageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPackageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPackageMutation, { data, loading, error }] = useAddPackageMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddPackageMutation(baseOptions?: Apollo.MutationHookOptions<AddPackageMutation, AddPackageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPackageMutation, AddPackageMutationVariables>(AddPackageDocument, options);
      }
export type AddPackageMutationHookResult = ReturnType<typeof useAddPackageMutation>;
export type AddPackageMutationResult = Apollo.MutationResult<AddPackageMutation>;
export type AddPackageMutationOptions = Apollo.BaseMutationOptions<AddPackageMutation, AddPackageMutationVariables>;
export const RemovePackageForSimDocument = gql`
    mutation removePackageForSim($data: RemovePackageFormSimInputDto!) {
  removePackageForSim(data: $data) {
    packageId
  }
}
    `;
export type RemovePackageForSimMutationFn = Apollo.MutationFunction<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>;

/**
 * __useRemovePackageForSimMutation__
 *
 * To run a mutation, you first call `useRemovePackageForSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRemovePackageForSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [removePackageForSimMutation, { data, loading, error }] = useRemovePackageForSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useRemovePackageForSimMutation(baseOptions?: Apollo.MutationHookOptions<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>(RemovePackageForSimDocument, options);
      }
export type RemovePackageForSimMutationHookResult = ReturnType<typeof useRemovePackageForSimMutation>;
export type RemovePackageForSimMutationResult = Apollo.MutationResult<RemovePackageForSimMutation>;
export type RemovePackageForSimMutationOptions = Apollo.BaseMutationOptions<RemovePackageForSimMutation, RemovePackageForSimMutationVariables>;
export const DeletePackageDocument = gql`
    mutation deletePackage($packageId: String!) {
  deletePackage(packageId: $packageId) {
    uuid
  }
}
    `;
export type DeletePackageMutationFn = Apollo.MutationFunction<DeletePackageMutation, DeletePackageMutationVariables>;

/**
 * __useDeletePackageMutation__
 *
 * To run a mutation, you first call `useDeletePackageMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeletePackageMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deletePackageMutation, { data, loading, error }] = useDeletePackageMutation({
 *   variables: {
 *      packageId: // value for 'packageId'
 *   },
 * });
 */
export function useDeletePackageMutation(baseOptions?: Apollo.MutationHookOptions<DeletePackageMutation, DeletePackageMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeletePackageMutation, DeletePackageMutationVariables>(DeletePackageDocument, options);
      }
export type DeletePackageMutationHookResult = ReturnType<typeof useDeletePackageMutation>;
export type DeletePackageMutationResult = Apollo.MutationResult<DeletePackageMutation>;
export type DeletePackageMutationOptions = Apollo.BaseMutationOptions<DeletePackageMutation, DeletePackageMutationVariables>;
export const GetPackagesForSimDocument = gql`
    query getPackagesForSim($data: GetPackagesForSimInputDto!) {
  getPackagesForSim(data: $data) {
    sim_id
    packages {
      ...SimPackages
    }
  }
}
    ${SimPackagesFragmentDoc}`;

/**
 * __useGetPackagesForSimQuery__
 *
 * To run a query within a React component, call `useGetPackagesForSimQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPackagesForSimQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPackagesForSimQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetPackagesForSimQuery(baseOptions: Apollo.QueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables> & ({ variables: GetPackagesForSimQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
      }
export function useGetPackagesForSimLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
        }
// @ts-ignore
export function useGetPackagesForSimSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>;
export function useGetPackagesForSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>): Apollo.UseSuspenseQueryResult<GetPackagesForSimQuery | undefined, GetPackagesForSimQueryVariables>;
export function useGetPackagesForSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>(GetPackagesForSimDocument, options);
        }
export type GetPackagesForSimQueryHookResult = ReturnType<typeof useGetPackagesForSimQuery>;
export type GetPackagesForSimLazyQueryHookResult = ReturnType<typeof useGetPackagesForSimLazyQuery>;
export type GetPackagesForSimSuspenseQueryHookResult = ReturnType<typeof useGetPackagesForSimSuspenseQuery>;
export type GetPackagesForSimQueryResult = Apollo.QueryResult<GetPackagesForSimQuery, GetPackagesForSimQueryVariables>;
export const AddPackagesToSimDocument = gql`
    mutation addPackagesToSim($data: AddPackagesToSimInputDto!) {
  addPackagesToSim(data: $data) {
    packages {
      packageId
    }
  }
}
    `;
export type AddPackagesToSimMutationFn = Apollo.MutationFunction<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>;

/**
 * __useAddPackagesToSimMutation__
 *
 * To run a mutation, you first call `useAddPackagesToSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddPackagesToSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addPackagesToSimMutation, { data, loading, error }] = useAddPackagesToSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddPackagesToSimMutation(baseOptions?: Apollo.MutationHookOptions<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>(AddPackagesToSimDocument, options);
      }
export type AddPackagesToSimMutationHookResult = ReturnType<typeof useAddPackagesToSimMutation>;
export type AddPackagesToSimMutationResult = Apollo.MutationResult<AddPackagesToSimMutation>;
export type AddPackagesToSimMutationOptions = Apollo.BaseMutationOptions<AddPackagesToSimMutation, AddPackagesToSimMutationVariables>;
export const DeleteSimDocument = gql`
    mutation deleteSim($data: DeleteSimInputDto!) {
  deleteSim(data: $data) {
    simId
  }
}
    `;
export type DeleteSimMutationFn = Apollo.MutationFunction<DeleteSimMutation, DeleteSimMutationVariables>;

/**
 * __useDeleteSimMutation__
 *
 * To run a mutation, you first call `useDeleteSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteSimMutation, { data, loading, error }] = useDeleteSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useDeleteSimMutation(baseOptions?: Apollo.MutationHookOptions<DeleteSimMutation, DeleteSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteSimMutation, DeleteSimMutationVariables>(DeleteSimDocument, options);
      }
export type DeleteSimMutationHookResult = ReturnType<typeof useDeleteSimMutation>;
export type DeleteSimMutationResult = Apollo.MutationResult<DeleteSimMutation>;
export type DeleteSimMutationOptions = Apollo.BaseMutationOptions<DeleteSimMutation, DeleteSimMutationVariables>;
export const UpdatePacakgeDocument = gql`
    mutation updatePacakge($packageId: String!, $data: UpdatePackageInputDto!) {
  updatePackage(packageId: $packageId, data: $data) {
    ...Package
  }
}
    ${PackageFragmentDoc}`;
export type UpdatePacakgeMutationFn = Apollo.MutationFunction<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>;

/**
 * __useUpdatePacakgeMutation__
 *
 * To run a mutation, you first call `useUpdatePacakgeMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdatePacakgeMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updatePacakgeMutation, { data, loading, error }] = useUpdatePacakgeMutation({
 *   variables: {
 *      packageId: // value for 'packageId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdatePacakgeMutation(baseOptions?: Apollo.MutationHookOptions<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>(UpdatePacakgeDocument, options);
      }
export type UpdatePacakgeMutationHookResult = ReturnType<typeof useUpdatePacakgeMutation>;
export type UpdatePacakgeMutationResult = Apollo.MutationResult<UpdatePacakgeMutation>;
export type UpdatePacakgeMutationOptions = Apollo.BaseMutationOptions<UpdatePacakgeMutation, UpdatePacakgeMutationVariables>;
export const UpdatePaymentDocument = gql`
    mutation UpdatePayment($data: UpdatePaymentInputDto!) {
  updatePayment(data: $data) {
    id
    itemId
    itemType
    amount
    currency
    paymentMethod
    depositedAmount
    paidAt
    payerName
    payerEmail
    payerPhone
    correspondent
    country
    description
    status
    failureReason
    createdAt
  }
}
    `;
export type UpdatePaymentMutationFn = Apollo.MutationFunction<UpdatePaymentMutation, UpdatePaymentMutationVariables>;

/**
 * __useUpdatePaymentMutation__
 *
 * To run a mutation, you first call `useUpdatePaymentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdatePaymentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updatePaymentMutation, { data, loading, error }] = useUpdatePaymentMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdatePaymentMutation(baseOptions?: Apollo.MutationHookOptions<UpdatePaymentMutation, UpdatePaymentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdatePaymentMutation, UpdatePaymentMutationVariables>(UpdatePaymentDocument, options);
      }
export type UpdatePaymentMutationHookResult = ReturnType<typeof useUpdatePaymentMutation>;
export type UpdatePaymentMutationResult = Apollo.MutationResult<UpdatePaymentMutation>;
export type UpdatePaymentMutationOptions = Apollo.BaseMutationOptions<UpdatePaymentMutation, UpdatePaymentMutationVariables>;
export const ProcessPaymentDocument = gql`
    mutation ProcessPayment($data: ProcessPaymentInputDto!) {
  processPayment(data: $data) {
    payment {
      id
      itemId
      itemType
      amount
      currency
      paymentMethod
      depositedAmount
      paidAt
      payerName
      payerEmail
      payerPhone
      correspondent
      country
      description
      status
      failureReason
      createdAt
    }
  }
}
    `;
export type ProcessPaymentMutationFn = Apollo.MutationFunction<ProcessPaymentMutation, ProcessPaymentMutationVariables>;

/**
 * __useProcessPaymentMutation__
 *
 * To run a mutation, you first call `useProcessPaymentMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useProcessPaymentMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [processPaymentMutation, { data, loading, error }] = useProcessPaymentMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useProcessPaymentMutation(baseOptions?: Apollo.MutationHookOptions<ProcessPaymentMutation, ProcessPaymentMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ProcessPaymentMutation, ProcessPaymentMutationVariables>(ProcessPaymentDocument, options);
      }
export type ProcessPaymentMutationHookResult = ReturnType<typeof useProcessPaymentMutation>;
export type ProcessPaymentMutationResult = Apollo.MutationResult<ProcessPaymentMutation>;
export type ProcessPaymentMutationOptions = Apollo.BaseMutationOptions<ProcessPaymentMutation, ProcessPaymentMutationVariables>;
export const GetPaymentDocument = gql`
    query GetPayment($paymentId: String!) {
  getPayment(paymentId: $paymentId) {
    ...payment
  }
}
    ${PaymentFragmentDoc}`;

/**
 * __useGetPaymentQuery__
 *
 * To run a query within a React component, call `useGetPaymentQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPaymentQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPaymentQuery({
 *   variables: {
 *      paymentId: // value for 'paymentId'
 *   },
 * });
 */
export function useGetPaymentQuery(baseOptions: Apollo.QueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables> & ({ variables: GetPaymentQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPaymentQuery, GetPaymentQueryVariables>(GetPaymentDocument, options);
      }
export function useGetPaymentLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPaymentQuery, GetPaymentQueryVariables>(GetPaymentDocument, options);
        }
// @ts-ignore
export function useGetPaymentSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>): Apollo.UseSuspenseQueryResult<GetPaymentQuery, GetPaymentQueryVariables>;
export function useGetPaymentSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>): Apollo.UseSuspenseQueryResult<GetPaymentQuery | undefined, GetPaymentQueryVariables>;
export function useGetPaymentSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentQuery, GetPaymentQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPaymentQuery, GetPaymentQueryVariables>(GetPaymentDocument, options);
        }
export type GetPaymentQueryHookResult = ReturnType<typeof useGetPaymentQuery>;
export type GetPaymentLazyQueryHookResult = ReturnType<typeof useGetPaymentLazyQuery>;
export type GetPaymentSuspenseQueryHookResult = ReturnType<typeof useGetPaymentSuspenseQuery>;
export type GetPaymentQueryResult = Apollo.QueryResult<GetPaymentQuery, GetPaymentQueryVariables>;
export const GetPaymentsDocument = gql`
    query GetPayments($data: GetPaymentsInputDto!) {
  getPayments(data: $data) {
    payments {
      ...payment
    }
  }
}
    ${PaymentFragmentDoc}`;

/**
 * __useGetPaymentsQuery__
 *
 * To run a query within a React component, call `useGetPaymentsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetPaymentsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetPaymentsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetPaymentsQuery(baseOptions: Apollo.QueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables> & ({ variables: GetPaymentsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetPaymentsQuery, GetPaymentsQueryVariables>(GetPaymentsDocument, options);
      }
export function useGetPaymentsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetPaymentsQuery, GetPaymentsQueryVariables>(GetPaymentsDocument, options);
        }
// @ts-ignore
export function useGetPaymentsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>): Apollo.UseSuspenseQueryResult<GetPaymentsQuery, GetPaymentsQueryVariables>;
export function useGetPaymentsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>): Apollo.UseSuspenseQueryResult<GetPaymentsQuery | undefined, GetPaymentsQueryVariables>;
export function useGetPaymentsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetPaymentsQuery, GetPaymentsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetPaymentsQuery, GetPaymentsQueryVariables>(GetPaymentsDocument, options);
        }
export type GetPaymentsQueryHookResult = ReturnType<typeof useGetPaymentsQuery>;
export type GetPaymentsLazyQueryHookResult = ReturnType<typeof useGetPaymentsLazyQuery>;
export type GetPaymentsSuspenseQueryHookResult = ReturnType<typeof useGetPaymentsSuspenseQuery>;
export type GetPaymentsQueryResult = Apollo.QueryResult<GetPaymentsQuery, GetPaymentsQueryVariables>;
export const GetReportsDocument = gql`
    query GetReports($data: GetReportsInputDto!) {
  getReports(data: $data) {
    reports {
      id
      ownerId
      ownerType
      networkId
      period
      type
      rawReport {
        ...rawReport
      }
      isPaid
      createdAt
    }
  }
}
    ${RawReportFragmentDoc}`;

/**
 * __useGetReportsQuery__
 *
 * To run a query within a React component, call `useGetReportsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetReportsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetReportsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetReportsQuery(baseOptions: Apollo.QueryHookOptions<GetReportsQuery, GetReportsQueryVariables> & ({ variables: GetReportsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetReportsQuery, GetReportsQueryVariables>(GetReportsDocument, options);
      }
export function useGetReportsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetReportsQuery, GetReportsQueryVariables>(GetReportsDocument, options);
        }
// @ts-ignore
export function useGetReportsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>): Apollo.UseSuspenseQueryResult<GetReportsQuery, GetReportsQueryVariables>;
export function useGetReportsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>): Apollo.UseSuspenseQueryResult<GetReportsQuery | undefined, GetReportsQueryVariables>;
export function useGetReportsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportsQuery, GetReportsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetReportsQuery, GetReportsQueryVariables>(GetReportsDocument, options);
        }
export type GetReportsQueryHookResult = ReturnType<typeof useGetReportsQuery>;
export type GetReportsLazyQueryHookResult = ReturnType<typeof useGetReportsLazyQuery>;
export type GetReportsSuspenseQueryHookResult = ReturnType<typeof useGetReportsSuspenseQuery>;
export type GetReportsQueryResult = Apollo.QueryResult<GetReportsQuery, GetReportsQueryVariables>;
export const GetReportDocument = gql`
    query GetReport($id: String!) {
  getReport(id: $id) {
    report {
      id
      ownerId
      ownerType
      networkId
      period
      type
      rawReport {
        ...rawReport
      }
      isPaid
      createdAt
    }
  }
}
    ${RawReportFragmentDoc}`;

/**
 * __useGetReportQuery__
 *
 * To run a query within a React component, call `useGetReportQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetReportQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetReportQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useGetReportQuery(baseOptions: Apollo.QueryHookOptions<GetReportQuery, GetReportQueryVariables> & ({ variables: GetReportQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetReportQuery, GetReportQueryVariables>(GetReportDocument, options);
      }
export function useGetReportLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetReportQuery, GetReportQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetReportQuery, GetReportQueryVariables>(GetReportDocument, options);
        }
// @ts-ignore
export function useGetReportSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetReportQuery, GetReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetReportQuery, GetReportQueryVariables>;
export function useGetReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportQuery, GetReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetReportQuery | undefined, GetReportQueryVariables>;
export function useGetReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetReportQuery, GetReportQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetReportQuery, GetReportQueryVariables>(GetReportDocument, options);
        }
export type GetReportQueryHookResult = ReturnType<typeof useGetReportQuery>;
export type GetReportLazyQueryHookResult = ReturnType<typeof useGetReportLazyQuery>;
export type GetReportSuspenseQueryHookResult = ReturnType<typeof useGetReportSuspenseQuery>;
export type GetReportQueryResult = Apollo.QueryResult<GetReportQuery, GetReportQueryVariables>;
export const GetSimPoolStatsDocument = gql`
    query GetSimPoolStats($data: GetSimsInput!) {
  getSimPoolStats(data: $data) {
    total
    available
    consumed
    failed
    esim
    physical
  }
}
    `;

/**
 * __useGetSimPoolStatsQuery__
 *
 * To run a query within a React component, call `useGetSimPoolStatsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimPoolStatsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimPoolStatsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimPoolStatsQuery(baseOptions: Apollo.QueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables> & ({ variables: GetSimPoolStatsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>(GetSimPoolStatsDocument, options);
      }
export function useGetSimPoolStatsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>(GetSimPoolStatsDocument, options);
        }
// @ts-ignore
export function useGetSimPoolStatsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>;
export function useGetSimPoolStatsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimPoolStatsQuery | undefined, GetSimPoolStatsQueryVariables>;
export function useGetSimPoolStatsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>(GetSimPoolStatsDocument, options);
        }
export type GetSimPoolStatsQueryHookResult = ReturnType<typeof useGetSimPoolStatsQuery>;
export type GetSimPoolStatsLazyQueryHookResult = ReturnType<typeof useGetSimPoolStatsLazyQuery>;
export type GetSimPoolStatsSuspenseQueryHookResult = ReturnType<typeof useGetSimPoolStatsSuspenseQuery>;
export type GetSimPoolStatsQueryResult = Apollo.QueryResult<GetSimPoolStatsQuery, GetSimPoolStatsQueryVariables>;
export const GetSimsFromPoolDocument = gql`
    query GetSimsFromPool($data: GetSimsInput!) {
  getSimsFromPool(data: $data) {
    sims {
      id
      qrCode
      iccid
      msisdn
      isAllocated
      isFailed
      simType
      smApAddress
      activationCode
      createdAt
      deletedAt
      updatedAt
      isPhysical
    }
  }
}
    `;

/**
 * __useGetSimsFromPoolQuery__
 *
 * To run a query within a React component, call `useGetSimsFromPoolQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimsFromPoolQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimsFromPoolQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimsFromPoolQuery(baseOptions: Apollo.QueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables> & ({ variables: GetSimsFromPoolQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>(GetSimsFromPoolDocument, options);
      }
export function useGetSimsFromPoolLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>(GetSimsFromPoolDocument, options);
        }
// @ts-ignore
export function useGetSimsFromPoolSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>;
export function useGetSimsFromPoolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsFromPoolQuery | undefined, GetSimsFromPoolQueryVariables>;
export function useGetSimsFromPoolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>(GetSimsFromPoolDocument, options);
        }
export type GetSimsFromPoolQueryHookResult = ReturnType<typeof useGetSimsFromPoolQuery>;
export type GetSimsFromPoolLazyQueryHookResult = ReturnType<typeof useGetSimsFromPoolLazyQuery>;
export type GetSimsFromPoolSuspenseQueryHookResult = ReturnType<typeof useGetSimsFromPoolSuspenseQuery>;
export type GetSimsFromPoolQueryResult = Apollo.QueryResult<GetSimsFromPoolQuery, GetSimsFromPoolQueryVariables>;
export const UploadSimsDocument = gql`
    mutation uploadSims($data: UploadSimsInputDto!) {
  uploadSims(data: $data) {
    iccid
  }
}
    `;
export type UploadSimsMutationFn = Apollo.MutationFunction<UploadSimsMutation, UploadSimsMutationVariables>;

/**
 * __useUploadSimsMutation__
 *
 * To run a mutation, you first call `useUploadSimsMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUploadSimsMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [uploadSimsMutation, { data, loading, error }] = useUploadSimsMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUploadSimsMutation(baseOptions?: Apollo.MutationHookOptions<UploadSimsMutation, UploadSimsMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UploadSimsMutation, UploadSimsMutationVariables>(UploadSimsDocument, options);
      }
export type UploadSimsMutationHookResult = ReturnType<typeof useUploadSimsMutation>;
export type UploadSimsMutationResult = Apollo.MutationResult<UploadSimsMutation>;
export type UploadSimsMutationOptions = Apollo.BaseMutationOptions<UploadSimsMutation, UploadSimsMutationVariables>;
export const AllocateSimDocument = gql`
    mutation allocateSim($data: AllocateSimInputDto!) {
  allocateSim(data: $data) {
    ...SimAllocation
  }
}
    ${SimAllocationFragmentDoc}`;
export type AllocateSimMutationFn = Apollo.MutationFunction<AllocateSimMutation, AllocateSimMutationVariables>;

/**
 * __useAllocateSimMutation__
 *
 * To run a mutation, you first call `useAllocateSimMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAllocateSimMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [allocateSimMutation, { data, loading, error }] = useAllocateSimMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAllocateSimMutation(baseOptions?: Apollo.MutationHookOptions<AllocateSimMutation, AllocateSimMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AllocateSimMutation, AllocateSimMutationVariables>(AllocateSimDocument, options);
      }
export type AllocateSimMutationHookResult = ReturnType<typeof useAllocateSimMutation>;
export type AllocateSimMutationResult = Apollo.MutationResult<AllocateSimMutation>;
export type AllocateSimMutationOptions = Apollo.BaseMutationOptions<AllocateSimMutation, AllocateSimMutationVariables>;
export const ToggleSimStatusDocument = gql`
    mutation toggleSimStatus($data: ToggleSimStatusInputDto!) {
  toggleSimStatus(data: $data) {
    simId
  }
}
    `;
export type ToggleSimStatusMutationFn = Apollo.MutationFunction<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>;

/**
 * __useToggleSimStatusMutation__
 *
 * To run a mutation, you first call `useToggleSimStatusMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useToggleSimStatusMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [toggleSimStatusMutation, { data, loading, error }] = useToggleSimStatusMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useToggleSimStatusMutation(baseOptions?: Apollo.MutationHookOptions<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>(ToggleSimStatusDocument, options);
      }
export type ToggleSimStatusMutationHookResult = ReturnType<typeof useToggleSimStatusMutation>;
export type ToggleSimStatusMutationResult = Apollo.MutationResult<ToggleSimStatusMutation>;
export type ToggleSimStatusMutationOptions = Apollo.BaseMutationOptions<ToggleSimStatusMutation, ToggleSimStatusMutationVariables>;
export const GetSimDocument = gql`
    query getSim($data: GetSimInputDto!) {
  getSim(data: $data) {
    ...Sim
  }
}
    ${SimFragmentDoc}`;

/**
 * __useGetSimQuery__
 *
 * To run a query within a React component, call `useGetSimQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimQuery(baseOptions: Apollo.QueryHookOptions<GetSimQuery, GetSimQueryVariables> & ({ variables: GetSimQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
      }
export function useGetSimLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimQuery, GetSimQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
        }
// @ts-ignore
export function useGetSimSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimQuery, GetSimQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimQuery, GetSimQueryVariables>;
export function useGetSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimQuery, GetSimQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimQuery | undefined, GetSimQueryVariables>;
export function useGetSimSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimQuery, GetSimQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimQuery, GetSimQueryVariables>(GetSimDocument, options);
        }
export type GetSimQueryHookResult = ReturnType<typeof useGetSimQuery>;
export type GetSimLazyQueryHookResult = ReturnType<typeof useGetSimLazyQuery>;
export type GetSimSuspenseQueryHookResult = ReturnType<typeof useGetSimSuspenseQuery>;
export type GetSimQueryResult = Apollo.QueryResult<GetSimQuery, GetSimQueryVariables>;
export const GetSimsDocument = gql`
    query GetSims($data: ListSimsInput!) {
  getSims(data: $data) {
    sims {
      ...Sim
    }
  }
}
    ${SimFragmentDoc}`;

/**
 * __useGetSimsQuery__
 *
 * To run a query within a React component, call `useGetSimsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSimsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSimsQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSimsQuery(baseOptions: Apollo.QueryHookOptions<GetSimsQuery, GetSimsQueryVariables> & ({ variables: GetSimsQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
      }
export function useGetSimsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
        }
// @ts-ignore
export function useGetSimsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsQuery, GetSimsQueryVariables>;
export function useGetSimsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>): Apollo.UseSuspenseQueryResult<GetSimsQuery | undefined, GetSimsQueryVariables>;
export function useGetSimsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSimsQuery, GetSimsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSimsQuery, GetSimsQueryVariables>(GetSimsDocument, options);
        }
export type GetSimsQueryHookResult = ReturnType<typeof useGetSimsQuery>;
export type GetSimsLazyQueryHookResult = ReturnType<typeof useGetSimsLazyQuery>;
export type GetSimsSuspenseQueryHookResult = ReturnType<typeof useGetSimsSuspenseQuery>;
export type GetSimsQueryResult = Apollo.QueryResult<GetSimsQuery, GetSimsQueryVariables>;
export const AddSubscriberDocument = gql`
    mutation addSubscriber($data: SubscriberInputDto!) {
  addSubscriber(data: $data) {
    ...Subscriber
  }
}
    ${SubscriberFragmentDoc}`;
export type AddSubscriberMutationFn = Apollo.MutationFunction<AddSubscriberMutation, AddSubscriberMutationVariables>;

/**
 * __useAddSubscriberMutation__
 *
 * To run a mutation, you first call `useAddSubscriberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddSubscriberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addSubscriberMutation, { data, loading, error }] = useAddSubscriberMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddSubscriberMutation(baseOptions?: Apollo.MutationHookOptions<AddSubscriberMutation, AddSubscriberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddSubscriberMutation, AddSubscriberMutationVariables>(AddSubscriberDocument, options);
      }
export type AddSubscriberMutationHookResult = ReturnType<typeof useAddSubscriberMutation>;
export type AddSubscriberMutationResult = Apollo.MutationResult<AddSubscriberMutation>;
export type AddSubscriberMutationOptions = Apollo.BaseMutationOptions<AddSubscriberMutation, AddSubscriberMutationVariables>;
export const GetSubscriberDocument = gql`
    query getSubscriber($subscriberId: String!) {
  getSubscriber(subscriberId: $subscriberId) {
    ...Subscriber
  }
}
    ${SubscriberFragmentDoc}`;

/**
 * __useGetSubscriberQuery__
 *
 * To run a query within a React component, call `useGetSubscriberQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSubscriberQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSubscriberQuery({
 *   variables: {
 *      subscriberId: // value for 'subscriberId'
 *   },
 * });
 */
export function useGetSubscriberQuery(baseOptions: Apollo.QueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables> & ({ variables: GetSubscriberQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscriberQuery, GetSubscriberQueryVariables>(GetSubscriberDocument, options);
      }
export function useGetSubscriberLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscriberQuery, GetSubscriberQueryVariables>(GetSubscriberDocument, options);
        }
// @ts-ignore
export function useGetSubscriberSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables>): Apollo.UseSuspenseQueryResult<GetSubscriberQuery, GetSubscriberQueryVariables>;
export function useGetSubscriberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables>): Apollo.UseSuspenseQueryResult<GetSubscriberQuery | undefined, GetSubscriberQueryVariables>;
export function useGetSubscriberSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSubscriberQuery, GetSubscriberQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSubscriberQuery, GetSubscriberQueryVariables>(GetSubscriberDocument, options);
        }
export type GetSubscriberQueryHookResult = ReturnType<typeof useGetSubscriberQuery>;
export type GetSubscriberLazyQueryHookResult = ReturnType<typeof useGetSubscriberLazyQuery>;
export type GetSubscriberSuspenseQueryHookResult = ReturnType<typeof useGetSubscriberSuspenseQuery>;
export type GetSubscriberQueryResult = Apollo.QueryResult<GetSubscriberQuery, GetSubscriberQueryVariables>;
export const UpdateSubscriberDocument = gql`
    mutation updateSubscriber($subscriberId: String!, $data: UpdateSubscriberInputDto!) {
  updateSubscriber(subscriberId: $subscriberId, data: $data) {
    success
  }
}
    `;
export type UpdateSubscriberMutationFn = Apollo.MutationFunction<UpdateSubscriberMutation, UpdateSubscriberMutationVariables>;

/**
 * __useUpdateSubscriberMutation__
 *
 * To run a mutation, you first call `useUpdateSubscriberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSubscriberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSubscriberMutation, { data, loading, error }] = useUpdateSubscriberMutation({
 *   variables: {
 *      subscriberId: // value for 'subscriberId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateSubscriberMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSubscriberMutation, UpdateSubscriberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateSubscriberMutation, UpdateSubscriberMutationVariables>(UpdateSubscriberDocument, options);
      }
export type UpdateSubscriberMutationHookResult = ReturnType<typeof useUpdateSubscriberMutation>;
export type UpdateSubscriberMutationResult = Apollo.MutationResult<UpdateSubscriberMutation>;
export type UpdateSubscriberMutationOptions = Apollo.BaseMutationOptions<UpdateSubscriberMutation, UpdateSubscriberMutationVariables>;
export const DeleteSubscriberDocument = gql`
    mutation deleteSubscriber($subscriberId: String!) {
  deleteSubscriber(subscriberId: $subscriberId) {
    success
  }
}
    `;
export type DeleteSubscriberMutationFn = Apollo.MutationFunction<DeleteSubscriberMutation, DeleteSubscriberMutationVariables>;

/**
 * __useDeleteSubscriberMutation__
 *
 * To run a mutation, you first call `useDeleteSubscriberMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteSubscriberMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteSubscriberMutation, { data, loading, error }] = useDeleteSubscriberMutation({
 *   variables: {
 *      subscriberId: // value for 'subscriberId'
 *   },
 * });
 */
export function useDeleteSubscriberMutation(baseOptions?: Apollo.MutationHookOptions<DeleteSubscriberMutation, DeleteSubscriberMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteSubscriberMutation, DeleteSubscriberMutationVariables>(DeleteSubscriberDocument, options);
      }
export type DeleteSubscriberMutationHookResult = ReturnType<typeof useDeleteSubscriberMutation>;
export type DeleteSubscriberMutationResult = Apollo.MutationResult<DeleteSubscriberMutation>;
export type DeleteSubscriberMutationOptions = Apollo.BaseMutationOptions<DeleteSubscriberMutation, DeleteSubscriberMutationVariables>;
export const GetSubscribersByNetworkDocument = gql`
    query getSubscribersByNetwork($networkId: String!) {
  getSubscribersByNetwork(networkId: $networkId) {
    subscribers {
      ...Subscriber
    }
  }
}
    ${SubscriberFragmentDoc}`;

/**
 * __useGetSubscribersByNetworkQuery__
 *
 * To run a query within a React component, call `useGetSubscribersByNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSubscribersByNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSubscribersByNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetSubscribersByNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables> & ({ variables: GetSubscribersByNetworkQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>(GetSubscribersByNetworkDocument, options);
      }
export function useGetSubscribersByNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>(GetSubscribersByNetworkDocument, options);
        }
// @ts-ignore
export function useGetSubscribersByNetworkSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>;
export function useGetSubscribersByNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<GetSubscribersByNetworkQuery | undefined, GetSubscribersByNetworkQueryVariables>;
export function useGetSubscribersByNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>(GetSubscribersByNetworkDocument, options);
        }
export type GetSubscribersByNetworkQueryHookResult = ReturnType<typeof useGetSubscribersByNetworkQuery>;
export type GetSubscribersByNetworkLazyQueryHookResult = ReturnType<typeof useGetSubscribersByNetworkLazyQuery>;
export type GetSubscribersByNetworkSuspenseQueryHookResult = ReturnType<typeof useGetSubscribersByNetworkSuspenseQuery>;
export type GetSubscribersByNetworkQueryResult = Apollo.QueryResult<GetSubscribersByNetworkQuery, GetSubscribersByNetworkQueryVariables>;
export const GetSubscriberMetricsByNetworkDocument = gql`
    query getSubscriberMetricsByNetwork($networkId: String!) {
  getSubscriberMetricsByNetwork(networkId: $networkId) {
    total
    active
    inactive
    terminated
  }
}
    `;

/**
 * __useGetSubscriberMetricsByNetworkQuery__
 *
 * To run a query within a React component, call `useGetSubscriberMetricsByNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSubscriberMetricsByNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSubscriberMetricsByNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetSubscriberMetricsByNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables> & ({ variables: GetSubscriberMetricsByNetworkQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>(GetSubscriberMetricsByNetworkDocument, options);
      }
export function useGetSubscriberMetricsByNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>(GetSubscriberMetricsByNetworkDocument, options);
        }
// @ts-ignore
export function useGetSubscriberMetricsByNetworkSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>;
export function useGetSubscriberMetricsByNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<GetSubscriberMetricsByNetworkQuery | undefined, GetSubscriberMetricsByNetworkQueryVariables>;
export function useGetSubscriberMetricsByNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>(GetSubscriberMetricsByNetworkDocument, options);
        }
export type GetSubscriberMetricsByNetworkQueryHookResult = ReturnType<typeof useGetSubscriberMetricsByNetworkQuery>;
export type GetSubscriberMetricsByNetworkLazyQueryHookResult = ReturnType<typeof useGetSubscriberMetricsByNetworkLazyQuery>;
export type GetSubscriberMetricsByNetworkSuspenseQueryHookResult = ReturnType<typeof useGetSubscriberMetricsByNetworkSuspenseQuery>;
export type GetSubscriberMetricsByNetworkQueryResult = Apollo.QueryResult<GetSubscriberMetricsByNetworkQuery, GetSubscriberMetricsByNetworkQueryVariables>;
export const GetGeneratedPdfReportDocument = gql`
    query getGeneratedPdfReport($Id: String!) {
  getGeneratedPdfReport(id: $Id) {
    contentType
    filename
    downloadUrl
  }
}
    `;

/**
 * __useGetGeneratedPdfReportQuery__
 *
 * To run a query within a React component, call `useGetGeneratedPdfReportQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetGeneratedPdfReportQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetGeneratedPdfReportQuery({
 *   variables: {
 *      Id: // value for 'Id'
 *   },
 * });
 */
export function useGetGeneratedPdfReportQuery(baseOptions: Apollo.QueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables> & ({ variables: GetGeneratedPdfReportQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>(GetGeneratedPdfReportDocument, options);
      }
export function useGetGeneratedPdfReportLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>(GetGeneratedPdfReportDocument, options);
        }
// @ts-ignore
export function useGetGeneratedPdfReportSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>;
export function useGetGeneratedPdfReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>): Apollo.UseSuspenseQueryResult<GetGeneratedPdfReportQuery | undefined, GetGeneratedPdfReportQueryVariables>;
export function useGetGeneratedPdfReportSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>(GetGeneratedPdfReportDocument, options);
        }
export type GetGeneratedPdfReportQueryHookResult = ReturnType<typeof useGetGeneratedPdfReportQuery>;
export type GetGeneratedPdfReportLazyQueryHookResult = ReturnType<typeof useGetGeneratedPdfReportLazyQuery>;
export type GetGeneratedPdfReportSuspenseQueryHookResult = ReturnType<typeof useGetGeneratedPdfReportSuspenseQuery>;
export type GetGeneratedPdfReportQueryResult = Apollo.QueryResult<GetGeneratedPdfReportQuery, GetGeneratedPdfReportQueryVariables>;
export const WhoamiDocument = gql`
    query Whoami {
  whoami {
    user {
      ...User
    }
    ownerOf {
      ...Org
    }
    memberOf {
      ...Org
    }
  }
}
    ${UserFragmentDoc}
${OrgFragmentDoc}`;

/**
 * __useWhoamiQuery__
 *
 * To run a query within a React component, call `useWhoamiQuery` and pass it any options that fit your needs.
 * When your component renders, `useWhoamiQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useWhoamiQuery({
 *   variables: {
 *   },
 * });
 */
export function useWhoamiQuery(baseOptions?: Apollo.QueryHookOptions<WhoamiQuery, WhoamiQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<WhoamiQuery, WhoamiQueryVariables>(WhoamiDocument, options);
      }
export function useWhoamiLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<WhoamiQuery, WhoamiQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<WhoamiQuery, WhoamiQueryVariables>(WhoamiDocument, options);
        }
// @ts-ignore
export function useWhoamiSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<WhoamiQuery, WhoamiQueryVariables>): Apollo.UseSuspenseQueryResult<WhoamiQuery, WhoamiQueryVariables>;
export function useWhoamiSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<WhoamiQuery, WhoamiQueryVariables>): Apollo.UseSuspenseQueryResult<WhoamiQuery | undefined, WhoamiQueryVariables>;
export function useWhoamiSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<WhoamiQuery, WhoamiQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<WhoamiQuery, WhoamiQueryVariables>(WhoamiDocument, options);
        }
export type WhoamiQueryHookResult = ReturnType<typeof useWhoamiQuery>;
export type WhoamiLazyQueryHookResult = ReturnType<typeof useWhoamiLazyQuery>;
export type WhoamiSuspenseQueryHookResult = ReturnType<typeof useWhoamiSuspenseQuery>;
export type WhoamiQueryResult = Apollo.QueryResult<WhoamiQuery, WhoamiQueryVariables>;
export const GetUserDocument = gql`
    query GetUser($userId: String!) {
  getUser(userId: $userId) {
    ...User
  }
}
    ${UserFragmentDoc}`;

/**
 * __useGetUserQuery__
 *
 * To run a query within a React component, call `useGetUserQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUserQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUserQuery({
 *   variables: {
 *      userId: // value for 'userId'
 *   },
 * });
 */
export function useGetUserQuery(baseOptions: Apollo.QueryHookOptions<GetUserQuery, GetUserQueryVariables> & ({ variables: GetUserQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
      }
export function useGetUserLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetUserQuery, GetUserQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
        }
// @ts-ignore
export function useGetUserSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetUserQuery, GetUserQueryVariables>): Apollo.UseSuspenseQueryResult<GetUserQuery, GetUserQueryVariables>;
export function useGetUserSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetUserQuery, GetUserQueryVariables>): Apollo.UseSuspenseQueryResult<GetUserQuery | undefined, GetUserQueryVariables>;
export function useGetUserSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetUserQuery, GetUserQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetUserQuery, GetUserQueryVariables>(GetUserDocument, options);
        }
export type GetUserQueryHookResult = ReturnType<typeof useGetUserQuery>;
export type GetUserLazyQueryHookResult = ReturnType<typeof useGetUserLazyQuery>;
export type GetUserSuspenseQueryHookResult = ReturnType<typeof useGetUserSuspenseQuery>;
export type GetUserQueryResult = Apollo.QueryResult<GetUserQuery, GetUserQueryVariables>;
export const GetNetworksDocument = gql`
    query getNetworks {
  getNetworks {
    networks {
      ...UNetwork
    }
  }
}
    ${UNetworkFragmentDoc}`;

/**
 * __useGetNetworksQuery__
 *
 * To run a query within a React component, call `useGetNetworksQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworksQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworksQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetNetworksQuery(baseOptions?: Apollo.QueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworksQuery, GetNetworksQueryVariables>(GetNetworksDocument, options);
      }
export function useGetNetworksLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworksQuery, GetNetworksQueryVariables>(GetNetworksDocument, options);
        }
// @ts-ignore
export function useGetNetworksSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>): Apollo.UseSuspenseQueryResult<GetNetworksQuery, GetNetworksQueryVariables>;
export function useGetNetworksSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>): Apollo.UseSuspenseQueryResult<GetNetworksQuery | undefined, GetNetworksQueryVariables>;
export function useGetNetworksSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworksQuery, GetNetworksQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNetworksQuery, GetNetworksQueryVariables>(GetNetworksDocument, options);
        }
export type GetNetworksQueryHookResult = ReturnType<typeof useGetNetworksQuery>;
export type GetNetworksLazyQueryHookResult = ReturnType<typeof useGetNetworksLazyQuery>;
export type GetNetworksSuspenseQueryHookResult = ReturnType<typeof useGetNetworksSuspenseQuery>;
export type GetNetworksQueryResult = Apollo.QueryResult<GetNetworksQuery, GetNetworksQueryVariables>;
export const GetNetworkDocument = gql`
    query getNetwork($networkId: String!) {
  getNetwork(networkId: $networkId) {
    ...UNetwork
  }
}
    ${UNetworkFragmentDoc}`;

/**
 * __useGetNetworkQuery__
 *
 * To run a query within a React component, call `useGetNetworkQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetNetworkQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetNetworkQuery({
 *   variables: {
 *      networkId: // value for 'networkId'
 *   },
 * });
 */
export function useGetNetworkQuery(baseOptions: Apollo.QueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables> & ({ variables: GetNetworkQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
      }
export function useGetNetworkLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
        }
// @ts-ignore
export function useGetNetworkSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<GetNetworkQuery, GetNetworkQueryVariables>;
export function useGetNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>): Apollo.UseSuspenseQueryResult<GetNetworkQuery | undefined, GetNetworkQueryVariables>;
export function useGetNetworkSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetNetworkQuery, GetNetworkQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetNetworkQuery, GetNetworkQueryVariables>(GetNetworkDocument, options);
        }
export type GetNetworkQueryHookResult = ReturnType<typeof useGetNetworkQuery>;
export type GetNetworkLazyQueryHookResult = ReturnType<typeof useGetNetworkLazyQuery>;
export type GetNetworkSuspenseQueryHookResult = ReturnType<typeof useGetNetworkSuspenseQuery>;
export type GetNetworkQueryResult = Apollo.QueryResult<GetNetworkQuery, GetNetworkQueryVariables>;
export const AddNetworkDocument = gql`
    mutation AddNetwork($data: AddNetworkInputDto!) {
  addNetwork(data: $data) {
    ...UNetwork
  }
}
    ${UNetworkFragmentDoc}`;
export type AddNetworkMutationFn = Apollo.MutationFunction<AddNetworkMutation, AddNetworkMutationVariables>;

/**
 * __useAddNetworkMutation__
 *
 * To run a mutation, you first call `useAddNetworkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddNetworkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addNetworkMutation, { data, loading, error }] = useAddNetworkMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddNetworkMutation(baseOptions?: Apollo.MutationHookOptions<AddNetworkMutation, AddNetworkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddNetworkMutation, AddNetworkMutationVariables>(AddNetworkDocument, options);
      }
export type AddNetworkMutationHookResult = ReturnType<typeof useAddNetworkMutation>;
export type AddNetworkMutationResult = Apollo.MutationResult<AddNetworkMutation>;
export type AddNetworkMutationOptions = Apollo.BaseMutationOptions<AddNetworkMutation, AddNetworkMutationVariables>;
export const SetDefaultNetworkDocument = gql`
    mutation SetDefaultNetwork($data: SetDefaultNetworkInputDto!) {
  setDefaultNetwork(data: $data) {
    success
  }
}
    `;
export type SetDefaultNetworkMutationFn = Apollo.MutationFunction<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>;

/**
 * __useSetDefaultNetworkMutation__
 *
 * To run a mutation, you first call `useSetDefaultNetworkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSetDefaultNetworkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [setDefaultNetworkMutation, { data, loading, error }] = useSetDefaultNetworkMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useSetDefaultNetworkMutation(baseOptions?: Apollo.MutationHookOptions<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>(SetDefaultNetworkDocument, options);
      }
export type SetDefaultNetworkMutationHookResult = ReturnType<typeof useSetDefaultNetworkMutation>;
export type SetDefaultNetworkMutationResult = Apollo.MutationResult<SetDefaultNetworkMutation>;
export type SetDefaultNetworkMutationOptions = Apollo.BaseMutationOptions<SetDefaultNetworkMutation, SetDefaultNetworkMutationVariables>;
export const GetSiteDocument = gql`
    query getSite($siteId: String!) {
  getSite(siteId: $siteId) {
    ...USite
  }
}
    ${USiteFragmentDoc}`;

/**
 * __useGetSiteQuery__
 *
 * To run a query within a React component, call `useGetSiteQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSiteQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSiteQuery({
 *   variables: {
 *      siteId: // value for 'siteId'
 *   },
 * });
 */
export function useGetSiteQuery(baseOptions: Apollo.QueryHookOptions<GetSiteQuery, GetSiteQueryVariables> & ({ variables: GetSiteQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSiteQuery, GetSiteQueryVariables>(GetSiteDocument, options);
      }
export function useGetSiteLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSiteQuery, GetSiteQueryVariables>(GetSiteDocument, options);
        }
// @ts-ignore
export function useGetSiteSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetSiteQuery, GetSiteQueryVariables>;
export function useGetSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>): Apollo.UseSuspenseQueryResult<GetSiteQuery | undefined, GetSiteQueryVariables>;
export function useGetSiteSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSiteQuery, GetSiteQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSiteQuery, GetSiteQueryVariables>(GetSiteDocument, options);
        }
export type GetSiteQueryHookResult = ReturnType<typeof useGetSiteQuery>;
export type GetSiteLazyQueryHookResult = ReturnType<typeof useGetSiteLazyQuery>;
export type GetSiteSuspenseQueryHookResult = ReturnType<typeof useGetSiteSuspenseQuery>;
export type GetSiteQueryResult = Apollo.QueryResult<GetSiteQuery, GetSiteQueryVariables>;
export const AddSiteDocument = gql`
    mutation addSite($data: AddSiteInputDto!) {
  addSite(data: $data) {
    ...USite
  }
}
    ${USiteFragmentDoc}`;
export type AddSiteMutationFn = Apollo.MutationFunction<AddSiteMutation, AddSiteMutationVariables>;

/**
 * __useAddSiteMutation__
 *
 * To run a mutation, you first call `useAddSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useAddSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [addSiteMutation, { data, loading, error }] = useAddSiteMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useAddSiteMutation(baseOptions?: Apollo.MutationHookOptions<AddSiteMutation, AddSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<AddSiteMutation, AddSiteMutationVariables>(AddSiteDocument, options);
      }
export type AddSiteMutationHookResult = ReturnType<typeof useAddSiteMutation>;
export type AddSiteMutationResult = Apollo.MutationResult<AddSiteMutation>;
export type AddSiteMutationOptions = Apollo.BaseMutationOptions<AddSiteMutation, AddSiteMutationVariables>;
export const GetSitesDocument = gql`
    query GetSites($data: SitesInputDto!) {
  getSites(data: $data) {
    sites {
      ...USite
    }
  }
}
    ${USiteFragmentDoc}`;

/**
 * __useGetSitesQuery__
 *
 * To run a query within a React component, call `useGetSitesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetSitesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetSitesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetSitesQuery(baseOptions: Apollo.QueryHookOptions<GetSitesQuery, GetSitesQueryVariables> & ({ variables: GetSitesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
      }
export function useGetSitesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
        }
// @ts-ignore
export function useGetSitesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>): Apollo.UseSuspenseQueryResult<GetSitesQuery, GetSitesQueryVariables>;
export function useGetSitesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>): Apollo.UseSuspenseQueryResult<GetSitesQuery | undefined, GetSitesQueryVariables>;
export function useGetSitesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetSitesQuery, GetSitesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetSitesQuery, GetSitesQueryVariables>(GetSitesDocument, options);
        }
export type GetSitesQueryHookResult = ReturnType<typeof useGetSitesQuery>;
export type GetSitesLazyQueryHookResult = ReturnType<typeof useGetSitesLazyQuery>;
export type GetSitesSuspenseQueryHookResult = ReturnType<typeof useGetSitesSuspenseQuery>;
export type GetSitesQueryResult = Apollo.QueryResult<GetSitesQuery, GetSitesQueryVariables>;
export const UpdateSiteDocument = gql`
    mutation updateSite($siteId: String!, $data: UpdateSiteInputDto!) {
  updateSite(siteId: $siteId, data: $data) {
    name
  }
}
    `;
export type UpdateSiteMutationFn = Apollo.MutationFunction<UpdateSiteMutation, UpdateSiteMutationVariables>;

/**
 * __useUpdateSiteMutation__
 *
 * To run a mutation, you first call `useUpdateSiteMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSiteMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSiteMutation, { data, loading, error }] = useUpdateSiteMutation({
 *   variables: {
 *      siteId: // value for 'siteId'
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateSiteMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSiteMutation, UpdateSiteMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateSiteMutation, UpdateSiteMutationVariables>(UpdateSiteDocument, options);
      }
export type UpdateSiteMutationHookResult = ReturnType<typeof useUpdateSiteMutation>;
export type UpdateSiteMutationResult = Apollo.MutationResult<UpdateSiteMutation>;
export type UpdateSiteMutationOptions = Apollo.BaseMutationOptions<UpdateSiteMutation, UpdateSiteMutationVariables>;
export const GetComponentByIdDocument = gql`
    query getComponentById($componentId: String!) {
  getComponentById(componentId: $componentId) {
    ...UComponent
  }
}
    ${UComponentFragmentDoc}`;

/**
 * __useGetComponentByIdQuery__
 *
 * To run a query within a React component, call `useGetComponentByIdQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetComponentByIdQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetComponentByIdQuery({
 *   variables: {
 *      componentId: // value for 'componentId'
 *   },
 * });
 */
export function useGetComponentByIdQuery(baseOptions: Apollo.QueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables> & ({ variables: GetComponentByIdQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetComponentByIdQuery, GetComponentByIdQueryVariables>(GetComponentByIdDocument, options);
      }
export function useGetComponentByIdLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetComponentByIdQuery, GetComponentByIdQueryVariables>(GetComponentByIdDocument, options);
        }
// @ts-ignore
export function useGetComponentByIdSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetComponentByIdQuery, GetComponentByIdQueryVariables>;
export function useGetComponentByIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetComponentByIdQuery | undefined, GetComponentByIdQueryVariables>;
export function useGetComponentByIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentByIdQuery, GetComponentByIdQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetComponentByIdQuery, GetComponentByIdQueryVariables>(GetComponentByIdDocument, options);
        }
export type GetComponentByIdQueryHookResult = ReturnType<typeof useGetComponentByIdQuery>;
export type GetComponentByIdLazyQueryHookResult = ReturnType<typeof useGetComponentByIdLazyQuery>;
export type GetComponentByIdSuspenseQueryHookResult = ReturnType<typeof useGetComponentByIdSuspenseQuery>;
export type GetComponentByIdQueryResult = Apollo.QueryResult<GetComponentByIdQuery, GetComponentByIdQueryVariables>;
export const GetComponentsByUserIdDocument = gql`
    query GetComponentsByUserId($data: ComponentTypeInputDto!) {
  getComponentsByUserId(data: $data) {
    components {
      ...UComponent
    }
  }
}
    ${UComponentFragmentDoc}`;

/**
 * __useGetComponentsByUserIdQuery__
 *
 * To run a query within a React component, call `useGetComponentsByUserIdQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetComponentsByUserIdQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetComponentsByUserIdQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetComponentsByUserIdQuery(baseOptions: Apollo.QueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables> & ({ variables: GetComponentsByUserIdQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>(GetComponentsByUserIdDocument, options);
      }
export function useGetComponentsByUserIdLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>(GetComponentsByUserIdDocument, options);
        }
// @ts-ignore
export function useGetComponentsByUserIdSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>;
export function useGetComponentsByUserIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>): Apollo.UseSuspenseQueryResult<GetComponentsByUserIdQuery | undefined, GetComponentsByUserIdQueryVariables>;
export function useGetComponentsByUserIdSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>(GetComponentsByUserIdDocument, options);
        }
export type GetComponentsByUserIdQueryHookResult = ReturnType<typeof useGetComponentsByUserIdQuery>;
export type GetComponentsByUserIdLazyQueryHookResult = ReturnType<typeof useGetComponentsByUserIdLazyQuery>;
export type GetComponentsByUserIdSuspenseQueryHookResult = ReturnType<typeof useGetComponentsByUserIdSuspenseQuery>;
export type GetComponentsByUserIdQueryResult = Apollo.QueryResult<GetComponentsByUserIdQuery, GetComponentsByUserIdQueryVariables>;
export const CreateInvitationDocument = gql`
    mutation CreateInvitation($data: CreateInvitationInputDto!) {
  createInvitation(data: $data) {
    ...Invitation
  }
}
    ${InvitationFragmentDoc}`;
export type CreateInvitationMutationFn = Apollo.MutationFunction<CreateInvitationMutation, CreateInvitationMutationVariables>;

/**
 * __useCreateInvitationMutation__
 *
 * To run a mutation, you first call `useCreateInvitationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useCreateInvitationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [createInvitationMutation, { data, loading, error }] = useCreateInvitationMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useCreateInvitationMutation(baseOptions?: Apollo.MutationHookOptions<CreateInvitationMutation, CreateInvitationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<CreateInvitationMutation, CreateInvitationMutationVariables>(CreateInvitationDocument, options);
      }
export type CreateInvitationMutationHookResult = ReturnType<typeof useCreateInvitationMutation>;
export type CreateInvitationMutationResult = Apollo.MutationResult<CreateInvitationMutation>;
export type CreateInvitationMutationOptions = Apollo.BaseMutationOptions<CreateInvitationMutation, CreateInvitationMutationVariables>;
export const GetInvitationsDocument = gql`
    query GetInvitations {
  getInvitations {
    invitations {
      ...Invitation
    }
  }
}
    ${InvitationFragmentDoc}`;

/**
 * __useGetInvitationsQuery__
 *
 * To run a query within a React component, call `useGetInvitationsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetInvitationsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetInvitationsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetInvitationsQuery(baseOptions?: Apollo.QueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetInvitationsQuery, GetInvitationsQueryVariables>(GetInvitationsDocument, options);
      }
export function useGetInvitationsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetInvitationsQuery, GetInvitationsQueryVariables>(GetInvitationsDocument, options);
        }
// @ts-ignore
export function useGetInvitationsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>): Apollo.UseSuspenseQueryResult<GetInvitationsQuery, GetInvitationsQueryVariables>;
export function useGetInvitationsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>): Apollo.UseSuspenseQueryResult<GetInvitationsQuery | undefined, GetInvitationsQueryVariables>;
export function useGetInvitationsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsQuery, GetInvitationsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetInvitationsQuery, GetInvitationsQueryVariables>(GetInvitationsDocument, options);
        }
export type GetInvitationsQueryHookResult = ReturnType<typeof useGetInvitationsQuery>;
export type GetInvitationsLazyQueryHookResult = ReturnType<typeof useGetInvitationsLazyQuery>;
export type GetInvitationsSuspenseQueryHookResult = ReturnType<typeof useGetInvitationsSuspenseQuery>;
export type GetInvitationsQueryResult = Apollo.QueryResult<GetInvitationsQuery, GetInvitationsQueryVariables>;
export const DeleteInvitationDocument = gql`
    mutation DeleteInvitation($deleteInvitationId: String!) {
  deleteInvitation(id: $deleteInvitationId) {
    id
  }
}
    `;
export type DeleteInvitationMutationFn = Apollo.MutationFunction<DeleteInvitationMutation, DeleteInvitationMutationVariables>;

/**
 * __useDeleteInvitationMutation__
 *
 * To run a mutation, you first call `useDeleteInvitationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useDeleteInvitationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [deleteInvitationMutation, { data, loading, error }] = useDeleteInvitationMutation({
 *   variables: {
 *      deleteInvitationId: // value for 'deleteInvitationId'
 *   },
 * });
 */
export function useDeleteInvitationMutation(baseOptions?: Apollo.MutationHookOptions<DeleteInvitationMutation, DeleteInvitationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<DeleteInvitationMutation, DeleteInvitationMutationVariables>(DeleteInvitationDocument, options);
      }
export type DeleteInvitationMutationHookResult = ReturnType<typeof useDeleteInvitationMutation>;
export type DeleteInvitationMutationResult = Apollo.MutationResult<DeleteInvitationMutation>;
export type DeleteInvitationMutationOptions = Apollo.BaseMutationOptions<DeleteInvitationMutation, DeleteInvitationMutationVariables>;
export const UpdateInvitationDocument = gql`
    mutation UpdateInvitation($data: UpateInvitationInputDto!) {
  updateInvitation(data: $data) {
    id
  }
}
    `;
export type UpdateInvitationMutationFn = Apollo.MutationFunction<UpdateInvitationMutation, UpdateInvitationMutationVariables>;

/**
 * __useUpdateInvitationMutation__
 *
 * To run a mutation, you first call `useUpdateInvitationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateInvitationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateInvitationMutation, { data, loading, error }] = useUpdateInvitationMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateInvitationMutation(baseOptions?: Apollo.MutationHookOptions<UpdateInvitationMutation, UpdateInvitationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateInvitationMutation, UpdateInvitationMutationVariables>(UpdateInvitationDocument, options);
      }
export type UpdateInvitationMutationHookResult = ReturnType<typeof useUpdateInvitationMutation>;
export type UpdateInvitationMutationResult = Apollo.MutationResult<UpdateInvitationMutation>;
export type UpdateInvitationMutationOptions = Apollo.BaseMutationOptions<UpdateInvitationMutation, UpdateInvitationMutationVariables>;
export const GetInvitationsByEmailDocument = gql`
    query GetInvitationsByEmail($email: String!) {
  getInvitationsByEmail(email: $email) {
    invitations {
      ...Invitation
    }
  }
}
    ${InvitationFragmentDoc}`;

/**
 * __useGetInvitationsByEmailQuery__
 *
 * To run a query within a React component, call `useGetInvitationsByEmailQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetInvitationsByEmailQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetInvitationsByEmailQuery({
 *   variables: {
 *      email: // value for 'email'
 *   },
 * });
 */
export function useGetInvitationsByEmailQuery(baseOptions: Apollo.QueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables> & ({ variables: GetInvitationsByEmailQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>(GetInvitationsByEmailDocument, options);
      }
export function useGetInvitationsByEmailLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>(GetInvitationsByEmailDocument, options);
        }
// @ts-ignore
export function useGetInvitationsByEmailSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>): Apollo.UseSuspenseQueryResult<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>;
export function useGetInvitationsByEmailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>): Apollo.UseSuspenseQueryResult<GetInvitationsByEmailQuery | undefined, GetInvitationsByEmailQueryVariables>;
export function useGetInvitationsByEmailSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>(GetInvitationsByEmailDocument, options);
        }
export type GetInvitationsByEmailQueryHookResult = ReturnType<typeof useGetInvitationsByEmailQuery>;
export type GetInvitationsByEmailLazyQueryHookResult = ReturnType<typeof useGetInvitationsByEmailLazyQuery>;
export type GetInvitationsByEmailSuspenseQueryHookResult = ReturnType<typeof useGetInvitationsByEmailSuspenseQuery>;
export type GetInvitationsByEmailQueryResult = Apollo.QueryResult<GetInvitationsByEmailQuery, GetInvitationsByEmailQueryVariables>;
export const GetCountriesDocument = gql`
    query GetCountries {
  getCountries {
    countries {
      name
      code
    }
  }
}
    `;

/**
 * __useGetCountriesQuery__
 *
 * To run a query within a React component, call `useGetCountriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetCountriesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCountriesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetCountriesQuery(baseOptions?: Apollo.QueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetCountriesQuery, GetCountriesQueryVariables>(GetCountriesDocument, options);
      }
export function useGetCountriesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetCountriesQuery, GetCountriesQueryVariables>(GetCountriesDocument, options);
        }
// @ts-ignore
export function useGetCountriesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>): Apollo.UseSuspenseQueryResult<GetCountriesQuery, GetCountriesQueryVariables>;
export function useGetCountriesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>): Apollo.UseSuspenseQueryResult<GetCountriesQuery | undefined, GetCountriesQueryVariables>;
export function useGetCountriesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCountriesQuery, GetCountriesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetCountriesQuery, GetCountriesQueryVariables>(GetCountriesDocument, options);
        }
export type GetCountriesQueryHookResult = ReturnType<typeof useGetCountriesQuery>;
export type GetCountriesLazyQueryHookResult = ReturnType<typeof useGetCountriesLazyQuery>;
export type GetCountriesSuspenseQueryHookResult = ReturnType<typeof useGetCountriesSuspenseQuery>;
export type GetCountriesQueryResult = Apollo.QueryResult<GetCountriesQuery, GetCountriesQueryVariables>;
export const GetCurrencySymbolDocument = gql`
    query GetCurrencySymbol($code: String!) {
  getCurrencySymbol(code: $code) {
    code
    symbol
    image
  }
}
    `;

/**
 * __useGetCurrencySymbolQuery__
 *
 * To run a query within a React component, call `useGetCurrencySymbolQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetCurrencySymbolQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetCurrencySymbolQuery({
 *   variables: {
 *      code: // value for 'code'
 *   },
 * });
 */
export function useGetCurrencySymbolQuery(baseOptions: Apollo.QueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables> & ({ variables: GetCurrencySymbolQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>(GetCurrencySymbolDocument, options);
      }
export function useGetCurrencySymbolLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>(GetCurrencySymbolDocument, options);
        }
// @ts-ignore
export function useGetCurrencySymbolSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>): Apollo.UseSuspenseQueryResult<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>;
export function useGetCurrencySymbolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>): Apollo.UseSuspenseQueryResult<GetCurrencySymbolQuery | undefined, GetCurrencySymbolQueryVariables>;
export function useGetCurrencySymbolSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>(GetCurrencySymbolDocument, options);
        }
export type GetCurrencySymbolQueryHookResult = ReturnType<typeof useGetCurrencySymbolQuery>;
export type GetCurrencySymbolLazyQueryHookResult = ReturnType<typeof useGetCurrencySymbolLazyQuery>;
export type GetCurrencySymbolSuspenseQueryHookResult = ReturnType<typeof useGetCurrencySymbolSuspenseQuery>;
export type GetCurrencySymbolQueryResult = Apollo.QueryResult<GetCurrencySymbolQuery, GetCurrencySymbolQueryVariables>;
export const GetTimezonesDocument = gql`
    query GetTimezones {
  getTimezones {
    timezones {
      value
      abbr
      offset
      isdst
      text
      utc
    }
  }
}
    `;

/**
 * __useGetTimezonesQuery__
 *
 * To run a query within a React component, call `useGetTimezonesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetTimezonesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetTimezonesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetTimezonesQuery(baseOptions?: Apollo.QueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetTimezonesQuery, GetTimezonesQueryVariables>(GetTimezonesDocument, options);
      }
export function useGetTimezonesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetTimezonesQuery, GetTimezonesQueryVariables>(GetTimezonesDocument, options);
        }
// @ts-ignore
export function useGetTimezonesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>): Apollo.UseSuspenseQueryResult<GetTimezonesQuery, GetTimezonesQueryVariables>;
export function useGetTimezonesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>): Apollo.UseSuspenseQueryResult<GetTimezonesQuery | undefined, GetTimezonesQueryVariables>;
export function useGetTimezonesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetTimezonesQuery, GetTimezonesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetTimezonesQuery, GetTimezonesQueryVariables>(GetTimezonesDocument, options);
        }
export type GetTimezonesQueryHookResult = ReturnType<typeof useGetTimezonesQuery>;
export type GetTimezonesLazyQueryHookResult = ReturnType<typeof useGetTimezonesLazyQuery>;
export type GetTimezonesSuspenseQueryHookResult = ReturnType<typeof useGetTimezonesSuspenseQuery>;
export type GetTimezonesQueryResult = Apollo.QueryResult<GetTimezonesQuery, GetTimezonesQueryVariables>;
export const UpdateNotificationDocument = gql`
    mutation UpdateNotification($isRead: Boolean!, $updateNotificationId: String!) {
  updateNotification(isRead: $isRead, id: $updateNotificationId) {
    id
  }
}
    `;
export type UpdateNotificationMutationFn = Apollo.MutationFunction<UpdateNotificationMutation, UpdateNotificationMutationVariables>;

/**
 * __useUpdateNotificationMutation__
 *
 * To run a mutation, you first call `useUpdateNotificationMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateNotificationMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateNotificationMutation, { data, loading, error }] = useUpdateNotificationMutation({
 *   variables: {
 *      isRead: // value for 'isRead'
 *      updateNotificationId: // value for 'updateNotificationId'
 *   },
 * });
 */
export function useUpdateNotificationMutation(baseOptions?: Apollo.MutationHookOptions<UpdateNotificationMutation, UpdateNotificationMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateNotificationMutation, UpdateNotificationMutationVariables>(UpdateNotificationDocument, options);
      }
export type UpdateNotificationMutationHookResult = ReturnType<typeof useUpdateNotificationMutation>;
export type UpdateNotificationMutationResult = Apollo.MutationResult<UpdateNotificationMutation>;
export type UpdateNotificationMutationOptions = Apollo.BaseMutationOptions<UpdateNotificationMutation, UpdateNotificationMutationVariables>;
export const GetDataUsagesDocument = gql`
    query GetDataUsages($data: SimUsagesInputDto!) {
  getDataUsages(data: $data) {
    usages {
      usage
      simId
    }
  }
}
    `;

/**
 * __useGetDataUsagesQuery__
 *
 * To run a query within a React component, call `useGetDataUsagesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetDataUsagesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetDataUsagesQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useGetDataUsagesQuery(baseOptions: Apollo.QueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables> & ({ variables: GetDataUsagesQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetDataUsagesQuery, GetDataUsagesQueryVariables>(GetDataUsagesDocument, options);
      }
export function useGetDataUsagesLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetDataUsagesQuery, GetDataUsagesQueryVariables>(GetDataUsagesDocument, options);
        }
// @ts-ignore
export function useGetDataUsagesSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>): Apollo.UseSuspenseQueryResult<GetDataUsagesQuery, GetDataUsagesQueryVariables>;
export function useGetDataUsagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>): Apollo.UseSuspenseQueryResult<GetDataUsagesQuery | undefined, GetDataUsagesQueryVariables>;
export function useGetDataUsagesSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetDataUsagesQuery, GetDataUsagesQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetDataUsagesQuery, GetDataUsagesQueryVariables>(GetDataUsagesDocument, options);
        }
export type GetDataUsagesQueryHookResult = ReturnType<typeof useGetDataUsagesQuery>;
export type GetDataUsagesLazyQueryHookResult = ReturnType<typeof useGetDataUsagesLazyQuery>;
export type GetDataUsagesSuspenseQueryHookResult = ReturnType<typeof useGetDataUsagesSuspenseQuery>;
export type GetDataUsagesQueryResult = Apollo.QueryResult<GetDataUsagesQuery, GetDataUsagesQueryVariables>;
export const GetAppsDocument = gql`
    query GetApps {
  getApps {
    apps {
      name
      space
      notes
      metricsKeys
    }
  }
}
    `;

/**
 * __useGetAppsQuery__
 *
 * To run a query within a React component, call `useGetAppsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetAppsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetAppsQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetAppsQuery(baseOptions?: Apollo.QueryHookOptions<GetAppsQuery, GetAppsQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetAppsQuery, GetAppsQueryVariables>(GetAppsDocument, options);
      }
export function useGetAppsLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetAppsQuery, GetAppsQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetAppsQuery, GetAppsQueryVariables>(GetAppsDocument, options);
        }
// @ts-ignore
export function useGetAppsSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetAppsQuery, GetAppsQueryVariables>): Apollo.UseSuspenseQueryResult<GetAppsQuery, GetAppsQueryVariables>;
export function useGetAppsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetAppsQuery, GetAppsQueryVariables>): Apollo.UseSuspenseQueryResult<GetAppsQuery | undefined, GetAppsQueryVariables>;
export function useGetAppsSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetAppsQuery, GetAppsQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetAppsQuery, GetAppsQueryVariables>(GetAppsDocument, options);
        }
export type GetAppsQueryHookResult = ReturnType<typeof useGetAppsQuery>;
export type GetAppsLazyQueryHookResult = ReturnType<typeof useGetAppsLazyQuery>;
export type GetAppsSuspenseQueryHookResult = ReturnType<typeof useGetAppsSuspenseQuery>;
export type GetAppsQueryResult = Apollo.QueryResult<GetAppsQuery, GetAppsQueryVariables>;
export const SoftwareDocument = gql`
    query Software($data: GetSoftwaresInput!) {
  getSoftwares(data: $data) {
    software {
      id
      releaseDate
      nodeId
      status
      changeLog
      currentVersion
      desiredVersion
      name
      space
      notes
      metricsKeys
      createdAt
      updatedAt
    }
  }
}
    `;

/**
 * __useSoftwareQuery__
 *
 * To run a query within a React component, call `useSoftwareQuery` and pass it any options that fit your needs.
 * When your component renders, `useSoftwareQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSoftwareQuery({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useSoftwareQuery(baseOptions: Apollo.QueryHookOptions<SoftwareQuery, SoftwareQueryVariables> & ({ variables: SoftwareQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<SoftwareQuery, SoftwareQueryVariables>(SoftwareDocument, options);
      }
export function useSoftwareLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<SoftwareQuery, SoftwareQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<SoftwareQuery, SoftwareQueryVariables>(SoftwareDocument, options);
        }
// @ts-ignore
export function useSoftwareSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<SoftwareQuery, SoftwareQueryVariables>): Apollo.UseSuspenseQueryResult<SoftwareQuery, SoftwareQueryVariables>;
export function useSoftwareSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SoftwareQuery, SoftwareQueryVariables>): Apollo.UseSuspenseQueryResult<SoftwareQuery | undefined, SoftwareQueryVariables>;
export function useSoftwareSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<SoftwareQuery, SoftwareQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<SoftwareQuery, SoftwareQueryVariables>(SoftwareDocument, options);
        }
export type SoftwareQueryHookResult = ReturnType<typeof useSoftwareQuery>;
export type SoftwareLazyQueryHookResult = ReturnType<typeof useSoftwareLazyQuery>;
export type SoftwareSuspenseQueryHookResult = ReturnType<typeof useSoftwareSuspenseQuery>;
export type SoftwareQueryResult = Apollo.QueryResult<SoftwareQuery, SoftwareQueryVariables>;
export const UpdateSoftwareDocument = gql`
    mutation UpdateSoftware($data: UpdateSoftwareInputDto!) {
  updateSoftware(data: $data) {
    message
  }
}
    `;
export type UpdateSoftwareMutationFn = Apollo.MutationFunction<UpdateSoftwareMutation, UpdateSoftwareMutationVariables>;

/**
 * __useUpdateSoftwareMutation__
 *
 * To run a mutation, you first call `useUpdateSoftwareMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useUpdateSoftwareMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [updateSoftwareMutation, { data, loading, error }] = useUpdateSoftwareMutation({
 *   variables: {
 *      data: // value for 'data'
 *   },
 * });
 */
export function useUpdateSoftwareMutation(baseOptions?: Apollo.MutationHookOptions<UpdateSoftwareMutation, UpdateSoftwareMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<UpdateSoftwareMutation, UpdateSoftwareMutationVariables>(UpdateSoftwareDocument, options);
      }
export type UpdateSoftwareMutationHookResult = ReturnType<typeof useUpdateSoftwareMutation>;
export type UpdateSoftwareMutationResult = Apollo.MutationResult<UpdateSoftwareMutation>;
export type UpdateSoftwareMutationOptions = Apollo.BaseMutationOptions<UpdateSoftwareMutation, UpdateSoftwareMutationVariables>;