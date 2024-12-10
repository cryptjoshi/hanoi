import React, { useEffect, useState } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { AddGame, GetGameStatus, GetGameById, UpdateGame } from '@/actions';
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

interface EditGamesPanelProps {
  gameId: number | null;
  lng: string;
  prefix: string;
  onClose: () => void;
  onCancel: () => void;
}

// Update the Promotion interface
interface Games {
  id: string;
  productCode: string;
  product: string;
  gameType: string;
  active: number;
  remark: string;
  position: string;
  urlimage: string;
  name: string;
  status: string;
}

// Update the form schema
const formSchema = z.object({
  name: z.string().optional(),
  productCode: z.string().optional(),
  product: z.string().optional(),
  gameType: z.string().optional(),
  active: z.union([z.number(), z.string()]).optional(),
  remark: z.string().optional(),
  position: z.string().optional(),
  urlimage: z.string().optional(),
  status: z.union([z.string(), z.object({
    id: z.string(),
    name: z.string()
  })]).optional(),
})

interface gameStatus {
  id: string;
  name: string;
}
 
function cleanJsonString(jsonString: string): gameStatus {
  if (!jsonString) {
    return { id: '', name: '' };
  }

  try {
    let cleanJsonString = jsonString.trim().replace(/^["']|["']$/g, '');
    cleanJsonString = cleanJsonString.replace(/\\"/g, '"');
    return JSON.parse(cleanJsonString) as gameStatus;
  } catch { 
    return { id: '', name: '' };
  }
}

export const EditGamesPanel: React.FC<EditGamesPanelProps> = ({ gameId, prefix, lng, onClose, onCancel }) => {
 
  const [game, setGame] = useState<Games>({
    id: '',
    name: '',
    productCode:'',
    product:'',
    gameType:'',
    active:0,
    remark: '',
    status: '0',
    position:'',
    urlimage:''
 
  });

  const {t} = useTranslation(lng,'translation',{keyPrefix:'games'})

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      productCode: '',
      product: '',
      gameType: '',
      active: 0,
      remark: '',
      position: '',
      urlimage: '',
      status: '',
    }
  })
  const fetchGame = async (prefix:string,id:number) => {
    const data = await GetGameById(prefix, id);
    form.reset(data.Data as z.infer<typeof formSchema>);
  };


  const [gameStatus, setGameStatus] = useState<gameStatus[]>([]);
  useEffect(() => {
    const fetchGameStatus = async () => {
      const data = await GetGameStatus(prefix);
      
     
      const formattedGameStatus: gameStatus[] = data.Data.map((item: any) => ({
        id: JSON.parse(item.status).id,
        name: JSON.parse(item.status).name
      }));
    //  console.log(formattedGameStatus)
      setGameStatus(formattedGameStatus);
    };
    fetchGameStatus();
    if (gameId) {
      
      fetchGame(prefix,gameId);
    } else {
      form.reset({} as z.infer<typeof formSchema>);
    }

    
  }, [prefix,gameId,gameStatus]);

 

  const handleSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      let processedValues = { ...values };

      // Handle status
      if (typeof processedValues.status === 'string') {
        try {
          processedValues.status = JSON.parse(processedValues.status);
        } catch (error) {
          console.error('Error parsing status:', error);
          // If parsing fails, keep the original string value
        }
      }

      // Handle active (ensure it's a number)
      processedValues.active = Number(processedValues.active);

      console.log('Processed values:', processedValues);
      console.log(gameId)

      if (gameId) {
        const data = await UpdateGame(prefix, gameId, processedValues);

        if (data.Status) {
          toast({
            title: t("edit.success"),
            description: t("edit.success_description"),
            variant: "default",
          })
          onClose();
        } else {
          toast({
            title: t("edit.error"),
            description: t("edit.error_description") + data.Message,
            variant: "destructive",
          })
        }
      } else {
        const data = await AddGame(prefix, processedValues)
        if (data.Status) {
          toast({
            title: t("add.success"),
            description: t("add.success_description"),
            variant: "default",
          })
          onClose();
        } else {
          toast({
            title: t("add.error"),
            description: t("add.error_description") + data.Message,
            variant: "destructive",
          })
        }
      }
    } catch (error) {
      console.error('Error processing values:', error);
      toast({
        title: t("edit.error"),
        description: t("edit.error_description"),
        variant: "destructive",
      });
    }
  };

  return (
    <div className="p-6 bg-white rounded-lg shadow-md md:max-w-md">
      <h2 className="text-2xl font-bold mb-4">{gameId ? t('edit.title') : t('add.title')}</h2>
      <p className="text-gray-600 mb-6">{t('edit.description')}</p>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('columns.name')}</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}  
          />
          <FormField
            control={form.control}
            name="productCode"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('columns.productCode')}</FormLabel>
                <FormControl>
                  <Textarea {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="product"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('columns.product')}</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
         
         
        
          <FormField
            control={form.control}
            name="position"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('columns.position')}</FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="urlimage"
            render={({ field }) => (
            <FormItem>
              <FormLabel>{t('columns.urlimage')}</FormLabel>
              <FormControl>
                <Input {...field}/>
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
         <FormField
            control={form.control}
            name="status"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('columns.gameType')}</FormLabel>
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
                        {gameStatus.map((item: gameStatus) => (
                          <SelectItem key={item.name} value={JSON.stringify(item)}>{item.name}</SelectItem>
                        ))}
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
            name="active"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('columns.active')}</FormLabel>
                <Select
            onValueChange={(value) => field.onChange(parseInt(value))}
                  value={field.value?.toString() || ''}
                >
                <FormControl>
                <SelectTrigger>
                <SelectValue placeholder={t('columns.selectStatus')} />
                </SelectTrigger>
              </FormControl>  
              <SelectContent>
                <SelectItem value="1">{t('active')}</SelectItem>
              <SelectItem value="0">{t('inactive')}</SelectItem>
              <SelectItem value="-1">{t('maintenance')}</SelectItem>
            </SelectContent>
          </Select>
        
          <FormMessage />
        </FormItem>
            )}
          />
          <div className="flex justify-end space-x-2 mt-6">
            <Button type="submit">{t('save')}</Button>
            <Button type="button" variant="outline" onClick={onCancel}>{t('cancel')}</Button>
          </div>
        </form>
      </Form>
    </div>
  );
};

export default EditGamesPanel;
