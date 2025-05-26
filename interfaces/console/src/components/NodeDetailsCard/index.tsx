/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NODE_IMAGES } from '@/constants';
import { colors } from '@/theme';
import { Chip, Link, Paper, Stack, Typography } from '@mui/material';
import Grid from '@mui/material/Grid2';
import DeviceModalView from '../DeviceModalView';
import LoadingWrapper from '../LoadingWrapper';

interface INodeDetailsCard {
  loading: boolean;
  nodeTitle: string;
  nodeType?: any;
  isUpdateAvailable: boolean;
  handleUpdateNode: () => void;
  getNodeUpdateInfos: () => void;
}

const NodeDetailsCard = ({
  loading,
  nodeTitle,
  isUpdateAvailable,
  getNodeUpdateInfos,
  nodeType = 'HOME',
}: INodeDetailsCard) => {
  return (
    <LoadingWrapper
      width="100%"
      radius={'small'}
      isLoading={loading}
      height="fit-content"
    >
      <Paper sx={{ p: 2, gap: 1 }}>
        <Grid container>
          <Grid size={{ xs: 5 }}>
            <Typography variant="h6">{nodeTitle}</Typography>
          </Grid>
          {isUpdateAvailable && (
            <Grid container size={{ xs: 7 }} justifyContent="flex-end">
              <Chip
                variant="outlined"
                sx={{
                  color: colors.primaryMain,
                  border: `1px solid ${colors.primaryMain}`,
                }}
                label={
                  <Stack spacing={'4px'} direction="row" alignItems="center">
                    <Typography variant="body2">
                      Software update available — view
                    </Typography>
                    <Link
                      onClick={() => getNodeUpdateInfos()}
                      sx={{
                        cursor: 'pointer',
                        typography: 'body2',
                        color: colors.primaryDark,
                      }}
                    >
                      notes
                    </Link>
                  </Stack>
                }
              />
            </Grid>
          )}
          <Grid size={{ xs: 12 }} height={'fit-content'} my={4}>
            {NODE_IMAGES[nodeType as 'hnode' | 'anode' | 'tnode'] && (
              <DeviceModalView
                nodeType={nodeType}
                image={NODE_IMAGES[nodeType as 'hnode' | 'anode' | 'tnode']}
              />
            )}
          </Grid>
        </Grid>
      </Paper>
    </LoadingWrapper>
  );
};

export default NodeDetailsCard;
