'use client';

import React from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from '@/app/i18n/client';
import { DeletePromotion } from '@/actions';
import { Button } from '@/components/ui/button';
import { DialogFooter } from '@/components/ui/dialog';

interface DeletePromotionDialogProps {
  prefix: string;
  lng: string;
  promotionId: string;
  setIsOpen: (isOpen: boolean) => void;
  onDeleteSuccess?: () => void;
}

export default function DeletePromotionDialog({
  prefix,
  lng,
  promotionId,
  setIsOpen,
  onDeleteSuccess
}: DeletePromotionDialogProps) {
  const { t } = useTranslation(lng, 'translation', undefined);
  const { handleSubmit } = useForm();

  const onSubmit = async () => {
    try {
     
       await DeletePromotion(prefix, promotionId);
      setIsOpen(false);
      if (onDeleteSuccess) {
        onDeleteSuccess();
      }
    } catch (error) {
      console.error('Error deleting promotion:', error);
      // You might want to show an error message to the user here
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <p>{t('promotion.edit.delete_confirm')}</p>
      <DialogFooter>
        <Button type="button" variant="secondary" onClick={() => setIsOpen(false)}>
          {t('promotion.edit.cancel')}
        </Button>
        <Button type="submit" variant="destructive">
          {t('promotion.edit.delete')}
        </Button>
      </DialogFooter>
    </form>
  );
}
