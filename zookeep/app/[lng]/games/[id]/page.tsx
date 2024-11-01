import { useTranslation } from "@/app/i18n"
import { ContentLayout } from "@/components/admin-panel/content-layout"
import GameList from "@/components/games/Gamelist"


export default async function GamePage({ params: { lng,id } }: { params: { lng: string,id:string } }) {

    const {t} = await useTranslation(lng,"translation",undefined)

    return (
        <ContentLayout title="Games">
        <div>
            <h1>{`${t("menu.games")} No. ${id}`}</h1>
             <GameList lng={lng} id={id}/>
        </div>
        </ContentLayout>
    )
}