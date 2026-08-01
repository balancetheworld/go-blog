// 引入Node.js内置process模块，用于读取环境变量（如.env里的后端地址）
import process from 'node:process'
// 引入axios，用于发起HTTP网络请求
import axios from 'axios'
// 引入字段命名转换工具函数：驼峰 ↔ 下划线互相转换
// camelToSnakeObj：对象所有key 驼峰转下划线
// camelToSnakeStr：单个字符串驼峰转下划线
// snakeToCamelObj：对象所有key 下划线转驼峰
import { camelToSnakeObj, camelToSnakeStr, snakeToCamelObj } from 'field-conv'

// 后端基础地址常量：优先读取环境变量BACKEND_URL，无环境变量则默认本地8888端口
export const BACKEND_URL = process.env.BACKEND_URL ?? 'http://localhost:8888'

/**
 * 转换FormData表单数据的key命名：驼峰字段名转下划线
 * @param data 原始前端FormData表单对象
 * @returns 转换key为下划线的全新FormData
 */
function convertFormData(data: FormData) {
  // 创建新FormData存储转换后数据
  const convertedData = new FormData()
  // 遍历原始表单每一对key-value
  data.forEach((value, key) => {
    // 把前端驼峰key（userName）转为后端下划线key（user_name），再存入新表单
    convertedData.append(camelToSnakeStr(key), value)
  })
  return convertedData
}

// 创建全局统一axios请求实例，所有接口统一使用该客户端，统一拦截处理
export const axiosClient = axios.create({
  // 接口统一前缀 /api/v1，所有请求自动拼接
  baseURL: '/api/v1',
})

// ===================== 请求拦截器（发请求之前自动执行） =====================
axiosClient.interceptors.request.use((config) => {
  // 场景1：请求体是FormData表单（上传文件/表单提交）
  if (config.data instanceof FormData) {
    // 转换表单所有key为下划线，适配后端Go接收规范
    config.data = convertFormData(config.data)
  }
  // 场景2：普通JSON请求体（POST/PUT JSON参数）
  else if (config.data && typeof config.data === 'object') {
    // JSON对象所有驼峰key统一转下划线
    config.data = camelToSnakeObj(config.data)
  }

  // 场景3：URL查询参数（GET请求 ?userName=xxx）
  if (config.params && typeof config.params === 'object') {
    // 查询参数驼峰key转下划线
    config.params = camelToSnakeObj(config.params)
  }

  // 返回修改完成的请求配置，正式发送请求
  return config
})

// ===================== 响应拦截器（收到后端返回数据后自动执行） =====================
axiosClient.interceptors.response.use((response) => {
  // 如果后端返回JSON对象
  if (response.data && typeof response.data === 'object') {
    // 后端返回的下划线字段（user_name）自动转回前端JS标准驼峰（userName）
    response.data = snakeToCamelObj(response.data)
  }
  // 返回转换完成的响应数据，页面拿到直接使用驼峰字段，无需手动转换
  return response
})
