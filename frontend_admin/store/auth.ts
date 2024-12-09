import { create } from 'zustand';
import { persist, PersistOptions } from 'zustand/middleware';
//import { useRouter } from 'next/router';

export interface AuthStore {
    isLoggedIn: boolean;
    accessToken: string | null;
    accessTokenData: string | null;
    refreshToken: string | null;
    customerCurrency: string | null;
    Signin: (body: User) => Promise<any>;   
    Logout: () => void;
    setIsLoggedIn: (isLoggedIn: boolean | false) => void;
    setAccessToken: (accessToken: string | null) => void;
    setRefreshToken: (refreshToken: string | null) => void;
    setCustomerCurrency: (customerCurrency: string | null) => void;
    init: () => void;
    clearTokens: () => void;
    lng: string;
    setLng: (lng: string) => void;
}

export type User = {
    username: string;
    password: string;
};

const endpoint = `${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/login`//   `${"http://localhost"}:4006/api/v1/db/login`//`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/login`// process.env.NEXT_PUBLIC_BACKEND_ENDPOINT +"api/v1/users/login"

// เพิ่มฟังก์ชันสำหรับอ่านค่า cookie
function getCookie(name: string): string | null {
  if (typeof document === 'undefined') return null; // ตรวจสอบว่าอยู่ใน browser environment
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()?.split(';').shift() || null;
  return null;
}

const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      isLoggedIn: false,
      accessToken: null,
      accessTokenData: null,
      refreshToken: null,
      customerCurrency: "THB",
      Signin: async (body: User) => {
      //  const router = useRouter()
        try {
          const response = await fetch(endpoint, {
            method: 'POST',
            headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username: body.username, password: body.password, prefix: "psc" }),
          });

          const data = await response.json();
          
          if (data.Status) {
            set({
              isLoggedIn: true,
              accessToken: data?.token,
            });
            localStorage.setItem('isLoggedIn', JSON.stringify(true));
            document.cookie = "isLoggedIn=true; path=/";
           // router.redirect("/");
            // location.replace("/dashboard"); // หากต้องการ redirect ควรพิจารณาให้แน่ใจว่าใช้งานใน context ที่ถูกต้อง
            return data
          } else {
            set({ isLoggedIn: false, accessToken: null });
            localStorage.setItem('isLoggedIn', JSON.stringify(false));
            document.cookie = "isLoggedIn=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
            return data
          }
         // return false
        } catch (error) {
          console.error(error);
          return false;
        }
      },
      Logout: () => {
       // const router = useRouter()
        set({ isLoggedIn: false, accessToken: null });
        document.cookie = "isLoggedIn=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC;";
       //  router.push("/"); 
        // location.replace("/"); // แนะนำให้ใช้ใน context ที่ปลอดภัย เช่นใน useEffect หรือ handle event
      },
      setIsLoggedIn:(isLoggedIn:boolean)=>{
       // const isLoggedData = isLoggedIn || false;
        set({isLoggedIn})
      },
      setAccessToken: (accessToken: string | null) => {
        const accessTokenData = accessToken || null;
        set({ accessToken, accessTokenData });
      },
      setRefreshToken: (refreshToken: string | null) => set({ refreshToken }),
      setCustomerCurrency: (customerCurrency: string | null) => set({ customerCurrency }),
      init: () => {
        const { setAccessToken, setRefreshToken, setIsLoggedIn, setLng, setCustomerCurrency } = get();
        const isloggedIn = localStorage.getItem('isLoggedIn') == 'true';
        const accessToken = localStorage.getItem('accessToken');
        const refreshToken = localStorage.getItem('refreshToken');
        const lng = getCookie('lng') || 'en'; // Get language from cookie or use default
        setIsLoggedIn(isloggedIn);
        setAccessToken(accessToken);
        setRefreshToken(refreshToken);
        setLng(lng);
      },
      clearTokens: () => {
        set({
          accessToken: null,
          accessTokenData: null,
          refreshToken: null,
        });
      },
      lng: 'en', // Default language
      setLng: (lng: string) => {
        set({ lng });
        document.cookie = `lng=${lng}; path=/`;
      },
    }),
    {
      name: 'userLoginStatus',
      storage: {
          getItem: (name) => (typeof window !== 'undefined' ? localStorage.getItem(name) : null),
          setItem: (name, value) => {
              if (typeof window !== 'undefined') {
                  localStorage.setItem(name, JSON.stringify(value));
              }
          },
          removeItem: (name) => {
              if (typeof window !== 'undefined') {
                  localStorage.removeItem(name);
              }
          },
      },
  } as PersistOptions<AuthStore>
    // {
    //   name: 'userLoginStatus',
    //   storage: localStorage, // ใช้ localStorage โดยตรง ไม่ต้องใส่เป็น function
    // } as PersistOptions<AuthStore> // ไม่จำเป็นต้องแปลงชนิดข้อมูลเป็น unknown
  )
);

export default useAuthStore;
