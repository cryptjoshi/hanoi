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
import { LucideEyeOff,LucideEye } from "lucide-react";
 
import { useTranslation } from '@/app/i18n/client'
 
  
export default function Login({lng}) {
  const router = useRouter();
  const { Signin, isLoggedIn } = useAuthStore();
  const [showing,setShowing] = React.useState(false)
  const { t } = useTranslation(lng, 'login')
  const {toast} = useToast()
  const {
    register,
    handleSubmit
  } = useForm<User>()
 
 
 
  const onSubmit: SubmitHandler<User> = async (data:User) => {
    try {

    
     const response = await Signin(data);  
    
     
     if (response) {
  
        router.push('/dashboard');
    } else {
  
      toast({
        variant: "destructive",
        title: "มีข้อผิดพลาด",
        description: "เข้าสู่ระบบผิดพลาด!",
        action: <ToastAction altText="ผิดพลาด">{"ผิดพลาด"}</ToastAction>,
      })
    }
    }
    catch
    {}
  }

  const redirect = ()=>{
    location.replace("/register")
}

  React.useEffect(() => {
  if (isLoggedIn) {
      router.push('/dashboard');  
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
                className="mt-2 rounded w-full px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300"
                placeholder=""
                required
                defaultValue=""  {...register("username", { required: true })} 
              />
            </div>
            <div className="mt-4 ">
              <label className="block text-gray-700" htmlFor="password">
                Password
              </label>
              <div className="flex items-center justify-between gap-2 ">
              <Input
                type={showing ? "text" : "password"}
                id="password"
                className="mt-2 rounded  px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300"
                required
                defaultValue="" {...register("password", { required: true })} 
              />
            <button type="button" className=" px-3 py-2 mt-2 bg-gray-700 text-white rounded hover:bg-gray-600" onClick={() => setShowing(!showing)}>{showing ? <LucideEye className="w-3 h-4"/> : <LucideEyeOff className="w-3 h-4"/>}</button>
            </div>
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