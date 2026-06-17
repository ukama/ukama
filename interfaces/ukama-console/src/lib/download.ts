/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

/**
 * Client-side download helpers. The BFF's getGeneratedPdfReport returns the
 * file as base64 (optionally as a `data:` URL), so we decode it to a Blob and
 * trigger a browser download via a temporary anchor.
 */

/** Decode a base64 string (raw or `data:<mime>;base64,...`) into a Blob. */
export function base64ToBlob(input: string, mime = 'application/octet-stream'): Blob {
  const base64 = input.startsWith('data:') ? input.slice(input.indexOf(',') + 1) : input;
  const binary = atob(base64);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i += 1) {
    bytes[i] = binary.charCodeAt(i);
  }
  return new Blob([bytes], { type: mime });
}

/** Save a Blob to the user's machine under `filename`. */
export function downloadBlob(blob: Blob, filename: string): void {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  window.URL.revokeObjectURL(url);
}
