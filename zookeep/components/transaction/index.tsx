'use client'
import React, { useState } from 'react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form"
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '../ui/select'
import { useTranslation } from '@/app/i18n/client';
//import Footer from '../footer';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { AddStatement } from '@/actions';
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
            channel:"1stpay",
            status: "101"
             };
    
            

             if(accessToken){

                const response = await AddStatement(accessToken,formattedValues)
             //console.log(response)
                if(response.Status){

                    toast({
                        title: t("promotion.edit.success"),
                        description: response.Message,
                        variant: "default",
                      })
        
                      router.push(`/${lng}/home`)
                }  else {
                    
                    toast({
                        title: t("promotion.edit.error"),
                        description: response.Message,
                        variant: "destructive",
                      })
        
                }   
            } else {
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
    
    return (
        <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4 p-4 max-w-md mx-auto">
      
            <h2 className="text-xl font-bold mb-4">{t(`${slug}`)}</h2>
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
               
                <FormField 
                    control={form.control}
                    name="transactionType"
                    render={({field})=>{
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
                    }}
                />
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
             
            </form>
            </Form>
        
    );
};

export default TransactionForm;