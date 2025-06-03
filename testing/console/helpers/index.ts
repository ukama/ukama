/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */
import { faker } from '@faker-js/faker';
import fs from 'fs';
import { Parser } from 'json2csv';
import path from 'path';

interface SimRecord {
  'Unnamed: 0': number;
  ICCID: string;
  MSISDN: string;
  SmDpAddress: string;
  ActivationCode: number;
  IsPhysical: boolean;
  QrCode: string;
}

function generateSimRecord(index: number): SimRecord {
  return {
    'Unnamed: 0': index,
    ICCID: faker.string.numeric(19),
    MSISDN: faker.string.numeric(15),
    SmDpAddress: `1001.9.0.0.${faker.number.int({ min: 1, max: 9 })}`,
    ActivationCode: faker.number.int({ min: 1000, max: 9999 }),
    IsPhysical: faker.datatype.boolean(),
    QrCode: faker.string.uuid(),
  };
}

export async function createFakeSimCSV(
  recordCount: number,
  filePath: string,
): Promise<void> {
  const data: SimRecord[] = Array.from({ length: recordCount }, (_, i) =>
    generateSimRecord(i),
  );

  const fields: (keyof SimRecord)[] = [
    'Unnamed: 0',
    'ICCID',
    'MSISDN',
    'SmDpAddress',
    'ActivationCode',
    'IsPhysical',
    'QrCode',
  ];

  const parser = new Parser<SimRecord>({ fields });
  const csv = parser.parse(data);

  await fs.promises.mkdir(path.dirname(filePath), { recursive: true });
  await fs.promises.writeFile(filePath, csv, 'utf8');
  return;
}
