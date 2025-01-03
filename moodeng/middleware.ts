import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import acceptLanguage from 'accept-language';
import { fallbackLng, languages, cookieName } from './app/i18n/settings';
import { getSession } from './actions';

acceptLanguage.languages(languages);

export async function middleware(request: NextRequest) {
  const url = request.nextUrl.clone();
  let lng = url.pathname.split('/')[1];
  const session = await getSession()
  // ตรวจสอบว่าเป็น root path หรือไม่
 
  if (url.pathname === '/') {
    lng = request.cookies.get('NEXT_LOCALE')?.value || fallbackLng;
    const isLoggedIn = session.isLoggedIn !== null;

    // if(isLoggedIn){
    //   return NextResponse.redirect(new URL(`/${lng}/dashboard/overview`, request.url));
    // }else{
    //   return NextResponse.redirect(new URL(`/${lng}/login`, request.url));
    // }
    return NextResponse.redirect(new URL(`/${lng}/dashboard`, request.url));
  }

  if (!languages.includes(lng)) {
    lng = acceptLanguage.get(request.headers.get('Accept-Language')) || fallbackLng;
  }

  //const isLoggedIn = request.cookies.get('isLoggedIn')?.value;

  // ตรวจสอบเส้นทางเมื่อเข้าหน้าแรกของภาษานั้นๆ
  if (url.pathname === `/${lng}` || url.pathname === `/${lng}/`) {
    
    //const isLoggedIn = session.isLoggedIn !== null;
   // console.log("session 18:",session.isLoggedIn)
    const redirectPath = session.isLoggedIn ? `/${lng}/dashboard/overview` : `/${lng}/login`;
    return NextResponse.redirect(new URL(redirectPath, request.url));
  }

  // จัดการกับ referer
  if (request.headers.has('referer')) {
    const refererUrl = new URL(request.headers.get('referer')!);
    const lngInReferer = languages.find((l) => refererUrl.pathname.startsWith(`/${l}`));
    const response = NextResponse.next();
    if (lngInReferer) {
      response.cookies.set(cookieName, lngInReferer);
    }
    return response;
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/:lng((?!api|_next/static|_next/image|assets|favicon.ico|sw.js|site.webmanifest).*)',
  ],
};
