import { z } from "zod";

export const gameStatusSchema = z.object({
    id: z.string(),
    name: z.string(),
});

export type GameStatus = z.infer<typeof gameStatusSchema>;