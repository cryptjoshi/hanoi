'use client'
import React, { useState } from 'react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form"
import { Dialog, DialogTrigger, DialogContent, DialogTitle, DialogDescription } from "@/components/ui/dialog"; // นำเข้า Dialog

import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '../ui/select'
import { useTranslation } from '@/app/i18n/client';
import { formatNumber } from '@/lib/utils';
//import Footer from '../footer';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Deposit, GetPromotion, GetUserInfo, Withdraw } from '@/actions';
import type { User } from '@/store/auth';
import useAuthStore from '@/store/auth';
 import { useRouter } from 'next/navigation';
 import { useToast } from "@/hooks/use-toast"

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
    const [loading, setLoading] = React.useState(true);
    const [balance, setBalance] = React.useState(0);
    const [turnover,setTurnOver] =React.useState(0);
    const [user, setUser] = React.useState<any | null>(null);
    const [currency, setCurrency] = React.useState('USD');
    const {t} = useTranslation(lng,"home",undefined)
    const { toast } = useToast()
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            transactionType:"deposit"
        } as z.infer<typeof formSchema>
      })
    const router = useRouter()
    const {accessToken} = useAuthStore()
  //  console.log(accessToken)
   // console.log(lng)
    //if(!accessToken)
    //    router.push(`${lng}/login`)
    const [promotions, setPromotions] = React.useState<any>();
    const [isLoading, setIsLoading] = React.useState<boolean>(false);
    const [isDialogOpen, setIsDialogOpen] = useState(false); // สถานะสำหรับเปิดปิด Dialog
    const [qrCodeLink, setQrCodeLink] = useState<string | null>(null); // สถานะสำหรับเก็บลิงก์ QR Code

    //const [minTurnover, setMinTurnover] = React.useState<number>(0);
    const [bonus, setBonus] = React.useState<number>(0);

    React.useEffect(() => {
     
        
        const fetchBalance = async (accessToken:string) => {
    
          try {
          setLoading(true);
         
          
    
        //  if (userLoginStatus.state) {
            // if(userLoginStatus.state.isLoggedIn && userLoginStatus.state.accessToken) {
                  const user = await GetUserInfo(accessToken);
                
                        if(user.Status){
                          setBalance(user.Data.balance);
                          setUser(user.Data);
                         
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
        const fetchPromotion = async (accessToken: string) => {
     
          //if(userLoginStatus.state.isLoggedIn && userLoginStatus.state.accessToken){
          const promotion = await GetPromotion(accessToken);
          console.log(promotion.Data)
          // if (promotion.Status) {
          //   // กรองโปรโมชั่นที่มี ID ไม่ตรงกับ user.pro_status
          //   let pro_use = promotion.Data.Promotions.find((promo:any) => promo.ID.toString() == promotion.Data.Prostatus)
          //   setPromotions(pro_use)
          //   //console.log(pro_use)
          //   setBonus(promotion.Data.Promotions.find((promo:any) => promo.ID.toString() == promotion.Data.Prostatus)?.minSpendType=="deposit"?0:user?.lastproamount)
          //  // setPromotions(promotion.Data.Promotions);
          //  //console.log(promotion.Data.Promotions.find((promo:any) => promo.ID.toString() == promotion.Data.Prostatus))
            
          // } else {
            
          //   //router.push(`/${lng}/login`)toast({
          // toast({title: t('unsuccess'),
          //   description: promotion.Message,
          //   variant: "destructive",
          // });
             
          // } 
        //} 
       }
    
       const userLoginStatus = JSON.parse(localStorage.getItem('userLoginStatus') || '{}');
        setCurrency(userLoginStatus.state.customerCurrency);
        fetchBalance(userLoginStatus.state.accessToken);
        fetchPromotion(userLoginStatus.state.accessToken);
        setLoading(false);
       
         
        //const lastbonus = promotion.Data.Promotions.find((promo:any) => promo.ID.toString() == user?.pro_status?.toString())?.minSpendType=="deposit"?0:user?.lastproamount
        //const minTurnoverValue = minTurnover * ((user?.lastdeposit)+lastbonus/100)
       
       //setMinTurnover(minTurnover)
        // ถ้า filtered เป็น array ว่าง ให้สร้าง promotion เริ่มต้น

      }, [lng, router]);
   
    const handleSubmit = async (values: z.infer<typeof formSchema>) => {
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

             if(accessToken){

                const response = await (slug === "deposit" ? Deposit(accessToken, formattedValues) : Withdraw(accessToken, formattedValues));
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
            <p className="text-xs sm:text-sm text-muted-foreground">{t('your_balance')}</p>
            <h2 className="text-xl sm:text-2xl font-bold mt-1">{formatNumber(balance)}</h2>
            { slug === "withdraw" && (
              <div>
                <p className="text-xs sm:text-sm text-muted-foreground">{ promotions?.turntype=="turncredit"?t('your_credit'):t('your_turnover')}</p>
                <h2 className="text-xl sm:text-2xl font-bold mt-1">{promotions?.turntype=="turncredit"?formatNumber(user?.pro_balance):formatNumber(turnover || 0)}</h2>
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
            <p className="text-xs sm:text-sm text-muted-foreground">
             
              {  // Display selected promotion name if available
                promotions?.name || t('No Promotion') 
              }   
            </p>
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
               
                {/* <FormField 
                    control={form.control}
                    name="transactionType"
                      render={({field})=>{
                        return (
                        <FormItem>
                            <FormLabel>
                            {t('transactionType')}
                            </FormLabel>
                            <FormControl>
                            <Select
                                    value={transactionType}
                                    onValueChange={setTransactionType}
                                    className="mt-1"
                                >
                                    <SelectTrigger>
                                        <SelectValue placeholder="เลือกประเภท" />
                                    </SelectTrigger>
                                    <SelectContent>
                                        <SelectItem value="deposit">ฝากเงิน</SelectItem>
                                        <SelectItem value="withdraw">ถอนเงิน</SelectItem>
                                    </SelectContent>
                                </Select>
                            </FormControl>
                        </FormItem>
                        )
                    }}
                /> */}
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