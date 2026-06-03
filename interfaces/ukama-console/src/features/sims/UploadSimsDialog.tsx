/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';
import Meter from '@/components/Meter';

/** Upload SIMs — dropzone → file chip → progress → toast (form-dialogs.jsx). */
import { useEffect, useRef, useState } from 'react';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import CloseRounded from '@mui/icons-material/CloseRounded';
import CloudUploadRounded from '@mui/icons-material/CloudUploadRounded';
import DescriptionRounded from '@mui/icons-material/DescriptionRounded';
import UploadFileRounded from '@mui/icons-material/UploadFileRounded';
import AppModal from '@/components/AppModal';
import { useToast } from '@/components/ToastProvider';

type Stage = 'idle' | 'ready' | 'uploading';

export default function UploadSimsDialog({ onClose }: { onClose: () => void }) {
  const toast = useToast();
  const [file, setFile] = useState<{ name: string; size: string } | null>(null);
  const [stage, setStage] = useState<Stage>('idle');
  const [pct, setPct] = useState(0);
  const [drag, setDrag] = useState(false);
  const timer = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(
    () => () => {
      if (timer.current) clearInterval(timer.current);
    },
    [],
  );

  const pick = () => {
    setFile({ name: 'kwacha-sims-2026.csv', size: '248 SIMs' });
    setStage('ready');
  };

  const start = () => {
    setStage('uploading');
    setPct(8);
    timer.current = setInterval(() => {
      setPct((p) => {
        if (p >= 100) {
          if (timer.current) clearInterval(timer.current);
          setTimeout(() => {
            onClose();
            toast('SIMs uploaded successfully!');
          }, 350);
          return 100;
        }
        return p + 11;
      });
    }, 160);
  };

  return (
    <AppModal
      title="Upload SIMs"
      width={520}
      onClose={onClose}
      footer={
        stage === 'uploading' ? (
          <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={onClose}>
            Cancel
          </Button>
        ) : (
          <>
            <Button color="inherit" sx={{ color: 'var(--uk-ink-3)' }} onClick={onClose}>
              Cancel
            </Button>
            <Button
              variant="contained"
              startIcon={<UploadFileRounded />}
              disabled={!file}
              onClick={start}
            >
              Upload
            </Button>
          </>
        )
      }
    >
      <p style={{ fontSize: 13.5, color: 'var(--uk-ink-2)', lineHeight: 1.6, margin: '0 0 16px', textWrap: 'pretty' }}>
        Upload the SIM CSV file you received so that you can digitally assign SIMs to your
        customers, and authorize them to use your network.
      </p>
      {stage === 'uploading' ? (
        <div style={{ padding: '8px 0 4px' }}>
          <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 8 }}>Uploading file…</div>
          <Meter value={pct} height={8} />
          <div style={{ fontSize: 12, color: 'var(--uk-ink-3)', marginTop: 6 }}>{file?.name}</div>
        </div>
      ) : file ? (
        <div className="sim-file">
          <DescriptionRounded sx={{ fontSize: 22, color: 'var(--uk-ac)' }} />
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ fontSize: 13.5, fontWeight: 600 }}>{file.name}</div>
            <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>{file.size}</div>
          </div>
          <IconButton
            size="small"
            aria-label="Remove file"
            sx={{ color: 'var(--uk-ink-3)' }}
            onClick={() => {
              setFile(null);
              setStage('idle');
            }}
          >
            <CloseRounded sx={{ fontSize: 18 }} />
          </IconButton>
        </div>
      ) : (
        <button
          type="button"
          className={`dropzone${drag ? ' drag' : ''}`}
          onClick={pick}
          onDragOver={(e) => {
            e.preventDefault();
            setDrag(true);
          }}
          onDragLeave={() => setDrag(false)}
          onDrop={(e) => {
            e.preventDefault();
            setDrag(false);
            pick();
          }}
        >
          <CloudUploadRounded sx={{ fontSize: 34, color: 'var(--uk-ink-3)' }} />
          <div style={{ fontSize: 13.5, color: 'var(--uk-ink-2)', marginTop: 8 }}>
            <b style={{ color: 'var(--uk-ac-dark)' }}>Browse</b> or drag &amp; drop to upload a file
          </div>
          <div style={{ fontSize: 12, color: 'var(--uk-ink-3)', marginTop: 3 }}>CSV up to 10 MB</div>
        </button>
      )}
    </AppModal>
  );
}
