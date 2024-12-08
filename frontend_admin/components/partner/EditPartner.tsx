'use client'
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
//import { useQuery } from '@tanstack/react-query';

import { partnerSchema, Partner } from "@/lib/zod/partner";
import { Input } from "@/components/ui/input";
//import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
//import { GameStatus } from "@/lib/zod/gameStatus";
import { useTranslation } from "@/app/i18n/client";
import { AddPartner,UpdatePartner,GetPartnerById,GetPartnerSeed} from "@/actions";
import { useEffect,useState,useRef } from "react";
import { z } from "zod";
import { toast } from "@/hooks/use-toast";
import { formatNumber } from "@/lib/utils";
import { cn } from "@/lib/utils";
// import { CalendarIcon } from "lucide-react";
// import { th } from "date-fns/locale";
// import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
// import { Calendar } from "@/components/ui/calendar";
import useAuthStore from "@/store/auth";

// const gametype = [
//     { id: '1', name: 'Slot' },
//     { id: '2', name: 'Live Casino' },
// ];

const formSchema = z.object({
  ID:z.number().optional(),     
  Username:z.string(),    
  Password:z.string(),      
  Fullname:z.string(),    
  Bankname:z.string(),    
  Banknumber:z.string(),    
  Status:z.string(),  
  ProStatus:z.string().optional(),   
  RefferalCode:z.string(), 
  RefferedCode:z.string().optional(),
})


