import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';
import acceptLanguage from 'accept-language';
import { fallbackLng, languages, cookieName } from './app/i18n/settings';
//import { i18n } from '@/app/i18n/'; // ปรับ path ตามโครงสร้างโปรเจคของคุณ

// กำหนดภาษาที่รองรับ
acceptLanguage.languages(languages);

export function middleware(request: NextRequest) {
  const url = request.nextUrl.clone();
  const { pathname } = url;

  // ตรวจสอบว่าเป็น root path หรือไม่
  if (pathname === '/' || pathname === '') {
    // ดึงภาษาจาก cookie หรือใช้ภาษาเริ่มต้น
    const lng = request.cookies.get('NEXT_LOCALE')?.value ||fallbackLng;
    // Redirect ไปยัง URL ที่มีภาษากำกับ
    return NextResponse.redirect(new URL(`/${lng}`, request.url));
  }

  // ดึงค่า lng จาก URL path
  let lng: string = url.pathname.split('/')[1];
  

  
  const cookieValue = request.cookies.get(cookieName)?.value;
  if (request.cookies.has(cookieName) && cookieValue) {
    lng = acceptLanguage.get(cookieValue) || lng;
  }
  if (!languages.includes(lng)) {
    const acceptLang = request.headers.get('Accept-Language');
    lng = acceptLang ? acceptLanguage.get(acceptLang) || fallbackLng : fallbackLng;
  }
  // ดึงค่าจาก cookies ที่บันทึกสถานะล็อกอิน
  const isLoggedIn = request.cookies.get('isLoggedIn')?.value;
  // ตรวจสอบว่า lng อยู่ในรายการภาษาที่รองรับหรือไม่
  if (!languages.includes(lng)) {
    lng = fallbackLng;
  }
  // ตรวจสอบเส้นทางเมื่อเข้าหน้าแรก
  if (url.pathname === `/${lng}` || url.pathname === `/${lng}/`) {
    if (!isLoggedIn) {
      url.pathname = `/${lng}/login`; // ถ้ายังไม่ได้ล็อกอินให้ไปที่หน้า login
    } else {
      url.pathname = `/${lng}/dashboard`; // ถ้าล็อกอินแล้วให้ไปที่หน้า dashboard
    }
    return NextResponse.redirect(url);
  } else {
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



  // else {
  //   // if (!isLoggedIn) {
  //   //   url.pathname = `/${lng}/login`; // ถ้ายังไม่ได้ล็อกอินให้ไปที่หน้า login
  //   // } else {
  //   //   url.pathname = `/${lng}/dashboard`; // ถ้าล็อกอินแล้วให้ไปที่หน้า dashboard
  //   // }
  //  // return NextResponse.redirect(url);
 


  // // ดึงค่า lng จาก cookie (client-side)
  // const clientLng = request.cookies.get(cookieName)?.value;
 
  // // ถ้า lng จาก URL ไม่ตรงกับ lng จาก cookie
  // if (clientLng && clientLng !== lng) {
   
  //     // ตรวจสอบว่าเส้นทางเป็น /login และหากยังไม่ได้ล็อกอิน ไม่ทำการ redirect ซ้ำ
  //   if (url.pathname.startsWith('/login')) {
  //     if (!isLoggedIn) {
  //       return NextResponse.next(); // อนุญาตให้ดำเนินการต่อไปที่หน้า login
  //     } else {
  //       // สร้าง URL ใหม่โดยใช้ lng จาก cookie
  //       url.pathname = `/${clientLng}${url.pathname.substring(3)}/dashboard`;
        
  //       return NextResponse.redirect(url);
  //     }
  //   }

  //   // // การจัดการ locale (i18n)

  


  //   // if (url.pathname.startsWith('/login')) {
  //   //   if (!isLoggedIn) {
  //   //     return NextResponse.next(); // อนุญาตให้ดำเนินการต่อไปที่หน้า login
  //   //   } else {
  //   //     // หากผู้ใช้ล็อกอินแล้วและพยายามเข้า /login จะเปลี่ยนเส้นทางไปที่ /dashboard
  //   //     url.pathname = `/${url.pathname.split('/')[1]}/dashboard`;
  //   //     return NextResponse.redirect(url);
  //   //   }
  //   // }

  //   // สร้าง response เพื่อ redirect
  //   const response = NextResponse.redirect(url);
    
  //   // อัปเดต cookie ให้ตรงกับ lng ใหม่
  //   response.cookies.set(cookieName, clientLng);
    
  //   return response;
  // } else {
    
  //   // ตรวจสอบการล็อกอิน
  //   if (!isLoggedIn && !url.pathname.startsWith(`/${lng}/login`)) {
  //     url.pathname = `/${lng}/login`;  // ถ้ายังไม่ได้ล็อกอินให้ไปที่หน้า login
  //     return NextResponse.redirect(url);
  //   }
    
  //   // ถ้า url.pathname เป็น / ให้เปลี่ยนเป็น /{lng}/dashboard
  //   if (url.pathname === '/') {
  //     url.pathname = `/${lng}/dashboard`;
  //   }
   
  //   return NextResponse.redirect(url);   
  // }
  // }
  // ถ้า lng จาก URL ตรงกับ lng จาก cookie หรือไม่มี cookie
  // ให้อัปเดตหรือตั้งค่า cookie ใหม่
  // const response = NextResponse.next();
  // response.cookies.set(cookieName, lng);

  
  //  let lng: string | undefined = url.pathname.split('/')[1]; // ดึงค่า locale จาก path
  
  
  
  // // ตรวจสอบว่าเส้นทางเป็น /login และหากยังไม่ได้ล็อกอิน ไม่ทำการ redirect ซ้ำ
  // if (url.pathname.startsWith('/login')) {
  //   if (!isLoggedIn) {
  //     return NextResponse.next(); // อนุญาตให้ดำเนินการต่อไปที่หน้า login
  //   } else {
  //     // หากผู้ใช้ล็อกอินแล้วและพยายามเข้า /login จะเปลี่ยนเส้นทางไปที่ /dashboard
  //     url.pathname = `/${url.pathname.split('/')[1]}/dashboard`;
  //     return NextResponse.redirect(url);
  //   }
  // }

  // // การจัดการ locale (i18n)


  // // Redirect if lng in path is not supported
  // if (
  //   !languages.some(loc => request.nextUrl.pathname.startsWith(`/${loc}`)) &&
  //   !request.nextUrl.pathname.startsWith('/_next')
  // ) {
  //   return NextResponse.redirect(new URL(`/${lng}${request.nextUrl.pathname}`, request.url))
  // }

  
  // if (!languages.includes(lng)) {
  //   // ถ้า locale ไม่ถูกต้อง ให้ตั้งเป็นค่า fallback
  //   lng = fallbackLng;
  // }

  



  //ตรวจสอบการเปลี่ยนภาษาผ่าน referer
 

  // return NextResponse.next();
}

// การตั้งค่า matcher
export const config = {
  matcher: [
    '/:lng((?!api|_next/static|_next/image|assets|favicon.ico|sw.js|site.webmanifest).*)', // สำหรับเส้นทางที่ไม่ต้องการตรวจสอบ
  ],
};
