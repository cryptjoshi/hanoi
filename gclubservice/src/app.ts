
//@ts-nocheck

import express from "express"
import crypto from 'crypto';
import expressPino from "express-pino-logger";
import pino from "pino";
import moment from "moment"
import fetch from "cross-fetch"
import jwt from "jsonwebtoken"
import cors from "cors"
import bodyParser from 'body-parser'
import request from 'request'
var CryptoJS = require("crypto-js");


const app = express()
const CONFIG = {
  CLIENT_ID: "6342e1be-fa03-456f-8d2d-8e1c9513c351",
  CLIENT_SECRET:	"6d83ac42",
  DESKEY: "9c62a148",
  DESIV:	"8e014099",
  SYSTEMCODE: "ckthb01",
  WEBID: "ckdthb",
  API_URL_G: "http://rcgapiv2.rcg666.com/",
  API_URL_PROXY: "http://api.tsxbet.info:8001"
}
const levels = {
  http: 10,
  debug: 20,
  info: 30,
  warn: 40,
  error: 50,
  fatal: 60,
};
const logger = pino({
  customLevels: levels, // our defined levels
  useOnlyCustomLevels: true,
  level: 'http',
  formatters: {
      level: (label) => {
        return { level: label.toUpperCase() };
      },
    },
    timestamp: () => `,"timestamp":"${moment().format('MMMM Do YYYY, h:mm:ss a')}"`
});

const logRequest = (enabled: boolean) =>
expressPino({
  level: "info",
  enabled,
});
const gencryption =  (data:string)=> {
 

  const algorithm = "des-cbc"
  const cipher = crypto.createCipheriv(algorithm, Buffer.from(CONFIG.DESKEY),Buffer.from(CONFIG.DESIV));
  cipher.setAutoPadding(true);
  const hash_data = cipher.update(data,"utf8","base64") + cipher.final("base64")

     return hash_data
}
const ghash = (data:string,unixtime:string) =>{
   
  const source = CONFIG.CLIENT_ID+CONFIG.CLIENT_SECRET+unixtime+data;
  const hasstr = crypto.createHash('md5').update(source).digest("hex").toString();
return {source,hasstr}
}
const gdecryption =  (data:string)=> {
 

  const algorithm = "des-cbc"
  const decipher = crypto.createDecipheriv(algorithm, CONFIG.DESKEY, CONFIG.DESIV).setAutoPadding(true);
  const hash_data = decipher.update(data,"base64","utf8") + decipher.final("utf8")

     return hash_data
}
const encryption = (data: string) => {
  try {
    let unx = moment().valueOf();
    let des = CryptoJS.TripleDES.encrypt(data, CryptoJS.enc.Utf8.parse(CONFIG.DESKEY), {
        iv: CryptoJS.enc.Utf8.parse(CONFIG.DESIV),
        mode: CryptoJS.mode.CBC
    }).toString();
    let md5 = CryptoJS.MD5(`${CONFIG.CLIENT_ID}${CONFIG.CLIENT_SECRET}${unx}${des}`);
    let enc = CryptoJS.enc.Base64.stringify(md5);
    return { enc, des, unx }
  }

catch(err){
  logger.error(err)
}
} 
const decryption = (data:string) => {
    let decipher = crypto.createDecipheriv('des-cbc', CONFIG.DESKEY, CONFIG.DESIV);
    let decrypted = decipher.update(data, 'base64', 'utf8');
    decrypted += decipher.final('utf8');
    
    return decrypted
  }

 const login = (account: string) => {
    return new Promise((resolve,reject) => {
     
        const JSONString = JSON.stringify({
            "SystemCode": CONFIG.SYSTEMCODE,
            "WebId": CONFIG.WEBID,
            "MemberAccount": account,
            "ItemNo": "1",
            "BackUrl": "https://tsx.bet/",
            "GroupLimitID": "1,4,12",
            "Lang": "th-TH",
        })
      
        const encrypt: any = encryption(JSONString);
        
        logger.info(encrypt)
        request({
            'method': 'POST',
            'url': `${CONFIG.API_URL_PROXY}/api/Player/Login`,
            'headers': {
                'X-API-ClientID': CONFIG.CLIENT_ID,
                'X-API-Signature': encrypt.enc,
                'X-API-Timestamp': encrypt.unx,
                'Content-Type': 'application/json'
            },
            body: encodeURIComponent(encrypt.des)

        }, function (error, response) {
            if (error) throw new Error(error);
            let result = decryption(response.body)
            resolve(JSON.parse(result))
            //resolve(response.body)
        });
    })
}
const CreateOrUpdate =  (account: string, name: string) => {
  return new Promise((resolve,reject) => {
 const JSONString = JSON.stringify({
     "SystemCode": CONFIG.SYSTEMCODE,
     "WebId": CONFIG.WEBID,
     "MemberAccount": account,
     "MemberName": name,
     "StopBalance": -1,
     "BetLimitGroup": "1,4,12",
     "Currency": "THB",
     "Language": "th-TH",
     "OpenGameList": "ALL"
 })

  const encrypt: any = encryption(JSONString);
 try {

 
  request({
  'method': 'POST',
  'url': `${CONFIG.API_URL_PROXY}/api/Player/CreateOrSetUser`,
  'headers': {
      'X-API-ClientID': CONFIG.CLIENT_ID,
      'X-API-Signature': encrypt.enc,
      'X-API-Timestamp': encrypt.unx,
      'Content-Type': 'application/json'
  },
  body: encodeURIComponent(encrypt.des)

  }, function (error, response) {
      if (error) throw new Error(error);
      let result = decryption(response.body)
      resolve(JSON.parse(result))
      //return result;
     // console.log(result, response.statusCode);
  });
}
catch(error){
  if (error) throw new Error(error);
}
})
}


