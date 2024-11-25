'use client'
import Link from "next/link";

import { useEffect, useState } from "react";

import { useTranslation } from "@/app/i18n/client";
import useGameStore from "@/store/gameStore";
import { useRouter } from "next/navigation";
import { GetGameByProvide, getGameUrl } from "@/actions";
import useAuthStore from "@/store/auth";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { useToast } from "@/hooks/use-toast"

export default function GamesList({ product,id,lng }: { product:string,id:string,lng:string }) {
   //const { gameStatus, fetchGameStatus } = useGameStore()
  const {t} = useTranslation(lng,'translation',undefined);
  //const [gameid,setGameId] =useState(id)
  const router = useRouter()
  const [gamelist,setGameList] = useState<any[]>([{}])
  const {toast} = useToast()
  const {accessToken,user,customerCurrency} = useAuthStore()
  user
  if(accessToken){
    useEffect(()=>{
        const fetchgame = async (id:string) =>{
            let provider = "pg"
            if(id =='9999')
              provider = "gc"
            else if(id =="8888"){
              provider = "pgsoft"
            } else {
              provider = "ef"
            }
            
            const body = {
              "ProductID": id,
              "GameType": product,
              "LanguageCode": "4",
              "Platform": "0"
            }
            
            const response = await GetGameByProvide(accessToken,provider,body)
          
             if(response.Status){
                 setGameList(response.Data?.games)
             } else {

             }
        }
      //  console.log(product)
        fetchgame(id)
    },[])
  } else {
    router.push(`/${lng}/login`)
  }
//   const playgame = (ID:string) =>{
//       router.push(`/${lng}/games/${ID}`)
//   }
const openInNewTab = (url:string) => {
  const newWindow = window.open(url, '_blank', 'noopener,noreferrer')
  if (newWindow) newWindow.opener = null
}
   let url = "";
   let data
   const playgame = (code:string) =>{
    //console.log("product:"+product+"id:"+id)
    switch(id){
      case "8888":
         data  = {  "currency": customerCurrency || "USD", "productId": code, "username": user.username,"password":user.password, "sessionToken": accessToken }
        getGameUrl("http://152.42.185.164:4007/api/v1/pg/launchgame",data).then((gameurl)=>{
        if(gameurl.Status){
         // url=gameurl.Data.url
        openInNewTab(gameurl.Data.url)
        }else {
          toast({
            title: t("common.error"),
            description: t("common.error"),
            variant: "destructive",
          });
        }
        })
      
     break;
     case "9999":
      break;
      default:
       
        
     
         data  = {  "currency": customerCurrency || "USD", "username": user.username,  
         "ProductID": id,
         "GameType": product,
         "LanguageCode": "4",
         "Platform": "1",
         "sessionToken": accessToken
        }

         getGameUrl("http://152.42.185.164:4007/api/v1/ef/launchgame",data).then((gameurl)=>{
          
          if(gameurl.Data.errorcode == "0"){
           // url=gameurl.Data.url
         
          openInNewTab(gameurl.Data.url)
          }else{
            toast({
              title: t("common.error"),
              description: t("common.error"),
              variant: "destructive",
            });
          }
          })

        break;
    } 
  
   
  }

  //useEffect(() => {
   //fetchGameStatus(prefix)
    //console.log('fetchGameStatus',gameStatus)
  //}, [prefix, fetchGameStatus])
  //href={`/${lng}/games/list/${id}/play/${item.code}`}

  if (!gamelist) {
    return <div>Game Maintainance status...</div>
  }


    return (
        <div className="grid grid-cols-4 gap-2 sm:gap-4 p-4 sm:p-6">
        {gamelist.map((item: any, index: any) => (
           <div key={index} onClick={()=>playgame(item.code || item.GameCode)}
           className="cursor-pointer hover:opacity-80 transition-opacity"> {/* เพิ่ม hover effect */}
         <div  className="flex flex-col items-center">
          
           <Avatar>
            <AvatarImage src={item.img || item.ImageUrl} />
            <AvatarFallback>{item.name || item.GameName}</AvatarFallback>
          </Avatar>
        
          
           <span className="text-[10px] sm:text-xs text-center" >{`${item.name || item.GameName}`}</span>
          
         </div>
         </div>
       ))}
     </div>
    )
}