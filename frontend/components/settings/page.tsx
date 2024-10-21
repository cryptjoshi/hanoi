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
import { AppearanceForm } from "./appearance/appearance-form"
import About from "./about"

export default function EditSettings({ params: { lng,id } }: { params: { lng: string,id: string } }) {
  const {t} =  useTranslation(lng,'translation','')
  return (
    <div className="space-y-6 md:container md:mx-auto md:px-4">
      <Tabs defaultValue="general" className="w-full h-auto md:h-full">
        <TabsList className="w-full flex-wrap justify-start">
          <TabsTrigger value="general" className="flex-grow md:flex-grow-0">{t('settings.general')}</TabsTrigger>
          <TabsTrigger value="about" className="flex-grow md:flex-grow-0">{t('settings.about')}</TabsTrigger>
        </TabsList>
        <TabsContent value="general" className="h-auto md:h-[calc(100vh-120px)]">
          <AppearanceForm lng={lng} id={id} />
        </TabsContent>
        <TabsContent value="about">
          <About lng={lng} id={id} />
        </TabsContent>
      </Tabs>
      {/* <SettingsLayout>
        <ProfileEdit lng={lng} id={id} />
      </SettingsLayout> */}
    </div>
  )
}
