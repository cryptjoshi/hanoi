'use client'
import { Navbar } from "@/components/admin-panel/navbar";
import React from 'react'
import type { ReactElement } from 'react'
import Footer from "../footer";
import useAuthStore from "@/store/auth";
import { getSession } from "@/actions";


interface ContentLayoutProps {
  title: string;
  children: React.ReactNode;
}
 


export function ContentLayout({ title, children }: ContentLayoutProps):ReactElement {
 //const {lng} = useAuthStore()
const [lng,setLng] = React.useState("")

//  React.useEffect(()=>{
//   const fetchSession = async ()=>{
//   const session = await getSession()
//  // console.log(session)
//   //setLng(session.lng)
//   }
//   fetchSession()
//  },[lng])

 return (
    <div className="min-h-screen flex flex-col">
      <Navbar title={title} />
      <div className="container flex-grow pt-8 pb-8 px-4 sm:px-8">{children}</div>
      <Footer />
    </div>
  );
}
