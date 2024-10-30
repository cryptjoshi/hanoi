import Register from '@/components/authen/register'


export default  function RegisterPage({ params: { lng } }: { params: { lng: string } }) {
   

    return (
        <>
        <Register lng={lng} />
        </>
    )
}