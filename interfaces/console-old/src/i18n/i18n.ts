import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import { TRANSLATIONS_EN } from "./translations/en";
import { TRANSLATIONS_FR } from "./translations/fr";
const defaultLocale = localStorage["i18nextLng"]
    ? localStorage["i18nextLng"]
    : "en";
i18n.use(LanguageDetector)
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
