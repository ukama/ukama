/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import {
  Box,
  Button,
  Card,
  CardContent,
  Skeleton,
  Tooltip,
  Typography,
} from '@mui/material';
import React from 'react';

interface UnassignedNodeProps {
  id: string;
  name: string;
  loading?: boolean;
  handleConfigureNode: (nodeId: string) => void;
}

const UnassignedNodeCard: React.FC<UnassignedNodeProps> = ({
  id,
  name,
  loading = false,
  handleConfigureNode,
}) => {
  const handleConfigureClick = (e: React.MouseEvent) => {
    e.stopPropagation();
    handleConfigureNode(id);
  };

  return (
    <Card
      sx={{
        border: `1px solid ${colors.darkGradient}`,
        borderRadius: 2,
        marginBottom: 2,
        backgroundColor: colors.lightGray,
        borderStyle: 'dashed',
        transition: 'transform 0.2s, box-shadow 0.2s',
        '&:hover': {
          transform: 'translateY(-2px)',
          boxShadow: '0 4px 8px rgba(0,0,0,0.1)',
        },
      }}
    >
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Box>
            <Typography
              variant="h6"
              sx={{
                display: 'inline-block',
                mb: 1,
                fontWeight: 'bold',
              }}
            >
              {loading ? <Skeleton width={150} /> : name}
            </Typography>

            <Tooltip title={id} placement="top-start">
              <Typography
                color="textSecondary"
                variant="body2"
                sx={{
                  whiteSpace: 'nowrap',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  maxWidth: '200px',
                }}
              >
                {loading ? (
                  <Skeleton width={200} />
                ) : (
                  `ID: ${id.substring(0, 15)}...`
                )}
              </Typography>
            </Tooltip>
          </Box>

          <Button
            variant="contained"
            color="primary"
            onClick={handleConfigureClick}
            disabled={loading}
            sx={{ textTransform: 'none' }}
          >
            Configure
          </Button>
        </Box>

        <Box display="flex" mt={3}>
          <Typography
            variant="body2"
            sx={{ fontStyle: 'italic', color: colors.darkGray }}
          >
            {loading ? (
              <Skeleton width={250} />
            ) : (
              'This node is not assigned to any site. Configure it to add to a site.'
            )}
          </Typography>
        </Box>
      </CardContent>
    </Card>
  );
};

export default UnassignedNodeCard;
