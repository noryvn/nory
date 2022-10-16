export * as prompt from "https://deno.land/x/cliffy@v0.25.2/prompt/mod.ts";
import { createClient } from "https://esm.sh/@supabase/supabase-js@2";

export const supabase = createClient(
  mustGetEnv("SUPABASE_URL"),
  mustGetEnv("SUPABASE_KEY"),
);

export function mustGetEnv(name: string) {
  const val = Deno.env.get(name);
  if (!name) {
    const msg = `missing "${name}" in environment variable`;
    console.error(msg);
    throw new Error(msg);
  }
  return val as string;
}
