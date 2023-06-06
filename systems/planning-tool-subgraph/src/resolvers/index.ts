import { Mutation } from './Mutation'
import { Query } from './Query'

const resolvers = {
    ...Query,
    ...Mutation,
}

export default resolvers
