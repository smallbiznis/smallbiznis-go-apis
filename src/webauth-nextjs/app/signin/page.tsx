'use client'

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";
import { useToast } from "@/components/ui/use-toast";
import { PasswordCheckList } from "@/components/password_rules";
import { useApp } from "../context/appContext";
import { api } from "@/lib/axios";
import { Metadata, ResolvingMetadata } from "next";

export default function Page() {

  const { toast } = useToast();
  const query = useSearchParams();
  const { replace } = useRouter();

  // @ts-ignore
  const { app, rules } = useApp()

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [valid, setValid] = useState(false);

  const handleSubmit = async (event: any) => {
    try {
      event.preventDefault();
      setLoading(true)

      const { status, data } = await api.post(`/accounts/signInWithPassword`, {
        email,
        password,
      });

      if (status > 200) {
        setError(data.error);
        return
      }

      replace(`/oauth/authorize?${query.toString()}`);
    } catch (error) {
      
    } finally {
      setLoading(false)
    }
  };

  return (
    <>
      <div className="grid gap-2 text-center">
        <h1 className="text-3xl font-bold">Login</h1>
        <p className="text-balance text-muted-foreground">
          Continue to <span className={'font-bold text-primary'}>{app?.display_name}</span>
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
        <form>
          <div className="grid gap-2 mb-4">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="your@example.com"
              autoComplete="email"
              required
              onChange={(e: any) => setEmail(e.target.value)}
            />
          </div>
          <div className="grid gap-2 mb-4">
            <div className="flex items-center">
              <Label htmlFor="password">Password</Label>
              <Link
                href="/forgot-password"
                className="ml-auto inline-block text-sm underline"
              >
                Forgot your password?
              </Link>
            </div>
            <Input
              id="password"
              type="password"
              autoComplete="current-password"
              required
              onChange={(e: any) => setPassword(e.target.value)}
            />
          </div>
          <div className="grid gap-2 mb-4">
            <PasswordCheckList
              value={password}
              rules={rules}
              onChange={setValid}
            />
          </div>
          <Button
            type={'submit'}
            className={`w-full ${loading ? 'bg-primary' : 'bg-secondary-foreground'}`}
            onClick={handleSubmit}
            disable={loading}
          >
            Login
          </Button>
        </form>
        {/* <Button variant="outline" className="w-full">
          Login with Google
        </Button> */}
      </div>
      <div className="mt-4 text-center text-sm">
        Don&apos;t have an account?{" "}
        <Link href={query.size > 0 ? `/signup?${query.toString()}` : '/signup'} className="underline">
          Sign up
        </Link>
      </div>
    </>
  );
}
