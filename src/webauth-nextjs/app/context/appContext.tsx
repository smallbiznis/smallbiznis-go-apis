'use client'

import { api } from "@/lib/axios";
import { IApplication } from "@/types/application";
import { IPasswordRules } from "@/types/password";
import { Metadata, ResolvingMetadata } from "next";
import { useSearchParams } from "next/navigation";
import { createContext, useContext, useEffect, useState } from "react"

// @ts-ignore
export const AppContext = createContext();

// @ts-ignore
export const AppProvider = ({ children }) => {

  const searchParam = useSearchParams()
  const [loading, setLoading] = useState<Boolean>(true);
  const [organization, setOrganization] = useState(undefined)
  const [app, setApp] = useState(undefined)

  const [rules, setRules] = useState<IPasswordRules[]>([
      {
        label: "Lower case letters (a-Z)",
        pattern: "(?=.*[a-z])",
      },
      {
        label: "Upper case letters (A-Z)",
        pattern: "(?=.*[A-Z])",
      },
      {
        label: "Number (i.e. 0-9)",
        pattern: "(?=.*[0-9])",
      },
      {
        label: "At least 8 character",
        pattern: ".{8,}",
      },
      {
        label: "Special character (@$!%*?&)",
        pattern: "(?=.*[@$!%*?&])",
      },
    ]);

  const fetchApp = async () => {
    try {
      const q = new URLSearchParams()
      // @ts-ignore
      q.append("client_id", searchParam.get("client_id"))
      q.append("exludes", "client_id, client_secret")
      const response = await api.get(`/applications?${q.toString()}`)
      if (response.status > 200) {}
      const {total_data, data} = await response.data
      setApp(data[0])
    } catch (error) {
      throw error
    }
  }

  const fetchPasswodRule = () => {
    try {
      fetch(`/api/password/rules`)
    } catch (error) {
      
    }
  }

  useEffect(() => {
    fetchApp().finally(() => setLoading(false))
  }, [])

  const LoaderMarkup = (
    <div className={'h-screen flex items-center justify-center'} style={{
      height: '"calc(100% - 60px);"'
    }}>
      <svg aria-hidden="true" className="w-8 h-8 text-gray-200 animate-spin dark:text-gray-600 fill-blue-600" viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" fill="currentColor"/>
          <path d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" fill="currentFill"/>
      </svg>
      <span className="sr-only">Loading...</span>
    </div>
  )

  return (
    <AppContext.Provider value={{ loading, app, rules }}>
      {loading ? LoaderMarkup : children}
    </AppContext.Provider>
  )

}

export const useApp = () => useContext(AppContext)