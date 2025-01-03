 
 
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

const url = "http://backendservice"// process.env.NEXT_PUBLIC_BACKEND_ENDPOINT

export const Signin = async (body:User) =>{
     
 
       // const state = useAuthStore()

        const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/users/login`, { method: 'POST',
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
export const GetPartnerList = async (prefix:string) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/partner/all`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix})
})
return response.json()
}

export const  GetPartnerSeed = async (prefix:string) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/partner/checkseed`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix})
})
return response.json()
}
export const GetPartnerById = async (prefix:string,id:number) =>{
  
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/partner/byid`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"ID":id})
})
return response.json()
}
export const AddPartner = async (prefix:string,body:any) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/partner/create`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"body":body})
    })
    return response.json()
} 

export const UpdatePartner = async (prefix:string,id:any,body:any) =>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/partner/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"id":id,"body":body})
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

export const GetGameStatus = async (prefix:string) =>{
 
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/game/status`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix})
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

    export const UpdatePromotion = async (dbname: string, promotionId: any, values: any) =>{
 
   //console.log(JSON.stringify({"prefix":dbname,"promotionId":promotionId,"body":values}))
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
export const GetPromotion = async (dbname:string) =>{
 try{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/promotion/all`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":dbname})
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
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/exchange/rate`,{method:'POST',
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

export const GetDBMode = async (dbname:string) =>{
  try{
    const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/setting`,{method:'POST',
    headers:{
      'Accept':'application/json',
      'Content-Type':'application/json'
    },
    body: JSON.stringify({"prefix":dbname})
  })
    return response.json()
  }catch(error){
    console.log(error)
    return error
  }
}


export const UpdateMaster = async (prefix:string,id:any,body:any) =>{
  
  //console.log(body)

  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/master/update`, { method: 'POST',
    headers: {   
      'Accept': 'application/json',
      'Content-Type': 'application/json',
      },
      body: JSON.stringify({"prefix":prefix,"id":id,"body":body})
    })
    return response.json()
}

export const GetCommission = async (prefix:string) =>{
  const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/db/master/commission`,{method: 'Post',
  headers: {   
    'Accept': 'application/json',
    'Content-Type': 'application/json',
    },
    body: JSON.stringify({"prefix":prefix})
})
  return response.json()
}

export async function navigate(path:string) {
  redirect(path)
}