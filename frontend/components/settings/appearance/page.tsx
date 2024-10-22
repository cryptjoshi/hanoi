import { Separator } from "@/components/ui/separator"
import { AppearanceForm } from "@/app/forms/appearance/appearance-form"

export default function SettingsAppearancePage({ params: { lng,prefix } }: { params: { lng: string,prefix: string } }) {
  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium">Appearance</h3>
        <p className="text-sm text-muted-foreground">
          Customize the appearance of the app. Automatically switch between day
          and night themes.
        </p>
      </div>
      <Separator />
      <AppearanceForm lng={lng} prefix={prefix} />
    </div>
  )
}
