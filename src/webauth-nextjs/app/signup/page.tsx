"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import { Suspense, useContext, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useToast } from "@/components/ui/use-toast";
import { PasswordCheckList } from "@/components/password_rules";
import { IPasswordRules } from "@/types/password.d";
import { api } from "@/lib/axios";
import { AppContext, AppProvider } from "../context/appContext";
import Head from "next/head";

export default function Page() {
  const { toast } = useToast();
  const query = useSearchParams();
  const { replace } = useRouter();

  // @ts-ignore
  const { app, rules } = useContext(AppContext)

  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const [provider, setProvider] = useState("password");
  const [providers, setProviders] = useState([
    {
      id: "password",
    },
    {
      id: "phone_number",
    },
    {
      id: "google",
    },
    {
      id: "facebook",
    },
  ]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const [valid, setValid] = useState(false);

  const onSubmit = async (event: any) => {
    event.preventDefault();
    setLoading(true);

    const { data } = await api.post(`/accounts/signup`, {
      // @ts-ignore
      client_id: query.get("client_id"),
      provider,
      first_name: firstName,
      last_name: lastName,
      email,
      password,
    });
    if (data.error) {
      setLoading(false);
      setError(data.error);
      return;
    }

    replace(`/signin?${query.toString()}`);
  };

  return (
    <>
      <Head>
        <title>{app?.display_name}</title>
      </Head>
      <div className="grid gap-2 text-center">
        <h1 className="text-3xl font-bold">Sign Up</h1>
        <p className="text-balance text-muted-foreground">
          Create a <span className={'font-bold text-primary'}>{app?.display_name}</span> account
        </p>
      </div>
      {error != "" && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      <div className="grid gap-4">
        <div className="grid grid-cols-2 gap-4">
          <div className="grid gap-2">
            <Label htmlFor="first-name">First name</Label>
            <Input
              id="first-name"
              placeholder="Max"
              autoComplete="given-name"
              required
              onChange={(e: any) => setFirstName(e.target.value)}
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="last-name">Last name</Label>
            <Input
              id="last-name"
              placeholder="Robinson"
              autoComplete="family-name"
              required
              onChange={(e: any) => setLastName(e.target.value)}
            />
          </div>
        </div>
        <div className="grid gap-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            placeholder="m@example.com"
            autoComplete="email"
            required
            onChange={(e: any) => setEmail(e.target.value)}
          />
        </div>
        <div className="grid gap-2">
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            type="password"
            autoComplete="new-password"
            required
            onChange={(e: any) => setPassword(e.target.value)}
          />
        </div>
        <div className="grid gap-2">
          <PasswordCheckList
            value={password}
            rules={rules}
            onChange={setValid}
          />
        </div>
        <Button
          type="submit"
          className="w-full"
          onClick={onSubmit}
          disabled={loading}
        >
          Create an account
        </Button>
        {/* <Button variant="outline" className="w-full">
          Sign up with Google
        </Button> */}
      </div>
      <div className="mt-4 text-center text-sm">
        Already have an account?{" "}
        <Link href={query.size > 0 ? `/signin?${query.toString()}` : `/signin`} className="underline">
          Sign in
        </Link>
      </div>
    </>
  );
}
