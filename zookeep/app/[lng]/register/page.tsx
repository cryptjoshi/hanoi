import Register from '@/components/authen/register'
 

export default function RegisterPage({ params: { lng }, searchParams }: { params: { lng: string }, searchParams: { code?: string } }) {
    // ดึง referralCode จาก searchParams
    const { code } = searchParams;
 
    return (
        <>
        <Register lng={lng} refferedcode={code} />
        </>
    )
}