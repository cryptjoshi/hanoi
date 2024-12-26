// import { NextResponse } from 'next/server';
// import { withIronSession } from 'next-iron-session';
// import { User } from '@/store/auth';
// import { Signin } from '@/actions';

// const loginHandler = async (req: Request) => {
//     // ตรวจสอบว่ารับเป็น POST request
//     if (req.method !== 'POST') {
//         return NextResponse.json({ message: 'Method not allowed' }, { status: 405 });
//     }

//     const { username, password, prefix } = await req.json();

//     // ส่งคำร้องไปยัง backend เพื่อเข้าสู่ระบบ
//     const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/users/login`, {
//         method: 'POST',
//         headers: {
//             'Accept': 'application/json',
//             'Content-Type': 'application/json',
//         },
//         body: JSON.stringify({ username, password, prefix }),
//     });

//     const data = await response.json();

//     // ตรวจสอบสถานะการเข้าสู่ระบบ
//     if (response.ok && data.token) {
//         // สมมุติว่าคุณได้รับ token หลังจากเข้าสู่ระบบสำเร็จ
//         req.session.set('user', { username, token: data.token }); // เก็บข้อมูลที่จำเป็นใน session
//         await req.session.save();
//         return NextResponse.json({ message: 'Logged in' }, { status: 200 });
//     } else {
//         return NextResponse.json({ message: data.message || 'Invalid credentials' }, { status: 401 });
//     }
// };

// export default withIronSession(loginHandler, {
//     cookieName: 'zookeep_cookies',
//     password: process.env.PASSWORD_SECRET as string, // รหัสลับที่เข้มแข็ง
// });