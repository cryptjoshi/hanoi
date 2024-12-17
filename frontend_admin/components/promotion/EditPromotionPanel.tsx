import React, { useEffect, useState } from 'react';
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { AddPromotion, GetPromotionById, UpdatePromotion,GetGameStatus } from '@/actions';
import { useTranslation } from '@/app/i18n/client';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form"
import { useForm } from "react-hook-form"
import TestPromotion, {TransactionForm} from "./TestPromotion" 
import { toast } from "@/hooks/use-toast"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { cn } from "@/lib/utils"
import { format, parse,isValid } from "date-fns"
import { th } from "date-fns/locale"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { CalendarIcon } from 'lucide-react';
import { Calendar } from '@/components/ui/calendar';
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useQuery } from '@tanstack/react-query';
import { useToast } from "@/hooks/use-toast"
import { useQueryClient } from '@tanstack/react-query';

interface EditPromotionPanelProps {
  promotionId: number | null;
  lng: string;
  prefix: string;
  onClose: () => void;
  onCancel: () => void;
}
interface SpecificTime {
  type?: string;
  daysOfWeek?: string[];
  hour?: string;
  minute?: string;
}
interface Formular {
  amount: number;
}
// Update the Promotion interface
// interface Promotion {
//   id: string;
//   ID: string;
//   name: string;
//   description: string;
//   percentDiscount: number;
//   startDate: string;
//   endDate: string;
//   maxDiscount: number;
//   usageLimit: number;
//   specificTime: string;
//   paymentMethod: string;
//   minSpend: number;
//   maxSpend: number;
//   termsAndConditions: string;
//   example: string;
//   status: string;
//   includegames: string;
//   excludegames: string;
// }

// Update the form schema
const formSchema = z.object({
  name: z.string().min(1, { message: "Name is required" }),
  description: z.string(),
  percentDiscount: z.coerce.number(),
  unit: z.string(),
  startDate: z.string().refine((val) => isValid(parse(val, 'dd-MM-yyyy', new Date())), {
    message: "Invalid date format"
  }),
  endDate: z.string().refine((val) => isValid(parse(val, 'dd-MM-yyyy', new Date())), {
    message: "Invalid date format"
  }),
  maxDiscount: z.coerce.number(),
  usageLimit: z.coerce.number(),
  minDept:z.coerce.number(),
  minCredit:z.string().optional(),
  turnType:z.string(),
  minwithdrawal:z.string(),
  minSpend: z.string().optional(),
  minSpendType: z.string().optional(),
  maxSpend: z.coerce.number(),
  termsAndConditions: z.string().optional(),
  status: z.coerce.number(),
  example: z.string().optional(),
  includegames: z.string(),
  excludegames: z.string(),
  specificTime: z.string().optional(),
  paymentMethod: z.string().optional(),
  Zerobalance:z.coerce.number().default(0)
})

 

