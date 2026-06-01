/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NodeTypeEnum, SiteDto } from '@/client/graphql/generated';
import { Card, CardContent, Stack, Typography } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { format } from 'date-fns';
import React from 'react';

interface SiteInfoProps {
  selectedSite: SiteDto;
  address?: string | null;
  nodeIds?: string[];
  createdDate?: string;
}

interface InfoRowProps {
  label: string;
  children: React.ReactNode;
}

const InfoRow = ({ label, children }: InfoRowProps) => {
  return (
    <>
      <Grid size={{ xs: 12, md: 4 }}>
        <Typography variant="body2" color="text.secondary" fontWeight="medium">
          {label}
        </Typography>
      </Grid>
      <Grid size={{ xs: 12, md: 8 }}>{children}</Grid>
    </>
  );
};

const SiteInfo: React.FC<SiteInfoProps> = ({
  selectedSite,
  address,
  nodeIds = [],
  createdDate,
}) => {
  const formatMaybeDate = (value?: string | null) => {
    if (!value) return null;
    const parsed = new Date(value);
    if (Number.isNaN(parsed.getTime())) return null;
    return format(parsed, 'MMMM d, yyyy');
  };

  const formattedDate =
    formatMaybeDate(createdDate) ??
    formatMaybeDate(selectedSite.installDate ?? null) ??
    'Not available';

  const toNumberOrNull = (value: unknown) => {
    if (value == null) return null;
    if (typeof value === 'number') return Number.isFinite(value) ? value : null;
    if (typeof value === 'string') {
      const parsed = Number.parseFloat(value);
      return Number.isFinite(parsed) ? parsed : null;
    }
    return null;
  };

  const latitude = toNumberOrNull(selectedSite.latitude);
  const longitude = toNumberOrNull(selectedSite.longitude);
  const hasCoordinates = latitude != null && longitude != null;
  const latitudeHemisphere = latitude != null && latitude >= 0 ? 'N' : 'S';
  const longitudeHemisphere = longitude != null && longitude >= 0 ? 'E' : 'W';
  const formattedCoordinates = hasCoordinates
    ? `(${Math.abs(latitude)}° ${latitudeHemisphere}, ${Math.abs(
        longitude,
      )}° ${longitudeHemisphere})`
    : null;

  const locationLabel = address || selectedSite.location || 'Not available';

  return (
    <Card
      sx={{
        borderRadius: 2,
        boxShadow: '0px 2px 6px rgba(0, 0, 0, 0.05)',
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <CardContent sx={{ padding: 2, flexGrow: 1, minHeight: 0 }}>
        <Typography variant="h6" sx={{ mb: { xs: 2, md: 4 } }}>
          Site information
        </Typography>

        <Grid container spacing={{ xs: 1, md: 3 }}>
          <InfoRow label="Tower Node:">
            <Typography variant="body2">
              {nodeIds.map((nodeId) =>
                nodeId.includes(NodeTypeEnum.Tnode) ? nodeId : '',
              )}
            </Typography>
          </InfoRow>
          <InfoRow label="Controller Node:">
            <Typography variant="body2">
              {nodeIds.map((nodeId) =>
                nodeId.includes(NodeTypeEnum.Cnode) ? nodeId : '',
              )}
            </Typography>
          </InfoRow>
          <InfoRow label="Amplifier Node:">
            <Typography variant="body2">
              {nodeIds.map((nodeId) =>
                nodeId.includes(NodeTypeEnum.Anode) ? nodeId : '',
              )}
            </Typography>
          </InfoRow>

          <InfoRow label="Date created:">
            <Typography variant="body2">{formattedDate}</Typography>
          </InfoRow>

          <InfoRow label="Location:">
            <Stack spacing={0.5}>
              <Typography variant="body2">{locationLabel}</Typography>
              {formattedCoordinates ? (
                <Typography variant="body2" color="text.secondary">
                  {formattedCoordinates}
                </Typography>
              ) : null}
            </Stack>
          </InfoRow>
        </Grid>
      </CardContent>
    </Card>
  );
};

export default SiteInfo;
