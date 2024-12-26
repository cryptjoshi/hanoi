
//import useAuthStore from "@/store/auth"
'use server'

import { redirect } from 'next/navigation'
import { SessionData } from "@/lib";
import { defaultSession, sessionOptions } from "@/lib";
import { getIronSession } from "iron-session";
import { cookies } from "next/headers";
 

type User = {
    username: string;
    fullname:string;
    password: string;
    prefix:string;
    referred_by:string;
    banknumber:string;
    bankname:string;
}
type Dbstruct = {
  dbname:string;
  prefix:string;
  username:string;
  dbnames:string[];
}

type ProBody = {
  prefix:string
  pro_status:string
}
export async function getSession() {
  const session = await getIronSession<SessionData>(cookies(), sessionOptions);

  if (!session.isLoggedIn) {
    session.isLoggedIn = defaultSession.isLoggedIn;
  }

  return session;
}

const port = ":4002"

export const Login = async (body:User) =>{
      
  const session = await getSession();

   // const state = useAuthStore()

    const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/users/login`, { method: 'POST',
        headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        //'Authorization': 'Bearer ' +  token
        },
        body: JSON.stringify({"username":body.username,password:body.password,prefix:body.prefix})
      })
      const data = await response.json();
       
      if(data.Status){
       
        session.isLoggedIn = data.Status;
        session.token = data.Token;
        session.username = data.Data.Username;
        session.userId = data.Data.ID
        session.prefix = data.Data.Prefix
        session.customerCurrency= data.Data.Currency
        session.lng = "en"
        await session.save();
      }
      return data
      //return response.json()
}
export async function Logout() {
  const session = await getSession();
  const lng = session.lng
  console.log(lng)
  session.destroy();
  redirect(`/${lng}/`)
}
export const Signin = async (body:User) =>{
      
      const session = await getSession();

       // const state = useAuthStore()

        const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/users/login`, { method: 'POST',
            headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            //'Authorization': 'Bearer ' +  token
            },
            body: JSON.stringify({"username":body.username,password:body.password,prefix:body.prefix})
          })
          const data = await response.json();
           
          if(data.Status){
           
            session.isLoggedIn = data.Status;
            session.token = data.Token;
            session.username = data.Data.Username;
            session.userId = data.Data.ID
            session.prefix = data.Data.Prefix
            
            await session.save();
          }
          
          return response.json()
}
export const GetDatabaseList = async () =>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/list`, { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
   // body: JSON.stringify(body)
  })
  return response.json()
}
export const CreateUser = async (body:Dbstruct) =>{
 
const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/create`, { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
    body: JSON.stringify({"dbname":body.prefix,"prefix":body.prefix,"username":body.username,"dbnames":body.dbnames})
  })
  return response.json()
}
export const UpdateDatabaseListByPrefix = async (body:Dbstruct) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"dbname":body.dbname,"prefix":body.prefix,"username":body.username,"dbnames":body.dbnames})
    })
    return response.json()
}  
export const GetDatabaseListByPrefix = async () =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/prefix`, { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
   body: JSON.stringify({"prefix":session.prefix})
  })
  return response.json()
}
export const GetMemberList = async () =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/member/all`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":session.prefix})
})
return response.json()
}
export const GetMemberById = async (id:number) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/member/byid`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":session.prefix,"ID":id})
})
return response.json()
}
export const GetUserInfo = async () =>{
  try {
    const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/users/info`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
  })
  return response.json()
}catch(error){
  console.log(error)
  return error
}
}
export const AddMember = async (body:any) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/member/create`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":session.prefix,"body":body})
    })
    return response.json()
} 
export const UpdateMember = async (id:any,body:any) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/member/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":session.prefix,"id":id,"body":body})
    })
    return response.json()
}
export const GetGameList = async () =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/game/all`, { method: 'POST',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
   body: JSON.stringify({"prefix":session.prefix})
  })
  return response.json()
}
export const GetGameStatus = async () =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/game/status`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      body: JSON.stringify({"prefix":session.prefix})
})
return response.json()
}
export const GetGameByType = async (id:string) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/game/bytype`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + session.token
      },
      body: JSON.stringify({"id":id})
})
return response.json()
}
export const GetGameGC = async () =>{
 
  const session = await getSession()
  const response = await fetch(`http://152.42.185.164:4005/api/Auth/LaunchGame`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      body: JSON.stringify({username:session.username})
})
return response.json()
}
export const GetGameByProvide = async (provider:string,body:any) =>{
 
  const session = await getSession()
 
    const response = await fetch(`http://152.42.185.164:4007/callback/${provider}/gamelist`, { method: 'POST',
      headers: {   
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' +  session.token
        },
        body: JSON.stringify(body)
  })
  return response.json()
}
export const GetGameById = async (id:number) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/game/byid`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":session.prefix,"id":id})
})
return response.json()
}
export const getEFGameUrl = async (url:string,data:any)=>{ //ProductID:string,username:string,currency:string) => {
  const session = await getSession()
    //const data  = {  "currency": currency, "productId": ProductID, "username": username, "sessionToken": token }
  // const data  = {  "currency": session.customerCurrency || "USD", "productId": code, "username": session.username,"password":session.password, "sessionToken": session.token,"callbackUrl":"http://128.199.92.45:4002/en/games/list/1/8888" }
     data.sessionToken = session.token
     data.currency = session.customerCurrency
     data.username = session.username
    const response = await fetch(url, { method: 'POST',
      headers: {   
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' +  session.token
        },
        body: JSON.stringify(data)
  })
  return response.json()
  }
  
export const getGameUrl = async (url:string,code:string)=>{ //ProductID:string,username:string,currency:string) => {
const session = await getSession()
  //const data  = {  "currency": currency, "productId": ProductID, "username": username, "sessionToken": token }
 const data  = {  "currency": session.customerCurrency || "USD", "productId": code, "username": session.username,"password":session.password, "sessionToken": session.token,"callbackUrl":"http://128.199.92.45:4002/en/games/list/1/8888" }
   
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
export const AddGame = async (body:any) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/game/create`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":session.prefix,"body":body})
    })  
    return response.json()
}
export const UpdateGame = async (id:any,body:any) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/game/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json', 
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":session.prefix,"id":id,"body":body})
    })
    return response.json()
}
export const AddPromotion = async (body:any) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/promotion/create`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":session.prefix,"body":body})
    })
    return response.json()
}
export const UpdatePromotion = async (dbname: string, promotionId: any, values: { name: string; description: string; percentDiscount: string; startDate: string; endDate: string; maxDiscount: string; usageLimit: string; specificTime: string; paymentMethod: string; minSpend: string; maxSpend: string; termsAndConditions: string; status: string; }) =>{
 
  // console.log(JSON.stringify({"prefix":dbname,"promotionId":promotionId,"body":values}))
    const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/promotion/update`, { method: 'POST',
      headers: {   
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        },
        body: JSON.stringify({"prefix":dbname,"promotionId":promotionId,"body":values})
  })
  return response.json()
}
export const GetPromotionById = async (dbname:string,id:any) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/promotion/byid`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":dbname,"promotionId":id})
})
return response.json()
}
export const GetPromotionByUser = async (dbname:string,) =>{
  const session = await getSession()
  try{
   const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/promotion/all`, { method: 'POST',
     headers: {   
       'Accept': 'application/json',
       'Content-Type': 'application/json',
       'Authorization': 'Bearer ' +  session.token
       },
       body: JSON.stringify({"prefix":dbname.toLowerCase()})
 })
 return response.json()
 }catch(error){
 
   console.log(error)
   return error
 }
} 
export const GetPromotionLog = async () =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/promotion/all`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      body: JSON.stringify({"prefix":session.prefix})
})
return response.json()
}
export const GetPromotion = async () =>{
  const session = await getSession()
 try{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/promotion/byuser`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      //body: JSON.stringify({"prefix":dbname.toLowerCase()})
})
return response.json()
}catch(error){

  console.log(error)
  return error
}
}
export const GetUserPromotion = async () =>{
  const session = await getSession()
  try{
   const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/users/promotions/byuser`, { method: 'POST',
     headers: {   
       'Accept': 'application/json',
       'Content-Type': 'application/json',
       'Authorization': 'Bearer ' +  session.token
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
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/promotion/delete`, { method: 'POST',
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
  const response = await fetch(`http://152.42.185.164${port}/api/v2/db/exchange/rate`,{method:'POST',
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
export const UpdateMaster = async (id:any,body:any) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/db/master/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":session.prefix,"id":id,"body":body})
    })
    return response.json()
}
export async function navigate(path:string) {
  redirect(path)
}
export const UpdateUserPromotion = async (body:ProBody) =>{
  
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/users/promotions`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      body: JSON.stringify({"proID":body.pro_status})
    })
    return response.json()
}
// export const UpdateUserPromotion = async (body:ProBody) =>{
//   console.log(body)
//   const session = await getSession()
//   const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/users/update/pro`, { method: 'POST',
//     headers: {   
//       'Accept': 'application/json',
//       'Content-Type': 'application/json',
//       'Authorization': 'Bearer ' +  session.token
//       },
//       body: JSON.stringify({"prefix":body.prefix,"pro_status":body.pro_status})
//     })
//     return response.json()
// }
export const UpdateUser = async (body:any) =>{
 
 const session = await getSession()
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/users/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      body: JSON.stringify(body)
    })
    return response.json()
}
export const Deposit = async (body:any)=>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/users/promotions/deposit`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      body: JSON.stringify(body)
    })
    return response.json()
}
export const Withdraw = async (body:any)=>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/statement/withdraw`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + session.token
      },
      body: JSON.stringify(body)
    })
    return response.json()
}
export const createTransaction = async (body:any) =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/transaction/add`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      body: JSON.stringify({"Body":body})
})
return response.json()
}
export const GetHistory = async () =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/statement/all`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      body: JSON.stringify({"prefix":session.prefix})
})
return response.json()
}
export const GetTransaction = async () =>{
  const session = await getSession()
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/transaction/all`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' +  session.token
      },
      body: JSON.stringify({"prefix":session.prefix})
})
return response.json()
}
export const RegisterUser = async (prefix:string,body:User)=>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/users/register`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body:JSON.stringify({"username":body.username,"password":body.password,"fullname":body.fullname,"preferredname":body.username,"banknumber":body.banknumber,"bankname":body.bankname,"role":"user","prefix":prefix,"referred_by":body.referred_by})
 
    })
    return response.json()
}
export const Webhoook = async ( uid:string,username: string,isexpired:string,isverify: string, prefix: string,method: string) =>{
 
  

  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/statement/webhook`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body:
      JSON.stringify({
            "TransactionID":uid,
            "isExpired":isexpired,
            "verify":isverify,
            "ref":username,
            "merchantID":prefix,
            "type":method /* payin,payout */ 
        })
    })
    
    return response.json()


  // const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}${port}/api/v2/statement/webhook`, { method: 'POST',
  //     headers: {   
  //       'Accept': 'application/json',
  //       'Content-Type': 'application/json',
  //       },
        // return 
      
  // })
  // return response.json()
}