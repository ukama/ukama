/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { HorizontalContainer } from '@/styles/global';
import { TVariant } from '@/types';
import { formatKPIValue } from '@/utils';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import { IconButton, Tooltip, Typography } from '@mui/material';
import Grid from '@mui/material/Grid2';
import Zoom from '@mui/material/Zoom';
import { useEffect, useState } from 'react';

interface IThreshold {
  min: number;
  normal: number;
  max: number;
}

interface INodeStatItem {
  id?: string;
  unit?: string;
  name: string;
  value: string;
  format?: string;
  variant?: TVariant;
  threshold?: IThreshold | null;
  nameInfo?: string; //Tooltip info about stat
  valueInfo?: string; //Tooltip info about stat value
  showAlertInfo?: boolean; //Pass true if its an alert value
}

interface ITextWithToolTip {
  title: string;
  isAlert?: boolean;
  tooltipText: string;
  showToottip?: boolean;
}

const TextWithToolTip = ({
  title,
  tooltipText,
  isAlert = false,
  showToottip = false,
}: ITextWithToolTip) => (
  <HorizontalContainer>
    <Typography
      fontWeight={500}
      variant="body1"
      color={isAlert ? 'error' : 'default'}
    >
      {title}
    </Typography>
    {showToottip && (
      <Tooltip
        arrow
        placement="right"
        title={tooltipText}
        TransitionComponent={Zoom}
      >
        <IconButton
          color={isAlert ? 'error' : 'default'}
          sx={{
            '&:hover svg path': {
              fill: 'inherit',
            },
          }}
        >
          <InfoOutlinedIcon
            sx={{
              width: '16px',
              height: '16px',
            }}
          />
        </IconButton>
      </Tooltip>
    )}
  </HorizontalContainer>
);

const variants = (variant: TVariant, key: string) => {
  switch (key) {
    case 'NG':
      return variant === 'small' ? 4 : variant === 'large' ? 8 : 5;
    case 'VG':
      return variant === 'small' ? 8 : variant === 'large' ? 4 : 7;
  }
};

const NodeStatItem = ({
  name,
  value,
  id = '',
  unit = '',
  threshold,
  nameInfo = '',
  valueInfo = '',
  format = undefined,
  variant = 'medium',
}: INodeStatItem) => {
  const [v, setV] = useState<string>('');

  useEffect(() => {
    setV(format ? formatKPIValue(value, format) : value);
  }, [value, format]);

  useEffect(() => {
    if (id) {
      const token = PubSub.subscribe(`stat-${id}`, (_, data) => {
        if (data.length > 0) {
          setV(format ? formatKPIValue(data[1], format) : data[1]);
        }
      });
      return () => {
        PubSub.unsubscribe(token);
      };
    }
  }, [id, unit]);

  const isAlert = (): boolean => {
    if (id && threshold && parseFloat(v) > threshold.normal) {
      return true;
    }
    return false;
  };

  return (
    <Grid container spacing={2}>
      <Grid size={{ xs: variants(variant, 'NG') }}>
        <TextWithToolTip
          title={name}
          showToottip={!!nameInfo}
          tooltipText={nameInfo}
        />
      </Grid>
      <Grid size={{ xs: variants(variant, 'VG') }}>
        <TextWithToolTip
          title={`${v} ${unit}`}
          isAlert={isAlert()}
          tooltipText={valueInfo}
          showToottip={!!valueInfo}
        />
      </Grid>
    </Grid>
  );
};

export default NodeStatItem;
