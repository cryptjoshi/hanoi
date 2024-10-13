import { ProfileForm } from "@/app/forms/profile-form";

export default function SomePage({ params: { lng } }: { params: { lng: string } }) {
  return <ProfileForm lng={lng} />;
}
