'use client'
import React, { useState } from 'react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form"
import { Dialog, DialogTrigger, DialogContent, DialogTitle, DialogDescription } from "@/components/ui/dialog"; // นำเข้า Dialog

//import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '../ui/select'
import { useTranslation } from '@/app/i18n/client';
import { formatNumber } from '@/lib/utils';
//import Footer from '../footer';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Deposit, GetUserPromotion, GetUserInfo, Withdraw } from '@/actions';
//import type { User } from '@/store/auth';
//import useAuthStore from '@/store/auth';
 import { useRouter } from 'next/navigation';
 import { cn } from "@/lib/utils"
 import { useToast } from "@/hooks/use-toast"
import { getSession } from '@/actions';
interface TransProps {
    lng:string
    slug:string
}

// interface Statement {
     
//         userid:string
//         walletid:string
//         uid:string
//         betamount:number
//         transactionamount:number
//         channel:string
//         status: string
 
// }

const formSchema = z.object({
    userid:z.string().optional(),
    walletid:z.string().optional(),
    uid:z.string().optional(),
    betamount:z.coerce.number().optional(),
    transactionamount:z.coerce.number(),
    channel:z.string().optional(),
    status: z.string().optional(),
    transactionType:z.string().optional()
    
})

 

function TransactionForm({lng,slug}:TransProps) {
    //const [amount, setAmount] = useState('');
    const [transactionType, setTransactionType] = useState('deposit'); // 'deposit' or 'withdraw'
    const [loading, setLoading] = React.useState(false);
    const [balance, setBalance] = React.useState(0);
    const [turnover,setTurnOver] =React.useState(0);
    const [user, setUser] = React.useState<any | null>(null);
    const [currency, setCurrency] = React.useState('USD');
    const [isBlinking, setIsBlinking] = useState(false);
    const {t} = useTranslation(lng,"home",undefined)
    const { toast } = useToast()
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            transactionType:"deposit"
        } as z.infer<typeof formSchema>
      })
    const router = useRouter()
    //const {accessToken} = useAuthStore()
  //  console.log(accessToken)
   // console.log(lng)
    //if(!accessToken)
    //    router.push(`${lng}/login`)
    const [promotions, setPromotions] = React.useState<any>();
    const [isLoading, setIsLoading] = React.useState<boolean>(false);
    const [isDialogOpen, setIsDialogOpen] = useState(false); // สถานะสำหรับเปิดปิด Dialog
    const [qrCodeLink, setQrCodeLink] = useState<string | null>(null); // สถานะสำหรับเก็บลิงก์ QR Code
    const transactionAmountWatch = form.watch('transactionamount', 0);
    //const [minTurnover, setMinTurnover] = React.useState<number>(0);
    const [bonus, setBonus] = React.useState<number>(0);

    React.useEffect(() => {
     
        
        const fetchBalance = async () => {
          //const session = await getSession()

          try {
        
         
          
          //setCurrency(session.customerCurrency);
        //  if (userLoginStatus.state) {
            // if(userLoginStatus.state.isLoggedIn && userLoginStatus.state.accessToken) {
                  const user = await GetUserInfo();
                
                        if(user.Status){
                          setBalance(user.Data.balance);
                          setUser(user.Data);
                         setCurrency(user.Data.currency)
                          setTurnOver(user.Data.turnover)
                        //  setPrefix(user.Data.prefix);
                          
                        } else {
                          // Redirect to login page if token is null
                        //router.push(`/${lng}/login`);
                        // console.log(user)
                        return;
                        }
                  // } else {
                  //   router.push(`/${lng}/login`);
                  //   return;
                  //   }
            // } else {
            //   router.push(`/${lng}/login`);
            //   return;
            // }
          
          } catch (error) {
           // router.push(`/${lng}/login`);
           console.log(error)
          }
         
        };
        //const userLoginStatus = JSON.parse(localStorage.getItem('userLoginStatus') || '{}');
        const fetchPromotion = async () => {
          setLoading(true);
         // const session = await getSession()
          //if(userLoginStatus.state.isLoggedIn && userLoginStatus.state.accessToken){
          const promotion = await GetUserPromotion();
         
          if (promotion.Status) {
          //   // กรองโปรโมชั่นที่มี ID ไม่ตรงกับ user.pro_status
          //   let pro_use = promotion.Data.Promotions.find((promo:any) => promo.ID.toString() == promotion.Data.Prostatus)
          //   setPromotions(pro_use)
          //   //console.log(pro_use)
          //   setBonus(promotion.Data.Promotions.find((promo:any) => promo.ID.toString() == promotion.Data.Prostatus)?.minSpendType=="deposit"?0:user?.lastproamount)
          setPromotions(promotion.Data);
          //  //console.log(promotion.Data.Promotions.find((promo:any) => promo.ID.toString() == promotion.Data.Prostatus))
          if(promotion.Data.status!="0")
          setIsBlinking(true)  
         } else {
         console.log(promotion.Message)
            //router.push(`/${lng}/login`)toast({
          // toast({title: t('unsuccess'),
          //   description: promotion.Message,
          //   variant: "destructive",
          // });
             
          } 
        //} 
        setLoading(false);
        
       }
    
        
        
        fetchBalance();
        fetchPromotion();
       
       
         
        //const lastbonus = promotion.Data.Promotions.find((promo:any) => promo.ID.toString() == user?.pro_status?.toString())?.minSpendType=="deposit"?0:user?.lastproamount
        //const minTurnoverValue = minTurnover * ((user?.lastdeposit)+lastbonus/100)
       
       //setMinTurnover(minTurnover)
        // ถ้า filtered เป็น array ว่าง ให้สร้าง promotion เริ่มต้น

      }, [lng, router]);
   
    const handleSubmit = async (values: z.infer<typeof formSchema>) => {
      const session = await getSession()
        //e.preventDefault();
        //console.log(values)
        // const result = await form.trigger();
        // if (!result) {
        //     // If validation fails, show errors in toast
        //     const errors = form.formState.errors;
        //     let errorMessage = t('form.validationError');
        //     Object.keys(errors).forEach((key) => {
        //       // @ts-ignore
        //       errorMessage += `\n${t(`promotion.${key}`)}: ${errors[key]?.message}`;
        //     });
        //     toast({
        //         title: t('form.error'),
        //         description: errorMessage,
        //         variant: "destructive",
        //       })
        //       return; // Stop the submission
        //     }
       try {
       
            const formattedValues = {
            ...values,
           // walletid:"",
           // uid:"",
            turnover:turnover,
            transactionType:slug,
            transactionamount:slug==="deposit"?values.transactionamount:values.transactionamount*(-1),
            channel:"1stpay",
            status: "101"
             };
    
            
            //console.log(formattedValues)

             if(session.isLoggedIn){

                const response = await (slug === "deposit" ? Deposit(formattedValues) : Withdraw(formattedValues));
                //console.log(response)
                if(response.Status  ){
                  //const link = response.Data.link;
                  // toast({
                  //       title: t("promotion.edit.success"),
                  //       description: response.Message,
                  //       variant: "default",
                  //     })
                  if( slug === "deposit"){
                      setQrCodeLink(response.Data.link);
                      setIsDialogOpen(true);
                  } else {
                    toast({
                        title: t("promotion.edit.success"),
                        description: response.Message,
                        variant: "default",
                      })
                      router.push(`/${lng}/home`)
                  }
              
                     // router.push(`/${lng}/home`)
                }  else {
                   // console.log(response)
                    toast({
                        title: t("promotion.edit.error"),
                        description: response.Message,
                        variant: "destructive",
                      })
        
                }   
            } else {
               // 
               toast({
                title: t('form.error'),
                description: "User not logged",
                variant: "destructive",
              })
              router.push(`/${lng}/login`)
             }
        
            }
            catch (error:any){
               console.log(error)
            toast({
                title: t('form.error'),
                description: error.Message,
                variant: "destructive",
              })
            }
    };
    
     return loading ? <div>Loading...</div> : (
 
        <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4 p-4 max-w-md mx-auto">
       
            <h2 className="text-xl font-bold mb-4">{t(`${slug}`)}</h2>
            
            <div>
           
            <p className="text-xs sm:text-sm text-muted-foreground">{promotions?.turntype!="turncredit"?t('your_balance'):t('your_credit')}</p>
            <h2 className="text-xl sm:text-2xl font-bold mt-1">{formatNumber(balance)}</h2>
         
        
            { slug === "withdraw" && (
              <div>
                <p className="text-xs sm:text-sm text-muted-foreground">{ promotions?.turntype!="turncredit"?t('your_turnover'):""}</p>
                <h2 className="text-xl sm:text-2xl font-bold mt-1">{promotions?.turntype!="turncredit"?formatNumber(user?.turnover):""}</h2>
            </div>
            ) }
            {promotions?.turntype=="turnover"? (
            <p className="text-xs sm:text-sm text-muted-foreground mt-1"> 
           {  promotions ? `≈ Min Turnover ${(parseFloat(user?.lastdeposit)+parseFloat(promotions?.minSpendType=="deposit"?0:user?.lastproamount))*promotions?.minSpend}  ${currency}`:''}</p>
          ):(
              <p className="text-xs sm:text-sm text-muted-foreground mt-1"> 
            
              {promotions ? 
               promotions.minCreditType === "deposit" 
               ?
                `   Min Credit ${user?.lastdeposit} x ${promotions.MinCredit} ≈ ${
                  parseFloat(user?.lastdeposit) * (promotions.MinCredit?.toString().includes("%") 
                      ? (100 + parseFloat(promotions.MinCredit.toString().replace("%",""))) / 100 
                      : parseFloat(promotions.MinCredit || 0)) } ${currency}
                `
                : 
                `  Min Credit ${(parseFloat(user?.lastdeposit)).toFixed(2)} + ${(parseFloat(user?.lastproamount)).toFixed(2)} x ${promotions.MinCredit} ≈ ${ ((parseFloat(user?.lastdeposit) + parseFloat(user?.lastproamount)) * (promotions.MinCredit?.toString().includes("%") 
                      ? (100 + parseFloat(promotions.MinCredit.toString().replace("%",""))) / 100
                      : parseFloat(promotions.MinCredit || 0))).toFixed(2) } ${currency}
                `
                : ''
              }
              
            </p>

              // <p className="text-xs sm:text-sm text-muted-foreground mt-1"> {  promotions ? `≈ Min Credit ${parseFloat(promotions?.MinCredit)}  ${currency}`:''}</p>
            )}
            </div>
            <div className="mt-2">
            <p className="text-xs sm:text-sm font-semibold">{t('promotionStatus')}:</p>
          
            <div className="flex items-center gap-2">
            <div
                className={cn(
                  "h-3 w-3 rounded-full transition-all duration-300",
                   promotions?.status != "0" ? "bg-green-500" : "bg-red-500",
                  isBlinking && "animate-pulse"
                )}
              />
               <span className={cn(
                "transition-opacity",
                promotions?.status != "0" ? "text-green-500" : "text-red-500",
                 isBlinking && "animate-pulse"
              )}> 
                {  // Display selected promotion name if available
                  promotions?.Name || t('No Promotion') 
                }   
                </span>
              </div>
           
            </div>
            <FormField
                    control={form.control}
                    name="transactionamount"
                    render={({ field }) => (
                    <FormItem>
                        <FormLabel>{t(`${slug}`)}</FormLabel>
                        <FormControl>
                        <Input {...field} />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                    )}
                />
             
            <p>{ slug=="deposit" && transactionAmountWatch?`Result ≈ ${ (eval(promotions?.Example?.replace("deposit", isNaN(transactionAmountWatch)?0:transactionAmountWatch)))?.toFixed(2)}`:""}</p>
                    
            <Button type="submit" onClick={async () => {
              const result = await form.trigger();
              if (!result) {
                const errors = form.formState.errors;
                let errorMessage = t('form.validationError');
                Object.keys(errors).forEach((key) => {
                  // @ts-ignore
                  errorMessage += `\n${t(`promotion.${key}`)}: ${errors[key]?.message}`;
                });

                toast({
                  title: t('form.error'),
                  description: errorMessage,
                  variant: "destructive",
                })
              }
            }}>{t(`${slug}`)}</Button>

            <Dialog open={isDialogOpen} onOpenChange={(open) => {
                setIsDialogOpen(open);
                if (!open) {
                    router.push(`/${lng}/home`); // เปลี่ยนเส้นทางเมื่อ Dialog ถูกปิด
                }}}>
            <DialogContent 
                    closeButton={false}
                    >
                <DialogDescription>
                {qrCodeLink && (
                    <iframe 
                    src={qrCodeLink}  
                    style={{ width: '100%', height: '100vh',border:'none' }}
                    ></iframe>
                )}
              </DialogDescription>
              </DialogContent>
            </Dialog>
             
            </form>
            </Form>
         
         
     
    );
};

export default TransactionForm;