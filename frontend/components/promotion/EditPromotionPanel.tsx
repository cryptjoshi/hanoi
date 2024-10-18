import React, { useEffect, useState } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { AddPromotion, GetPromotionById, UpdatePromotion } from '@/actions';
import { useTranslation } from '@/app/i18n/client';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage, FormDescription } from "@/components/ui/form"
import { useForm } from "react-hook-form"
import { toast } from "@/hooks/use-toast"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { cn } from "@/lib/utils"
import { format, parse } from "date-fns"
import th from "date-fns/locale/th"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { CalendarIcon } from 'lucide-react';
import { Calendar } from '@/components/ui/calendar';
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";

interface EditPromotionPanelProps {
  promotionId: string | null;
  lng: string;
  prefix: string;
  onClose: () => void;
  onCancel: () => void;
}

// Update the Promotion interface
interface Promotion {
  id: string;
  name: string;
  description: string;
  percentDiscount: number;
  startDate: string;
  endDate: string;
  maxDiscount: number;
  usageLimit: number;
  specificTime: string;
  paymentMethod: string;
  minSpend: number;
  maxSpend: number;
  termsAndConditions: string;
  example: string;
  status: string;
}

// Update the form schema
const formSchema = z.object({
  name: z.string().min(1, { message: "Name is required" }),
  description: z.string(),
  percentDiscount: z.coerce.number(),
  startDate: z.string(),
  endDate: z.string(),
  maxDiscount: z.coerce.number(),
  usageLimit: z.coerce.number(),
  specificTime: z.string().refine((val) => {
    try {
      const parsed = JSON.parse(val);
      return parsed && typeof parsed === 'object';
    } catch {
      return false;
    }
  }, { message: "Invalid JSON string" }),
  paymentMethod: z.string(),
  minSpend: z.coerce.number(),
  maxSpend: z.coerce.number(),
  termsAndConditions: z.string(),
  status: z.coerce.number(),
  example: z.string(),
})

interface SpecificTime {
  type?: string;
  daysOfWeek?: string[];
  hour?: string;
  minute?: string;
}

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

export const EditPromotionPanel: React.FC<EditPromotionPanelProps> = ({ promotionId, prefix, lng, onClose, onCancel }) => {
 
  const [promotion, setPromotion] = useState<Promotion>({
    id: '',
    name: '',
    description: '',
    percentDiscount: 0,
    startDate: '',
    endDate: '',
    maxDiscount: 0,
    usageLimit: 0,
    specificTime: '',
    paymentMethod: '',
    minSpend: 0,
    maxSpend: 0,
    termsAndConditions: '',
    status: '0',
    example: '',
  });

  const {t} = useTranslation(lng,'translation',undefined)

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: promotion // Initialize with an empty object
  })
  const fetchPromotion = async (prefix:string,promotionId:string) => {
    const data = await GetPromotionById(prefix, promotionId);
    form.reset(data.Data as z.infer<typeof formSchema>);
  };
  useEffect(() => {
    if (promotionId) {
      
      fetchPromotion(prefix,promotionId);
    } else {
      form.reset({} as z.infer<typeof formSchema>);
    }
  }, [promotionId]);

  const handleSubmit = async (values: z.infer<typeof formSchema>) => {

    const specificTime = cleanJsonString(values.specificTime);
    const stringifiedValues = {
      ...values,
      percentDiscount: values.percentDiscount.toString(),
      maxDiscount: values.maxDiscount.toString(),
      usageLimit: values.usageLimit,
      minSpend: values.minSpend.toString(),
      maxSpend: values.maxSpend.toString(),
      status: values.status, // Convert status to string
      specificTime: JSON.stringify(specificTime),
    };
   
    if (promotionId) {
      const data = await UpdatePromotion(prefix, promotionId, stringifiedValues);

      if (data.Status) {
     
        toast({
          title: t("promotion.edit.success"),
          description: t("promotion.edit.success_description"),
          variant: "default",
        })
        onClose();
      } else {
        toast({
          title: t("promotion.edit.error"),
          description:  t("promotion.edit.error_description")+data.Message,
          variant: "destructive",
        })
      }
    } else {
    
      const data = await AddPromotion(prefix,stringifiedValues)
      if (data.Status) {
      
        toast({
          title: t("promotion.add.success"),
          description: t("promotion.add.success_description"),
          variant: "default",
        })
        onClose();
      } else {
        toast({
          title: t("promotion.add.error"),
          description:  t("promotion.add.error_description")+data.Message,
          variant: "destructive",
        })
      }
    }
  };

  return (
    <div className="p-6 bg-white rounded-lg shadow-md md:max-w-md">
      <h2 className="text-2xl font-bold mb-4">{promotionId ? t('promotion.edit.title') : t('promotion.add.title')}</h2>
      <p className="text-gray-600 mb-6">{t('promotion.edit.description')}</p>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
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
            name="percentDiscount"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.percentDiscount')}  %</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="startDate"
            render={({ field }) => (
              <FormItem className="flex flex-col">
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
                          format(field.value, "dd-MM-yyyy")
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
                      timezone="Asia/Bangkok"
                      locale={th}
                      selected={field.value}
                      onSelect={field.onChange}
                      disabled={(date) =>
                        date > new Date() || date < new Date("1900-01-01")
                      }
                      initialFocus
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
              <FormItem className="flex flex-col">
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
                        {field.value ? (
                          format(field.value, "dd-MM-yyyy")
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
                      timezone="Asia/Bangkok"
                      locale={th}
                      selected={field.value}
                      onSelect={field.onChange}
                      disabled={(date) =>
                        date > new Date() || date < new Date("1900-01-01")
                      }
                      initialFocus
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
                  <Input {...field} type="number" />
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
                        const newValue = { ...JSON.parse(field.value || '{}'), type: value };
                        field.onChange(JSON.stringify(newValue));
                      }}
                      value={cleanJsonString(field.value || '{}').type}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder={t('promotion.selectTimeType')} />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
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
                                    const newValue = { ...currentValue, daysOfWeek: updatedDays };
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
            name="minSpend"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.minSpend')} %</FormLabel>
                <FormControl>
                  <Input {...field} type="number" />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="maxSpend"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('promotion.maxSpend')}</FormLabel>
                <FormControl>
                  <Input {...field} type="number" />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="status"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                <div className="space-y-0.5">
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
          <div className="flex justify-end space-x-2 mt-6">
            <Button type="submit">{t('promotion.save')}</Button>
            <Button type="button" variant="outline" onClick={onCancel}>{t('promotion.cancel')}</Button>
          </div>
        </form>
      </Form>
    </div>
  );
};

export default EditPromotionPanel;
