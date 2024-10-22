'use client'

import { useEffect, useState } from 'react'
import i18next from 'i18next'
import { initReactI18next, useTranslation as useTranslationOrg } from 'react-i18next'
import { useCookies } from 'react-cookie'
import resourcesToBackend from 'i18next-resources-to-backend'
import LanguageDetector from 'i18next-browser-languagedetector'
import { getOptions, languages, cookieName } from './settings'

const runsOnServerSide = typeof window === 'undefined'

// 
i18next
  .use(initReactI18next)
  .use(LanguageDetector)
  .use(resourcesToBackend((language:string, namespace:string) => import(`./locales/${language}/${namespace}.json`)))
  .init({
    ...getOptions(),
    lng: undefined, // let detect the language on client side
    detection: {
      order: ['path', 'htmlTag', 'cookie', 'navigator'],
    },
    preload: runsOnServerSide ? languages : []
  })

export function useTranslation(lng:string, ns:string, options:any) {
  const [cookies, setCookie] = useCookies([cookieName])
  const ret = useTranslationOrg(ns, options)
  const { i18n } = ret

  // eslint-disable-next-line react-hooks/rules-of-hooks
  const [activeLng, setActiveLng] = useState(i18n.resolvedLanguage)

  // eslint-disable-next-line react-hooks/rules-of-hooks
  useEffect(() => {
    if (runsOnServerSide) return;
    
    if (lng && i18n.resolvedLanguage !== lng) {
      i18n.changeLanguage(lng).then(() => {
        setActiveLng(lng);
        if (cookies.i18next !== lng) {
          setCookie(cookieName, lng, { path: '/' });
        }
      });
    } else if (i18n.resolvedLanguage !== activeLng) {
      setActiveLng(i18n.resolvedLanguage);
    }
  }, [lng, i18n, activeLng, cookies.i18next, setCookie]);

  if (runsOnServerSide && lng && i18n.resolvedLanguage !== lng) {
    i18n.changeLanguage(lng);
  }

  return ret;
}
