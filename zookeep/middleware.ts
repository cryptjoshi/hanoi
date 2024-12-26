import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import acceptLanguage from 'accept-language';
import { fallbackLng, languages, cookieName } from './app/i18n/settings';
import { getSession } from './actions';
acceptLanguage.languages(languages);

export async function  middleware(request: NextRequest) {
  const url = request.nextUrl.clone();
  let lng = url.pathname.split('/')[1];

const session = await getSession()
  

  //const isLoggedIn = request.cookies.get('isLoggedIn')?.value;
  // ตรวจสอบว่าเป็น root path หรือไม่

 
  if (url.pathname === '/') {
    lng = request.cookies.get('NEXT_LOCALE')?.value || fallbackLng;
    // const sessionResponse = await fetch(`${request.url}/api/session`, {
    //   method: 'GET',
    //   headers: {
    //     // ส่ง cookie ไปด้วย
    //     Cookie: request.headers.get('cookie') || '',
    //   },
    // });
  
    // const sessionData = await sessionResponse.json();
    const isLoggedIn = session.isLoggedIn !== null;

    if(isLoggedIn){
      return NextResponse.redirect(new URL(`/${lng}/home`, request.url));
    }else{
      return NextResponse.redirect(new URL(`/${lng}/login`, request.url));
    }
  }

  if (!languages.includes(lng)) {
    lng = acceptLanguage.get(request.headers.get('Accept-Language')) || fallbackLng;
  }

  // ตรวจสอบเส้นทางเมื่อเข้าหน้าแรกของภาษานั้นๆ
  if (url.pathname === `/${lng}` || url.pathname === `/${lng}/` ) {
    // เรียก API เพื่อดึง session
    // const sessionResponse = await fetch(`${request.url}/api/session`, {
    //   method: 'GET',
    //   headers: {
    //     Cookie: request.headers.get('cookie') || '',
    //   },
    // });
    // const sessionData = await sessionResponse.json();
    // const isLoggedIn = sessionData.user !== null;
    const isLoggedIn = session.isLoggedIn !== null;
    const redirectPath = isLoggedIn ? `/${lng}/home` : `/${lng}/login`;
    return NextResponse.redirect(new URL(redirectPath, request.url));
  }
  // // ตรวจสอบเส้นทางเมื่อเข้าหน้าแรกของภาษานั้นๆ
  // if (url.pathname === `/${lng}` || url.pathname === `/${lng}/`) {

  //   const redirectPath = isLoggedIn ? `/${lng}/home` : `/${lng}/login`;
  //   return NextResponse.redirect(new URL(redirectPath, request.url));
  // }

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

// export const config = {
//   matcher: [
//     '/:lng((?!api|_next/static|_next/image|assets|favicon.ico|sw.js|site.webmanifest).*)',
//   ],
// };

export const config = {
  matcher: [
    '/:lng((?!api|_next/static|_next/image|assets|favicon.ico|sw.js|site.webmanifest).*)',
  ],
};
