import { z } from "zod";

export const partnerSchema = z.object({
      ID:z.number(),
      //Walletid:z.number().optional(),       
      Username:z.string().optional(),    
      Password:z.string().optional(),    
      //ProviderPassword:z.string(),    
      Fullname:z.string().optional(),    
      Bankname:z.string().optional(),    
      Banknumber:z.string().optional(),    
      //Balance:z.number().optional(),    
      //Beforebalance:z.number(),    
      //Token:z.string(),    
      //Role:z.string(),    
      //Salt:z.string(),    
      Status:z.number().optional(),
      //MinTurnoverDef:z.string().default('10%').optional(),    
      //Betamount:z.number(),    
      //Win:z.number(),    
      //Lose:z.number(),    
      //Turnover:z.number().optional(),    
      //ProID:z.string(),    
      //PartnersKey:z.string(),    
      // ProStatus:z.string().optional(),   
      ReferralCode:z.string().optional(), 
      RefferedCode:z.string().optional(),
      //ProActive:z.string()

});

export type Partner = z.infer<typeof partnerSchema>;