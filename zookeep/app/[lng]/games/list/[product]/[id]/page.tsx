'use client'
import { GetGameByProvide, getGameUrl } from "@/actions"
import { useTranslation } from "@/app/i18n/client"
import { ContentLayout } from "@/components/admin-panel/content-layout"
import GamesList from "@/components/games/Gameslist"
import {Options} from "@/components/games/options"
import useAuthStore from "@/store/auth"
import { useEffect } from "react"
//import { Route } from "lucide-react"


export default  function GamePage({ params: { lng,id,product } }: { params: { lng: string,id:string,product:string } }) {

   

    
    const {t} =  useTranslation(lng,"translation",undefined)
    const {accessToken,user,customerCurrency} = useAuthStore()

    useEffect(()=>{

        const fetchData = async ()=>{
            
            
            // const  gameslist = await GetGameByProvide(accessToken || "",product);

            // if(gameslist.Status){

            // }
            // const openInNewTab = (url:string) => {
            //     const newWindow = window.open(url, '_blank', 'noopener,noreferrer')
            //     if (newWindow) newWindow.opener = null
            // }
        
            // if(gameurl.Status){
            //     //console.log(gameurl.Data.url)
            //  //   window.open(gameurl.Data.url,'_blank','noopener,noreferrer')
            //    openInNewTab(gameurl.Data.url)
            // }   
        
        }
      //  fetchData()
        
    },[])

   
    return (
        <ContentLayout title="Games">
        <div>
        <h1>{`${t("menu.games")} No. ${product}`}</h1>
            <GamesList lng={lng} product={product} id={id} />
        </div>
        </ContentLayout>
    )
}