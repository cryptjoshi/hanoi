"use client";
import { Button } from "@/components/ui/button";
import useAuthStore from "@/store/auth";
import  Image  from "next/image"
import { RxAvatar,RxBell,RxMagnifyingGlass,RxHome,RxListBullet    } from "react-icons/rx";

export default function Dashboar() {
  const { Logout } = useAuthStore();
  return (
    <>
       
        <div className="flex justify-between items-center bg-white p-4 shadow-md">
          <div className="flex items-center space-x-2">
       
            <RxAvatar   className="w-6 h-6 rounded-full" />
            <span>0.00000000 BTC</span>
          </div>
          <div className="flex space-x-4">
            <button>
              <RxMagnifyingGlass  className="w-6 h-6" />
            </button>
            <button>
              <RxBell 
                className="w-6 h-6"
              />
            </button>
          </div>
        </div>

        <div className="p-4 bg-white shadow-md">
          <div className="text-center">
            <h2 className="text-2xl font-semibold">0.00000000 BTC</h2>
            <p className="text-sm text-gray-500">~$0.00</p>
          </div>
          <div className="flex justify-center mt-4 space-x-4">
            <Button size={"lg"} className="bg-yellow-400 text-white px-4 py-2 rounded-md w-full">
              ฝาก
            </Button>
            <Button size={"lg"} className="bg-black text-white px-4 py-2 rounded-md w-full">
              ซื้อ
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-4 gap-4 p-4 bg-white shadow-md mt-4">
          <div className="text-center">
            <Button className="bg-red-500 text-white px-2 py-1 rounded">
              HOT
            </Button>
            <p className="mt-2">6th Anniversary</p>
          </div>
          <div className="text-center">
            <Image
              src=""
              alt="Copy Trading"
              className="w-10 h-10 mx-auto"
            />
            <p className="mt-2">Copy Trading</p>
          </div>
          <div className="text-center">
            <Image src="" alt="Spot" className="w-10 h-10 mx-auto" />
            <p className="mt-2">Spot</p>
          </div>
          <div className="text-center">
            <Image
              src=""
              alt="Futures Grid"
              className="w-10 h-10 mx-auto"
            />
            <p className="mt-2">Futures Grid</p>
          </div>
        </div>

        <div className="p-4 mt-4">
          <Image
            src=""
            alt="Promo Banner"
            className="w-full rounded-md"
          />
        </div>

        <div className="grid grid-cols-3 gap-4 p-4 bg-white shadow-md mt-4">
          <div className="text-center">
            <p className="text-lg font-bold">7D ROI</p>
            <p className="text-green-500 font-semibold">+203.99%</p>
          </div>
          <div className="text-center">
            <p className="text-lg font-bold">7D ROI</p>
            <p className="text-green-500 font-semibold">+101.79%</p>
          </div>
          <div className="text-center">
            <p className="text-lg font-bold">7D ROI</p>
            <p className="text-green-500 font-semibold">+57.17%</p>
          </div>
        </div>

        <div className="fixed bottom-0 left-0 right-0 bg-white shadow-md flex justify-around p-4">
          <Button size={"icon"} variant={"ghost"} className="flex flex-col items-center w-10 h-10">
            <RxHome className="w-6 h-6" />
            <span className="text-sm">Home</span>
          </Button>
          <Button  size={"icon"} variant={"ghost"} className="flex flex-col items-center w-10 h-10">
            <RxListBullet  className="w-6 h-6" />
            <span className="text-sm">Market</span>
          </Button>
          <Button size={"icon"} variant={"ghost"} className="flex flex-col items-center w-10 h-10">
            <Image src="" alt="Buy/Sell" className="w-6 h-6" />
            <span className="text-sm">ซื้อขาย</span>
          </Button>
          <Button size={"icon"} variant={"ghost"} className="flex flex-col items-center w-10 h-10">
            <Image src="" alt="Futures" className="w-6 h-6" />
            <span className="text-sm">ฟิวเจอร์ส</span>
          </Button>
          <Button size={"icon"} variant={"ghost"} className="flex flex-col items-center w-10 h-10">
            <Image src="" alt="Assets" className="w-6 h-6" />
            <span className="text-sm">ทรัพย์สิน</span>
          </Button>
          <Button size={"icon"} variant={"ghost"} className="flex flex-col items-center w-10 h-10"
            onClick={() => {
              Logout();
            }}
          >
            Logout
          </Button>
        </div>
       
    </>
  );
}
