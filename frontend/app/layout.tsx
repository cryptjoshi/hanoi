import type { Metadata } from "next";
import localFont from "next/font/local";
import AdminPanelLayout from "@/components/admin-panel/admin-panel-layout";
import { Toaster } from "@/components/ui/toaster"
import "./globals.css";
import { kanit } from "@/lib/fonts";
 

import { dir } from 'i18next'
import { languages } from '@/app/i18n/settings'

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
  title: "Create Next App",
  description: "Generated by create next app",
  
};

 function RootLayout({
  children, params: {
    lng
  }
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html  lang={lng} dir={dir(lng)}>
      <body className={`${kanit.className} ${geistSans.variable} ${geistMono.variable} antialiased`}
      >
       {children} 
        <Toaster />
      </body>
    </html>
  );
}

export default RootLayout;