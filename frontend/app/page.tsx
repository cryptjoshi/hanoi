// import { Button } from "@/components/ui/button"
// import Login from "@/components/login"
import Link from "next/link";
import Image from "next/image";
import { PanelsTopLeft } from "lucide-react";
import { ArrowRightIcon, GitHubLogoIcon } from "@radix-ui/react-icons";

import { Button } from "@/components/ui/button";
import { ModeToggle } from "@/components/mode-toggle";
import { useTranslation } from '@/app/i18n'

export default async function Home({ params: { lng } }) {
    const { t } = await useTranslation(lng)
    return (
        <></>
            );
        }
