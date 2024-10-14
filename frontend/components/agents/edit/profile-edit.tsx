"use client"

import Link from "next/link"
import { zodResolver } from "@hookform/resolvers/zod"
import { useFieldArray, useForm,SubmitHandler } from "react-hook-form"
import { z } from "zod"
import { cn } from "@/lib/utils"
import { toast } from "@/hooks/use-toast"
import { Button } from "@/components/ui/button"
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useTranslation } from "@/app/i18n/client"
import { Textarea } from "@/components/ui/textarea"
import { CreateUser, GetDatabaseListByPrefix, UpdateDatabaseListByPrefix } from "@/actions"
import { useState, useEffect, useMemo } from "react"
import { Checkbox } from "@/components/ui/checkbox"
import { useRouter } from 'next/navigation'
//import { useQuery } from '@tanstack/react-query'

 

const profileFormSchema = z.object({
  username: z
    .string()
    .min(2, {
      message: "ต้องไม่น้อยกว่า 3 ตัวอักษร.",
    })
    .max(30, {
      message:  "ต้องไม่เกิน 30 ตัวอักษร.",
    }),
  prefix: z
    .string()
    .min(3, {
      message: "ต้องไม่น้อยกว่า 3 ตัวอักษร.",
    })
    .max(5, {
      message: "ต้องไม่เกิน 5 ตัวอักษร.",
    }),
  dbname: z.string().min(3, {
    message: "ต้องไม่น้อยกว่า 3 ตัวอักษร.",
  }),
  // email: z
  //   .string({
  //     required_error: "Please select an email to display.",
  //   })
  //   .email(),
  // bio: z.string().max(160).min(4),
  // urls: z
  //   .array(
  //     z.object({
  //       value: z.string().url({ message: "Please enter a valid URL." }),
  //     })
  //   )
  //   .optional(),
})

type ProfileFormValues = z.infer<typeof profileFormSchema>

// This can come from your database or API.
const defaultValues: Partial<ProfileFormValues> = {
  bio: "I own a computer.",
  urls: [
    { value: "https://shadcn.com" },
    { value: "http://twitter.com/shadcn" },
  ],
}

interface ModeSelection {
  development: boolean;
  production: boolean;
}

interface ProfileEditProps {
  lng: string;
  id: string;
}

