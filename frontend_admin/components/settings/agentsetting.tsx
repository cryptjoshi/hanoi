"use client"
import { Separator } from "@/components/ui/separator"
import { ProfileEdit } from "@/components/agents/edit/profile-edit"
import SettingsLayout from "@/components/agents/layout"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import PromotionList from "@/components/promotion/list"
import GameList from "@/components/games/list"
import MemberList from "@/components/member/list"
import { useTranslation } from "@/app/i18n/client"
import { General } from "./general"
import {AccountForm} from "@/components/settings/account/account-form"
import { About } from "./about"

export default function AgentSettings({ lng,id }: { lng: string,id: string  }) {
  const {t} =  useTranslation(lng,'translation','')
  return (
    <div className="space-y-6 md:container md:mx-auto md:px-4">
    <AccountForm lng={lng} prefix={id}/>
    </div>
  )
}
