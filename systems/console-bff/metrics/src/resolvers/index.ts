import WebSocket = require('ws')
import { Resolvers } from '../../types'

const metricResolvers: Resolvers = {
 Query: {
  getMetrics: async (_, { input }, { pubSub }) => {
   const { type, orgId, userId } = input
   const ws = new WebSocket('ws://localhost:8080')

   ws.on('error', (e) => console.log('Error: ', e))

   ws.on('open' as any, function open() {
    ws.send('something')
   })

   ws.on('message' as any, function message(data) {
    console.log('received: %s', data)
    pubSub.publish(`metric-${input.type}`, `${orgId}/${userId}`, {
     value: data.toString(),
    })
   })

   return [
    {
     value: '1',
    },
   ]
  },
 },
 Subscription: {
  getMetricEvent: {
   subscribe: (_, { input }, { pubSub }) =>
    pubSub.subscribe(`metric-${input.type}`, `${input.orgId}/${input.userId}`),
   resolve: (payload) => {
    console.log(payload)
    return payload
   },
  },
 },
}

export default metricResolvers
