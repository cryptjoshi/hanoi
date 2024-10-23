'use client'
import { usePathname, useRouter } from 'next/navigation'
import { languages } from '@/app/i18n/settings'
import { useState, useEffect } from 'react'
import { TbLanguage } from 'react-icons/tb'

import { hasFlag } from 'country-flag-icons'
import * as CountryFlags from 'country-flag-icons/react/3x2'
import { Tooltip, TooltipProvider, TooltipTrigger } from './ui/tooltip'
import { Button } from "@/components/ui/button";
const languageToCountry: { [key: string]: string } = {
  en: 'US',
  th: 'TH',
  // เพิ่มการแมปภาษาเป็นรหัสประเทศสำหรับภาษาอื่นๆ ตามต้องการ
}

export default function LanguageSwitcher() {
  const pathname = usePathname()
  const router = useRouter()
  const [currentLang, setCurrentLang] = useState('')

  useEffect(() => {
    const lang = pathname?.split('/')[1] || languages[0]
    setCurrentLang(lang)
  }, [pathname])

  const handleLanguageChange = () => {
    const currentIndex = languages.indexOf(currentLang)
    const nextIndex = (currentIndex + 1) % languages.length
    const nextLang = languages[nextIndex]
    const newPathname = pathname!.replace(/^\/[^\/]+/, `/${nextLang}`)
    router.push(newPathname)
  }

  const countryCode = languageToCountry[currentLang]
  const FlagComponent = countryCode && hasFlag(countryCode) ? CountryFlags[countryCode as keyof typeof CountryFlags] : null

  return (
    <TooltipProvider disableHoverableContent>
    <Tooltip delayDuration={100}>
      <TooltipTrigger asChild>
    <Button 
      onClick={handleLanguageChange}
      className="h-5 w-5 sm:h-5 sm:w-5 mt-2 shadow"
      variant="ghost"
      size="icon"
      title={`Current language: ${currentLang.toUpperCase()}. Click to change.`}
    >
      {FlagComponent ? (
        <FlagComponent className="w-4 h-4" />
      ) : (
        <TbLanguage className="w-4 h-4" />
      )}
    </Button>
    </TooltipTrigger>
    </Tooltip>
    </TooltipProvider>
  )
}