export function ProfileEdit({ lng, id }: ProfileEditProps) {
  
  //const { lng, setLng } = useAuthStore()
 
  const {t} = useTranslation(lng,'translation',undefined)
  
  const [isDbNameSameAsPrefix, setIsDbNameSameAsPrefix] = useState(true)
  const [modeSelection, setModeSelection] = useState<ModeSelection>({
    development: true,
    production: false
  })

  const router = useRouter();

  // const { data: userData, isLoading, error } = useQuery({
  //   queryKey: ["prefix", id],
  //   queryFn: async () => {
  //     const response = await GetDatabaseListByPrefix(id);
  //     return response.data;
  //   }
  // });
 

  const form = useForm<ProfileFormValues>({
    resolver: zodResolver(profileFormSchema),
    defaultValues,
    mode: "onChange",
  })

  useEffect(() => {
    const fetchData = async () => {
      const response = await GetDatabaseListByPrefix(id);
     
      if (response.Status) {
        // Response: {Databases: ["ckd_dev","ckd_prod"], Message: "ดึงรายชื่อฐานข้อมูลสำเร็จ", Status: true}
        const databases = response.Databases || [];
        
        if (databases.length > 0) {
          const [prefix, ...rest] = databases[0].split('_');
          
          form.reset({
            username: id, // Assuming id is the username
            prefix: prefix,
            dbname: databases.join(', '),
          });

          // Check if all database names are the same as prefix + mode
          const isDbNameSameAsPrefix = databases.every(db => 
            db === `${prefix}_dev` || db === `${prefix}_prod`
          );
          setIsDbNameSameAsPrefix(isDbNameSameAsPrefix);

          // Set mode selection based on database names
          setModeSelection({
            development: databases.some(db => db.endsWith('_dev')),
            production: databases.some(db => db.endsWith('_prod')),
          });
        }
      }
    }
    fetchData();
  }, [id, form]);

  const prefixValue = form.watch("prefix")
  const customDbName = form.watch("customDbName")

  const generateDbNames = useMemo(() => {
    const names: string[] = [];
    if (modeSelection.development) {
      names.push(isDbNameSameAsPrefix ? `${prefixValue}_development` : `${prefixValue}_${customDbName || prefixValue}_development`);
    }
    if (modeSelection.production) {
      names.push(isDbNameSameAsPrefix ? `${prefixValue}_production` : `${prefixValue}_${customDbName || prefixValue}_production`);
    }
    return names.join(', ');
  }, [isDbNameSameAsPrefix, modeSelection, prefixValue, customDbName]);

  useEffect(() => {
    form.setValue("dbname", generateDbNames, { shouldValidate: true });
  }, [generateDbNames, form]);

  const handleCheckboxChange = (checked: boolean) => {
    if (checked && prefixValue.length < 3) {
      toast({
        title: "ข้อผิดพลาด",
        description: "ชื่อย่อต้องมีความยาวอย่างน้อย 3 ตัวอักษรเพื่อใช้เป็นชื่อฐานข้อมูล",
        variant: "destructive",
      })
      return;
    }
    setIsDbNameSameAsPrefix(checked);
  }

  const onSubmit: SubmitHandler<ProfileFormValues> = async (data: ProfileFormValues) => {
    // สร้าง array ของชื่อฐานข้อมูลจาก dbname
    const dbNamesArray = data.dbname.split(', ');
    
    const submitData = {
      ...data,
      dbnames: dbNamesArray,
    };

    try {
      const response = await UpdateDatabaseListByPrefix(submitData);  
      toast({
        title: "อัพเดทข้อมูลสำเร็จ",
        description: (
          <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
            <code className="text-white">{JSON.stringify(response.data, null, 2)}</code>
          </pre>
        ),
      });
      router.push(`/${lng}/dashboard/agents`);
    } catch (error) {
      toast({
        title: "การอัพเดทข้อมูลล้มเหลว",
        description: "กรุณาตรวจสอบข้อมูลอีกครั้ง",
        variant: "destructive",
      });
    }
  }

  // if (isLoading) return <div>กำลังโหลด...</div>;
  // if (error) return <div>เกิดข้อผิดพลาดในการโหลดข้อมูล</div>;

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"ชื่อผู้ใช้"}</FormLabel>
              <FormControl>
                <Input placeholder="" {...field} />
              </FormControl>
              <FormDescription>
                ชื่อผู้ใช้ เพื่อเข้าสู่ระบบจัดการ
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
         <FormField
          control={form.control}
          name="prefix"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"ชื่อย่อ"}</FormLabel>
              <FormControl>
                <Input placeholder="" {...field} />
              </FormControl>
              <FormDescription>
               ชื่อย่อ เพื่อกำหนดในฐานข้อมูล 
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="samePrefix"
          render={({ field }) => (
            <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
              <FormControl>
                <Checkbox
                  checked={isDbNameSameAsPrefix}
                  onCheckedChange={handleCheckboxChange}
                />
              </FormControl>
              <div className="space-y-1 leading-none">
                <FormLabel>
                  ใช้ชื่อย่อเป็นชื่อฐานข้อมูล
                </FormLabel>
                <FormDescription>
                  เมื่อเลือกตัวเลือกนี้ ชื่อฐานข้อมูลจะใช้เฉพาะชื่อย่อ (ชื่อย่อต้องมีความยาวอย่างน้อย 5 ตัวอักษร)
                </FormDescription>
              </div>
            </FormItem>
          )}
        />
        <div className="space-y-2">
          <FormLabel>โหมดการสร้างฐานข้อมูล</FormLabel>
          <div className="flex space-x-4">
            <FormField
              control={form.control}
              name="developmentMode"
              render={({ field }) => (
                <FormItem className="flex flex-row items-start space-x-3 space-y-0">
                  <FormControl>
                    <Checkbox
                      checked={modeSelection.development}
                      onCheckedChange={(checked) => {
                        setModeSelection(prev => ({ ...prev, development: checked as boolean }));
                      }}
                    />
                  </FormControl>
                  <FormLabel>Development</FormLabel>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="productionMode"
              render={({ field }) => (
                <FormItem className="flex flex-row items-start space-x-3 space-y-0">
                  <FormControl>
                    <Checkbox
                      checked={modeSelection.production}
                      onCheckedChange={(checked) => {
                        setModeSelection(prev => ({ ...prev, production: checked as boolean }));
                      }}
                    />
                  </FormControl>
                  <FormLabel>Production</FormLabel>
                </FormItem>
              )}
            />
          </div>
          <FormDescription>
            เลือกโหมดที่ต้องการสร้างฐานข้อมูล (สามารถเลือกได้ทั้งสองโหมด)
          </FormDescription>
        </div>

        {!isDbNameSameAsPrefix && (
          <FormField
            control={form.control}
            name="customDbName"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{"ชื่อฐานข้อมูลเพิ่มเติม"}</FormLabel>
                <FormControl>
                  <Input
                    placeholder="ระบุชื่อฐานข้อมูลเพิ่มเติม (ถ้ามี)"
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  ระบุชื่อฐานข้อมูลเพิ่มเติมหากไม่ต้องการใช้ชื่อย่อ
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        )}

        <FormField
          control={form.control}
          name="dbname"
          render={({ field }) => (
            <FormItem>
              <FormLabel>{"ชื่อฐานข้อมูลที่มีอยู่"}</FormLabel>
              <FormControl>
                <Input
                  {...field}
                  disabled={true}
                  readOnly={true}
                />
              </FormControl>
              <FormDescription>
                ชื่อฐานข้อมูลที่มีอยู่ในระบบ
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        {/* <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select a verified email to display" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem value="m@example.com">m@example.com</SelectItem>
                  <SelectItem value="m@google.com">m@google.com</SelectItem>
                  <SelectItem value="m@support.com">m@support.com</SelectItem>
                </SelectContent>
              </Select>
              <FormDescription>
                You can manage verified email addresses in your{" "}
                <Link href="/examples/forms">email settings</Link>.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        /> */}
        {/* <FormField
          control={form.control}
          name="bio"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Bio</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="Tell us a little bit about yourself"
                  className="resize-none"
                  {...field}
                />
              </FormControl>
              <FormDescription>
                You can <span>@mention</span> other users and organizations to
                link to them.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        /> */}
        {/* <div>
          {fields.map((field, index) => (
            <FormField
              control={form.control}
              key={field.id}
              name={`urls.${index}.value`}
              render={({ field }) => (
                <FormItem>
                  <FormLabel className={cn(index !== 0 && "sr-only")}>
                    URLs
                  </FormLabel>
                  <FormDescription className={cn(index !== 0 && "sr-only")}>
                    Add links to your website, blog, or social media profiles.
                  </FormDescription>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          ))}
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="mt-2"
            onClick={() => append({ value: "" })}
          >
            Add URL
          </Button>
        </div> */}
        <Button type="submit">{t("agents.settings.edit.submit")}</Button>
      </form>
    </Form>
  )
}