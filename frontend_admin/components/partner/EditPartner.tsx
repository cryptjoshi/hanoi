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

function EditPartner({ partnerId, lng, prefix, onClose, onCancel, isAdd }: { partnerId: number, lng: string, prefix: string, onClose: () => void, onCancel: () => void, isAdd: boolean }) {
  
  const { t } = useTranslation(lng, 'translation', undefined);
  const { customerCurrency } = useAuthStore();
  //const [isSeedFetched, setIsSeedFetched] = useState(false);
  const isSeedFetchedRef = useRef(false);
  // const { data: memberStatus, isLoading: memberStatusLoading } = useQuery({
  //   queryKey: ['memberStatus'],
  //   queryFn: async () => await GetMemberStatus(prefix),
  // });
  const {
    register,
    handleSubmit,
    formState: { errors },
    form
  } = useForm<Partner>({
    resolver: zodResolver(partnerSchema),
  });

  // const form = useForm<Partner>({
  //   resolver: zodResolver(partnerSchema),
  //   defaultValues: {},
  // });

  const fetchPartner = async (prefix:string,id:number) => {
    const data = await GetPartnerById(prefix, id);
    form.reset(data.Data as z.infer<typeof partnerSchema>);
  };

  const fetchSeed = async (prefix:string) => {
    const seed = await GetPartnerSeed(prefix)
    form.setValue("RefferalCode",seed.Data.affiliatekey)
  }

  useEffect(() => {
    //console.log(partnerId)
    if (isAdd && !form.getValues("RefferalCode") && !isSeedFetchedRef.current) {
      fetchSeed(prefix);
      isSeedFetchedRef.current = true; //// ตั้งค่าให้เป็น true หลังจากเรียก fetchSeed
  }
  if (partnerId && !isSeedFetchedRef.current) {
      fetchPartner(prefix, partnerId);
    }
  }, [partnerId, prefix]);

  const handlesubmit = async (data: Partner) => {
    console.log(data)
    const errors = form.formState.errors;
    if (Object.keys(errors).length > 0) {
      // แสดงข้อผิดพลาดหรือทำการจัดการตามที่ต้องการ
      toast({
        title: t("edit.error"),
        description: t("edit.error_description"),
        variant: "destructive",
      });
      return; // หยุดการทำงานหากมีข้อผิดพลาด
    }

    if (isAdd) {
      // Combine prefix and username when saving
      data.Username = `${prefix}${data.Username}`;
    }
   // console.log(data)
    data.Status = JSON.parse(data.Status?.toString());//JSON.parse(data.Status?.toString() || '{}').name
 
   const result = !isAdd ? await UpdatePartner(prefix, partnerId, data) : await AddPartner(prefix, data)
   if (!isAdd) {
 

    if (result.Status) {
      toast({
        title: t("partner.edit.success"),
        description: t("partner.edit.success_description"),
        variant: "default",
      })
      onClose();
    } else {
      toast({
        title: t("edit.error"),
        description: t("edit.error_description") + result.Message,
        variant: "destructive",
      })
    }
  } else {
  
    if (result.Status) {
      toast({
        title: t("partner.add.success"),
        description: t("partner.add.success_description"),
        variant: "default",
      })
      onClose();
    } else {
      toast({
        title: t("partner.add.error"),
        description: t("partner.add.error_description") + result.Message,
        variant: "destructive",
      })
    }
  }
  };
    return (
        <div className="p-6 bg-white rounded-lg shadow-md md:max-w-md">
          <h2 className="text-2xl font-bold mb-4">{partnerId ? t('partner.edit.title') : t('partner.add.title')}</h2>
          <p className="text-gray-600 mb-6">{t('partner.edit.description')}</p>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handlesubmit)} className="space-y-4">
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
                <Button type="submit">{t('common.save')}</Button>
                <Button type="button" variant="outline" onClick={onCancel}>{t('common.cancel')}</Button>
              </div>
            </form>
          </Form>
        </div>
      );
}

export default EditPartner;
