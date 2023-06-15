require('dotenv').config()
const express = require('express')
const { ApolloServer } = require('apollo-server-express')
const {
 ApolloGateway,
 RemoteGraphQLDataSource,
 IntrospectAndCompose,
} = require('@apollo/gateway')
const { ApolloServerPluginInlineTrace } = require('apollo-server-core')

class CustomDataSource extends RemoteGraphQLDataSource {
 willSendRequest({ request, context }) {
  request.http.headers.set('Authorization', context.authToken)
 }
}

const gateway = new ApolloGateway({
 supergraphSdl: new IntrospectAndCompose({
  subgraphs: [{ name: 'subgraph', url: process.env.SUBGRAPH_PLANNING_TOOL }],
  introspectionHeaders: {
   Authorization: 'Bearer abc123',
  },
 }),
 buildService({ url }) {
  return new CustomDataSource({ url })
 },
})

const app = express()

async function startApolloServer() {
 const server = new ApolloServer({
  gateway,
  plugins: [ApolloServerPluginInlineTrace()],
  context: ({ req }) => {
   const authToken = req.headers.authorization || ''

   return { authToken }
  },
 })
 await server.start()

 server.applyMiddleware({ app })

 return server
}

startApolloServer().then((server) => {
 const port = 4001
 app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}${server.graphqlPath}`)
 })
})
