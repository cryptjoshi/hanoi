import React, { useState, useEffect } from 'react';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription, SheetFooter } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { AddPromotion, GetPromotionById, UpdatePromotion } from '@/actions';
import { useTranslation } from '@/app/i18n/client';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Form, FormField, FormItem, FormLabel, FormControl, FormDescription, FormMessage } from "@/components/ui/form"
import { useForm } from "react-hook-form"
import { toast } from "@/hooks/use-toast"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"

interface EditPromotionPanelProps {
  isOpen: boolean;
  onClose: () => void;
  promotionId: string | null;
  lng: string;
  prefix: string;
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

// Update the form schema to include all fields
const formSchema = z.object({
  name: z.string().min(1, { message: "Name is required" }),
  description: z.string(),
  percentDiscount: z.coerce.number(),
  startDate: z.string(),
  endDate: z.string(),
  maxDiscount: z.coerce.number(),
  usageLimit: z.coerce.number(),
  specificTime: z.string(),
  paymentMethod: z.string(),
  minSpend: z.coerce.number(),
  maxSpend: z.coerce.number(),
  termsAndConditions: z.string(),
  status: z.coerce.number(),
  example: z.string(),
})

export const EditPromotionPanel: React.FC<EditPromotionPanelProps> = ({ isOpen, onClose, promotionId, prefix, lng }) => {
 
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
    if (isOpen && promotionId) {
      
      fetchPromotion(prefix,promotionId);
    } else if (isOpen) {
      form.reset({} as z.infer<typeof formSchema>);
    }
  }, [isOpen, promotionId]);

  const handleSubmit = async (values: z.infer<typeof formSchema>) => {
    // TODO: Implement save logic
    const stringifiedValues = {
      ...values,
      percentDiscount: values.percentDiscount.toString(),
      maxDiscount: values.maxDiscount.toString(),
      usageLimit: values.usageLimit,
      minSpend: values.minSpend.toString(),
      maxSpend: values.maxSpend.toString(),
      status: values.status,
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
          description: t("promotion.edit.error_description"),
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
          description: t("promotion.add.error_description"),
          variant: "destructive",
        })
      }
    }
  };

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent side="bottom" className="h-[90vh] overflow-y-auto max-w-xl mx-auto">
        <SheetHeader>
          <SheetTitle>{promotionId ? t('promotion.edit.title') : t('promotion.add.title')}</SheetTitle>
          <SheetDescription>{t('promotion.edit.description')}</SheetDescription>
        </SheetHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4 py-4">
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
                  <FormLabel>% {t('promotion.percentDiscount')}</FormLabel>
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
                <FormItem>
                  <FormLabel>{t('promotion.startDate')}</FormLabel>
                  <FormControl>
                    <Input {...field} type="date" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="endDate"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('promotion.endDate')}</FormLabel>
                  <FormControl>
                    <Input {...field} type="date" />
                  </FormControl>
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
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            {/* <FormField
              control={form.control}
              name="paymentMethod"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('promotion.paymentMethod')}</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            /> */}
            <FormField
              control={form.control}
              name="minSpend"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>%{t('promotion.minSpend')}</FormLabel>
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
            {/* <FormField
              control={form.control}
              name="example"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('promotion.example')}</FormLabel>
                  <FormControl>
                    <Textarea {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            /> */}
            {/* <FormField
              control={form.control}
              name="termsAndConditions"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('promotion.termsAndConditions')}</FormLabel>
                  <FormControl>
                    <Textarea {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            /> */}
            <FormField
              control={form.control}
              name="status"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('promotion.status')}</FormLabel>
                  <Select onValueChange={field.onChange} value={field.value} defaultValue={field.value?.toString() ?? '0'}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder={t('promotion.selectStatus')} />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="1">{t('promotion.status_active')}</SelectItem>
                      <SelectItem value="0">{t('promotion.status_inactive')}</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
            <SheetFooter>
              <div className="flex justify-end gap-2 w-full">
                <Button className="bg-gray-300 text-black" variant="outline" onClick={onClose}>{t('promotion.cancel')}</Button>
                <Button className="bg-blue-500 text-white" type="submit">{t('promotion.save')}</Button>
              </div>
            </SheetFooter>
          </form>
        </Form>
      </SheetContent>
    </Sheet>
  );
};

export default EditPromotionPanel;
