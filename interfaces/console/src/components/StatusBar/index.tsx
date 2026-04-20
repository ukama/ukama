/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, SiteDto as Site } from '@/client/graphql/generated';
import { Button } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { styled } from '@mui/material/styles';
import LoadingWrapper from '../LoadingWrapper';
import NodeDropDown from '../NodeDropDown';
import SplitButton from '../SplitButton';

const StyledBtn = styled(Button)({
  whiteSpace: 'nowrap',
  minWidth: 'max-content',
});

interface INodeStatus {
  objs: any;
  uptime: number;
  loading: boolean;
  ActionOptions: any[];
  handleSelected: (obj: Node | Site) => void;
  handleEditClick?: (obj: Node | Site) => void;
  selected: Node | Site | undefined;
  handleActionClick: (id: string) => void;
}

const StatusBar = ({
  objs,
  uptime,
  selected,
  loading = false,
  ActionOptions,
  handleSelected,
  handleEditClick,
  handleActionClick,
}: INodeStatus) => {
  const handleUpdateNode = () =>
    handleEditClick?.(objs.find((item: any) => item.id === selected) as Node | Site);

  return (
    <Grid container alignItems={'center'}>
      <Grid size={{ xs: 12, md: 9 }}>
        <NodeDropDown
          objs={objs}
          uptime={uptime}
          loading={loading}
          selected={selected}
          isReady={uptime > 0}
          onSelected={handleSelected}
        />
      </Grid>
      <Grid
        container
        columnSpacing={2}
        size={{ xs: 12, md: 3 }}
        justifyContent="flex-end"
        visibility={uptime > 0 ? 'visible' : 'hidden'}
      >
        {handleEditClick && 
          <Grid>
            <LoadingWrapper isLoading={loading} height={40}>
              <StyledBtn variant="contained" onClick={handleUpdateNode}>
                Edit NODE
              </StyledBtn>
            </LoadingWrapper>
          </Grid>
        }

        <Grid>
          <LoadingWrapper isLoading={loading} height={40} width={'100%'}>
            <SplitButton
              options={ActionOptions}
              handleSplitActionClick={handleActionClick}
            />
          </LoadingWrapper>
        </Grid>
      </Grid>
    </Grid>
  );
};

export default StatusBar;
