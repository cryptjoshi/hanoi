 
 
//import useAuthStore from "@/store/auth"
'use server'
import { redirect } from 'next/navigation'

type User = {
    username: string;
    password: string;
}


export const Signin = async (body:User) =>{
     
 
       // const state = useAuthStore()

        const response = await fetch("https://backend.paribrand.shop/api/v1/users/login", { method: 'POST',
            headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            //'Authorization': 'Bearer ' +  token
            },
            body: JSON.stringify({"username":body.username,password:body.password})
          })
       return response.json()
}


 

 
export async function navigate(path:string) {
  redirect(path)
}