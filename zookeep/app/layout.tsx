import type { Metadata } from "next";
import localFont from "next/font/local";
import AdminPanelLayout from "@/components/admin-panel/admin-panel-layout";
import { Toaster } from "@/components/ui/toaster"
import "./globals.css";
import { kanit } from "@/lib/fonts";
 
import clsx from 'clsx';
import { dir } from 'i18next'
import { languages } from '@/app/i18n/settings'
import LanguageSwitcher from "@/components/LanguageSwitcher"
import Providers from '@/app/providers'
export async function generateStaticParams() {
  return languages.map((lng) => ({ lng }))
}


const geistSans = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-geist-sans",
  weight: "100 900",
});
const geistMono = localFont({
  src: "./fonts/GeistMonoVF.woff",
  variable: "--font-geist-mono",
  weight: "100 900",
});

export const metadata: Metadata = {
  metadataBase: new URL(
    process.env.APP_URL
    ? `${process.env.APP_URL}`
    : process.env.VERCEL_URL
    ? `https://${process.env.VERCEL_URL}`
    : `http://localhost:${process.env.PORT || 3001}`
  ),
  title: "PLAY WITH ME",
  description: "play with me",
  
};

export default function RootLayout({
  children,
  params: { lng }
}: Readonly<{
  children: React.ReactNode;
  params: { lng: string };
}>) {
  const isRTL = lng === 'ar';
  return (
    <html lang={lng} dir={isRTL ? 'rtl' : 'ltr'}>
      <body 
        className={clsx(
          kanit.className,
          geistSans.variable,
          geistMono.variable,
          'antialiased relative min-h-screen',
          isRTL && 'rtl'
        )}
      >
        <Providers>
          {children}
          <Toaster />
        </Providers>
      </body>
    </html>
  );
}
