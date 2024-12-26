'use client'
import Link from "next/link";

import { useEffect, useState } from "react";

import { useTranslation } from "@/app/i18n/client";
import useGameStore from "@/store/gameStore";
import { useRouter } from "next/navigation";
import { GetGameByProvide, getGameUrl,GetGameGC, getSession,getEFGameUrl } from "@/actions";
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
  //const {accessToken,user,customerCurrency} = useAuthStore()

 

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
            if(id!='9999'){
            const body = {
              "ProductID": id,
              "GameType": product,
              "LanguageCode": "4",
              "Platform": "0"
            }
            
            const response = await GetGameByProvide(provider,body)
          
             if(response.Status){
                 setGameList(response.Data?.games)
             } 
            }else {
               
              
              const responseg = await GetGameGC()
          

               if(responseg.status){
                const gameTokens: { token: string; url: string }[] = []; 
                const { token, url } = responseg.data; // ดึง token และ url จาก response
                gameTokens.push({ token, url }); 
                   setGameList(gameTokens)
               } 
            }

        }
      //  console.log(product)
    
     
        fetchgame(id)
    },[])
   
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
        // data  = {  "currency": session.customerCurrency || "USD", "productId": code, "username": session.username,"password":session.password, "sessionToken": session.token,"callbackUrl":"http://128.199.92.45:4002/en/games/list/1/8888" }
        getGameUrl("http://152.42.185.164:4007/api/v1/pg/launchgame",code).then((gameurl)=>{
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
      // data  = {  "currency": customerCurrency || "USD", "productId": code, "username": user.username,"password":user.password, "sessionToken": accessToken,"callbackUrl":"http://128.199.92.45:4002/en/games/list/1/8888" }
      // getGameUrl("http://152.42.185.164:4005/api/Auth/LaunchGame",data).then((gameurl)=>{
      // if(gameurl.Status){
       // url=gameurl.Data.url
      // console.log(gamelist[0].url)
      openInNewTab(gamelist[0].url)
      // }else {
      //   toast({
      //     title: t("common.error"),
      //     description: t("common.error"),
      //     variant: "destructive",
      //   });
      // }
      //})
      break;
      default:
       
        
     
         data  = {  
        "currency": "", 
        "username": "",  
         "ProductID": id,
         "GameType": product,
         "LanguageCode": "4",
         "Platform": "1",
         "sessionToken": ""
        }

         getEFGameUrl("http://152.42.185.164:4007/api/v1/ef/launchgame",data).then((gameurl)=>{
          console.log(gameurl)
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