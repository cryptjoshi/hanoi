'use client'
//import { Signin } from "@/actions"
import React, { useState } from 'react'
import { Input } from "@/components/ui/input"
import {Button} from "@/components/ui/button"
//import useAuthStore from "@/store/auth" 
import { useForm, SubmitHandler } from "react-hook-form"
//import { Login } from '@/components/authen/login';
import { Separator } from "@/components/ui/separator"
// import AlertComponent from "@/components/AlertComponent"
// import { ExclamationTriangleIcon } from "@radix-ui/react-icons"
// import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
//import { Toaster } from "@/components/ui/toaster"
import {RegisterUser} from "@/actions"
import { useToast } from "@/hooks/use-toast"
import { ToastAction } from "@/components/ui/toast"
import { LucideEye, LucideEyeOff } from 'lucide-react'
type User = {
  username:string
  password:string
  repassword:string
  fullname:string
  prefix:string
  referred_by:string
  
  
}

 


//  const Register = async (body:User) =>{

  
//     // const state = useAuthStore()
//     try {
//     const response = await  fetch("http://167.71.100.123:3003/api/v1/users/register", { method: 'POST',
//           headers: {
//           'Accept': 'application/json',
//           'Content-Type': 'application/json',
//           //'Authorization': 'Bearer ' +  token
//           },
//           body: JSON.stringify({"username":body.username,"password":body.password,"fullname":body.fullname,"preferredname":body.username,"prefix":"psc"})
//         })
    
//        // if(!response.ok)
//       // return {"status":false}

//         return response.json()
//     }
//     catch(err:unknown) {

//         return err
//     }
//         // .then((response) => {
//         //     return new Promise((resolve) => response.json()
//         //       .then((json) => resolve({
//         //         status: response.status,
//         //         ok: response.ok,
//         //         json,
//         //       })));
//         //   }).catch((json) => resolve({json});
      
          
//         //   .then(({ status, json, ok }) => {
            
//         //     const message = json.message;
//         //     let color = 'black';
//         //     switch (status) {
//         //       case 400:
//         //         color = 'red';
//         //         break;
//         //       case 201:
//         //       case 200:
//         //         color = 'grey';
//         //         break;
//         //       case 500:
//         //       default:
//         //         handleUnexpected({ status, json, ok });
//         //     }
//         //   })
        
       
        
       
// }


export default function RegisterComponent({lng,refferedcode}:{lng:string,refferedcode:string}) {
 
 //const [iserror,setError] = React.useState(false)
 //const state = useAuthStore()
 const [showingA,setShowingA] = React.useState(false)
 const [showingB, setShowingB] = useState(false);
 const {toast} = useToast()
  const {
    register,
    reset,
    handleSubmit
  } = useForm<User>()

  const redirect = ()=>{
    location.replace(`/${lng}/login`)
}

  //const {login} = useAuthStore()
  const onSubmit: SubmitHandler<User> = (body:User) => {
    if (refferedcode) {
      body.referred_by = refferedcode; // ใส่ค่า refferedcode ลงใน referred_by
    }
   console.log(body)
   RegisterUser("ckd",body).then((response) =>{
 
    if(!response.Status){
        toast({
            variant: "destructive",
            title: "มีข้อผิดพลาด",
            description: response.message,
            action: <ToastAction altText="ผิดพลาด">{"ผิดพลาด"}</ToastAction>,
          })
          
    } else {
        toast({
            title: "สำเร็จ",
            description: response.message,
          })
       location.replace(`/${lng}/`)
    }
 
    
})
      
    
     
  //console.log('logged in')
  //  login(data)
 
  }
 
 

 
  // const handleLogOut() => {
  //     console.log('logged out')
  //     logout()
  //     }
  return (
    <>
    <form onSubmit={handleSubmit(onSubmit)}>
    <div className="bg-gray-100 min-h-screen flex items-center justify-center p-6">
      <div className="bg-white shadow-lg rounded-lg max-w-md mx-auto">
        <div className="px-6 py-4">
          <h2 className="text-gray-700 text-3xl font-semibold">Register</h2>
          <p className="mt-1 text-gray-600">Please register your account.</p>
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
            <div className="mt-4">
              <label className="block text-gray-700" htmlFor="username">
                Fullname
              </label>
              <Input
                type="text"
                id="fullname"
                className="mt-2 rounded w-full px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300"
                placeholder=""
                required
                defaultValue=""  
                {...register("fullname", { required: true })} 
              />
            </div>
            <div className="mt-4">
              <label className="block text-gray-700" htmlFor="referred_by">
                Referred By
              </label>
              <Input
                type="text"
                id="referred_by"
                className="mt-2 rounded w-full px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300"
                placeholder=""
                defaultValue={refferedcode}
                readOnly={!!refferedcode} 
                {...register("referred_by")} 
              />
            </div>
            <div className="mt-4">
              <label className="block text-gray-700" htmlFor="password">
                Password
              </label>
              <div className="flex items-center justify-between gap-2 ">
              <Input
                type={showingA ? "text" : "password"}
                id="password"
                className="mt-2 rounded w-full px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300"
                required
                defaultValue="" {...register("password",{required:true})} 
              />
              <button type="button" className="px-3 py-2 mt-2 bg-gray-700 text-white rounded hover:bg-gray-600" onClick={() => setShowingA(!showingA)}>
                  {showingA ? <LucideEye className="w-3 h-4"/> : <LucideEyeOff className="w-3 h-4"/>}
                </button> 
              </div>
            </div>
            <div className="mt-4">
              <label className="block text-gray-700" htmlFor="repassword">
                RePassword
              </label>
              <div className="flex items-center justify-between gap-2 ">
              <Input
                type={showingB ? "text" : "password"}
                id="repassword"
                className="mt-2 rounded w-full px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300"
                required
                defaultValue="" {...register("repassword", { required: true })} 
              />
              <button type="button" className="px-3 py-2 mt-2 bg-gray-700 text-white rounded hover:bg-gray-600" onClick={() => setShowingB(!showingB)}>
                  {showingB ? <LucideEye className="w-3 h-4"/> : <LucideEyeOff className="w-3 h-4"/>}
                </button>
              </div>
            </div>

            <div className="mt-6">
              <button type="submit" className="py-2 px-4 bg-gray-700 text-white rounded hover:bg-gray-600 w-full">
                Regiser
              </button>
            </div>
            <Separator className="my-4" />
            {/* <div className="mt-6">
              <button onClick={()=>{location.replace("/login")}} className="py-2 px-4 bg-gray-700 text-white rounded hover:bg-gray-600 w-full">
                Login
              </button>
            </div> */}
            <div className="mt-3">
            <Button type="button" onClick={redirect} className="py-2 px-4 bg-gray-700 text-white rounded hover:bg-gray-600 w-full">
                Cancel
              </Button>
            </div>
        </div>
      </div>
    </div>
    </form>
 
    </>
  )
}