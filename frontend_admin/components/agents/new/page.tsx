import { Separator } from "@/components/ui/separator"
import { ProfileRegister } from "@/components/agents/new/profile-register"

export default function SettingsProfilePage({ params: { lng,id } }: { params: { lng: string,id: string } }) {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium">สมัครสมาชิกใหม่</h3>
        <p className="text-sm text-muted-foreground">
         รายละเอียด
        </p>
      </div>
      <Separator />
      <ProfileRegister lng={lng} id={id} />
    </div>
  )
}
