import { z } from "zod";

export const gameSchema = z.object({
    name: z.string().optional(),
    productCode: z.string().optional(),
    product: z.string().optional(),
    gameType: z.string().optional(),
    active: z.number().optional(),
    remark: z.string().optional(),
    status: z.string().optional(),
});

export type Game = z.infer<typeof gameSchema>;