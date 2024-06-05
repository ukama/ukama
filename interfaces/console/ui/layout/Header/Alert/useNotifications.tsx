import { metricsClient } from "@/client/ApolloClient";
import { NotificationRes, useGetNotificationsQuery, useGetNotificationsSubSubscription } from "@/generated/metrics";
import { useState } from "react";
// custom hook to use graphql queries

const useNotifications = () =>{
 const [alerts, setAlerts] = useState<NotificationRes[] | undefined>(undefined);

 // Fetch initial notifications
 useGetNotificationsQuery({
   client: metricsClient,
   fetchPolicy: 'cache-and-network',
   variables: {
     data: {
       orgId: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b221',
       userId: 'da421ed5-0fba-4638-9661-a9204f49006a',
       networkId: 'da421ed5-0fba-4638-9661-a9204f490069',
       scopes: ['notifications'],
       siteId: 'da421ed5-0fba-4638-9661-a9204f490062',
       subscriberId: 'da421ed5-0fba-4638-9661-a9204f490065',
     },
   },
   onCompleted: (data) => {
     setAlerts(data.getNotifications.notifications);
   },
 });

 // Subscribe to notifications
 useGetNotificationsSubSubscription({
   client: metricsClient,
   variables: {
     orgId: 'fbc2a80d-339e-4d3c-acaa-329dd3d1b221',
     userId: 'da421ed5-0fba-4638-9661-a9204f49006a',
     networkId: 'da421ed5-0fba-4638-9661-a9204f490069',
     scopes: ['notifications'],
     siteId: 'da421ed5-0fba-4638-9661-a9204f490062',
     subscriberId: 'da421ed5-0fba-4638-9661-a9204f490065',
   },
   onData: ({ data: subscriptionData }) => {
     const newAlerts = subscriptionData.data?.getNotificationsSub;
     if (newAlerts) {
       setAlerts((prev) => (prev ? [newAlerts, ...prev] : [newAlerts]));
     }
   },
 });
 return {alerts, setAlerts}
}
export default useNotifications