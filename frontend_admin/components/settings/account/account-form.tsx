"use client"
import {useEffect,useState} from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { CalendarIcon, CaretSortIcon, CheckIcon } from "@radix-ui/react-icons"
import { format } from "date-fns"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { cn } from "@/lib/utils"
import { toast } from "@/hooks/use-toast"
import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { useTranslation } from "@/app/i18n/client"
import { GetCommission,UpdateMaster } from "@/actions"


const languages = [
  { label: "English", value: "en" },
  { label: "French", value: "fr" },
  { label: "German", value: "de" },
  { label: "Spanish", value: "es" },
  { label: "Portuguese", value: "pt" },
  { label: "Russian", value: "ru" },
  { label: "Japanese", value: "ja" },
  { label: "Korean", value: "ko" },
  { label: "Chinese", value: "zh" },
] as const

const accountFormSchema = z.object({
  partner_commission:z.string().default("5%"),
  user_commission:z.string().default("5%"),
})

type AccountFormValues = z.infer<typeof accountFormSchema>

// This can come from your database or API.
const defaultValues: Partial<AccountFormValues> = {
  partner_commission:"5%",
  user_commission:"5%",
}

export function AccountForm({ lng,prefix }: { lng: string,prefix:string }) {

  const { t } =  useTranslation(lng, "translation",undefined)
  
 
  
  const form = useForm<AccountFormValues>({
    resolver: zodResolver(accountFormSchema),
    defaultValues,
  })


  const [settings,setSettings] = useState([])


  function onSubmit(data: AccountFormValues) {

    const update = async (prefix:string) => {
      const response = await UpdateMaster(prefix,1,[{
        "key":   `${prefix}_partner_commission`,
        "value": data.partner_commission
    },
    {
        "key":   `${prefix}_user_commission`,
        "value": data.user_commission
    },
   ])
  
      if(response.Status){
        //setCustomerCurrency(data.targetCurrency)
        toast({
          title: t("settings.appearance.update_preferences"),
        description: (
          <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
            <code className="text-white">{JSON.stringify(data, null, 2)}</code>
          </pre>
          ),
        })
      
      }else{
        toast({
          title: t("settings.appearance.update_preferences"),
          description: response.Message,
        })
      }
    }
    update(prefix) 
  }

  useEffect(() => {
    const fetchSettings = async (prefix: string) => {
      try {
        const data = await  GetCommission(prefix)
        //console.log(data)
        if(data.Status){
         // form.Set()
        //setSettings(data.Data)
        // console.log(data.Data[1].key,data.Data[1].value)
        //z.Set("partner_commission",data.Data[1].value)
        // console.log(data.Data[2].key,data.Data[2].value)
       //z.Set("user_commission",data.Data[2].value)

        //console.log(data.Data.filter((obj:any)=>obj.key.indexOf("partner_commission"))[1])
        //console.log(data.Data.filter((obj:any)=>obj.key.indexOf("user_commission"))[1])
        
       // console.log(data.Data)
        
        form.setValue("partner_commission",data.Data[1].value)
        form.setValue("user_commission",data.Data[2].value)
        




        // Remove ID from formattedData before setting form values
        // const { ID, ...formData } = formattedData;
        // form.reset(formData as z.infer<typeof updatedFormSchema>);
        // } 

        // const usdData = await GetExchangeRate("USD")
        // const baseToUsd = base === "USD" ? 1 : 1 / usdData.rates[base]
        // const targetToUsd = target === "USD" ? 1 : 1 / usdData.rates[target]
        
        // const baseToTarget = baseToUsd / targetToUsd
        // const targetToBase = targetToUsd / baseToUsd

        // setExchangeRates({
        //   [base]: 1,
        //   [target]: baseToTarget,
        //   [`${target}To${base}`]: targetToBase
        // })
      }
      } catch (error) {
        console.error("Error fetching agent settigs:", error)
         
      }
    }
    
    fetchSettings(prefix)
  
  },[])




  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
      <FormField
          control={form.control}
          name="partner_commission"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("agents.settings.partner_commission")}</FormLabel>
              <FormControl>
                <Input placeholder={t("agents.settings.partner_commission")} {...field} />
              </FormControl>
              <FormDescription>
               {t("agents.settings.partner_descriptions")}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="user_commission"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("agents.settings.user_commission")}</FormLabel>
              <FormControl>
                <Input placeholder={t("agents.settings.user_commission")} {...field} />
              </FormControl>
              <FormDescription>
               {t("agents.settings.user_descriptions")}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        {/* <FormField
          control={form.control}
          name="language"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Language</FormLabel>
              <Popover>
                <PopoverTrigger asChild>
                  <FormControl>
                    <Button
                      variant="outline"
                      role="combobox"
                      className={cn(
                        "w-[200px] justify-between",
                        !field.value && "text-muted-foreground"
                      )}
                    >
                      {field.value
                        ? languages.find(
                            (language) => language.value === field.value
                          )?.label
                        : "Select language"}
                      <CaretSortIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className="w-[200px] p-0">
                  <Command>
                    <CommandInput placeholder="Search language..." />
                    <CommandList>
                      <CommandEmpty>No language found.</CommandEmpty>
                      <CommandGroup>
                        {languages.map((language) => (
                          <CommandItem
                            value={language.label}
                            key={language.value}
                            onSelect={() => {
                              form.setValue("language", language.value)
                            }}
                          >
                            <CheckIcon
                              className={cn(
                                "mr-2 h-4 w-4",
                                language.value === field.value
                                  ? "opacity-100"
                                  : "opacity-0"
                              )}
                            />
                            {language.label}
                          </CommandItem>
                        ))}
                      </CommandGroup>
                    </CommandList>
                  </Command>
                </PopoverContent>
              </Popover>
              <FormDescription>
                This is the language that will be used in the dashboard.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        /> */}
        <Button type="submit">{t("agents.settings.edit.submit")}</Button>
      </form>
    </Form>
  )
}
