const TOKEN_KEY = 'admin_token'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function removeToken() {
  localStorage.removeItem(TOKEN_KEY)
}

// api 封装：自动带 JWT token，401 时清除 token 并跳转登录页
export async function api(url, options = {}) {
  const token = getToken()
  const headers = { 'Content-Type': 'application/json', ...options.headers }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(url, { ...options, headers })

  if (res.status === 401) {
    removeToken()
    if (!window.location.pathname.startsWith('/login')) {
      window.location.href = '/login'
    }
    return { code: 401, message: '请重新登录' }
  }

  return res.json()
}
