import { AreaGraph } from './area-graph';
import { BarGraph } from './bar-graph';
import { PieGraph } from './pie-graph';
import { useEffect, useRef, useState } from 'react';
import { useRouter } from 'next/navigation'
import { CalendarDateRangePicker } from '@/components/date-range-picker';
import {CalendarDatePicker} from '@/components/date-picker'
import PageContainer from '@/components/layout/page-container';
import { RecentSales } from './recent-sales';
import { Button } from '@/components/ui/button';
import { toast } from "@/hooks/use-toast"
import { cn } from "@/lib/utils"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle
} from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useTranslation } from '@/app/i18n/client';
import {GetOverview} from '@/actions'
import useAuthStore from '@/store/auth'
import { TZDate } from "@date-fns/tz";
import { formatNumber } from '@/lib/utils'
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { z } from "zod"
import { io } from 'socket.io-client'; 
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"


export interface IProvider {
    Name: string;
    Total: string;  // ใช้ string สำหรับ decimal
}

export interface IFinancialStats {
    allmembers: string;      // ใช้ string สำหรับ decimal
    newcomer: string;        // ใช้ string สำหรับ decimal
    firstdept: string;       // ใช้ string สำหรับ decimal
    deposit: string;         // ใช้ string สำหรับ decimal
    withdrawl: string;       // ใช้ string สำหรับ decimal
    totaldeposit: string;         // ใช้ string สำหรับ decimal
    totalwithdrawl: string;       // ใช้ string สำหรับ decimal
    winlose: string;         // ใช้ string สำหรับ decimal
    totalprofit: string;     // ใช้ string สำหรับ decimal
    provider: IProviderp[];      // ใช้ interface IProvider
}


const FormSchema = z.object({
  startdate: z.date({
    required_error: "A date of Start is required.",
  }),
})


export default function OverViewPage({lng}:{lng:string}) {
const tzDate = new TZDate(new Date(), "Asia/Bangkok");
 const { t } =  useTranslation(lng,'dashboard' ,undefined);
 const [date, setDate] = useState<Date | undefined>()
 const [isLoading, setIsLoading] = useState(true)
 const [overview,setOverView] = useState<IFinancialStats>({})
 const [refreshTrigger, setRefreshTrigger] = useState(0);
 const { customerCurrency,accessToken } = useAuthStore();
 const [connected,setConncted] = useState(false)
 const router = useRouter()
 const [messages, setMessages] = useState([]);
 const [isBlinking, setIsBlinking] = useState(false);
 const socketRef = useRef(null);
 const [socketid,setSocketId] = useState("");
 const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
    defaultValues: {
        startdate: tzDate, // ตั้งค่าเริ่มต้นให้กับ startdate
    },
  })
  const connectWebSocket = () => {
    socketRef.current = io('https://report.tsxbet.net', {
      path: '/socket.io',
      transports: ['websocket'], // บังคับให้ใช้ WebSocket
      upgrade: false, // ป้องกันการ upgrade protocol
      reconnection: true,
      reconnectionAttempts: 5,
      reconnectionDelay: 1000
  });

    
    socketRef.current.on('connect', () => {
      console.log('Socket.IO is connected.');
      setSocketId(socketRef.current.id);
      setConncted(true)
      console.log('Socket ID:', socketRef.current.id);

  });

  socketRef.current.on('message', (newMessage) => {
      setMessages((prevMessages) => [...prevMessages, newMessage]);
      setIsBlinking(true);
      // ตั้งเวลาให้หยุดกระพริบหลังจาก 3 วินาที
      setTimeout(() => {
        setIsBlinking(false);
        fetchGames();
      }, 3000);

      console.log('Received message:', newMessage);
  });

  socketRef.current.on('disconnect', () => {
      console.log('Socket.IO is disconnected. Reconnecting...');
      setConncted(false)
      setTimeout(connectWebSocket, 5000); // reconnect after 5 seconds
  });

  socketRef.current.on('error', (error) => {
      console.error('Socket.IO error observed:', error);
  });
};
const fetchGames = async () => {
      
  setIsLoading(true);
  try {
  
    if(accessToken){
        //console.log(form.getValues("startdate").toLocaleDateString())
    const fetchedGames = await GetOverview(accessToken,form.getValues("startdate").toLocaleDateString());
    
    setOverView(fetchedGames.Data);
    console.log(fetchedGames)
    } else {
      router.replace(`/${lng}/login`)
    }
  } catch (error) {
    console.error('Error fetching games:', error);
  } finally {
    setIsLoading(false);
  }
};

  useEffect(() => {
    // สร้าง WebSocket connection
    connectWebSocket();

    // Cleanup function to close WebSocket connection
    return () => {
        socketRef.current.close();
    };
}, []);

   useEffect(() => {
    
    fetchGames();
  }, [ refreshTrigger])


 function onSubmit(data: z.infer<typeof FormSchema>) {
    setIsLoading(true);
    GetOverview(accessToken,data.startdate.toLocaleDateString()).then(response=>{
        if(response.Status){
            setOverView(response.Data)
        } else {
            toast({
                title: t("common.fetch.error"),
                description: t("common.fetch.error_description"),
                variant: "destructive",
              })
        }
    setIsLoading(false);
    })
   
  }

