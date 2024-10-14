import { Separator } from "@/components/ui/separator"
import { ProfileEdit } from "@/components/agents/edit/profile-edit"
import SettingsLayout from "@/components/agents/layout"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import PromotionList from "@/components/promotion/list"

export default function EditAgentSettings({ params: { lng,id } }: { params: { lng: string,id: string } }) {
  return (
    <div className="space-y-6">
      <Tabs defaultValue="account" className="w-full">
  <TabsList>
    <TabsTrigger value="account">Account</TabsTrigger>
    <TabsTrigger value="password">Promotion</TabsTrigger>
  </TabsList>
  <TabsContent value="account">
  <ProfileEdit lng={lng} id={id} />
  </TabsContent>
  <TabsContent value="password"><PromotionList lng={lng} id={id} /></TabsContent>
</Tabs>

        
      {/* <SettingsLayout>
        <ProfileEdit lng={lng} id={id} />
      </SettingsLayout> */}
    </div>
  )
}
