'use client'
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
  } from "@/components/ui/alert-dialog"
  import { MdOutlineError } from "react-icons/md";
type AlertDialogProps = {
    caption: string;
    title: string;
    description: string;
    canceltext:string;
    oktext:string;
    icon:string;

}

export default function AlertComponent (props: AlertDialogProps) {
    return (
    <AlertDialog>
        <AlertDialogTrigger>
        Open
           </AlertDialogTrigger>
            <AlertDialogContent>
                <AlertDialogHeader>
                <AlertDialogTitle> <MdOutlineError />
             {props.title}</AlertDialogTitle>
                <AlertDialogDescription>
                    {props.description}
                </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                <AlertDialogCancel>{props.canceltext}</AlertDialogCancel>
                <AlertDialogAction>{props.oktext}</AlertDialogAction>
            </AlertDialogFooter>
        </AlertDialogContent>
    </AlertDialog>
    )
}