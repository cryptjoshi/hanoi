import Login from "@/components/authen/login"
// import { fallbackLng, languages } from '@/app/i18n/settings';
// import { useTranslation } from '@/app/i18n';

// export async function generateStaticParams() {
//   return languages.map((lng) => ({ lng }));
// }
// export async function generateStaticParams() {
//    return languages.map((lng) => ({ lng }));
//  }
 
export default async function LoginPage({ params }: { params: { lng: string } }) {
 // const { t } = await useTranslation(params.lng,'login');
 return (
    <Login lng={params.lng} />
 )
}
 