import { NonEmptyArray } from 'type-graphql'
import { DraftResolver } from './draft/draftResolver'

const resolvers: NonEmptyArray<Function> = [DraftResolver]

export default resolvers
