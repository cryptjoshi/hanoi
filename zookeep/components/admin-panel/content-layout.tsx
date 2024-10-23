import { Navbar } from "@/components/admin-panel/navbar";
import React from 'react'
import type { ReactElement } from 'react'

interface ContentLayoutProps {
  title: string;
  children: React.ReactNode;
}
 

export function ContentLayout({ title, children }: ContentLayoutProps):ReactElement {
 
  return (
    <div  >
      <Navbar title={title} />
      <div className="container pt-8 pb-8 px-4 sm:px-8">{children}</div>
    </div>
  );
}
