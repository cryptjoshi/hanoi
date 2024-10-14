import { Separator } from "@/components/ui/separator"
import { ProfileEdit } from "@/components/agents/edit/profile-edit"
import SettingsLayout from "@/components/agents/edit/layout"

export default function EditAgentSettings({ params: { lng,id } }: { params: { lng: string,id: string } }) {
  return (
    <div className="space-y-6">
      
      <SettingsLayout>
        <ProfileEdit lng={lng} id={id} />
      </SettingsLayout>
    </div>
  )
}
