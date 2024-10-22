"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import { ChevronDownIcon } from "@radix-ui/react-icons"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { useState, useEffect } from "react"
import { GetExchangeRate, UpdateMaster } from "@/actions"

import { cn } from "@/lib/utils"
import { toast } from "@/hooks/use-toast"
import { Button, buttonVariants } from "@/components/ui/button"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { useTranslation } from "@/app/i18n/client"

const appearanceFormSchema = z.object({
  theme: z.enum(["light", "dark"], {
    required_error: "Please select a theme.",
  }).optional(),
  font: z.enum(["inter", "manrope", "system"], {
    invalid_type_error: "Select a font",
    required_error: "Please select a font.",
  }).optional(),
  baseCurrency: z.enum(["USD", "EUR", "THB"], {
    invalid_type_error: "Select a base currency",
    required_error: "Please select a base currency.",
  }).optional(),
  targetCurrency: z.enum(["USD", "EUR", "THB"], {
    invalid_type_error: "Select a target currency",
    required_error: "Please select a target currency.",
  }).optional(),
})

type AppearanceFormValues = z.infer<typeof appearanceFormSchema>

// This can come from your database or API.
const defaultValues: Partial<AppearanceFormValues> = {
  theme: "light",
  font: "inter",
  baseCurrency: "USD",
  targetCurrency: "THB",
}
interface Props {
  prefix: string;
}

