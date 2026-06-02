/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node } from '@/client/graphql/generated';
import { Grid, Link, Typography } from '@mui/material';
import LoadingWrapper from '@/components/ui/LoadingWrapper';
interface INodeGroup {
  nodes: Node[];
  loading: boolean;
  handleNodeAction: (id: string) => void;
}

const NodeGroup = ({ nodes, loading, handleNodeAction }: INodeGroup) => {
  return (
    <Grid container spacing={2} alignItems="center">
      <Grid size={5}>
        <Typography fontWeight={500} variant="body2">
          Node Group
        </Typography>
      </Grid>
      <Grid size={7}>
        <LoadingWrapper isLoading={loading} height={24} radius="small">
          {nodes.length > 0 ? (
            nodes.map((item) => (
              <Link
                variant="body2"
                fontWeight={500}
                key={item.id}
                underline="always"
                sx={{ textTransform: 'capitalize' }}
                onClick={() => handleNodeAction(item.id)}
              >
                {item.name}
              </Link>
            ))
          ) : (
            <Typography fontWeight={500} variant="body2">
              N/A
            </Typography>
          )}
        </LoadingWrapper>
      </Grid>
    </Grid>
  );
};

export default NodeGroup;
