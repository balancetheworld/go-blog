'use client'

import type { ReactNode } from 'react'
import type {
  CurrentRole,
  LoginReq,
  User,
} from '@/models/user'
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from 'react'
import {
  getCurrentUser,
  login as loginRequest,
  logout as logoutRequest,
} from '@/api/user'

interface AuthProviderProps {
  children: ReactNode
}

interface AuthContextValue {
  currentUser: User | null
  currentRole: CurrentRole
  isLoading: boolean
  login: (request: LoginReq) => Promise<void>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: AuthProviderProps) {
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const [currentRole, setCurrentRole] = useState<CurrentRole>('guest')
  const [isLoading, setIsLoading] = useState(true)

  // 封装一个把 currentUser 设为游客的函数
  const setGuest = useCallback(() => {
    setCurrentUser(null)
    setCurrentRole('guest')
  }, [])

  const refreshUser = useCallback(async () => {
    try {
      const result = await getCurrentUser()
      if (result.user == null) {
        setGuest()
        return
      }
      setCurrentUser(result.user)
      setCurrentRole(result.role)
    }
    catch (error) {
      setGuest()
      throw error
    }
  }, [setGuest])

  const login = useCallback(async (request: LoginReq) => {
    await loginRequest(request)
    await refreshUser()
  }, [refreshUser])

  const logout = useCallback(async () => {
    try {
      await logoutRequest()
    }
    finally {
      setGuest()
    }
  }, [setGuest])

  useEffect(() => {
    void refreshUser()
      .catch(() => undefined)
      .finally(() => setIsLoading(false))
  }, [refreshUser])

  return (
    <AuthContext.Provider
      value={{
        currentUser,
        currentRole,
        isLoading,
        login,
        refreshUser,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

// 封装读取登录全局状态的工具钩子，简化页面获取登录信息，同时做安全校验，防止开发者错误使用
export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useauth must be used within AuthProvider')
  }
  return context
}

// 关于封装：
// 减少重复请求、统一数据流，避免状态不一致
