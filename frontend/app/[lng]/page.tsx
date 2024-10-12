// import { Separator } from "@/components/ui/separator"
// import { ProfileForm } from "@/app/forms/profile-form"

import { redirect } from 'next/navigation';
export default function Page({ params: {
    lng
  }}) {
  redirect(`/${lng}/dashboard`);
}
