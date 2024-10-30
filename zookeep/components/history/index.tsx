'use client'
import { BellIcon, EyeNoneIcon, PersonIcon } from "@radix-ui/react-icons"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { useEffect, useState } from "react"
 
import { GetHistory } from "@/actions"
import useAuthStore from "@/store/auth"
import { useRouter } from "next/navigation"
import { toast } from "@/hooks/use-toast"
import { useTranslation } from "@/app/i18n/client"

export function History({lng}:{lng:string}) {
    const [history,setHistory] = useState<any[]>([])
    const { accessToken } = useAuthStore()  
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
       } else {
        router.push(`/${lng}/login`)
       }
    }, [])


  return (
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
            <p className="text-sm font-medium leading-none">{t('bankstatement.bankName')}: {item.Bankname}</p>
            <p className="text-sm font-medium leading-none">{t('bankstatement.beforeBalance')}: {item.Beforebalance}</p>
            <p className="text-sm font-medium leading-none">{t('bankstatement.transactionAmount')}: {item.Transactionamount}</p>
            <p className="text-sm font-medium leading-none">{t('bankstatement.proAmount')}: {item.Proamount}</p>
            <p className="text-sm font-medium leading-none">{t('bankstatement.balance')}: {item.Balance}</p>
          
          </div>
        </div>
        ))}
    
      </CardContent>
    </Card>
  )
}