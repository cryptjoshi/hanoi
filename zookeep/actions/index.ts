 
 
//import useAuthStore from "@/store/auth"
'use server'
import { AnyMxRecord } from 'dns';
import { redirect } from 'next/navigation'
 

type User = {
    username: string;
    fullname:string;
    password: string;
    prefix:string;
    referred_by:string;
}

type Dbstruct = {
  dbname:string;
  prefix:string;
  username:string;
  dbnames:string[];
}
export const Signin = async (body:User) =>{
     
 
       // const state = useAuthStore()

        const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/users/login`, { method: 'POST',
            headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            //'Authorization': 'Bearer ' +  token
            },
            body: JSON.stringify({"username":body.username,password:body.password,prefix:body.prefix})
          })
       return response.json()
}
export const GetDatabaseList = async () =>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/list`, { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
   // body: JSON.stringify(body)
  })
  return response.json()
}
export const CreateUser = async (body:Dbstruct) =>{
 
const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/create`, { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
    body: JSON.stringify({"dbname":body.prefix,"prefix":body.prefix,"username":body.username,"dbnames":body.dbnames})
  })
  return response.json()
}
export const UpdateDatabaseListByPrefix = async (body:Dbstruct) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"dbname":body.dbname,"prefix":body.prefix,"username":body.username,"dbnames":body.dbnames})
    })
    return response.json()
}  
export const GetDatabaseListByPrefix = async (prefix:string) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/prefix`, { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
   body: JSON.stringify({"prefix":prefix})
  })
  return response.json()
}
export const GetMemberList = async (prefix:string) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/member/all`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix})
})
return response.json()
}
export const GetMemberById = async (prefix:string,id:number) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/member/byid`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"ID":id})
})
return response.json()
}
export const GetUserInfo = async (token:string) =>{
  try {
  
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/users/info`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
  })
  return response.json()
}catch(error){
  console.log(error)
  return error
}
}
export const AddMember = async (prefix:string,body:any) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/member/create`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"body":body})
    })
    return response.json()
} 
export const UpdateMember = async (prefix:string,id:any,body:any) =>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/member/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"id":id,"body":body})
    })
    return response.json()
}
export const GetGameList = async (prefix:string) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/game/all`, { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
   body: JSON.stringify({"prefix":prefix})
  })
  return response.json()
}
export const GetGameStatus = async (prefix:string,token:string) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/game/status`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
      body: JSON.stringify({"prefix":prefix})
})
return response.json()
}
export const GetGameByType = async (token:string,id:string) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/game/bytype`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
      body: JSON.stringify({"id":id})
})
return response.json()
}
export const GetGameByProvide = async (token:string,provider:string,body:any) =>{
 
 
    const response = await fetch(`http://152.42.185.164:4007/callback/${provider}/gamelist`, { method: 'POST',
      headers: {   
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' +  token
        },
        body: JSON.stringify(body)
  })
  return response.json()
}
export const GetGameById = async (prefix:string,id:number) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/game/byid`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"id":id})
})
return response.json()
}
export const getGameUrl = async (url:string,data:any)=>{ //token:string,ProductID:string,username:string,currency:string) => {

  //const data  = {  "currency": currency, "productId": ProductID, "username": username, "sessionToken": token }
  console.log(data)
  const response = await fetch(url, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  data.sessionToken
      },
      body: JSON.stringify(data)
})
return response.json()
}
export const AddGame = async (prefix:string,body:any) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/game/create`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"body":body})
    })  
    return response.json()
}
export const UpdateGame = async (prefix:string,id:any,body:any) =>{
  
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/game/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json', 
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"id":id,"body":body})
    })
    return response.json()
}
export const AddPromotion = async (prefix:string,body:any) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/promotion/create`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"body":body})
    })
    return response.json()
}
export const UpdatePromotion = async (dbname: string, promotionId: any, values: { name: string; description: string; percentDiscount: string; startDate: string; endDate: string; maxDiscount: string; usageLimit: string; specificTime: string; paymentMethod: string; minSpend: string; maxSpend: string; termsAndConditions: string; status: string; }) =>{
 
  // console.log(JSON.stringify({"prefix":dbname,"promotionId":promotionId,"body":values}))
    const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/promotion/update`, { method: 'POST',
      headers: {   
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        },
        body: JSON.stringify({"prefix":dbname,"promotionId":promotionId,"body":values})
  })
  return response.json()
}
export const GetPromotionById = async (dbname:string,id:any) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/promotion/byid`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":dbname,"promotionId":id})
})
return response.json()
}
export const GetPromotionByUser = async (dbname:string,token:string) =>{
  try{
   const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/promotion/all`, { method: 'POST',
     headers: {   
       'Accept': 'application/json',
       'Content-Type': 'application/json',
       'Authorization': 'Bearer ' +  token
       },
       body: JSON.stringify({"prefix":dbname.toLowerCase()})
 })
 return response.json()
 }catch(error){
 
   console.log(error)
   return error
 }
} 
export const GetPromotion = async (token:string) =>{
 try{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/promotion/byuser`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
      //body: JSON.stringify({"prefix":dbname.toLowerCase()})
})
return response.json()
}catch(error){

  console.log(error)
  return error
}
}
export const DeletePromotion = async (dbname:string,id:any) =>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/promotion/delete`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":dbname,"promotionId":id})
})
return response.json()
}
export const GetExchangeRate = async (currency:string) =>{
  try{
  const response = await fetch(`http://152.42.185.164:4006/api/v1/db/exchange/rate`,{method:'POST',
  headers:{
    'Accept':'application/json',
    'Content-Type':'application/json'
  },
  body:JSON.stringify({"currency":currency})
})
  return response.json()
}catch(error){
  console.log(error)
  return error
}
}
export const UpdateMaster = async (prefix:string,id:any,body:any) =>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/master/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"id":id,"body":body})
    })
    return response.json()
}
export async function navigate(path:string) {
  redirect(path)
}
export const UpdateUserPromotion = async (token:string,body:any) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/users/update/pro`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
      body: JSON.stringify(body)
    })
    return response.json()
}
export const UpdateUser = async (token:string,body:any) =>{
 //console.log(JSON.stringify(body))
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/users/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
      body: JSON.stringify(body)
    })
    return response.json()
}
export const Deposit = async (token:string,body:any)=>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/statement/deposit`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
      body: JSON.stringify(body)
    })
    return response.json()
}
export const Withdraw = async (token:string,body:any)=>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/statement/withdraw`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
      body: JSON.stringify(body)
    })
    return response.json()
}
export const createTransaction = async (accessToken:string,body:any) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/transaction/add`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  accessToken
      },
      body: JSON.stringify({"Body":body})
})
return response.json()
}
export const GetHistory = async (token:string,prefix:string) =>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/statement/all`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
      body: JSON.stringify({"prefix":prefix})
})
return response.json()
}
export const GetTransaction = async (token:string,prefix:string) =>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/transaction/all`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  token
      },
      body: JSON.stringify({"prefix":prefix})
})
return response.json()
}
export const RegisterUser = async (prefix:string,body:User)=>{
 console.log(JSON.stringify(body))
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/users/register`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body:JSON.stringify({"username":body.username,"password":body.password,"fullname":body.fullname,"preferredname":body.username,"role":"user","prefix":prefix,"referred_by":body.referred_by})
 
    })
    return response.json()
}