'use client'

import type { ReactNode } from 'react'
import type {
  LoginRequest,
} from '@/lib/api/user'
import type {
  CurrentUser,
  UserRole,
} from '@/models/user'
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from 'react'
import {
  fetchCurrentUser,
  login as loginRequest,
  logout as logoutRequest,
} from '@/lib/api/user'

interface AuthProviderProps {
  children: ReactNode
}

type CurrentRole = UserRole | 'guest'
interface AuthContextValue {
  currentUser: CurrentUser | null
  currentRole: CurrentRole
  isLoading: boolean
  login: (request: LoginRequest) => Promise<void>
  logout: () => Promise<void>
  refreshCurrentUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: AuthProviderProps) {
  const [currentUser, setCurrentUser] = useState<CurrentUser | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const currentRole: CurrentRole = currentUser?.role ?? 'guest'

  const refreshCurrentUser = useCallback(
    async () => {
      const user = await fetchCurrentUser()
      setCurrentUser(user)
    },
    // 依赖为空 = 完全不依赖任何 state /props，组件整个生命周期只生成一次函数，永远复用同一个函数地址，不会重新创建。
    // 给 useEffect 做依赖，防止无限循环
    [],
  )

  // 前面缓存了这个函数，现在就只会在挂载的时候执行一次了
  useEffect(() => {
    // refreshCurrentUser()是异步函数，返回 promise， js/ts 规定凡是返回 promise 的异步操作，必须处理异常，否则视为代码隐患
    // 用 try catch 相对繁琐
    // `void xxx` 运算符：执行表达式，然后强制丢弃返回值，返回 undefined
    // 作用：告诉 TS：我明确知道这个函数返回 Promise，但我不需要接收它的结果，且已经手动处理异常，不用给我报未处理 Promise 警告。
    void refreshCurrentUser()
      .catch(() => setCurrentUser(null))
      .finally(() => setIsLoading(false))
  }, [refreshCurrentUser])

  async function login(request: LoginRequest) {
    await loginRequest(request)
    await refreshCurrentUser()
  }
  async function logout() {
    await logoutRequest()
    setCurrentUser(null)
  }

  return (
    <AuthContext.Provider
      value={{
        currentUser,
        currentRole,
        isLoading,
        login,
        logout,
        refreshCurrentUser,
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
