/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import EmptyView from '@/ui/molecules/EmptyView';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import PlanCard from '@/ui/molecules/PlanCard';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';
import { Grid, Paper } from '@mui/material';

interface IDataPlan {
  data: any;
  handleActionButon: Function;
  handleOptionMenuItemAction: Function;
}

const DataPlan = ({
  data,
  handleActionButon,
  handleOptionMenuItemAction,
}: IDataPlan) => {
  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        borderRadius: '10px',
        height: 'calc(100vh - 200px)',
      }}
    >
      <PageContainerHeader
        showSearch={false}
        title={'Data plans'}
        buttonTitle={'CREATE DATA PLAN'}
        handleButtonAction={handleActionButon}
        // warningIcon={true}
      />
      <br />
      {data.length === 0 ? (
        <EmptyView icon={UpdateIcon} title="No data plan created yet!" />
      ) : (
        <Grid container rowSpacing={2} columnSpacing={2}>
          {data.map(
            ({
              uuid,
              name,
              duration,
              users,
              currency,
              dataVolume,
              dataUnit,
              amount,
            }: any) => (
              <Grid item xs={12} sm={6} md={4} key={uuid}>
                <PlanCard
                  uuid={uuid}
                  name={name}
                  users={users}
                  amount={amount}
                  dataUnit={dataUnit}
                  duration={duration}
                  currency={currency}
                  dataVolume={dataVolume}
                  handleOptionMenuItemAction={handleOptionMenuItemAction}
                />
              </Grid>
            ),
          )}
        </Grid>
      )}
    </Paper>
  );
};

export default DataPlan;