const testUser = (req, res) => { 
  const myresponse:ResponseUser = {status:true,message:""}
  const data =   `{"SystemCode":"116688tsxbet","WebId":"tsxdemo","MaxId":0,"Rows":100}`
    const utf8EncodeText = new TextEncoder();

    const str = JSON.stringify(data);

   const result =  gencrypt(str)//encrypt(str,CONFIG.DESKEY,CONFIG.DESIV) //gencrypt(message)// encryptWithDES(message,key,iv)
 
    const {source,hasstr} = ghash(result)
    myresponse.message = { "DES": result, "source": source,"hash": hasstr, "base64": Buffer.from(hasstr).toString('base64'),"url_encode": encodeURIComponent(result)};

  return res.status(200).json(myresponse);
}

const launchGame = async (req:Request,res:Response) => {
  const myresponse:UserResponseX = {status:false, data:{}}

   try {

       let response:any = await  login(req.body.account,req.body.encrypt);
      // const userx = await getBnbUser(req.body.username.toUpperCase())
       logger.info(response);
       if(response.message=="OK")
       { 
          //  const [users, created] = await Users.findOrCreate({where:{username:req.body.username.toUpperCase()}}).catch(err =>logger.error(err))
          //  users.token = login.data.token;
          //  users.wallet_id = userx.id
          //  users.balance = (1*userx.balance)
          //  users.save()
           myresponse.status= true
           myresponse.data = response.data

          
       }else 
           if(response.Message=='MEMBER_NOT_EXISTS'){
                
               let response = await  CreateOrUpdate(req.body.account,req.body.account);
               // myesponse.status = true
               // myresponse.data= response.data
               response = await  login(req.body.account);
               
              // const [users, created] = await Users.findOrCreate({where:{username:req.body.username.toUpperCase()}}).catch(err =>logger.error(err))
               myresponse.status = true
               myresponse.data= response.data
              //  users.token = login.data.token;
              //  users.wallet_id = userx.id
              //  users.balance = (1*userx.balance)
              //  users.save()
               
           }
           //const [users,created] = await Users.findOrCreate({where:{username:req.body.username.toUpperCase()}})
          
           //const [user, created] = await Users.findOrCreate({where:{username:req.body.MemberName}})
           //logger.debug(created)
           //logger.info(users)
         
       } catch (e:any) {
           logger.error(e);
           myresponse.status = false;
           myresponse.data = e;
           
       }
       //logger.info(myresponse)
return res.status(200).json(myresponse)
}

var compression = require('compression')

const shouldCompress = (req, res) => {
    if (req.headers["x-no-compression"]) {
      return false;
    }
    return compression.filter(req, res);
  };

const port: number = 9003
const corsOptions = {
    origin: "*"
}
 


app.use(cors(corsOptions))
app.use(express.json())
app.use(compression({ filter: shouldCompress }));

app.use(bodyParser.urlencoded({ extended: true }))
app.use(bodyParser.json())

app.use(logRequest(false));

app.get("/",(req,res)=>{

  res.send("Hello I am G Service!")
})

app.post("/encryption",(req,res)=>{

  const {data,key,iv} = req.body

  const encrypted = encryption(data,key,iv)

  res.json({"data":encrypted})
})

app.post("/decryption",(req,res)=>{

  const {data,key,iv} = req.body

  const encrypted = decryption(data,key,iv)

  res.json({"data":encrypted})
})
app.post("/test" ,testUser)
app.post("/LaunchGame" ,launchGame)

app.post("/login",async (req,res)=>{
  const {account} = req.body

  let result = await login(account)
 
  // if(result.message=="MEMBER_NOT_EXISTS"){
  //   const created = await CreateOrUpdate(account,account)
  //   result = await login(account)
  // }
  res.json(result)

})

app.listen(port,()=>{
    console.log(` GService running on ${port}`)
})