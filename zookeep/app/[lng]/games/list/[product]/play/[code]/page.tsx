'use client'
import { getGameUrl } from "@/actions"
import { useTranslation } from "@/app/i18n/client"
import { ContentLayout } from "@/components/admin-panel/content-layout"
import GamesList from "@/components/games/Gameslist"
import {Options} from "@/components/games/options"
import useAuthStore from "@/store/auth"
import { useEffect } from "react"
//import { Route } from "lucide-react"


export default  function GamePage({ params: { lng,code } }: { params: { lng: string,code:string } }) {

   

    
    const {t} =  useTranslation(lng,"translation",undefined)
    const {accessToken,user,customerCurrency} = useAuthStore()

    useEffect(()=>{

        const fetchData = async ()=>{
            
            
            const  gameurl = await getGameUrl(accessToken || "",code,user.username,customerCurrency || "USD");

           
            const openInNewTab = (url:string) => {
                const newWindow = window.open(url, '_blank', 'noopener,noreferrer')
                if (newWindow) newWindow.opener = null
            }
        
            if(gameurl.Status){
                //console.log(gameurl.Data.url)
             //   window.open(gameurl.Data.url,'_blank','noopener,noreferrer')
               openInNewTab(gameurl.Data.url)
            }   
        
        }
        fetchData()
        
    },[])

   
    return (
        <ContentLayout title="Games">
        <div>
        <h1>{`${t("menu.games")} No. ${code}`}</h1>
            {/* <GamesList lng={lng} id={code}/> */}
            {/* {id == "all" ? (           
             <Options lng={lng} id={id}/>
             ) : (
             <GameList lng={lng} id={id}/>
             )} 
               */}
        </div>
        </ContentLayout>
    )
}