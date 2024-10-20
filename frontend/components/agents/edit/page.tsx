"use client"
import { Separator } from "@/components/ui/separator"
import { ProfileEdit } from "@/components/agents/edit/profile-edit"
import SettingsLayout from "@/components/agents/layout"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import PromotionList from "@/components/promotion/list"
import GameList from "@/components/games/list"
import MemberList from "@/components/member/list"
import { useTranslation } from "@/app/i18n/client"

export default function EditAgentSettings({ params: { lng,id } }: { params: { lng: string,id: string } }) {
  const {t} =  useTranslation(lng,'translation','')
  return (
    <div className="space-y-6">
      <Tabs defaultValue="account" className="w-full">
  <TabsList>
    <TabsTrigger value="account">{t('promotion.account')}</TabsTrigger>
    <TabsTrigger value="games">{t('games.title')}</TabsTrigger>
    <TabsTrigger value="promotion">{t('promotion.title')}</TabsTrigger>
    <TabsTrigger value="member">{t('member.title')}</TabsTrigger>
  </TabsList>
  <TabsContent value="account">
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
</Tabs>
      {/* <SettingsLayout>
        <ProfileEdit lng={lng} id={id} />
      </SettingsLayout> */}
    </div>
  )
}
