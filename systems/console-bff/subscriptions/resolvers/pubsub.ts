import { createPubSub } from "@graphql-yoga/subscription";

export const enum Topic {
  NOTIFICATIONS = "NOTIFICATIONS",
  DYNAMIC_ID_TOPIC = "DYNAMIC_ID_TOPIC",
}

export const pubSub = createPubSub();
