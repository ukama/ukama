import { Draft, Resolvers } from '../__generated__/resolvers-types'

export const Query: Resolvers = {
    Query: {
        async getDraft(_parent, { id }, _context): Promise<Draft> {
            return {
                id: id.toString(),
                name: 'Draft 1',
                lastSaved: 1686084992,
                site: {
                    name: 'My Site',
                    height: 0,
                    apOption: 'one-to-one',
                    solarUptime: 90,
                    isSetlite: false,
                    location: {
                        lat: '0.123',
                        lng: '0.456',
                        address: '123 Main Street',
                    },
                },
                events: [
                    {
                        id: 'event-id',
                        operation: 'attributeName',
                        value: 'attributeValue',
                    },
                ],
            }
        },
    },
}
