/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { IP_API_BASE_URL, IPFY_URL } from '@/constants';

const getMetaInfo = async () => {
  return await fetch(IPFY_URL, {
    method: 'GET',
  })
    .then((response) => response.text())
    .then((data) => JSON.parse(data))
    .then((data) =>
      fetch(`${IP_API_BASE_URL}/${data.ip}/json/`, {
        method: 'GET',
      }),
    )
    .then((response) => response.text())
    .then((data) => JSON.parse(data))
    .catch((err) => {
      console.log(err);
      return {};
    });
};

export { getMetaInfo };
