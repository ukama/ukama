import { ApolloServer, ContextFunction } from '@apollo/server'
import {
    StandaloneServerContextFunctionArgument,
    startStandaloneServer,
} from '@apollo/server/standalone'
import { buildSubgraphSchema } from '@apollo/subgraph'
import { readFileSync } from 'fs'
import { GraphQLError } from 'graphql'
import gql from 'graphql-tag'
import resolvers from './resolvers'
import { DataSourceContext } from './types/DataSourceContext'

const port = process.env.PORT ?? '4001'
const routerSecret = process.env.ROUTER_SECRET

const context: ContextFunction<
    [StandaloneServerContextFunctionArgument],
    DataSourceContext
> = async ({ req }) => {
    if (routerSecret && req.headers['router-authorization'] !== routerSecret) {
        throw new GraphQLError('Missing router authentication', {
            extensions: {
                code: 'UNAUTHENTICATED',
                http: { status: 401 },
            },
        })
    }

    return {
        auth: req.headers.authorization,
    }
}

async function main() {
    let typeDefs = gql(
        readFileSync('schema.graphql', {
            encoding: 'utf-8',
        })
    )
    const server = new ApolloServer({
        schema: buildSubgraphSchema({ typeDefs, resolvers }),
    })
    const { url } = await startStandaloneServer(server, {
        context,
        listen: { port: Number.parseInt(port) },
    })

    console.log(`Ukama Planning Tool Subgraph 🚀 ready at ${url}`)
}

main()
