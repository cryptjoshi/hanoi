import Link from "next/link";
import Image from "next/image";

import { Card, CardContent } from "@/components/ui/card";
interface PlaceholderContentProps {
  children: React.ReactNode;
  fullWidthOnLargeScreen?: boolean;
}
export default function PlaceholderContent({children, fullWidthOnLargeScreen = false}: PlaceholderContentProps) {
  const widthClass = fullWidthOnLargeScreen ? 'w-full' : 'w-full md:w-auto';
  return (
    <Card className={`${widthClass} rounded-lg border-none mt-6`}>
      {/* <CardContent className="p-6">
        <div className="flex justify-start items-center min-h-[calc(100vh-56px-64px-20px-24px-56px-48px)]">
          <div className="flex flex-col relative w-full">
          {children}
          </div>
        </div>
      </CardContent> */}
      <CardContent className="p-6">
        {children}
      </CardContent>
    </Card>
  );
}
