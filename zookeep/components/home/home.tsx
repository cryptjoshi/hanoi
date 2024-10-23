'use client'
import React from 'react'
import type { ReactElement } from 'react'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { ChevronLeft, Menu, Home, DollarSign, Bell, Coin, Plus, Activity, CreditCard, PieChart, Settings } from "lucide-react"
import useAuthStore from '@/store/auth'


function HomePage(): ReactElement {
    const {user} = useAuthStore()
    
  return (
    <div className="max-w-md mx-auto p-4 bg-indigo-600 min-h-screen">
      <Card className="bg-white shadow-lg rounded-3xl">
        <CardHeader className="flex justify-between items-center">
          <Avatar className="h-12 w-12">
            <AvatarImage src="/placeholder-avatar.jpg" alt="Profile picture" />
            <AvatarFallback>N</AvatarFallback>
          </Avatar>
          <div className="text-right">
            <CardTitle className="text-lg">Nasim</CardTitle>
            <p className="text-sm text-muted-foreground">10 April, 2020</p>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="text-center">
            <p className="text-sm text-muted-foreground">Current Balance</p>
            <p className="text-4xl font-bold text-indigo-600">$4,239.98</p>
          </div>
          
          <div className="space-y-2">
            <div className="flex items-center justify-between bg-gray-100 p-3 rounded-lg">
              <div>
                <p className="font-semibold">Deposit</p>
                <p className="text-sm text-muted-foreground">10 April, 2020</p>
              </div>
              <div className="text-right">
                <p className="font-semibold text-green-600">$321.00</p>
                <p className="text-sm text-muted-foreground">6% of Current Balance</p>
              </div>
            </div>
            <div className="flex items-center justify-between bg-gray-100 p-3 rounded-lg">
              <div>
                <p className="font-semibold">Expenses</p>
                <p className="text-sm text-muted-foreground">10 April, 2020</p>
              </div>
              <div className="text-right">
                <p className="font-semibold text-red-600">$142.89</p>
                <p className="text-sm text-muted-foreground">3.5% of Current Balance</p>
              </div>
            </div>
          </div>

          <div>
            <h3 className="font-semibold mb-2">Send Money To</h3>
            <div className="flex space-x-4">
              {['Kathryn', 'Eleanor', 'Diana'].map((name) => (
                <div key={name} className="flex flex-col items-center">
                  <Avatar className="h-12 w-12 mb-1">
                    <AvatarFallback>{name[0]}</AvatarFallback>
                  </Avatar>
                  <p className="text-xs">{name}</p>
                </div>
              ))}
              <Button variant="outline" size="icon" className="h-12 w-12 rounded-full">
                <Plus className="h-6 w-6" />
              </Button>
            </div>
          </div>

          <div>
            <div className="flex justify-between items-center mb-2">
              <h3 className="font-semibold">Expenses</h3>
              <Button variant="ghost" size="sm">See All</Button>
            </div>
            <div className="flex items-center justify-between bg-gray-100 p-3 rounded-lg">
              <div className="flex items-center">
                <div className="bg-red-500 p-2 rounded-full mr-3">
                  <Activity className="h-4 w-4 text-white" />
                </div>
                <div>
                  <p className="font-semibold">Strava</p>
                  <p className="text-xs text-muted-foreground">10:05, 12 April 2020</p>
                </div>
              </div>
              <p className="font-semibold text-red-600">-$124.00</p>
            </div>
          </div>
        </CardContent>
      </Card>
      
      <div className="fixed bottom-0 left-0 right-0 bg-white p-4">
        <div className="flex justify-around">
          <Button variant="ghost" size="icon">
            <Home className="h-6 w-6" />
          </Button>
          <Button variant="ghost" size="icon">
            <CreditCard className="h-6 w-6" />
          </Button>
          <Button variant="ghost" size="icon">
            <PieChart className="h-6 w-6" />
          </Button>
          <Button variant="ghost" size="icon">
            <Settings className="h-6 w-6" />
          </Button>
        </div>
      </div>
    </div>
  )
}


export default HomePage
 
