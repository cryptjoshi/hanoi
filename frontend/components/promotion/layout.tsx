'use client'
//import { Metadata } from "next"
import Image from "next/image"

import { Separator } from "@/components/ui/separator"
import { SidebarNav } from "@/app/forms/components/sidebar-nav"

import { useTranslation } from "@/app/i18n/client"
import useAuthStore from "@/store/auth"

import { useParams } from 'next/navigation'
// export const metadata: Metadata = {
//   title: "Forms",
//   description: "Advanced form example using react-hook-form and Zod.",
// }




interface PromotionLayoutProps {
  children: React.ReactNode
}

export default function PromotionLayout({ children }: PromotionLayoutProps) {
 
  const { lng, setLng } = useAuthStore()
 
  const {t} = useTranslation(lng,'translation',undefined)
  const params = useParams()
  // หรือ const router = useRouter() สำหรับ Next.js เวอร์ชันเก่ากว่า

 
  const prefix = params.prefix as string
 
  const sidebarNavItems = [
    {
      title: "Profile",
      href: `/${lng}/dashboard/agents/${prefix}`,
    },
    {
      title: "Account",
      href: `/${lng}/dashboard/agents/${prefix}/account`,
    },
    {
      title: "Appearance",
      href: `/${lng}/dashboard/agents/${prefix}/appearance`
    },
    {
      title: "Notifications",
      href: `/${lng}/dashboard/agents/${prefix}/notifications`,
    },
    {
      title: "Promotion",
      href: `/${lng}/dashboard/agents/${prefix}/promotion`,
    },
  ]

  return (
    <>
      <div className="md:hidden">
        {/* <Image
          src="/forms-light.png"
          width={1280}
          height={791}
          alt="Forms"
          className="block dark:hidden"
        />
        <Image
          src="/forms-dark.png"
          width={1280}
          height={791}
          alt="Forms"
          className="hidden dark:block"
        /> */}
      </div>
      <div className="hidden space-y-6 p-10 pb-16 md:block">
        {/* <div className="space-y-0.5">
          <h2 className="text-2xl font-bold tracking-tight">{t("agents.settings.title")}</h2>
          <p className="text-muted-foreground">
            {t("agents.settings.description")}
          </p>
        </div> */}
        <Separator className="my-6" />
        <div className="flex flex-col space-y-8 lg:flex-row lg:space-x-12 lg:space-y-0">
          <aside className="-mx-4 lg:w-1/5">
            <SidebarNav items={sidebarNavItems} />
          </aside>
          <div className="flex-1 lg:max-w-2xl">{children}</div>
        </div>
      </div>
    </>
  )
}