export function AppearanceForm({ lng, prefix }: { lng: string, prefix: string }) {

  const {t} = useTranslation(lng,'translation',undefined)
  const form = useForm<AppearanceFormValues>({
    resolver: zodResolver(appearanceFormSchema),
    defaultValues,
  })

  const [exchangeRates, setExchangeRates] = useState<{ [key: string]: number } | null>(null)

  function onSubmit(data: AppearanceFormValues) {

    const updateData = {
      baseCurrency: data.baseCurrency,
      customerCurrency: data.targetCurrency,
      baseRate: exchangeRates[data.baseCurrency || "USD"],
      customerRate: exchangeRates[data.targetCurrency || "USD"]
    }

    const update = async (prefix:string) => {
      const response = await UpdateMaster(prefix,1,updateData)
  
      if(response.Status){
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
    const fetchExchangeRates = async (base: string, target: string) => {
      try {
        const usdData = await GetExchangeRate("USD")
        const baseToUsd = base === "USD" ? 1 : 1 / usdData.rates[base]
        const targetToUsd = target === "USD" ? 1 : 1 / usdData.rates[target]
        
        const baseToTarget = baseToUsd / targetToUsd
        const targetToBase = targetToUsd / baseToUsd

        setExchangeRates({
          [base]: 1,
          [target]: baseToTarget,
          [`${target}To${base}`]: targetToBase
        })
      } catch (error) {
        console.error("Error fetching exchange rates:", error)
        setExchangeRates(null)
      }
    }

    const baseCurrency = form.watch("baseCurrency")
    const targetCurrency = form.watch("targetCurrency")

    if (baseCurrency && targetCurrency && baseCurrency !== targetCurrency) {
      fetchExchangeRates(baseCurrency, targetCurrency)
    } else {
      setExchangeRates(null)
    }
  }, [form.watch("baseCurrency"), form.watch("targetCurrency")])

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <div className="flex space-x-4">
          <FormField
            control={form.control}
            name="baseCurrency"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("settings.appearance.base_currency")}</FormLabel>
                <div className="relative w-max">
                  <FormControl>
                    <select
                      className={cn(
                        buttonVariants({ variant: "outline" }),
                        "w-[200px] appearance-none font-normal"
                      )}
                      {...field}
                    >
                      <option value="USD">USD</option>
                      <option value="EUR">EUR</option>
                      <option value="THB">THB</option>
                    </select>
                  </FormControl>
                  <ChevronDownIcon className="absolute right-3 top-2.5 h-4 w-4 opacity-50" />
                </div>
                <FormMessage />
                <p>1 {form.watch("baseCurrency")}</p>
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="targetCurrency"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("settings.appearance.target_currency")}</FormLabel>
                <div className="relative w-max">
                  <FormControl>
                    <select
                      className={cn(
                        buttonVariants({ variant: "outline" }),
                        "w-[200px] appearance-none font-normal"
                      )}
                      {...field}
                    >
                      <option value="USD">USD</option>
                      <option value="EUR">EUR</option>
                      <option value="THB">THB</option>
                    </select>
                  </FormControl>
                  <ChevronDownIcon className="absolute right-3 top-2.5 h-4 w-4 opacity-50" />
                </div>
                <FormMessage />
                <p>{exchangeRates && exchangeRates[form.watch("targetCurrency") || "USD"] 
                    ? exchangeRates[form.watch("targetCurrency") || "USD"].toFixed(4) 
                    : '1'} {form.watch("targetCurrency")}</p>
              </FormItem>
            )}
          />
        </div>
        {/* <div className="mt-2">
          {form.watch("baseCurrency") === form.watch("targetCurrency") ? (
            <p>1 {form.watch("baseCurrency")} = 1 {form.watch("targetCurrency")} (1:1)</p>
          ) : exchangeRates && exchangeRates[form.watch("targetCurrency")] ? (
            <p>1 {form.watch("baseCurrency")} = {exchangeRates[form.watch("targetCurrency")].toFixed(4)} {form.watch("targetCurrency")}</p>
          ) : (
            <p>Loading...</p>
          )}
        </div> */}
        <FormDescription>
          {/* {t("settings.appearance.currency_description")} */}
        </FormDescription>
        
        <FormField
          control={form.control}
          name="font"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{t("settings.appearance.font")}</FormLabel>
              <div className="relative w-max">
                <FormControl>
                  <select
                    className={cn(
                      buttonVariants({ variant: "outline" }),
                      "w-[200px] appearance-none font-normal"
                    )}
                    {...field}
                  >
                    <option value="inter">Inter</option>
                    <option value="manrope">Manrope</option>
                    <option value="system">System</option>
                  </select>
                </FormControl>
                <ChevronDownIcon className="absolute right-3 top-2.5 h-4 w-4 opacity-50" />
              </div>
              <FormDescription>
                {t("settings.appearance.font_description")}
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="theme"
          render={({ field }) => (
            <FormItem className="space-y-1">
              <FormLabel>{t("settings.appearance.theme")}</FormLabel>
              <FormDescription>
                {t("settings.appearance.theme_description")}
              </FormDescription>
              <FormMessage />
              <RadioGroup
                onValueChange={field.onChange}
                defaultValue={field.value}
                className="grid max-w-md grid-cols-2 gap-8 pt-2"
              >
                <FormItem>
                  <FormLabel className="[&:has([data-state=checked])>div]:border-primary">
                    <FormControl>
                      <RadioGroupItem value="light" className="sr-only" />
                    </FormControl>
                    <div className="items-center rounded-md border-2 border-muted p-1 hover:border-accent">
                      <div className="space-y-2 rounded-sm bg-[#ecedef] p-2">
                        <div className="space-y-2 rounded-md bg-white p-2 shadow-sm">
                          <div className="h-2 w-[80px] rounded-lg bg-[#ecedef]" />
                          <div className="h-2 w-[100px] rounded-lg bg-[#ecedef]" />
                        </div>
                        <div className="flex items-center space-x-2 rounded-md bg-white p-2 shadow-sm">
                          <div className="h-4 w-4 rounded-full bg-[#ecedef]" />
                          <div className="h-2 w-[100px] rounded-lg bg-[#ecedef]" />
                        </div>
                        <div className="flex items-center space-x-2 rounded-md bg-white p-2 shadow-sm">
                          <div className="h-4 w-4 rounded-full bg-[#ecedef]" />
                          <div className="h-2 w-[100px] rounded-lg bg-[#ecedef]" />
                        </div>
                      </div>
                    </div>
                    <span className="block w-full p-2 text-center font-normal">
                      Light
                    </span>
                  </FormLabel>
                </FormItem>
                <FormItem>
                  <FormLabel className="[&:has([data-state=checked])>div]:border-primary">
                    <FormControl>
                      <RadioGroupItem value="dark" className="sr-only" />
                    </FormControl>
                    <div className="items-center rounded-md border-2 border-muted bg-popover p-1 hover:bg-accent hover:text-accent-foreground">
                      <div className="space-y-2 rounded-sm bg-slate-950 p-2">
                        <div className="space-y-2 rounded-md bg-slate-800 p-2 shadow-sm">
                          <div className="h-2 w-[80px] rounded-lg bg-slate-400" />
                          <div className="h-2 w-[100px] rounded-lg bg-slate-400" />
                        </div>
                        <div className="flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-sm">
                          <div className="h-4 w-4 rounded-full bg-slate-400" />
                          <div className="h-2 w-[100px] rounded-lg bg-slate-400" />
                        </div>
                        <div className="flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-sm">
                          <div className="h-4 w-4 rounded-full bg-slate-400" />
                          <div className="h-2 w-[100px] rounded-lg bg-slate-400" />
                        </div>
                      </div>
                    </div>
                    <span className="block w-full p-2 text-center font-normal">
                      Dark
                    </span>
                  </FormLabel>
                </FormItem>
              </RadioGroup>
            </FormItem>
          )}
        />

        <Button type="submit">{t("settings.appearance.update_preferences")}</Button>
      </form>
    </Form>
  )
}
