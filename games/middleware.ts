import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import acceptLanguage from 'accept-language';
import { fallbackLng, languages, cookieName } from './app/i18n/settings';

acceptLanguage.languages(languages);

export function middleware(request: NextRequest) {
  // const url = request.nextUrl.clone();
  // let lng = url.pathname.split('/')[1];

  // // ตรวจสอบว่าเป็น root path หรือไม่
  // if (url.pathname === '/') {
  //   lng = request.cookies.get('NEXT_LOCALE')?.value || fallbackLng;
  //   return NextResponse.redirect(new URL(`/${lng}/games`, request.url));
  // }

  // if (!languages.includes(lng)) {
  //   lng = acceptLanguage.get(request.headers.get('Accept-Language')) || fallbackLng;
  // }

  // const isLoggedIn = request.cookies.get('isLoggedIn')?.value;

  // // ตรวจสอบเส้นทางเมื่อเข้าหน้าแรกของภาษานั้นๆ
  // if (url.pathname === `/${lng}` || url.pathname === `/${lng}/`) {
  //   const redirectPath = isLoggedIn ? `/${lng}/games` : `/${lng}/login`;
  //   return NextResponse.redirect(new URL(redirectPath, request.url));
  // }

  // // จัดการกับ referer
  // if (request.headers.has('referer')) {
  //   const refererUrl = new URL(request.headers.get('referer')!);
  //   const lngInReferer = languages.find((l) => refererUrl.pathname.startsWith(`/${l}`));
  //   const response = NextResponse.next();
  //   if (lngInReferer) {
  //     response.cookies.set(cookieName, lngInReferer);
  //   }
  //   return response;
  // }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/:lng((?!api|_next/static|_next/image|assets|favicon.ico|sw.js|site.webmanifest).*)',
  ],
};
