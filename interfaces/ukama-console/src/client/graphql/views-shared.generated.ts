import * as Types from './types';

import { gql } from '@apollo/client';
export type SectionErrorFieldsFragment = { __typename?: 'SectionError', section: string, code: Types.SectionErrorCode, message: string };

export type ViewNodeFragment = { __typename?: 'Node', id: string, name: string, type: Types.NodeTypeEnum, site: { __typename?: 'NodeSite', siteId?: string | null, networkId?: string | null }, status: { __typename?: 'NodeStatus', connectivity: string, state: string } };

export type ViewSiteFragment = { __typename?: 'SiteDto', id: string, name: string, networkId: string, latitude: string, longitude: string, location: string, isDeactivated: boolean, installDate: string, createdAt: string };

export const SectionErrorFieldsFragmentDoc = gql`
    fragment SectionErrorFields on SectionError {
  section
  code
  message
}
    `;
export const ViewNodeFragmentDoc = gql`
    fragment ViewNode on Node {
  id
  name
  type
  site {
    siteId
    networkId
  }
  status {
    connectivity
    state
  }
}
    `;
export const ViewSiteFragmentDoc = gql`
    fragment ViewSite on SiteDto {
  id
  name
  networkId
  latitude
  longitude
  location
  isDeactivated
  installDate
  createdAt
}
    `;