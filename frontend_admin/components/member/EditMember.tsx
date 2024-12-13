'use client'
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQuery } from '@tanstack/react-query';

import { memberSchema, Member } from "@/lib/zod/member";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
//import { GameStatus } from "@/lib/zod/gameStatus";
import { useTranslation } from "@/app/i18n/client";
import { GetMemberById,AddMember,UpdateMember} from "@/actions";
import { useEffect } from "react";
import { z } from "zod";
import { toast } from "@/hooks/use-toast";
import { formatNumber } from "@/lib/utils";
import { cn } from "@/lib/utils";
import { CalendarIcon } from "lucide-react";
import { th } from "date-fns/locale";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import useAuthStore from "@/store/auth";

// const gametype = [
//     { id: '1', name: 'Slot' },
//     { id: '2', name: 'Live Casino' },
// ];

function EditMember({ memberId, lng, prefix, onClose, onCancel, isAdd }: { memberId: number, lng: string, prefix: string, onClose: () => void, onCancel: () => void, isAdd: boolean }) {
  
  const { t } = useTranslation(lng, 'translation', undefined);
  const { customerCurrency } = useAuthStore();
  // const { data: memberStatus, isLoading: memberStatusLoading } = useQuery({
  //   queryKey: ['memberStatus'],
  //   queryFn: async () => await GetMemberStatus(prefix),
  // });
 

  const form = useForm<Member>({
    resolver: zodResolver(memberSchema),
    defaultValues: {

    },
  });

  const fetchMember = async (prefix:string,id:number) => {
    const data = await GetMemberById(prefix, id);
  // console.log(data)
    form.reset(data.Data as z.infer<typeof memberSchema>);
    form.setValue("RefferedCode",data.Data.referred_by)
  };

  useEffect(() => {
    if (memberId) {
      fetchMember(prefix, memberId);
    }
  }, [memberId, prefix]);

  const handleSubmit = async (data: Member) => {
 
    //console.log(isAdd)

    if (isAdd) {
      // Combine prefix and username when saving
      data.Username = `${prefix}${data.Username}`;
    }
    //console.log(data)
    data.Status = JSON.parse(data.Status?.toString());//JSON.parse(data.Status?.toString() || '{}').name
 
   const result = !isAdd ? await UpdateMember(prefix, memberId, data) : await AddMember(prefix, data)
   if (!isAdd) {
 

    if (result.Status) {
      toast({
        title: t("member.edit.success"),
        description: t("member.edit.success_description"),
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
        title: t("member.add.success"),
        description: t("member.add.success_description"),
        variant: "default",
      })
      onClose();
    } else {
      toast({
        title: t("member.add.error"),
        description: t("member.add.error_description") + result.Message,
        variant: "destructive",
      })
    }
  }
  };
    return (
        <div className="p-6 bg-white rounded-lg shadow-md md:max-w-md">
         
          <h2 className="text-2xl font-bold mb-4">{memberId ? t('member.edit.title') : t('member.add.title')}</h2>
          <p className="text-gray-600 mb-6">{memberId ?t('member.edit.description'):t('member.add.description')}</p>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
            <FormField
                control={form.control}
                name="ReferralCode"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('member.columns.refferalcode')}</FormLabel>
                    <FormControl>
                      <div className="flex">
                       
                        <Input
                          {...field}
                          className={cn(
                            "rounded-r-1 border-r-0 w-[20ch]",
                            "bg-muted" 
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
                    <FormLabel>{t('member.columns.username')}</FormLabel>
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
                    <FormLabel>{t('member.columns.password')}</FormLabel>
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
                    <FormLabel>{t('member.columns.fullname')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
             
                  {/* <FormField
                control={form.control}
                name="dob"
                render={({ field }) => (
                  <FormItem className="flex flex-col">
                    <FormLabel>Date of birth</FormLabel>
                    <Popover>
                      <PopoverTrigger asChild>
                        <FormControl>
                          <Button
                            variant={"outline"}
                            className={cn(
                              "w-[240px] pl-3 text-left font-normal",
                              !field.value && "text-muted-foreground"
                            )}
                          >
                            {field.value ? (
                              format(field.value, "PPP")
                            ) : (
                              <span>Pick a date</span>
                            )}
                            <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                          </Button>
                        </FormControl>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0" align="start">
                        <Calendar
                          mode="single"
                          selected={field.value}
                          onSelect={field.onChange}
                          disabled={(date) =>
                            date > new Date() || date < new Date("1900-01-01")
                          }
                          initialFocus
                          locale={th}
                          weekStartsOn={0}
                          dir="ltr"
                          className="ltr-calendar"
                        />
                      </PopoverContent>
                    </Popover>
                    <FormDescription>
                      Your date of birth is used to calculate your age.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              /> */}
            
              <FormField
                control={form.control}
                name="Bankname"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('member.columns.bankname')}</FormLabel>
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
                  <FormLabel>{t('member.columns.banknumber')}</FormLabel>
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
                <FormLabel>{t('member.columns.balance')}</FormLabel>
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
                name="ProStatus"
                render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('member.columns.prostatus')}</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />  
              <FormField
                control={form.control}
                name="MinTurnoverDef"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('member.columns.minturnover_def')}</FormLabel>
                    <FormControl>
                      <Input {...field} placeholder={"10%"} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              {/* <FormField
                control={form.control}
                name="gameType"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('columns.gameType')}</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              /> */}
             {/* <FormField
                control={form.control}
                name="Role"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('columns.role')}</FormLabel>
                    <FormControl>
                      <div className="space-y-4">
                        <Select
                          onValueChange={(value) => {
                            field.onChange(value); // Just set the value directly
                          }}
                          value={field.value?.toString() || ''}
                        >
                          <FormControl>
                            <SelectTrigger>
                                <SelectValue placeholder={t('selectStatus')} />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                             {gameStatusLoading ? (
                              <div>Loading...</div>
                            ) : (
                              gameStatus.Data.map((item:any  ) => {
                                const status = JSON.parse(item.status)
                              
                                return (
                                <SelectItem key={status.name} value={JSON.stringify(status)}>{t(status.name)}</SelectItem>
                              )
                            })
                            )}
                          </SelectContent>
                        </Select>
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              /> */}
            
           
              <FormField
                control={form.control}
                name="Status"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('member.columns.status')}</FormLabel>
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
              <FormField
                control={form.control}
                name="RefferedCode"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('member.columns.refferedcode')}</FormLabel>
                    <FormControl>
                      <div className="flex">
                       
                        <Input
                          {...field}
                          className={cn(
                            isAdd && "rounded-l-none","rounded-r-none",
                            "flex-1"
                          )}
                          disabled={!isAdd}
                        />
                         
                      </div>
                    </FormControl>
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
                        errorMessage += `\n${t(`member.${key}`)}: ${errors[key]?.message}`;
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
            </form>
          </Form>
        </div>
      );
}

export default EditMember;
