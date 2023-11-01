/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import Router from 'next/router';
import { useEffect } from 'react';

const Ping = () => {
  useEffect(() => {
    const { pathname } = Router;
    if (pathname == '/ping') {
      Router.push('/home');
    }
  }, []);
  return null;
};

export default Ping;