function cleanJsonString(jsonString: string): SpecificTime {
  if (!jsonString) {
    return {};
  }

  try {
    let cleanJsonString = jsonString.trim().replace(/^["']|["']$/g, '');
    cleanJsonString = cleanJsonString.replace(/\\"/g, '"');
    return JSON.parse(cleanJsonString) as SpecificTime;
  } catch { 
    return {}
  }
}

// Helper function to check if a string is a valid date
const isValidDateString = (dateString: string) => {
  return !isNaN(Date.parse(dateString));
};

export const EditPromotionPanel: React.FC<EditPromotionPanelProps> = ({ promotionId, prefix, lng, onClose, onCancel }) => {
  const { toast } = useToast()
  const {t} = useTranslation(lng, 'translation', undefined)
  const queryClient = useQueryClient();
  const { data: gameTypes, isLoading: gameStatusLoading } = useQuery({
    queryKey: ['gameTypes'],
    queryFn: async () => await GetGameStatus(prefix),
  });
  

  // Update the schema with translated messages
  const updatedFormSchema = formSchema.extend({
    startDate: z.string().refine((val) => isValid(parse(val, 'dd-MM-yyyy', new Date())), {
      message: t('promotion.edit.invalid_date_format')
    }),
    endDate: z.string().refine((val) => isValid(parse(val, 'dd-MM-yyyy', new Date())), {
      message: t('promotion.edit.invalid_date_format')
    }),
  });

  const form = useForm<z.infer<typeof updatedFormSchema>>({
    resolver: zodResolver(updatedFormSchema),
    defaultValues: {} as z.infer<typeof updatedFormSchema>
  })

  const fetchPromotion = async () => {
    if (promotionId) {
      try {
        const data = await GetPromotionById(prefix, promotionId);
        if (data.Status) {
         // console.log(data.Data)
          const formattedData = {
            ...data.Data,
            percentDiscount: Number(data.Data.percentDiscount),
            unit: data.Data.unit,
            maxDiscount: data.Data.max_discount,
            usageLimit: Number(data.Data.usageLimit),
            minDept: Number(data.Data.minDept),
            minSpend: data.Data.minSpend,
            minwithdrawal:data.Data.Widthdrawmin,
            maxSpend: Number(data.Data.maxSpend),
            minSpendType: data.Data.minSpendType,
            minCredit: data.Data.MinCredit,
            turnType: data.Data.turntype,
            example: data.Data.example,
            startDate: data.Data.startDate ? format(new Date(data.Data.startDate), 'dd-MM-yyyy') : '',
            endDate: data.Data.endDate ? format(new Date(data.Data.endDate), 'dd-MM-yyyy') : '',
            Zerobalance: data.Data.zerobalance
          };


          // Remove ID from formattedData before setting form values
          const { ID, ...formData } = formattedData;
          form.reset(formData as z.infer<typeof updatedFormSchema>);
        } else {
          toast({
            title: t("promotion.fetch.error"),
            description: data.Message,
            variant: "destructive",
          })
        }
      } catch (error) {
        console.error("Error fetching promotion:", error);
        toast({
          title: t("promotion.fetch.error"),
          description: t("promotion.fetch.error_description"),
          variant: "destructive",
        })
      }
    }
  };

  useEffect(() => {
    fetchPromotion();
  }, [promotionId]);

  const handleSubmit = async (values: z.infer<typeof updatedFormSchema>) => {
    // Validate the form before submitting
    const result = await form.trigger();
    if (!result) {
      // If validation fails, show errors in toast
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
      return; // Stop the submission
    }

     let example_str
      if (values.unit !== 'percent') {
        //console.log(`deposit+((${percentDiscountValue}/100)*100)`)
          example_str = `deposit+((${values.percentDiscount}/100)*100)`; // Update example in form
      } else {
        //console.log(`(1+(${percentDiscountValue}/100))*deposit`)
         example_str =  `(1+(${values.percentDiscount}/100))*deposit`; // Update example in form
      }

    const formattedValues = {
      ...values,
      startDate: values.startDate ? format(parse(values.startDate, 'dd-MM-yyyy', new Date()), 'yyyy-MM-dd') : '',
      endDate: values.endDate ? format(parse(values.endDate, 'dd-MM-yyyy', new Date()), 'yyyy-MM-dd') : '',
      percentDiscount: values.percentDiscount.toString(),
      maxDiscount: values.maxDiscount.toString(),
      minwithdrawal:values.minwithdrawal,
      minDept: values.minDept.toString(),
      minSpend: values.minSpend?.toString(),
      minSpendType:values.minSpendType?.toString(),
      maxSpend: values.maxSpend?.toString(),
      example: example_str,
      minCredit: values.minCredit?.toString(),
      turntype: values.turnType.toString(),
      Zerobalance:values.Zerobalance
    };
   // console.log("format values:"+JSON.stringify(formattedValues))
    if (promotionId) {
      const data = await UpdatePromotion(prefix, promotionId, formattedValues);
      if (data.Status) {
        toast({
          title: t("promotion.edit.success"),
          description: t("promotion.edit.success_description"),
          variant: "default",
        })
        queryClient.invalidateQueries({ queryKey: ['promotions'] });
        onClose();
      } else {
        toast({
          title: t("promotion.edit.error"),
          description: t("promotion.edit.error_description") + data.Message,
          variant: "destructive",
        })
      }
    } else {
      const data = await AddPromotion(prefix, formattedValues);
      if (data.Status) {
        toast({
          title: t("promotion.add.success"),
          description: t("promotion.add.success_description"),
          variant: "default",
        })
        queryClient.invalidateQueries({ queryKey: ['promotions'] });
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
      <h2 className="text-2xl font-bold mb-4">{promotionId ? t('promotion.edit.title') : t('promotion.add.title')}</h2>
      <p className="text-gray-600 mb-6">{t('promotion.edit.description')}</p>
      <Tabs defaultValue="promotion">
        <TabsList>
          <TabsTrigger value="promotion">{t('promotion.promotion')}</TabsTrigger>
          <TabsTrigger value="games">{t('promotion.edit.games')}</TabsTrigger>
        </TabsList>
        <TabsContent value="promotion" >
     

          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.name')}</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.description')}</FormLabel>
                <FormControl>
                  <Textarea {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="minDept"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.minDept')}</FormLabel>
                <FormControl>
                  <Input {...field} type="text" onChange={(e) => field.onChange(Number(e.target.value))} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          
          <FormField
            control={form.control}
            name="percentDiscount"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.percentDiscount')}</FormLabel>
                <FormControl>
                  <Input {...field} type="text" onChange={(e) => field.onChange(Number(e.target.value))} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="unit"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.unit')}</FormLabel>
                <FormControl>
                  <div className="space-y-4">
                    <Select
                      onValueChange={(value) => {
                        // Get the current percentDiscount value
                        const percentDiscountValue = form.getValues('percentDiscount');
                        //console.log("PercentDiscountValue:"+percentDiscountValue)
                        if(percentDiscountValue)
                        if (value !== 'percent') {
                          //console.log(`deposit+((${percentDiscountValue}/100)*100)`)
                           form.setValue('example', `deposit+((${percentDiscountValue}/100)*100)`); // Update example in form
                        } else {
                          //console.log(`(1+(${percentDiscountValue}/100))*deposit`)
                          form.setValue('example', `(1+(${percentDiscountValue}/100))*deposit`); // Update example in form
                        }
                        
                        field.onChange(value); // Update the field value with the selected value
                      }}
                      value={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder={t('promotion.selectunit')} />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="percent">{t('promotion.percent')}</SelectItem>
                        <SelectItem value="nonepercent">{t('promotion.nonepercent')}</SelectItem>
                      </SelectContent>
                    </Select>
 
                  </div>
          
            </FormControl>
            <FormMessage />
          </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="startDate"
            render={({ field }) => (
              <FormItem className="flex flex-col mt-2">
                <FormLabel>{t('promotion.startDate')}</FormLabel>
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
                          format(parse(field.value, "dd-MM-yyyy", new Date()), "dd-MM-yyyy")
                        ) : (
                          <span>{t('promotion.selectDate')}</span>
                        )}
                        <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                      </Button>
                    </FormControl>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={field.value ? parse(field.value, "dd-MM-yyyy", new Date()) : undefined}
                      onSelect={(date) => field.onChange(date ? format(date, "dd-MM-yyyy") : '')}
                      initialFocus
                      locale={th}
                      weekStartsOn={0} 
                      dir="ltr"
                      className="ltr-calendar" 
                    />
                  </PopoverContent>
                </Popover>
                <FormMessage />
              </FormItem>
            )}
          />
          
          <FormField
            control={form.control}
            name="endDate"
            render={({ field }) => (
              <FormItem className="flex flex-col mt-2">
                <FormLabel>{t('promotion.endDate')}</FormLabel>
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
                        {field.value && isValid(parse(field.value, "dd-MM-yyyy", new Date())) ? (
                          format(parse(field.value, "dd-MM-yyyy", new Date()), "dd-MM-yyyy")
                        ) : (
                          <span>{t('promotion.selectDate')}</span>
                        )}
                        <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                      </Button>
                    </FormControl>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={field.value ? parse(field.value, "dd-MM-yyyy", new Date()) : undefined}
                      onSelect={(date) => field.onChange(date ? format(date, "dd-MM-yyyy") : '')}
                      initialFocus
                      locale={th}
                      weekStartsOn={0} 
                      dir="ltr"
                      className="ltr-calendar" // 0 represents Sunday
                    />
                  </PopoverContent>
                </Popover>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="maxDiscount"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.maxDiscount')}</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="usageLimit"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.usageLimit')}</FormLabel>
                <FormControl>
                <Input 
                    {...field} 
                    type="number" 
                    readOnly={cleanJsonString(form.getValues('specificTime') || '{}').type === 'first' || cleanJsonString(form.getValues('specificTime') || '{}').type === 'once'} 
                    value={cleanJsonString(form.getValues('specificTime') || '{}').type === 'first' || cleanJsonString(form.getValues('specificTime') || '{}').type === 'once' ? 1 : field.value} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="specificTime"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.specificTime')}</FormLabel>
                <FormControl>
                  <div className="space-y-4">
                    <Select
                      onValueChange={(value) => {
                        const currentValue = cleanJsonString(field.value || '{}');
                        const newValue: SpecificTime = { type: value };
                        if (value !== 'once' || value !== 'first') {
                          newValue.daysOfWeek = currentValue.daysOfWeek;
                        }
                        field.onChange(JSON.stringify(newValue));
                        const specificTimeType = form.watch('specificTime') ? cleanJsonString(form.watch('specificTime')).type : '';
                        if (specificTimeType === 'first' || specificTimeType === 'once') {
                          form.setValue('usageLimit', 1);
                        }
                      }}
                      value={cleanJsonString(field.value || '{}').type}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder={t('promotion.selectTimeType')} />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="first">{t('common.first')}</SelectItem>
                        <SelectItem value="once">{t('common.once')}</SelectItem>
                        <SelectItem value="weekly">{t('common.weekly')}</SelectItem>
                        <SelectItem value="monthly">{t('common.monthly')}</SelectItem>
                      </SelectContent>
                    </Select>
                    
                    {cleanJsonString(field.value || '{}').type === 'weekly' && (
                      <div className="flex flex-wrap gap-2">
                        {['mon', 'tue', 'wed','thu','fri','sat','sun'].map((day) => (
                          <FormField
                            key={day}
                            control={form.control}
                            name="specificTime"
                            render={({ field: innerField }) => (
                              <FormItem className="flex items-center space-x-2">
                                <Checkbox
                                  checked={(cleanJsonString(innerField.value || '{}').daysOfWeek || []).includes(day)}
                                  onCheckedChange={(checked) => {
                                    const currentValue = cleanJsonString(innerField.value || '{}');
                                    const updatedDays = checked
                                      ? [...(currentValue.daysOfWeek || []), day]
                                      : (currentValue.daysOfWeek || []).filter((d: string) => d !== day);
                                    const newValue: SpecificTime = { ...currentValue, daysOfWeek: updatedDays };
                                    innerField.onChange(JSON.stringify(newValue));
                                  }}
                                />
                                <FormLabel className="text-sm font-normal">
                                  {t(`common.${day}`)}
                                </FormLabel>
                              </FormItem>
                            )}
                          />
                        ))}
                      </div>
                    )}

                    <div className="grid grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="hour">{t('common.hour')}</Label>
                        <Select
                          onValueChange={(value) => {
                            const currentValue = cleanJsonString(field.value || '{}');
                            const newValue = { ...currentValue, hour: value };
                            field.onChange(JSON.stringify(newValue));
                          }}
                          value={cleanJsonString(field.value || '{}').hour}
                        >
                          <SelectTrigger id="hour" aria-label="Hour">
                            <SelectValue placeholder={t('common.selectHour')} />
                          </SelectTrigger>
                          <SelectContent>
                            {Array.from({ length: 24 }, (_, i) => i).map((hour) => (
                              <SelectItem key={hour} value={hour.toString()}>{hour < 10 ? `0${hour}` : hour}</SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="minute">{t('common.minute')}</Label>
                        <Select
                          onValueChange={(value) => {
                            const currentValue = cleanJsonString(field.value || '{}');
                            const newValue = { ...currentValue, minute: value };
                            field.onChange(JSON.stringify(newValue));
                          }}
                          value={cleanJsonString(field.value || '{}').minute}
                        >
                          <SelectTrigger id="minute" aria-label="Minute">
                            <SelectValue placeholder={t('common.selectMinute')} />
                          </SelectTrigger>
                          <SelectContent>
                            {Array.from({ length: 60 }, (_, i) => i).map((minute) => (
                              <SelectItem key={minute} value={minute.toString()}>{minute < 10 ? `0${minute}` : minute}</SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      </div>
                    </div>
                  </div>
          
            </FormControl>
            <FormMessage />
          </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="turnType"
            render={({ field }) => (
              <FormItem>
                      <Label htmlFor="hour">{t('promotion.turntype')}</Label>
                      <Select
                      onValueChange={(value) => {
                        field.onChange(value);
                      }}
                      value={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder={t('promotion.select_turntype')} />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="turnover">{t('promotion.turnover')}</SelectItem>
                        <SelectItem value="turncredit">{t('promotion.turncredit')}</SelectItem>
                      </SelectContent>
                    </Select>
                    </FormItem>
            )}
        
            />
            {form.watch("turnType")=="turnover"?
           <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2"> 
              <FormField
            control={form.control}
            name="minSpend"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.minSpend')} </FormLabel>
                <FormControl>
                  <Input {...field} type="text" placeholder={t('promotion.minturn_placeholder')} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
            />
            </div>
            <div className="space-y-2">
            <FormField
            control={form.control}
            name="minSpendType"
            render={({ field }) => (
              <FormItem>
                      <Label htmlFor="hour">{t('promotion.minSpendtype')}</Label>
                      <Select
                      onValueChange={(value) => {
                        field.onChange(value);
                      }}
                      value={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder={t('promotion.selectMinspendType')} />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="deposit_bonus">{t('promotion.deposit_bonus')}</SelectItem>
                        <SelectItem value="deposit">{t('promotion.deposit')}</SelectItem>
                      </SelectContent>
                    </Select>
                    </FormItem>
            )}
        
            />
              </div>
          </div>
          :
         
              <div className="space-y-2"> 
              <FormField
            control={form.control}
            name="minCredit"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.minCredit')} </FormLabel>
                <FormControl>
                  <Input {...field} type="text"  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
            />
          
           
          </div>
          }
           <div className="grid grid-cols-2 gap-4">
           <div className="space-y-2"> 
          <FormField
            control={form.control}
            name="minwithdrawal"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.minwithdrawal')}</FormLabel>
                <FormControl>
                  <Input {...field} type="text" />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
         </div>
         <div className="space-y-2">
          <FormField
            control={form.control}
            name="maxSpend"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.maxSpend')}</FormLabel>
                <FormControl>
                  <Input {...field} type="text" />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          </div>
          </div>
          <FormField
            control={form.control}
            name="Zerobalance"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between rounded-lg border mt-2 p-2">
                <div className="">
                  <FormLabel className="text-base">{t('promotion.zerobalance')}</FormLabel>
                  <FormDescription>
                    {field.value === 1 ? t('promotion.status_active') : t('promotion.status_inactive')}
                  </FormDescription>
                </div>
                <FormControl>
                   
                  <Switch
                    checked={field.value === 1}
                    onCheckedChange={(checked) => {
                      console.log(checked)
                      field.onChange(checked?1:0);
                    }}
                  />
                </FormControl>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="status"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between rounded-lg border mt-2 p-2">
                <div className="">
                  <FormLabel className="text-base">{t('promotion.status')}</FormLabel>
                  <FormDescription>
                    {field.value === 1 ? t('promotion.status_active') : t('promotion.status_inactive')}
                  </FormDescription>
                </div>
                <FormControl>
                  <Switch
                    checked={field.value === 1}
                    onCheckedChange={(checked) => {
                      field.onChange(checked ? 1 : 0);
                    }}
                  />
                </FormControl>
              </FormItem>
            )}
          />
      
      
    </TabsContent>
    <TabsContent value="games">
      <div className="flex flex-col justify-between space-y-4">
        <Button type="button"
          onClick={() => {
            const allGameIds = gameTypes?.Data?.map((item: any) => {
              try {
                return item.status.id.toString();
              } catch {
                return null;
              }
            }).filter(Boolean) || [];
           
            const currentIncludeGames = form.watch('includegames')?.split(',').filter(Boolean) || [];
            const currentExcludeGames = form.watch('excludegames')?.split(',').filter(Boolean) || [];
            
            // Check if all games are currently included
            if (currentIncludeGames.length === allGameIds.length && allGameIds.length > 0) {
              // If all games are currently included, exclude all
              form.setValue('includegames', '');
              form.setValue('excludegames', allGameIds.join(','));
            } else {
              // Otherwise, include all games
              form.setValue('includegames', allGameIds.join(','));
              form.setValue('excludegames', '');
            }
          }}
          variant="outline"
          className="w-full"
          disabled={!gameTypes?.Data?.length}
        >
          {((form.watch('includegames')?.split(',').filter(Boolean).length || 0) === (gameTypes?.Data?.length || 0)) && gameTypes?.Data?.length > 0
            ? t('promotion.deselectAll') 
            : t('promotion.selectAll')}
        </Button>

        <div className="space-y-4">
          {gameTypes?.Data?.map((item: any) => {
            const status = item.status;
            
            const isIncluded = form.watch('includegames')?.split(',').includes(status.id.toString());

            function handleGameTypeToggle(id: string): void {
              const currentIncludeGames = form.watch('includegames')?.split(',').filter(Boolean) || [];
              const currentExcludeGames = form.watch('excludegames')?.split(',').filter(Boolean) || [];

              if (isIncluded) {
                form.setValue('includegames', currentIncludeGames.filter(gameId => gameId !== id).join(','));
                form.setValue('excludegames', [...currentExcludeGames, id].join(','));
              } else {
                form.setValue('excludegames', currentExcludeGames.filter(gameId => gameId !== id).join(','));
                form.setValue('includegames', [...currentIncludeGames, id].join(','));
              }
            }

            return (
              <div key={status.name} className="flex items-center justify-between p-2 border rounded">
                <span>{t(`games.${status.name}`)}</span>
                <Switch
                  checked={isIncluded}
                  onCheckedChange={() => handleGameTypeToggle(status.id.toString())}
                />
              </div>
            );
          })}
        </div>
      </div>
    </TabsContent>
    <div className="flex justify-end space-x-2 mt-6">
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
            }}>{t('promotion.save')}</Button>
            <Button type="button" variant="outline" onClick={onCancel}>{t('promotion.cancel')}</Button>
            <Dialog>
            <DialogTrigger asChild>
              <Button type='button' variant="outline">{t('promotion.testing')}</Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
              <DialogHeader>
                <DialogTitle>{t('promotion.testing')}</DialogTitle>
                <DialogDescription>
                   <TestPromotion lng={lng} promotion={form} />
                </DialogDescription>
              </DialogHeader>
              <DialogFooter className="sm:justify-start">
          <DialogClose asChild>
            <Button type="button" variant="secondary">
              Close
            </Button>
          </DialogClose>
        </DialogFooter>
        </DialogContent>
      </Dialog>
            
          </div>
      </Tabs>
    </div>
    
    </form>
          
    </Form>

  );
};

export default EditPromotionPanel;

