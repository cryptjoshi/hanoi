 
 
//import useAuthStore from "@/store/auth"
'use server'
import { redirect } from 'next/navigation'

type User = {
    username: string;
    password: string;
}

type Dbstruct = {
  dbname:string;
  prefix:string;
  username:string;
  dbnames:string[];
}
export const Signin = async (body:User) =>{
     
 
       // const state = useAuthStore()

        const response = await fetch("http://152.42.185.164:4006/api/v1/users/login", { method: 'POST',
            headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            //'Authorization': 'Bearer ' +  token
            },
            body: JSON.stringify({"username":body.username,password:body.password})
          })
       return response.json()
}

export const GetDatabaseList = async () =>{
  const response = await fetch("http://152.42.185.164:4006/api/v1/db/list", { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
   // body: JSON.stringify(body)
  })
  return response.json()
}

export const CreateUser = async (body:Dbstruct) =>{
 
const response = await fetch("http://152.42.185.164:4006/api/v1/db/create", { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
    body: JSON.stringify({"dbname":body.prefix,"prefix":body.prefix,"username":body.username,"dbnames":body.dbnames})
  })
  return response.json()
}


 

 
export async function navigate(path:string) {
  redirect(path)
}