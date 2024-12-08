import { z } from "zod";

export const partnerSchema = z.object({
      ID:z.number(),     
      Username:z.string(),    
      Password:z.string(),      
      Fullname:z.string(),    
      Bankname:z.string(),    
      Banknumber:z.string(),    
      Status:z.number(),  
      ProStatus:z.string(),   
      RefferalCode:z.string(), 
      RefferedCode:z.string(),
});

export type Partner = z.infer<typeof partnerSchema>;