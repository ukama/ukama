/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { AllocateSimApiDto, SimPoolResDto } from '@/client/graphql/generated';
import colors from '@/theme/colors';
import styled from '@emotion/styled';
import CloseIcon from '@mui/icons-material/Close';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { IconButton } from '@mui/material';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Button,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from '@mui/material';
import QRCode from 'qrcode.react';
import { useState } from 'react';

const CloseButtonStyle = styled(IconButton)({
  position: 'absolute',
  right: 10,
  top: 14,
});

interface SuccessStepProps {
  subscriberName: string;
  submissionData: AllocateSimApiDto;
  selectedSim: SimPoolResDto | null;
  onClose: () => void;
}

/** Step 3 — confirms successful subscriber creation, shows QR code for eSIM. */
const SuccessStep: React.FC<SuccessStepProps> = ({
  subscriberName,
  submissionData,
  selectedSim,
  onClose,
}) => {
  const [showQrCode, setShowQrCode] = useState(false);

  return (
    <>
      <DialogTitle>
        Successfully added [{subscriberName}]
        <CloseButtonStyle aria-label="close" onClick={onClose}>
          <CloseIcon />
        </CloseButtonStyle>
      </DialogTitle>

      <DialogContent>
        {submissionData.is_physical ? (
          <>
            <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 3 }}>
              You have successfully added {subscriberName} as a subscriber to
              your network, and a unique ID has been generated for them, which
              must be used to create a Ukama subscriber app.
            </Typography>
            <Box sx={{ bgcolor: 'grey.50', p: 2, borderRadius: 1, mb: 3 }}>
              <Typography fontFamily="monospace" sx={{ mb: 1 }}>
                UID: {submissionData.subscriber_id}
              </Typography>
              <Typography fontFamily="monospace">
                SIM ICCID: {submissionData.iccid}
              </Typography>
            </Box>
          </>
        ) : (
          <>
            <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 3 }}>
              You have successfully added {subscriberName} as a subscriber to
              your network, and an eSIM installation invitation has been sent
              out to them. If they would rather install their eSIM now, have
              them scan the QR code below.
            </Typography>
            <Accordion
              sx={{ boxShadow: 'none', background: 'transparent' }}
              onChange={(_, isExpanded) => setShowQrCode(isExpanded)}
            >
              <AccordionSummary
                expandIcon={<ExpandMoreIcon color="primary" />}
                sx={{
                  p: 0,
                  m: 0,
                  justifyContent: 'flex-start',
                  '& .MuiAccordionSummary-content': { flexGrow: 0.02 },
                }}
              >
                <Typography fontWeight={500} variant="caption" color={colors.primaryMain}>
                  {showQrCode ? 'HIDE QR CODE' : 'SHOW QR CODE'}
                </Typography>
              </AccordionSummary>
              <AccordionDetails sx={{ p: 2, display: 'flex', justifyContent: 'center' }}>
                <QRCode
                  id="qrCodeId"
                  value={selectedSim?.qrCode || ''}
                  style={{ height: 180, width: 180 }}
                />
              </AccordionDetails>
            </Accordion>
          </>
        )}
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} variant="contained">
          Close
        </Button>
      </DialogActions>
    </>
  );
};

export default SuccessStep;
