/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, SiteDto as Site } from '@/client/graphql/generated';
import { ToggleOption, ToggleValue } from '@/components/ui/ToggleButtonsMenu';
import { SplitButtonOption } from '@/components/ui/SplitButton';
import { Button } from '@mui/material';
import Grid from '@mui/material/Grid';
import { styled } from '@mui/material/styles';
import LoadingWrapper from '@/components/ui/LoadingWrapper';
import NodeDropDown from '@/app/(main)/console/nodes/_components/NodeDropDown';
import SplitButton from '@/components/ui/SplitButton';
import ToggleButtonsMenu from '@/components/ui/ToggleButtonsMenu';

const StyledBtn = styled(Button)({
  whiteSpace: 'nowrap',
  minWidth: 'max-content',
});

interface INodeStatus {
  objs: Node[] | Site[] | [];
  uptime: number;
  loading?: boolean;
  actionOptions: ToggleOption[] | SplitButtonOption[];
  actionLoading?: boolean;
  type: 'split' | 'toggle';
  actionOptionValues?: ToggleValue[];
  handleSelected: (obj: Node | Site) => void;
  handleEditClick?: (obj: Node | Site) => void;
  selected: Node | Site | undefined;
  handleActionClick: (id: string, value: boolean) => void;
}

const StatusBar = ({
  objs,
  type,
  uptime,
  selected,
  loading = false,
  actionLoading = false,
  actionOptions,
  handleSelected,
  handleEditClick,
  handleActionClick,
  actionOptionValues,
}: INodeStatus) => {
  const handleUpdateNode = () =>
    handleEditClick?.(
      (objs as (Node | Site)[]).find((item) => item.id === selected?.id) as Node | Site,
    );

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
      >
        {handleEditClick && (
          <Grid>
            <LoadingWrapper isLoading={loading} height={40}>
              <StyledBtn variant="contained" onClick={handleUpdateNode}>
                Edit NODE
              </StyledBtn>
            </LoadingWrapper>
          </Grid>
        )}

        <Grid>
          {type === 'toggle' ? (
            <ToggleButtonsMenu
              title="Actions"
              options={actionOptions}
              isLoading={actionLoading}
              values={actionOptionValues ?? []}
              handleToggle={handleActionClick}
            />
          ) : (
            <SplitButton
              options={actionOptions}
              handleSplitActionClick={(id: string) => {
                handleActionClick(id, true);
              }}
            />
          )}
        </Grid>
      </Grid>
    </Grid>
  );
};

export default StatusBar;
