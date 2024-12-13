'use client';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import { CalendarIcon } from '@radix-ui/react-icons';
import { addDays, format } from 'date-fns';
import { TZDate } from "@date-fns/tz";
import th from "date-fns/locale/th"
//import en from "date-fns/locale/en"
import * as React from 'react';
//import { DateRange } from 'react-day-picker';

export function CalendarDatePicker({
  className,
  lng,
  onChange,
  initialDate 
}: React.HTMLAttributes<HTMLDivElement>  & { onChange: (date: Date | undefined) => void,
    initialDate?: Date
 }) {
    const tzDate = new TZDate(new Date(), "Asia/Bangkok");
    const [date, setDate] = React.useState<Date | undefined>(initialDate || tzDate)
    let locale;
    switch (lng) {
      case 'th':
        locale = th; // ภาษาไทย
        break;
      case 'en':
        locale = undefined; // ภาษาอังกฤษ (ใช้ค่าเริ่มต้น)
        break;
      // สามารถเพิ่มกรณีอื่น ๆ สำหรับภาษาที่รองรับได้ที่นี่
      default:
        locale = undefined; // ค่าเริ่มต้น
    }

    const handleDateChange = (newDate: Date | undefined) => {
        //console.log(newDate)
        setDate(newDate);
        onChange(newDate); // ส่งค่ากลับไปยังฟังก์ชันที่เรียกใช้
      };
  

  return (
    <div className={cn('grid gap-2', className)}>
      <Popover>
        <PopoverTrigger asChild>
          <Button
            id="date"
            variant={'outline'}
            className={cn(
              'w-[260px] justify-start text-left font-normal',
              !date && 'text-muted-foreground'
            )}
          >
            <CalendarIcon className="mr-2 h-4 w-4" />
            
            {date?  (
                 format(date, 'dd MMMM yyyy', { locale }) // แก้ไขให้แสดงเดือนเป็นภาษาไทย
     
              )
             : (
              <span>Pick a date</span>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="end">
          <Calendar
            initialFocus
            mode="single"
            selected={date}
            onSelect={handleDateChange} // เปลี่ยนจาก setDate เป็น handleDateChange
           locale={locale}
           weekStartsOn={0} 
           dir="ltr"
           className="ltr-calendar" 
          />
        </PopoverContent>
      </Popover>
    </div>
  );
}