function EditPartner({ partnerId, lng, prefix, onClose, onCancel, isAdd }: { partnerId: number, lng: string, prefix: string, onClose: () => void, onCancel: () => void, isAdd: boolean }) {
 
  const { t } = useTranslation(lng, 'translation', undefined);
  const { customerCurrency } = useAuthStore();
 
  const isSeedFetchedRef = useRef(false);
   

  const form = useForm<Partner>({
    resolver: zodResolver(formSchema),
    defaultValues: {} as z.infer<typeof formSchema>
  });

  const fetchPartner = async (prefix:string,id:any) => {
   
    
      try {
      //  console.log("prefix:",prefix,",id:",id)
    const data = await GetPartnerById(prefix, id);
   // console.log(data)
    if(data.Status){
      const formattedData = {
        ...data.Data,
        RefferalCode:data.Data.affiliatekey,
        Username:data.Data.username,
        Fullname:data.Data.name,
        Password:data.Data.password,
        Bankname:data.Data.bankname,
        Banknumber:data.Data.banknumber,
        Balance:data.Data.balance,
        Status:data.Data.status

      }
      const { ID, ...formData } = formattedData;
      form.reset(formData as z.infer<typeof formSchema>);
    } else {
      toast({
        title: t("partner.fetch.error"),
        description: data.Message,
        variant: "destructive",
      })
    }
    } catch (error) {
      //console.error("Error fetching promotion:", error);
      toast({
        title: t("partner.fetch.error"),
        description: t("partner.fetch.error_description"),
        variant: "destructive",
      })
    }
 
  };

  const fetchSeed = async (prefix:string) => {
    const seed = await GetPartnerSeed(prefix)
    form.setValue("RefferalCode",seed.Data.affiliatekey)
  }



  useEffect(() => {
    if (isAdd && !form.getValues("RefferalCode") && !isSeedFetchedRef.current) {
      fetchSeed(prefix);
      isSeedFetchedRef.current = true;
      //// ตั้งค่าให้เป็น true หลังจากเรียก fetchSeed
  }
  if (partnerId) {
    
     
      fetchPartner(prefix, partnerId);
      
    }
   
  }, [partnerId, prefix]);

  const handleSubmit = async (data: z.infer<typeof formSchema>) => {
     
    // const result = await form.trigger();
    // if (!result) {
    //   // If validation fails, show errors in toast
    //   const errors = form.formState.errors;
    //   let errorMessage = t('form.validationError');
    //   Object.keys(errors).forEach((key) => {
    //     // @ts-ignore
    //     errorMessage += `\n${t(`promotion.${key}`)}: ${errors[key]?.message}`;
    //   });

    //   toast({
    //     title: t('form.error'),
    //     description: errorMessage,
    //     variant: "destructive",
    //   })
    //   return; // Stop the submission
    // } else {
    //   toast({
    //     title: t("partner.add.success"),
    //     description: result,
    //     variant: "default",
    //   })
    // }
    //data.Status = JSON.parse(data.Status?.toString());
    if (isAdd) {
      // Combine prefix and username when saving
      data.Username = `${prefix}${data.Username}`;
    }  
   
   //console.log("isAdd:",isAdd)
   //data.Status = JSON.parse(data.Status?.toString());//JSON.parse(data.Status?.toString() || '{}').name
   
   const formattedValues = {
    //...data,
    affiliateKey:data.RefferalCode.toString(),
    username:data.Username.toString(),    
    password:data.Password.toString(),      
    name:data.Fullname.toString(),    
    bankname:data.Bankname.toString(),    
    banknumber:data.Banknumber.toString(),    
    status:data.Status.toString()
  };
 // console.log("format values:"+JSON.stringify(formattedValues))
   
  if (partnerId) {
    const data = await UpdatePartner(prefix, partnerId, formattedValues)
    if (data.Status) {
      toast({
        title: t("partner.edit.success"),
        description: t("partner.edit.success_description"),
        variant: "default",
      })
     // queryClient.invalidateQueries({ queryKey: ['promotions'] });
      onClose();
    } else {
      toast({
        title: t("promotion.edit.error"),
        description: t("promotion.edit.error_description") + data.Message,
        variant: "destructive",
      })
    }
  } else {
    const data = await AddPartner(prefix, formattedValues)
    if (data.Status) {
      toast({
        title: t("promotion.add.success"),
        description: data.Message,
        variant: "default",
      })
     // queryClient.invalidateQueries({ queryKey: ['promotions'] });
      onClose();
    } else {
      toast({
        title: t("promotion.add.error"),
        description: t("promotion.add.error_description") + data.Message,
        variant: "destructive",
      })
    }
  }
};
   
   
    
    return (
      <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
        <div className="p-6 bg-white rounded-lg shadow-md md:max-w-md">
          <h2 className="text-2xl font-bold mb-4">{partnerId ? t('partner.edit.title') : t('partner.add.title')}</h2>
          <p className="text-gray-600 mb-6">{t('partner.edit.description')}</p>
        
            <FormField
                control={form.control}
                name="RefferalCode"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('partner.columns.refferalcode')}</FormLabel>
                    <FormControl>
                      <div className="flex">
                       
                        <Input
                          {...field}
                          className={cn(
                            isAdd && "rounded-l-none","rounded-r-none",
                            "flex-1"
                          )}
                          readOnly
                        />
                         
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}  
              />
   
            <FormField
                control={form.control}
                name="Username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('partner.columns.username')}</FormLabel>
                    <FormControl>
                      <div className="flex">
                        <Input
                          value={prefix}
                          readOnly
                          className={cn(
                            "rounded-r-none border-r-0 w-[5ch]",
                            isAdd ? "bg-muted" : "hidden"
                          )}
                        />
                        <Input
                          {...field}
                          className={cn(
                            isAdd && "rounded-l-none","rounded-r-none",
                            "flex-1"
                          )}
                          disabled={!isAdd}
                        />
                          <Input
                          value={customerCurrency?.toLowerCase()}
                          readOnly
                          className={cn(
                            "rounded-l-none border-l-0 w-[7ch]",
                            isAdd ? "bg-muted" : "hidden"
                          )}
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}  
              />
              <FormField
                control={form.control}
                name="Password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('partner.columns.password')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="Fullname"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('partner.columns.fullname')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
             
               
            
              <FormField
                control={form.control}
                name="Bankname"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('partner.columns.bankname')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="Banknumber"
                render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('partner.columns.banknumber')}</FormLabel>
                  <FormControl>
                    <Input {...field}/>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
            control={form.control}
            name="Balance"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('partner.columns.balance')}</FormLabel>
                <FormControl>
                  <Input
                    {...field}
                    value={formatNumber(parseFloat(field.value?.toString() || '0'), 2)}
                    readOnly
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
              
           
              <FormField
                control={form.control}
                name="Status"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('partner.columns.status')}</FormLabel>
                    <Select
                onValueChange={(value) => field.onChange(value ? parseInt(value) : null)}
                      value={field.value?.toString() || ''}
                    >
                    <FormControl>
                    <SelectTrigger>
                    <SelectValue placeholder={t('common.selectStatus')} />
                    </SelectTrigger>
                  </FormControl>  
                  <SelectContent>
                    <SelectItem value="1">{t('common.active')}</SelectItem>
                  <SelectItem value="0">{t('common.inactive')}</SelectItem>
                </SelectContent>
              </Select>
            
              <FormMessage />
            </FormItem>
                )}
              />
               
              <div className="flex justify-end space-x-2 mt-6">
              <Button type="submit" onClick={async () => {
              const result = await form.trigger();
              if (!result) {
                const errors = form.formState.errors;
                let errorMessage = t('form.validationError');
                Object.keys(errors).forEach((key) => {
                  // @ts-ignore
                  errorMessage += `\n${t(`partner.${key}`)}: ${errors[key]?.message}`;
                });

                toast({
                  title: t('form.error'),
                  description: errorMessage,
                  variant: "destructive",
                })
              }
            }}>{t('common.save')}</Button>
                <Button type="button" variant="outline" onClick={onCancel}>{t('common.cancel')}</Button>
              </div>
           
        </div>
        </form>
          </Form>
      );
}

export default EditPartner;
