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
import { GetHistory,GetTransaction,Webhoook } from "@/actions"
import useAuthStore from "@/store/auth"
import { useRouter } from "next/navigation"
import { toast } from "@/hooks/use-toast"
import { useTranslation } from "@/app/i18n/client"
import { Button } from "../ui/button"

export function History({lng}:{lng:string}) {
    const [history,setHistory] = useState<any[]>([])
    const [statement,setStatement] = useState<any[]>([])
    const { accessToken,user } = useAuthStore()  
    const router = useRouter()
    const { t } = useTranslation(lng,'translation' ,undefined);
    useEffect(() => {
       if(accessToken){
       GetHistory(accessToken,lng).then((response:any) => {
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

        GetTransaction(accessToken,lng).then((response:any) => {
      
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
        console.log(resp.json())
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
        <TabsTrigger value="history">Statement</TabsTrigger>
        <TabsTrigger value="transaction">Transaction</TabsTrigger>
      </TabsList>
    <TabsContent value="history">
    <Card>
      <CardHeader className="pb-3">
        <CardTitle>{t('bankstatement.history')}</CardTitle>
        <CardDescription>
          {t('bankstatement.historyDescription')}
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-1">
        {history.map((item:any) => (
         
        <div key={item.ID} className="-mx-2 flex items-start space-x-4 rounded-md p-2 transition-all hover:bg-accent hover:text-accent-foreground">
          <BellIcon className="mt-px h-5 w-5" />
          <div className="space-y-1">
            <p className="text-sm font-medium leading-none">{t('bankstatement.transactionDate')}: {item.CreatedAt}</p>
            <p className="text-sm font-medium leading-none">{t('bankstatement.uid')}: <Button  variant="link"  disabled={item.Status !== "waiting"} onClick={() => callwebhook(item)}> {item.Uid} </Button></p>
            <p className="text-sm font-medium leading-none">{t('bankstatement.bankName')}: {item.Bankname}</p>
            <p className="text-sm font-medium leading-none">{t('bankstatement.beforeBalance')}: {item.Beforebalance}</p>
            <p className={cn("text-sm font-medium leading-none", item.Transactionamount>0 ? 'text-green-500' : 'text-red-500')}>{ item.Transactionamount>0?t('bankstatement.deposit'):t('bankstatement.withdraw')}: {item.Transactionamount}</p>
            <p className="text-sm font-medium leading-none">{t('bankstatement.proAmount')}: {item.Proamount}</p>
            <p className="text-sm font-medium leading-none">{t('bankstatement.balance')}: {item.Balance}</p>
          
          </div>
        </div>
        ))}
    
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
        {statement.map((item:any) => (
        <div key={item.ID} className="-mx-2 flex items-start space-x-4 rounded-md p-2 transition-all hover:bg-accent hover:text-accent-foreground">
          <BellIcon className="mt-px h-5 w-5" />
          <div className="space-y-1">
          <p className={cn("text-sm font-medium leading-none text-muted-foreground", item.Status=='100' ? 'text-green-500' : 'text-red-500')}>{t('transaction.Status')}: {item.Status=="100" ? t('transaction.bet') : t('transaction.result')}</p>
            <p className="text-sm font-medium leading-none">{t('transaction.transactionDate')}:{new Date(item.CreatedAt).toLocaleDateString()} {new Date(item.CreatedAt).toLocaleTimeString()}</p>
            <p className="text-sm font-medium leading-none">{t('transaction.gameprovide')}: {item.GameProvide}</p>
            <p className="text-sm font-medium leading-none">{t('transaction.beforeBalance')}: {item.BeforeBalance}</p>
            <p className="text-sm font-medium leading-none">{t('transaction.betamount')}: {item.BetAmount}</p>
            <p className="text-sm font-medium leading-none">{t('transaction.transactionAmount')}: {item.TransactionAmount}</p>
            <p className="text-sm font-medium leading-none">{t('transaction.balance')}: {item.Balance}</p>
            <p className="text-sm font-medium leading-none">{t('transaction.turover')}: {item.BetAmount}</p>
          
          </div>
        </div>
        ))}
    
      </CardContent>
    </Card>
    </TabsContent>
    </Tabs>
  )
}