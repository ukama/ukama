/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import i18n from 'i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import { initReactI18next } from 'react-i18next';

import { TRANSLATIONS_EN } from './translations/en';
import { TRANSLATIONS_FR } from './translations/fr';
const defaultLocale =
  typeof window !== 'undefined' &&
  (window.localStorage['i18nextLng']
    ? window.localStorage['i18nextLng']
    : 'en');
i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      en: {
        translation: TRANSLATIONS_EN,
      },
      fr: {
        translation: TRANSLATIONS_FR,
      },
    },
  });

i18n.changeLanguage(defaultLocale);
