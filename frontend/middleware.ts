// import { NextResponse } from 'next/server';
// import type { NextRequest } from 'next/server';
// import acceptLanguage from 'accept-language'
// import { fallbackLng, languages, cookieName } from './app/i18n/settings'

// export function middleware(request: NextRequest) {
//   const url = request.nextUrl.clone();

//   // ดึงค่าจาก cookies ที่บันทึกสถานะล็อกอิน
//   const isLoggedIn = request.cookies.get('isLoggedIn');

//   if (!isLoggedIn) {
//     url.pathname = '/login';
//     return NextResponse.redirect(url);
//   }
//   if (url.pathname === '/') {
//     if (!isLoggedIn) {
//       url.pathname = '/login';
//     } else {
//       url.pathname = '/dashboard';
//     }
//     return NextResponse.redirect(url);
//   }  
//   return NextResponse.next();
// }

// export const config = {
//   matcher: ['/', '/dashboard'],
// };
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import acceptLanguage from 'accept-language';
import { fallbackLng, languages, cookieName } from './app/i18n/settings';

// กำหนดภาษาที่รองรับ
acceptLanguage.languages(languages);

export function middleware(request: NextRequest) {
  const url = request.nextUrl.clone();
//  let lng: string | undefined = url.pathname.split('/')[1]; // ดึงค่า locale จาก path
  
  // ดึงค่าจาก cookies ที่บันทึกสถานะล็อกอิน
  const isLoggedIn = request.cookies.get('isLoggedIn');
  
  // ตรวจสอบว่าเส้นทางเป็น /login และหากยังไม่ได้ล็อกอิน ไม่ทำการ redirect ซ้ำ
  if (url.pathname.startsWith('/login')) {
    if (!isLoggedIn) {
      return NextResponse.next(); // อนุญาตให้ดำเนินการต่อไปที่หน้า login
    } else {
      // หากผู้ใช้ล็อกอินแล้วและพยายามเข้า /login จะเปลี่ยนเส้นทางไปที่ /dashboard
      url.pathname = `/${url.pathname.split('/')[1]}/dashboard`;
      return NextResponse.redirect(url);
    }
  }

  // การจัดการ locale (i18n)
  let lng
  let cookiename = request.cookies.get(cookieName)
  if (request.cookies.has(cookieName)) lng = acceptLanguage.get(cookiename?.value)
  if (!lng) lng = acceptLanguage.get(request.headers.get('Accept-Language'))
  if (!lng) lng = fallbackLng

  // Redirect if lng in path is not supported
  if (
    !languages.some(loc => request.nextUrl.pathname.startsWith(`/${loc}`)) &&
    !request.nextUrl.pathname.startsWith('/_next')
  ) {
    return NextResponse.redirect(new URL(`/${lng}${request.nextUrl.pathname}`, request.url))
  }

  
  if (!languages.includes(lng)) {
    // ถ้า locale ไม่ถูกต้อง ให้ตั้งเป็นค่า fallback
    lng = fallbackLng;
  }

  // ตรวจสอบการล็อกอิน
  if (!isLoggedIn && !url.pathname.startsWith(`/${lng}/login`)) {
    url.pathname = `/${lng}/login`;  // ถ้ายังไม่ได้ล็อกอินให้ไปที่หน้า login
    return NextResponse.redirect(url);
  }

  // ตรวจสอบเส้นทางเมื่อเข้าหน้าแรก
  if (url.pathname === `/${lng}` || url.pathname === `/${lng}/`) {
    if (!isLoggedIn) {
      url.pathname = `/${lng}/login`; // ถ้ายังไม่ได้ล็อกอินให้ไปที่หน้า login
    } else {
      url.pathname = `/${lng}/dashboard`; // ถ้าล็อกอินแล้วให้ไปที่หน้า dashboard
    }
    return NextResponse.redirect(url);
  }

  // ตรวจสอบการเปลี่ยนภาษาผ่าน referer
  if (request.headers.has('referer')) {
    const refererUrl = new URL(request.headers.get('referer')!);
    const lngInReferer = languages.find((l) => refererUrl.pathname.startsWith(`/${l}`));
    const response = NextResponse.next();
    if (lngInReferer) {
      response.cookies.set(cookieName, lngInReferer);  // บันทึกภาษาที่เปลี่ยนแปลง
    }
    return response;
  }

  return NextResponse.next();
}

// การตั้งค่า matcher
export const config = {
  matcher: [
    '/:lng((?!api|_next/static|_next/image|assets|favicon.ico|sw.js|site.webmanifest).*)', // สำหรับเส้นทางที่ไม่ต้องการตรวจสอบ
  ],
};

