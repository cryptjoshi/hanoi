"use client"
import useState from "react"
import { Separator } from "@/components/ui/separator"
import { ProfileEdit } from "@/components/agents/edit/profile-edit"
import SettingsLayout from "@/components/agents/layout"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import PromotionList from "@/components/promotion/list"
import GameList from "@/components/games/list"
import MemberList from "@/components/member/list"
import PartnerList from "@/components/partner/list"
import AgentSettings from "@/components/settings/agentsetting"
import { useTranslation } from "@/app/i18n/client"
import { iMember } from "@/components/member/list"
import { DataTableProps } from "@/components/member/list" 

import Profile from "../partner/Profile"

export default function EditMemberSettings({ params: { lng,id } }: { params: { lng: string,id: string } }) {
  const {t} =  useTranslation(lng,'translation','')
  
  //const [data, setData] = useState<DataTableProps<iMember>>({}); 
  const data: DataTableProps<iMember> = []; 
  const closeEditPanel = () => {
    // setParnerId(null);
    // setIsAddingGame(false);
    // setShowTable(true);
    // setRefreshTrigger((prev:any) => prev + 1);
  };

  return (
    <div className="space-y-6 md:container md:mx-auto md:px-4">
      <Tabs defaultValue="account" className="w-full h-auto md:h-full">
        <TabsList className="w-full flex-wrap justify-start">
          <TabsTrigger value="account" className="flex-grow md:flex-grow-0">{t('promotion.account')}</TabsTrigger>
          <TabsTrigger value="member" className="flex-grow md:flex-grow-0">{t('member.title')}</TabsTrigger>
          <TabsTrigger value="settings" className="flex-grow md:flex-grow-0" disabled>{t('settings.title')}</TabsTrigger>
        </TabsList>
        <TabsContent value="account">
        <Profile
          //  partnerId={}
            isAdd={false}
            lng={lng}
           
          />
        </TabsContent>
        <TabsContent value="member">
          <MemberList prefix={id} data={data} lng={lng} />
        </TabsContent>
        
        <TabsContent value="settings">
          <AgentSettings lng={lng} id={id} />
        </TabsContent>
      </Tabs>
      {/* <SettingsLayout>
        <ProfileEdit lng={lng} id={id} />
      </SettingsLayout> */}
    </div>
  )
}
 

