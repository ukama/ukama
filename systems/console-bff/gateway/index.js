import { ApolloGateway } from '@apollo/gateway'
import { ApolloServer } from '@apollo/server'
import { expressMiddleware } from '@apollo/server/express4'
import { ApolloServerPluginDrainHttpServer } from '@apollo/server/plugin/drainHttpServer'
import { ApolloServerPluginInlineTrace } from 'apollo-server-core'
import { json } from 'body-parser'
import cors from 'cors'
import dotenv from 'dotenv'
import http from 'http'
import { logger } from '../common/logger'
import { configureExpress } from './configureExpress'
dotenv.config({ path: '.env' })

const app = configureExpress(logger)
const subgraphUrls = [process.env.SUBGRAPH_PLANNING_TOOL]
const httpServer = http.createServer(app)
const gateway = new ApolloGateway({
 //  supergraphSdl: new IntrospectAndCompose({
 //   subgraphs: [{ name: 'subgraph', url: '' }],
 // }),
 serviceList: subgraphUrls.map((url) => ({ name: 'subgraph', url })),
})

const startServer = async () => {
 const server = new ApolloServer({
  gateway,
  plugins: [
   ApolloServerPluginInlineTrace(),
   ApolloServerPluginDrainHttpServer({ httpServer }),
  ],
 })

 await server.start()

 app.use(
  '/graphql',

  cors({
   origin: [
    process.env.AUTH_APP_URL,
    process.env.PLAYGROUND_URL,
    process.env.CONSOLE_APP_URL,
   ],
   credentials: true,
  }),

  json(),
  expressMiddleware(server)
 )
 await new Promise((resolve) =>
  httpServer.listen({ port: process.env.PORT }, resolve)
 )
 logger.info(`🚀 Server ready at http://localhost:${process.env.PORT}/graphql`)
}

startServer()
