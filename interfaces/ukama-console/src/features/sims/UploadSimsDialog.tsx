/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/** Upload SIMs — pick/drop a CSV, base64-encode it, and call uploadSims. */
import { useRef, useState } from 'react';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import IconButton from '@mui/material/IconButton';
import CloseRounded from '@mui/icons-material/CloseRounded';
import CloudUploadRounded from '@mui/icons-material/CloudUploadRounded';
import DescriptionRounded from '@mui/icons-material/DescriptionRounded';
import UploadFileRounded from '@mui/icons-material/UploadFileRounded';
import { useUploadSimsMutation } from '@/client/graphql/sims.generated';
import { Sim_Types } from '@/client/graphql/types';
import AppModal from '@/components/AppModal';
import { useToast } from '@/components/ToastProvider';

const MAX_BYTES = 10 * 1024 * 1024; // 10 MB

const formatSize = (bytes: number): string =>
  bytes < 1024
    ? `${bytes} B`
    : bytes < 1024 * 1024
      ? `${(bytes / 1024).toFixed(0)} KB`
      : `${(bytes / (1024 * 1024)).toFixed(1)} MB`;

/** Reads a file as base64 (without the data-URL prefix), matching the BFF. */
const fileToBase64 = (file: File): Promise<string> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onloadend = () => {
      const result = reader.result as string;
      resolve(result.split(',')[1] ?? '');
    };
    reader.onerror = () => reject(new Error('Could not read the file'));
    reader.readAsDataURL(file);
  });

export default function UploadSimsDialog({
  onClose,
  onUploaded,
}: {
  onClose: () => void;
  onUploaded?: () => void;
}) {
  const toast = useToast();
  const [file, setFile] = useState<File | null>(null);
  const [drag, setDrag] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const [uploadSims, { loading: uploading }] = useUploadSimsMutation({
    onCompleted: (res) => {
      toast(`${res.uploadSims.iccid.length} SIMs uploaded successfully!`);
      onUploaded?.();
      onClose();
    },
    onError: (err) => toast(err.message),
  });

  const accept = (f: File | undefined) => {
    if (!f) return;
    if (!f.name.toLowerCase().endsWith('.csv')) {
      toast('Please choose a .csv file.');
      return;
    }
    if (f.size > MAX_BYTES) {
      toast('File is too large (max 10 MB).');
      return;
    }
    setFile(f);
  };

  const start = async () => {
    if (!file) return;
    try {
      const data = await fileToBase64(file);
      void uploadSims({
        variables: { data: { data, simType: Sim_Types.UkamaData } },
      });
    } catch (e) {
      toast(e instanceof Error ? e.message : 'Could not read the file');
    }
  };

  return (
    <AppModal
      title="Upload SIMs"
      width={520}
      onClose={onClose}
      footer={
        <>
          <Button
            color="inherit"
            sx={{ color: 'var(--uk-ink-3)' }}
            onClick={onClose}
            disabled={uploading}
          >
            Cancel
          </Button>
          <Button
            variant="contained"
            startIcon={
              uploading ? (
                <CircularProgress size={16} color="inherit" />
              ) : (
                <UploadFileRounded />
              )
            }
            disabled={!file || uploading}
            onClick={() => void start()}
          >
            {uploading ? 'Uploading…' : 'Upload'}
          </Button>
        </>
      }
    >
      <p style={{ fontSize: 13.5, color: 'var(--uk-ink-2)', lineHeight: 1.6, margin: '0 0 16px', textWrap: 'pretty' }}>
        Upload the SIM CSV file you received so that you can digitally assign SIMs to your
        customers, and authorize them to use your network.
      </p>

      <input
        ref={inputRef}
        type="file"
        accept=".csv,text/csv"
        hidden
        onChange={(e) => {
          accept(e.target.files?.[0]);
          e.target.value = '';
        }}
      />

      {file ? (
        <div className="sim-file">
          <DescriptionRounded sx={{ fontSize: 22, color: 'var(--uk-ac)' }} />
          <div style={{ flex: 1, minWidth: 0 }}>
            <div style={{ fontSize: 13.5, fontWeight: 600 }}>{file.name}</div>
            <div style={{ fontSize: 12, color: 'var(--uk-ink-3)' }}>
              {formatSize(file.size)}
            </div>
          </div>
          <IconButton
            size="small"
            aria-label="Remove file"
            sx={{ color: 'var(--uk-ink-3)' }}
            disabled={uploading}
            onClick={() => setFile(null)}
          >
            <CloseRounded sx={{ fontSize: 18 }} />
          </IconButton>
        </div>
      ) : (
        <button
          type="button"
          className={`dropzone${drag ? ' drag' : ''}`}
          onClick={() => inputRef.current?.click()}
          onDragOver={(e) => {
            e.preventDefault();
            setDrag(true);
          }}
          onDragLeave={() => setDrag(false)}
          onDrop={(e) => {
            e.preventDefault();
            setDrag(false);
            accept(e.dataTransfer.files?.[0]);
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
