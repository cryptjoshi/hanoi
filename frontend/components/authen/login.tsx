'use client'
import React from "react"
import { useRouter } from 'next/navigation';
import { Input } from "@/components/ui/input"
import useAuthStore from "@/store/auth" 
import {User} from "@/store/auth"
import { useForm, SubmitHandler } from "react-hook-form"
import { Separator } from "@/components/ui/separator"
import { ToastAction } from "@/components/ui/toast"
import { useToast } from "@/hooks/use-toast"
import { Button } from "react-day-picker";
import { LucideEyeOff,LucideEye, EyeOffIcon, EyeIcon } from "lucide-react";
import { useState } from "react"
import { useTranslation } from '@/app/i18n/client'
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

// Define the schema
const loginSchema = z.object({
  username: z.string().min(3, "Username must be at least 3 characters long"),
  password: z.string().min(6, "Password must be at least 6 characters long"),
});

// Infer the type from the schema
type LoginFormData = z.infer<typeof loginSchema>;

export default function Login({lng}:{lng:string}) {
  const router = useRouter();
  const { Signin, isLoggedIn } = useAuthStore();
  const [showing,setShowing] = React.useState(false)
  const { t } = useTranslation(lng, 'login', undefined  )
  const {toast} = useToast()
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });
 
 
 
  const onSubmit: SubmitHandler<LoginFormData> = async (data:LoginFormData) => {
    setIsSubmitting(true);
    try {
      const response = await Signin(data);  
    
     
     if (response) {
  
      router.push(`/${lng}/dashboard`);
    } else {
  
      toast({
        variant: "destructive",
        title: "มีข้อผิดพลาด",
        description: "เข้าสู่ระบบผิดพลาด!",
        action: <ToastAction altText="ผิดพลาด">{"ผิดพลาด"}</ToastAction>,
      })
    }
    } catch (error) {
      console.error("Login error:", error);
      toast({
        variant: "destructive",
        title: "เกิดข้อผิดพลาด",
        description: "มีข้อผิดพลาดในการเข้าสู่ระบบ โปรดลองอีกครั้ง",
      });
    } finally {
      setIsSubmitting(false);
    }
  }

  const redirect = ()=>{
    location.replace(`/${lng}/register`)
}

  React.useEffect(() => {
  if (isLoggedIn) {
      router.push(`/${lng}/dashboard`);  
  }
}, [isLoggedIn, router]);

  return (
    <>
    <form onSubmit={handleSubmit(onSubmit)}>
    <div className= "bg-gray-100 min-h-screen flex items-center justify-center p-6">
      <div className=" grow bg-white shadow-lg rounded-lg max-w-md mx-auto">
        <div className="px-6 py-4">
          <h2 className="text-gray-700 text-3xl font-semibold">{t('login.title')}</h2>
          <p className="mt-1 text-gray-600">{t('login.description')}</p>
        </div>
        <div className="px-6 py-4">
            <div className="mt-4">
              <label className="block text-gray-700" htmlFor="username">
                Username
              </label>
              <Input
                type="text"
                id="username"
                className={`mt-2 rounded w-full px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300 ${errors.username ? 'border-red-500' : ''}`}
                placeholder=""
                disabled={isSubmitting}
                {...register("username")}
              />
              {errors.username && <p className="text-red-500 text-sm mt-1">{errors.username.message}</p>}
            </div>
            <div className="mt-4 ">
              <label className="block text-gray-700" htmlFor="password">
                Password
              </label>
              <div className="flex items-center justify-between gap-2 ">
                <Input
                  type={showing ? "text" : "password"}
                  id="password"
                  className={`mt-2 rounded px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300 ${errors.password ? 'border-red-500' : ''}`}
                  disabled={isSubmitting}
                  {...register("password")}
                />
                <button type="button" className="px-3 py-2 mt-2 bg-gray-700 text-white rounded hover:bg-gray-600" onClick={() => setShowing(!showing)}>
                  {showing ? <LucideEye className="w-3 h-4"/> : <LucideEyeOff className="w-3 h-4"/>}
                </button>
              </div>
              {errors.password && <p className="text-red-500 text-sm mt-1">{errors.password.message}</p>}
            </div>
            <div className="mt-6">
              <button type="submit" className="py-2 px-4 bg-gray-700 text-white rounded hover:bg-gray-600 w-full">
                Login
              </button>
            </div>
            <Separator className="my-4" />
            <div className="mt-3">
            <button onClick={redirect} className="py-2 px-4 bg-gray-700 text-white rounded hover:bg-gray-600 w-full">
                Register
              </button>
            </div>
        </div>
      </div>
    </div>
    </form>
    </>
  )
}
