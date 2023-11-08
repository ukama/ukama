/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { RoundedCard, SkeletonRoundedCard } from '@/styles/global';
import ApexLineChart from '../ApexLineChart';

type StatsCardProps = {
  loading: boolean;
  metricData: any;
};

const StatsCard = ({ loading, metricData }: StatsCardProps) => {
  return (
    <>
      {loading ? (
        <SkeletonRoundedCard variant="rectangular" height={337} />
      ) : (
        <RoundedCard sx={{ minHeight: 337, display: 'flex' }}>
          <ApexLineChart data={metricData['uptimetrx']} />
        </RoundedCard>
      )}
    </>
  );
};
export default StatsCard;