//   const sendMessage = () => {
//     if (socketRef.current) {
//         const messageData = {
//             id: socketid,
//             message: "Hello Redis from Client!"
//         };
//         socketRef.current.emit('demo.hello', messageData);
//         setMessages((prevMessages) => [...prevMessages, messageData.message]);
//     } else {
//         console.error('Socket.IO is not connected.');
//     }
// }

  return (
   
    <PageContainer scrollable>
      <div className="space-y-4 ">
      <Card className="bg-[#CFE2F3] text-black">
            <CardContent className="item-center pt-4">
         <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} >
        {/* แถบการค้นหา */}
            <FormField
          control={form.control}
          name="startdate"
          render={({ field }) => (
            <FormItem className="flex flex-col">
             
        <div className="flex items-center justify-between">
        
          <div className="flex items-center space-x-2">
            <CalendarDatePicker lng={lng}  
            onChange={(value) => { console.log(value); field.onChange(value)} }
            initialDate={date}/>
            <Button type="submit">{t('button.refresh')}</Button>
          </div>
          <div className="flex items-center gap-2">
          <div
            className={cn(
              "h-3 w-3 rounded-full transition-all duration-300",
              connected ? "bg-green-500" : "bg-red-500",
              isBlinking && "animate-pulse"
            )}
          />
         <span className={cn(
          "transition-opacity",
          isBlinking && "animate-pulse"
        )}>{connected ? "Online" : "Offline"}</span>
        </div>
        </div>

        </FormItem>
          )}
          />
            
          </form>
          </Form>
          </CardContent>
            {/* <CardFooter> 
            <button onClick={sendMessage}>Send Message</button>
            </CardFooter> */}
            </Card>

               
{/* เปรียบเทียบกับเดือนที่ผ่านมา */}
<div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Card className="bg-blue-600 text-white">
            <CardHeader>
              <CardTitle>สมาชิก</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.allmembers?.toString()), 2)}</div>
              <p className="text-md text-white">± 100%</p>
            </CardContent>
          </Card>
          <Card className="bg-blue-600 text-white">
            <CardHeader>
              <CardTitle>ฝากครั้งแรก</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.firstdept?.toString()), 2)}</div>
              <p className="text-md">± 0.00%</p>
            </CardContent>
          </Card>
          <Card className="bg-blue-600 text-white">
            <CardHeader>
              <CardTitle>ลูกค้าที่ถอน</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.withdrawl?.toString()), 2)}</div>
              <p className="text-md">± 0.00%</p>
            </CardContent>
          </Card>
          <Card className="bg-blue-600 text-white">
            <CardHeader>
              <CardTitle>ลูกค้าที่สมัคร</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.newcomer?.toString()), 2)}</div>
              <p className="text-md">± 0.00%</p>
            </CardContent>
          </Card>
          <Card className="bg-blue-600 text-white">
            <CardHeader>
              <CardTitle>ฝาก</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-green-500">{formatNumber(parseFloat(overview?.totaldeposit?.toString()), 2)}</div>
              <p className="text-md text-green-500">↑ 1,883.61%</p>
            </CardContent>
          </Card>
          <Card className="bg-blue-600 text-white">
            <CardHeader>
              <CardTitle>ถอน</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-red-500">-{formatNumber(parseFloat(overview?.totalwithdrawl?.toString()), 2)}</div>
              <p className="text-md text-green-500">↑ 376.67%</p>
            </CardContent>
          </Card>
          <Card className="bg-blue-600 text-white">
            <CardHeader>
              <CardTitle>ได้เสีย (WL)</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.winlose?.toString()), 2)}</div>
              <p className="text-md text-green-500">↑ 1,353.11%</p>
            </CardContent>
          </Card>
          <Card className="bg-blue-600 text-white">
            <CardHeader>
              <CardTitle>รวมรายได้</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.totalprofit?.toString()), 2)}</div>
              <p className="text-md text-green-500">↑ 2,326.99%</p>
            </CardContent>
          </Card>
        </div>
  

        {/* แสดงข้อมูลเพิ่มเติม */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Card>
            <CardHeader>
              <CardTitle>ข้อมูลการทำธุรกรรม</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.allmembers?.toString()), 2)}</div>
                  <p className="text-xs">สมาชิก</p>
                </div>
                <div>
                  <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.firstdept?.toString()), 2)}</div>
                  <p className="text-xs">ฝากครั้งแรก</p>
                </div>
                <div>
                  <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.withdrawl?.toString()), 2)}</div>
                  <p className="text-xs">ลูกค้าที่ถอน</p>
                </div>
                <div>
                  <div className="text-2xl font-bold">{formatNumber(parseFloat(overview?.newcomer?.toString()), 2)}</div>
                  <p className="text-xs">ลูกค้าที่สมัคร</p>
                </div>
              </div>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle>ผู้ให้บริการ</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1">
               {overview?.Provider?.map((item, index) => (
                <div key={index} className="flex justify-between">
                    <span>{item.name}</span>
                    <span>{formatNumber(parseFloat(item.total.toString()), 2)}</span>
                </div>
            ))}
              </div>
            </CardContent>
          </Card>
        </div>
        {/* แสดงกราฟ  
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="col-span-1">
            <BarGraph />
          </div>
          <div className="col-span-1">
            <PieGraph />
          </div>
        </div>
*/}
 

      </div>
    </PageContainer>
 
  );
}