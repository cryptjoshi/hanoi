'use client'
import { BellIcon, EyeNoneIcon, PersonIcon } from "@radix-ui/react-icons"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useEffect, useState } from "react"
import { cn } from "@/lib/utils"
import { GetHistory,GetTransaction,Webhoook,GetPromotionLog,GetPromotion } from "@/actions"
import useAuthStore from "@/store/auth"
import { useRouter } from "next/navigation"
import { toast } from "@/hooks/use-toast"
import { useTranslation } from "@/app/i18n/client"
import { Button } from "../ui/button"
import { TransactionTable } from "./transaction"
import { HistoryTable } from "./statement"
import { HistoryPromotion } from "./promotionlog"

export function History({lng}:{lng:string}) {
    const [history,setHistory] = useState<any[]>([])
    const [statement,setStatement] = useState<any[]>([])
    const [promotionlog,setPromotionlog] = useState<any[]>([])
    const { accessToken,user } = useAuthStore()  
    const router = useRouter()
    const { t } = useTranslation(lng,'translation' ,undefined);
    useEffect(() => {
       if(accessToken){
       GetHistory(accessToken,user.prefix).then((response:any) => {
       if(response.Status){
            setHistory(response.Data)
            
        } else {
                toast({
                    title: t('common.unsuccess'),
                    description: response.Message,
                    variant: "destructive",
                })
        }
        })

        GetTransaction(accessToken,user.prefix).then((response:any) => {
      
          if(response.Status){
              setStatement(response.Data)
              
          } else {
                  toast({
                      title: t('common.unsuccess'),
                      description: response.Message,
                      variant: "destructive",
                  })
          }
        })

        GetPromotion(accessToken).then((response:any) => {
        //console.log(response)
          if(response.Status){
               console.log(response.Data)
               setPromotionlog(response.Data)
               
           } else {
                  //  toast({
                  //      title: t('common.unsuccess'),
                  //      description: response.Message,
                  //      variant: "destructive",
                  //  })
           }
           })
       } else {
        router.push(`/${lng}/login`)
       }
    }, [])

    const callwebhook = (item:any) => {
     // console.log(uid,prefix,method)
   try {
    //console.log(item)
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_ENDPOINT}:4006/api/v1/statement/webhook`, { method: 'POST',
      headers: {   
        'Accept': 'application/json',
        'Content-Type': 'application/json',
        },
        body:
        JSON.stringify({
              "TransactionID":item.Uid,
              "isExpired":0,
              "verify":1,
              "ref":user.Username,
              "merchantID":item.Prefix,
              "type":"payin" /* payin,payout */ 
          })
      }).then(resp=>{
        //console.log(resp.json())
        router.push(`/${lng}/home`)
      })
    //Webhoook(item.Uid,user.username,"0","1",item.Prefix,item.StatementType=="Deposit"?"payin":"payout").then((resp)=>console.log(resp))
    //  if (response && typeof response === 'string') {
    //      const jsonResponse = JSON.parse(response); // แปลงเป็น JSON
    //      console.log(jsonResponse);
    //  } else {
    //      console.log("Response is not a valid JSON string");
    //  }
   } catch (error) {
     console.log(error);
   }
    }


  return (
    <Tabs defaultValue="history">
      <TabsList>
      <TabsTrigger value="promotion">Promotion</TabsTrigger>
        <TabsTrigger value="history">Statement</TabsTrigger>
        <TabsTrigger value="transaction">Transaction</TabsTrigger>
        
      </TabsList>
      <TabsContent value="promotion">
      <Card>
      <CardHeader className="pb-3">
        <CardTitle>{t('bankstatement.history')}</CardTitle>
        <CardDescription>
          {t('bankstatement.historyDescription')}
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-1">
          <HistoryPromotion history={promotionlog}  />
      </CardContent>
    </Card>
      </TabsContent>
    <TabsContent value="history">
    <Card>
      <CardHeader className="pb-3">
        <CardTitle>{t('bankstatement.history')}</CardTitle>
        <CardDescription>
          {t('bankstatement.historyDescription')}
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-1">
          <HistoryTable history={history} callwebhook={callwebhook} />
      </CardContent>
    </Card>
    </TabsContent>
    <TabsContent value="transaction">
    <Card>
      <CardHeader className="pb-3">
        <CardTitle>{t('transaction.transaction')}</CardTitle>
        <CardDescription>
          {t('transaction.transactionDescription')}
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-1">
      <TransactionTable statement={statement} />
      </CardContent>
    </Card>
    </TabsContent>
    </Tabs>
  )
}