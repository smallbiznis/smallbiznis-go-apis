import { Inter as FontSans } from "next/font/google";
import "./globals.css";

import Image from "next/image";
import { ThemeProvider } from "next-themes";
import { cn } from "@/lib/utils";
import { useEffect } from "react";
import { initSDKFaro } from "@/lib/faro-sdk";
import { AppProvider } from "./context/appContext";

const fontSans = FontSans({ subsets: ["latin"], variable: "--font-sans" });

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {

  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={cn(
          "min-h-screen bg-background font-sans antialiased",
          fontSans.variable,
        )}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="light"
          enableSystem
          disableTransitionOnChange
        >
          <AppProvider>
            <div className="w-full h-screen lg:grid lg:min-h-[600px] lg:grid-cols-2 xl:min-h-[800px]">
              <div className="flex items-center justify-center py-12">
                <div className="mx-auto grid w-[350px] gap-6">{children}</div>
              </div>
              <div className="hidden bg-muted lg:block">
                <Image
                  src={"/placeholder.svg"}
                  alt="Image"
                  width="1920"
                  height="1080"
                  className="h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
                />
              </div>
            </div>
          </AppProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
