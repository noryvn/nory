import { prompt, supabase } from "./deps.ts";

const email = await prompt.Input.prompt({
  message: "Email: ",
});
const password = await prompt.Input.prompt({
  message: "Password: ",
});

const { data, error } = await supabase.auth.signInWithPassword({
  email,
  password,
});

if (error !== null) {
  console.error(error.message);
  throw error;
}

console.log(JSON.stringify(data.session, null, "\t"));
Deno.exit(0);
