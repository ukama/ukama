import { createPubSub, createSchema, createYoga } from '@graphql-yoga/node'
import * as fs from 'fs'
import { createServer } from 'node:http'
import * as path from 'path'
import resolvers from './resolvers'

const typeDefs = fs.readFileSync(path.join(process.cwd(), 'schema.graphql'), {
 encoding: 'utf-8',
})

const pubSub = createPubSub<{
 newMessage: [payload: { from: string; body: string }]
}>()

const yoga = createYoga({
 schema: createSchema({
  typeDefs,
  resolvers,
 }),
 context: {
  pubSub,
 },
})

const server = createServer(yoga)

server.listen(4000, () => {
 console.info('Server is running on http://localhost:4000/graphql')
})
