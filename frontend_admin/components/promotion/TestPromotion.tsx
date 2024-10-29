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
        balance: z.coerce.number()
       ,//z.number().gte(mindeposit, { message: `Minimum ${mindeposit}` }),
    })
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            amount: 0,
            balance: 0
        }as z.infer<typeof formSchema>
      })
      const amount = form.watch("amount");

      React.useEffect(() => {
          const percentDiscount = promotion.getValues("percentDiscount");
          const calculatedBalance = (amount * percentDiscount / 100) + amount;
          form.setValue("balance", calculatedBalance);
      }, [amount, promotion, form]);
     // console.log(promotion)
    return (
      
        
        <div className="p-4 max-w-md mx-auto">
            <h2 className="text-xl font-bold mb-4">{transactionType === 'deposit' ? t('promotion.deposit') : t('promotion.withdrawal')}</h2>
            <Form {...form}> 
            <form   className="space-y-4">
                <div>
                     { `${t('promotion.minDept')} : ${promotion.getValues("minDept")} `  } 
                     {  `${t('promotion.percentDiscount')} :  ${promotion.getValues("percentDiscount")} %`  } 
                </div>
                <div>
                    {/* <label htmlFor="amount" className="block text-sm font-medium">{transactionType === 'deposit'?t('promotion.deposit'):t('promotion.withdrawal')}</label> */}
                    <FormField
            control={form.control}
            name="amount"
            render={({ field }) => (
              <FormItem className='mb-4'>
                <FormLabel>{t('promotion.amount')}</FormLabel>
                        <FormControl>
                        <Input {...field} type="number"  onChange={(e) => field.onChange(Number(e.target.value))}/>
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
                 
                </div>
                
                {/* <Button type="button" className="w-full" onClick={form.handleSubmit(handleSubmit)} >
                    {transactionType === 'deposit' ? t('promotion.deposit') : t('promotion.withdrawal')}
                </Button> */}
            </form>
            </Form>
        </div>
    
       
    );
};

 
