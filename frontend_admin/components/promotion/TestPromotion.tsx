'use client'
import React, { useState } from 'react';
import * as z from "zod"
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '../ui/select'
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form"
import { useTranslation } from '@/app/i18n/client';
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"

// interface Promotion {
//     id: string;
//     ID: string;
//     name: string;
//     description: string;
//     percentDiscount: number;
//     startDate: string;
//     endDate: string;
//     maxDiscount: number;
//     usageLimit: number;
//     specificTime: string;
//     paymentMethod: string;
//     minSpend: number;
//     maxSpend: number;
//     termsAndConditions: string;
//     example: string;
//     status: string;
//     includegames: string;
//     excludegames: string;
//   }
  
interface TransProps {
    lng:string,
    promotion:any
}

export default function TestPromotion({lng,promotion}:TransProps) {
    //const [amount, setAmount] = useState(0);
    const [transactionType, setTransactionType] = useState('deposit'); // 'deposit' or 'withdraw'

    const {t} =   useTranslation(lng, 'translation', undefined)
    const handleSubmit = async (values: z.infer<typeof formSchema>) => {
        console.log(`Transaction Type: ${transactionType}, Amount: ${values.amount}`);
        // Set the balance value
        //form.setValue('balance', values.amount);
    };
 //   console.log(promotion.getValues("minDept"))
    const mindeposit = parseInt(promotion.getValues("minDept"))
    const formSchema = z.object({
        amount:    z.coerce.number().gte(100,{message:`Minimum ${mindeposit}`}),
        balance: z.coerce.number(),
        mincredit:z.coerce.number(),
        minturnover: z.coerce.number(),
        maxwithdrawal: z.coerce.number()
       ,//z.number().gte(mindeposit, { message: `Minimum ${mindeposit}` }),
    })
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            amount: 0,
            balance: 0,
            mincredit:0,
            minturnover:0,
            maxwithdrawal: 0,
        }as z.infer<typeof formSchema>
      })
      const amount = form.watch("amount");

      React.useEffect(() => {
          const percentDiscount = promotion.getValues("percentDiscount");
          const maxDiscount = promotion.getValues("maxDiscount")
          const maxwithdrawal = promotion.getValues("maxSpend")
          const minCredit = promotion.getValues("MinCredit")

          let calculatedBalance = (amount * percentDiscount / 100) + amount;
          if(calculatedBalance>maxDiscount)
          calculatedBalance = (1*amount)+(1*maxDiscount)
          
          form.setValue("balance", calculatedBalance);

        let minturn  = 0
        let percentage = 0
        if(promotion.getValues("turntype")=="turnover"){
                if(promotion.getValues("minSpend").indexOf("%")>-1)
                {
                    percentage = parseFloat(promotion.getValues("minSpend").replace("%",""))

                    minturn = calculatedBalance*((100+percentage)/100)
                }else {
                    minturn = promotion.getValues("minSpend")*calculatedBalance
                }
            //minturn = form.getValue("balance")*percentDiscount

            form.setValue("minturnover",minturn)
            } else {
                if(promotion.getValues("MinCredit").indexOf("%")>-1)
                {
                    percentage = parseFloat(promotion.getValues("MinCredit").replace("%",""))
        
                    minturn = calculatedBalance*((100+percentage)/100)
                }else {
                    minturn = promotion.getValues("MinCredit")*amount
                }
                //minturn = form.getValue("balance")*percentDiscount
                form.setValue("mincredit",minturn)
            }
           // console.log("MinTurn:"+minturn)
        
        form.setValue("maxwithdrawal",maxwithdrawal)


      }, [amount, promotion, form]);
      console.log(promotion)
    return (
      
        
        <div className="p-4 max-w-md mx-auto">
            <h2 className="text-xl font-bold mb-4">{transactionType === 'deposit' ? t('promotion.deposit') : t('promotion.withdrawal')}</h2>
            <Form {...form}> 
            <form   className="space-y-4">
                <div>
                    {
                    `   Min Credit ${amount|| 0} x ${promotion.getValues("MinCredit")} ≈ ${
                    parseFloat(amount || 0) * (promotion.getValues("MinCredit")?.toString().includes("%") 
                        ? (100 + parseFloat(promotion.getValues("MinCredit").toString().replace("%",""))) / 100 
                        : parseFloat(promotion.getValues("MinCredit") || 0)) } 
                    `
                    }
                    <p>{ `${t('promotion.minDept')} : ${promotion.getValues("minDept")} `  } </p>
                    <p> {  `${t('promotion.percentDiscount')} :  ${promotion.getValues("percentDiscount")} %`  } </p>
                    <p> {  `${t('promotion.maxDiscount')} :  ${promotion.getValues("maxDiscount")} `  } </p>
                    { promotion.getValues("turntype")=="turnover"?<>
                    <p> {  `${t('promotion.minSpend')} :  ${promotion.getValues("minSpend").indexOf("%")>-1?promotion.getValues("minSpend"):"x "+promotion.getValues("minSpend")} of ${promotion.getValues("minSpendType").indexOf("_")>-1?promotion.getValues("minSpendType").replace("_","+"):promotion.getValues("minSpendType")}`   }</p>
                    <p> {  `${t('promotion.maxSpend')} :  ${promotion.getValues("maxSpend")} `   }</p>
                    </>
                    :
                    <p> {  `${t('promotion.turncredit')} :  ${promotion.getValues("MinCredit") } `}  </p>
                    }
                </div>
                <div>
                    {/* <label htmlFor="amount" className="block text-sm font-medium">{transactionType === 'deposit'?t('promotion.deposit'):t('promotion.withdrawal')}</label> */}
                    <FormField
            control={form.control}
            name="amount"
            render={({ field }) => (
              <FormItem className='mb-4'>
                <FormLabel>{t('promotion.deposit')}</FormLabel>
                        <FormControl>
                        <Input {...field} type="text"  onChange={(e) => field.onChange(Number(e.target.value))}/>
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                    )}
                />
                 <FormField
            control={form.control}
            name="balance"
            render={({ field }) => (
              <FormItem className='mb-4'>
                <FormLabel>{t('promotion.balance')}</FormLabel>
                        <FormControl>
                        <Input {...field} type="number" readOnly />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                    )}
                />
                  { promotion.getValues("turntype")=="turnover"? 
                 <FormField
            control={form.control}
            name="minturnover"
            render={({ field }) => (
              <FormItem className='mb-4'>
                <FormLabel>{t('promotion.minSpend')}</FormLabel>
                        <FormControl>
                        <Input {...field} type="number" readOnly />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                    )}
                />
                :
                <FormField
                control={form.control}
                name="mincredit"
                render={({ field }) => (
                  <FormItem className='mb-4'>
                    <FormLabel>{t('promotion.turncredit')}</FormLabel>
                            <FormControl>
                            <Input {...field} type="number" readOnly />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                        )}
                    />
            }
                <FormField
            control={form.control}
            name="maxwithdrawal"
            render={({ field }) => (
              <FormItem className='mb-4'>
                <FormLabel>{t('promotion.maxwithdrawal')}</FormLabel>
                        <FormControl>
                        <Input {...field} type="number" readOnly />
                        </FormControl>
                        <FormMessage />
                    </FormItem>
                    )}
                />
                </div>
                {/* <Button type="button" className="w-full" onClick={form.handleSubmit(handleSubmit)} >
                    {transactionType === 'deposit' ? t('promotion.deposit') : t('promotion.withdrawal')}
                </Button> */}
            </form>
            </Form>
        </div>
    
       
    );
};

 
