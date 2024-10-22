import { dir } from 'i18next'
import { languages } from '@/app/i18n/settings'
import { ReactNode } from 'react'

export async function generateStaticParams() {
  return languages.map((lng) => ({ lng }))
}

interface LayoutProps {
  children: ReactNode;
  params: {
    lng: string;
  };
}

export default function Layout({
  children,
  params: { lng }
}: LayoutProps) {
  return (
    <div lang={lng} dir={dir(lng)}>
      {children}
    </div>
  )
}
