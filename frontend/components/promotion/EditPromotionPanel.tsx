import React, { useState, useEffect } from 'react';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription, SheetFooter } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { GetPromotionById } from '@/actions';

interface EditPromotionPanelProps {
  promotionId: string;
  onClose: () => void;
}

interface Promotion {
  id: string;
  name: string;
  description: string;
  // เพิ่มฟิลด์อื่น ๆ ตามโครงสร้างข้อมูลโปรโมชันของคุณ
}

export const EditPromotionPanel: React.FC<EditPromotionPanelProps> = ({ promotionId, onClose }) => {
  const [promotion, setPromotion] = useState<Promotion | null>(null);

  useEffect(() => {
    // โหลดข้อมูลโปรโมชันจาก API
    const fetchPromotion = async () => {
      // สมมติว่ามีฟังก์ชัน getPromotionById ที่เรียก API
      const data = await GetPromotionById(promotionId);
      setPromotion(data);
    };

    fetchPromotion();
  }, [promotionId]);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    // ส่งข้อมูลที่แก้ไขไปยัง API
    // await updatePromotion(promotion);
    onClose();
  };

  if (!promotion) return null;

  return (
    <Sheet open={true} onOpenChange={onClose}>
      <SheetContent>
        <SheetHeader>
          <SheetTitle>แก้ไขโปรโมชัน</SheetTitle>
          <SheetDescription>แก้ไขรายละเอียดโปรโมชันที่นี่</SheetDescription>
        </SheetHeader>
        <form onSubmit={handleSubmit}>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="name" className="text-right">
                ชื่อ
              </Label>
              <Input
                id="name"
                value={promotion.name}
                onChange={(e) => setPromotion({ ...promotion, name: e.target.value })}
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="description" className="text-right">
                รายละเอียด
              </Label>
              <Textarea
                id="description"
                value={promotion.description}
                onChange={(e) => setPromotion({ ...promotion, description: e.target.value })}
                className="col-span-3"
              />
            </div>
            {/* เพิ่มฟิลด์อื่น ๆ ตามต้องการ */}
          </div>
          <SheetFooter>
            <Button type="submit">บันทึกการเปลี่ยนแปลง</Button>
          </SheetFooter>
        </form>
      </SheetContent>
    </Sheet>
  );
};

export default EditPromotionPanel;
