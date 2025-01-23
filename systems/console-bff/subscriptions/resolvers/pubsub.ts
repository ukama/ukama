import { createPubSub } from "@graphql-yoga/subscription";

export const enum Topic {
  NOTIFICATIONS = "SUBSCRIPTION",
  DYNAMIC_ID_TOPIC = "SUBSCRIPTION_TOPIC",
}

export const pubSub = createPubSub();
