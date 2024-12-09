'use client'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQuery } from '@tanstack/react-query';

import { gameSchema, Game } from "@/lib/zod/game";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
//import { GameStatus } from "@/lib/zod/gameStatus";
import { useTranslation } from "@/app/i18n/client";
import { GetGameById,GetGameStatus,AddGame,UpdateGame} from "@/actions";
import { useEffect } from "react";
import { z } from "zod";
import { toast } from "@/hooks/use-toast";

// const gametype = [
//     { id: '1', name: 'Slot' },
//     { id: '2', name: 'Live Casino' },
// ];

function EditGame({ gameId, lng, prefix, onClose, onCancel,isAdd  }: { gameId: number, lng: string, prefix: string, onClose: () => void, onCancel: () => void,isAdd:boolean }) {
  
  const { t } = useTranslation(lng, 'translation', { keyPrefix: 'games' });

  const { data: gameStatus, isLoading: gameStatusLoading } = useQuery({
    queryKey: ['gameStatus'],
    queryFn: async () => await GetGameStatus(prefix),
  });
 

  const form = useForm<Game>({
    resolver: zodResolver(gameSchema),
    defaultValues: {
      name: '',
      productCode: '',
      product: '',
      gameType: '',
      active: 0,
      remark: '',
      status: '',
    },
  });

  const fetchGame = async (prefix:string,id:number) => {
    const data = await GetGameById(prefix, id);
    form.reset(data.Data as z.infer<typeof gameSchema>);
  };

  useEffect(() => {
    if (gameId) {
      fetchGame(prefix, gameId);
    }
  }, [gameId, prefix]);

  const handleSubmit = async (data: Game) => {
   data.gameType = JSON.parse(data.status?.toString() || '{}').name
   //console.log(data)
   const result = !isAdd ? await UpdateGame(prefix, gameId, data) : await AddGame(prefix, data)
   if (!isAdd) {
 

    if (result.Status) {
      toast({
        title: t("edit.success"),
        description: t("edit.success_description"),
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
        title: t("add.success"),
        description: t("add.success_description"),
        variant: "default",
      })
      onClose();
    } else {
      toast({
        title: t("add.error"),
        description: t("add.error_description") + result.Message,
        variant: "destructive",
      })
    }
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
}

export default EditGame;
