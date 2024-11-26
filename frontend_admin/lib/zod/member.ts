import { z } from "zod";

export const memberSchema = z.object({
      ID:z.number(),
      Walletid:z.number(),       
      Username:z.string(),    
      Password:z.string(),    
      //ProviderPassword:z.string(),    
      Fullname:z.string(),    
      Bankname:z.string(),    
      Banknumber:z.string(),    
      //Balance:z.number().optional(),    
      //Beforebalance:z.number(),    
      //Token:z.string(),    
      //Role:z.string(),    
      //Salt:z.string(),    
      Status:z.number(),
      MinTurnoverDef:z.string().default('10%'),    
      //Betamount:z.number(),    
      //Win:z.number(),    
      //Lose:z.number(),    
      //Turnover:z.number().optional(),    
      //ProID:z.string(),    
      //PartnersKey:z.string(),    
      ProStatus:z.string().optional(),   
      ReferralCode:z.string().optional(), 
      RefferedCode:z.string().optional(),
      //ProActive:z.string()

});

export type Member = z.infer<typeof memberSchema>;