"use client"
import { Separator } from "@/components/ui/separator"
import { ProfileEdit } from "@/components/agents/edit/profile-edit"
import SettingsLayout from "@/components/agents/layout"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import PromotionList from "@/components/promotion/list"
import GameList from "@/components/games/list"
import MemberList from "@/components/member/list"
import PartnerList from "@/components/partner/list"
import { useTranslation } from "@/app/i18n/client"

export default function EditAgentSettings({ params: { lng,id } }: { params: { lng: string,id: string } }) {
  const {t} =  useTranslation(lng,'translation','')
  return (
    <div className="space-y-6 md:container md:mx-auto md:px-4">
      <Tabs defaultValue="account" className="w-full h-auto md:h-full">
        <TabsList className="w-full flex-wrap justify-start">
          <TabsTrigger value="account" className="flex-grow md:flex-grow-0">{t('promotion.account')}</TabsTrigger>
          <TabsTrigger value="games" className="flex-grow md:flex-grow-0">{t('games.title')}</TabsTrigger>
          <TabsTrigger value="promotion" className="flex-grow md:flex-grow-0">{t('promotion.title')}</TabsTrigger>
          <TabsTrigger value="member" className="flex-grow md:flex-grow-0">{t('member.title')}</TabsTrigger>
          <TabsTrigger value="partner" className="flex-grow md:flex-grow-0">{t('partner.title')}</TabsTrigger>
        </TabsList>
        <TabsContent value="account" className="h-auto md:h-[calc(100vh-120px)]">
          <ProfileEdit lng={lng} id={id} />
        </TabsContent>
        <TabsContent value="games">
          <GameList prefix={id} lng={lng} />
        </TabsContent>
        <TabsContent value="promotion">
          <PromotionList prefix={id} lng={lng} />
        </TabsContent>
        <TabsContent value="member">
          <MemberList prefix={id} lng={lng} />
        </TabsContent>
        <TabsContent value="partner">
          <PartnerList prefix={id} lng={lng} />
        </TabsContent>
      </Tabs>
      {/* <SettingsLayout>
        <ProfileEdit lng={lng} id={id} />
      </SettingsLayout> */}
    </div>
  )
}
