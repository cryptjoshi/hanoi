'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslation } from '@/app/i18n/client';
import Login from "@/components/authen/login";

export default function LoginPage({ params }: { params: { lng: string } }) {
  const router = useRouter();
  const { i18n } = useTranslation(params.lng, 'login', undefined);

  useEffect(() => {
 
    if (i18n.language !== params.lng) {
      i18n.changeLanguage(params.lng);
    }
  }, [i18n, params.lng]);
 
  return <Login lng={params.lng}  />;
}
