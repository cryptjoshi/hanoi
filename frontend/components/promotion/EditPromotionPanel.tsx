import React, { useState, useEffect } from 'react';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription, SheetFooter } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { GetPromotionById } from '@/actions';
import { useTranslation } from '@/app/i18n/client';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';

interface EditPromotionPanelProps {
  isOpen: boolean;
  onClose: () => void;
  promotionId: string | null;
  lng: string;
  prefix: string;
}

interface Promotion {
  percentDiscount: string | number | readonly string[] | undefined;
  startDate: string | number | readonly string[] | undefined;
  endDate: string | number | readonly string[] | undefined;
  id: string;
  name: string;
  description: string;
  maxDiscount: string | number | readonly string[] | undefined;
  usageLimit: string | number | readonly string[] | undefined;
  specificTime: string | number | readonly string[] | undefined;
  paymentMethod: string | number | readonly string[] | undefined;
  minSpend: string | number | readonly string[] | undefined;
  maxSpend: string | number | readonly string[] | undefined;
  termsAndConditions: string | number | readonly string[] | undefined;
  status: string | number | readonly string[] | undefined;
}

export const EditPromotionPanel: React.FC<EditPromotionPanelProps> = ({ isOpen, onClose, promotionId, prefix, lng }) => {
  const [promotion, setPromotion] = useState<Promotion>({
    id: '',
    name: '',
    description: '',
    percentDiscount: '',
    startDate: '',
    endDate: '',
    maxDiscount: '',
    usageLimit: '',
    specificTime: '',
    paymentMethod: '',
    minSpend: '',
    maxSpend: '',
    termsAndConditions: '',
    status: '',
  });

  const {t} = useTranslation(lng,'translation',undefined)

  useEffect(() => {
    if (isOpen && promotionId) {
      const fetchPromotion = async () => {
        const data = await GetPromotionById(prefix,promotionId);
        setPromotion(data);
      };
      fetchPromotion();
    } else if (isOpen) {
      setPromotion({
        id: '',
        name: '',
        description: '',
        percentDiscount: '',
        startDate: '',
        endDate: '',
        maxDiscount: '',
        usageLimit: '',
        specificTime: '',
        paymentMethod: '',
        minSpend: '',
        maxSpend: '',
        termsAndConditions: '',
        status: '',
      });
    }
  }, [isOpen, promotionId]);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
      // TODO: Implement save logic
    onClose();
  };

  return (
    <Sheet open={isOpen} onOpenChange={onClose}>
      <SheetContent side="bottom" className="h-[90vh] overflow-y-auto">
        <SheetHeader>
          <SheetTitle>{promotionId ? t('promotion.edit.title') : t('promotion.add.title')}</SheetTitle>
          <SheetDescription>{t('promotion.edit.description')}</SheetDescription>
        </SheetHeader>
        <form onSubmit={handleSubmit} className="space-y-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
           
            <Input
              id="name"
              value={promotion.name}
              onChange={(e) => setPromotion({ ...promotion, name: e.target.value })}
              className="col-span-3"
            />
             <Label htmlFor="name" className="text-lef">
              {t('promotion.name')}
            </Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
           
            <Textarea
              id="description"
              value={promotion.description}
              onChange={(e) => setPromotion({ ...promotion, description: e.target.value })}
              className="col-span-3"
            />
             <Label htmlFor="description" className="text-right">
              {t('promotion.description')}
            </Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            
            <Input
              id="percentDiscount"
              value={promotion.percentDiscount}
              onChange={(e) => setPromotion({ ...promotion, percentDiscount: e.target.value })}
              className="col-span-3"
            />
            <Label htmlFor="percentDiscount" className="text-right">
              {t('promotion.percentDiscount')}
            </Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
             
            <Input
              id="startDate"
              value={promotion.startDate}
              onChange={(e) => setPromotion({ ...promotion, startDate: e.target.value })}
              className="col-span-3"
            />
            <Label htmlFor="startDate" className="text-right">{t('promotion.startDate')}</Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
           
            <Input
              id="endDate"
              value={promotion.endDate}
              onChange={(e) => setPromotion({ ...promotion, endDate: e.target.value })}
              className="col-span-3"
            />
              <Label htmlFor="endDate" className="text-right">{t('promotion.endDate')}</Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
       
            <Input
              id="maxDiscount"
              value={promotion.maxDiscount}
              onChange={(e) => setPromotion({ ...promotion, maxDiscount: e.target.value })}
              className="col-span-3"
            />
                 <Label htmlFor="maxDiscount" className="text-right">
              {t('promotion.maxDiscount')}
            </Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            
            <Input
              id="usageLimit"
              value={promotion.usageLimit}
              onChange={(e) => setPromotion({ ...promotion, usageLimit: e.target.value })}
              className="col-span-3"
            /><Label htmlFor="usageLimit" className="text-right">
            {t('promotion.usageLimit')}
          </Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
           
            <Input
              id="specificTime"
              value={promotion.specificTime}
              onChange={(e) => setPromotion({ ...promotion, specificTime: e.target.value })}
              className="col-span-3"
            />
            <Label htmlFor="specificTime" className="text-right">
              {t('promotion.specificTime')}
            </Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
           
            <Input
              id="paymentMethod"
              value={promotion.paymentMethod}
              onChange={(e) => setPromotion({ ...promotion, paymentMethod: e.target.value })}
              className="col-span-3"
            />
            <Label htmlFor="paymentMethod" className="text-right">{t('promotion.paymentMethod')}</Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
           
            <Input
              id="minSpend"
              value={promotion.minSpend}
              onChange={(e) => setPromotion({ ...promotion, minSpend: e.target.value })}
              className="col-span-3"
            />
            <Label htmlFor="minSpend" className="text-right">{t('promotion.minSpend')}</Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">

            <Input
              id="maxSpend"
              value={promotion.maxSpend}
              onChange={(e) => setPromotion({ ...promotion, maxSpend: e.target.value })}
              className="col-span-3"
            />
            <Label htmlFor="maxSpend" className="text-right">{t('promotion.maxSpend')}</Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            
            <Textarea
              id="termsAndConditions"
              value={promotion.termsAndConditions}
              onChange={(e) => setPromotion({ ...promotion, termsAndConditions: e.target.value })}
              className="col-span-3"
            />
            <Label htmlFor="termsAndConditions" className="text-right">{t('promotion.termsAndConditions')}</Label>
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
           
            <Select
              value={promotion.status as string}
              onValueChange={(value: string) => setPromotion({ ...promotion, status: value })}
            >
              <SelectTrigger className="col-span-3">
                <SelectValue placeholder={t('promotion.selectStatus')} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="active">{t('promotion.statusActive')}</SelectItem>
                <SelectItem value="inactive">{t('promotion.statusInactive')}</SelectItem>
              </SelectContent>
            </Select>
            <Label htmlFor="status" className="text-right">{t('promotion.status')}</Label>
          </div>
          <SheetFooter>
            <Button type="submit">{t('promotion.save')}</Button>
          </SheetFooter>
        </form>
      </SheetContent>
    </Sheet>
  );
};

export default EditPromotionPanel;
