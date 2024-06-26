/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import colors from '@/theme/colors';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Grid,
  Typography,
} from '@mui/material';
import QRCode from 'qrcode.react';
import React, { useState } from 'react';

interface SubscriberDialogProps {
  qrCode: string;
  simType: string;
}

const Step4: React.FC<SubscriberDialogProps> = ({ qrCode, simType }) => {
  const [showQrCode, setShowQrCode] = useState(false);

  return (
    <Grid container>
      <Grid item xs={12}>
        {simType === 'eSim' ? (
          <Accordion
            sx={{ boxShadow: 'none', background: 'transparent' }}
            onChange={(_, isExpanded: boolean) => {
              setShowQrCode(isExpanded);
            }}
          >
            <AccordionSummary
              expandIcon={<ExpandMoreIcon color="primary" />}
              sx={{
                p: 0,
                m: 0,
                justifyContent: 'flex-start',
                '& .MuiAccordionSummary-content': {
                  flexGrow: 0.02,
                },
              }}
            >
              <Typography
                fontWeight={500}
                variant="caption"
                color={colors.primaryMain}
              >
                {showQrCode ? 'HIDE QR CODE' : 'SHOW QR CODE'}
              </Typography>
            </AccordionSummary>
            <AccordionDetails
              sx={{ p: 2, display: 'flex', justifyContent: 'center' }}
            >
              <QRCode
                id="qrCodeId"
                value={qrCode}
                style={{ height: 164, width: 164 }}
              />
            </AccordionDetails>
          </Accordion>
        ) : (
          <Box sx={{ pl: 2 }}>
            <Typography variant="body1">{`pSIM ICCID : ${qrCode}`}</Typography>
          </Box>
        )}
      </Grid>
    </Grid>
  );
};

export default Step4;
