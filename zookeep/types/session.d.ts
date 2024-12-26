import {User,AuthStore} from "@/store/auth"

declare module "next-iron-session" {
    interface IronSessionData {
        user?: User;
        auth?:AuthStore;
    }
}