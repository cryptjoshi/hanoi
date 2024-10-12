import { Separator } from "@/components/ui/separator"
import { ProfileForm } from "@/app/forms/profile-form"

export default function SettingsProfilePage() {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium">สมัครสมาชิกใหม่</h3>
        <p className="text-sm text-muted-foreground">
         รายละเอียด
        </p>
      </div>
      <Separator />
      <ProfileForm />
    </div>
  )
}
