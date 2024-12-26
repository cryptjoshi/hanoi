import { SessionOptions } from "iron-session";

export interface SessionData {
  userId?: string;
  username?: string;
  prefix?: string;
  isLoggedIn: boolean;
  token?:string;
  customerCurrency?:string;
  lng?:string;
}

export const defaultSession: SessionData = {
  isLoggedIn: false,
};

export const sessionOptions: SessionOptions = {
  // You need to create a secret key at least 32 characters long.
  password: process.env.PASSWORD_SECRET!,
  cookieName: "zookeep-session",
  cookieOptions: {
    httpOnly: true,
    // Secure only works in `https` environments. So if the environment is `https`, it'll return true.
    secure: process.env.NODE_ENV === "production",
  },
};