// @vitest-environment jsdom

import { afterEach, describe, expect, it, vi } from 'vitest'
import { currentUser, onAuthFailure, responseData, updateUser } from '#/api'

afterEach(() => {
  vi.unstubAllGlobals()
  document.cookie = 'csrf_token=; Max-Age=0; path=/'
})

describe('API client authentication', () => {
  it('shares one refresh across concurrent 401 responses', async () => {
    document.cookie = 'csrf_token=test-csrf; path=/'
    let refreshed = false
    let refreshCalls = 0
    vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const request = new Request(input, init)
      const path = new URL(request.url).pathname
      if (path === '/api/v1/auth/refresh') {
        refreshCalls++
        expect(request.headers.get('X-CSRF-Token')).toBe('test-csrf')
        await Promise.resolve()
        refreshed = true
        return Response.json({ expiresIn: 900 })
      }
      if (path === '/api/v1/auth/me' && !refreshed) {
        return Response.json({ code: 'UNAUTHORIZED', message: 'expired', requestId: 'test' }, { status: 401 })
      }
      return Response.json({ id: 1, username: 'admin', nickname: 'Administrator', status: 1 })
    }))

    const users = await Promise.all([
      responseData(currentUser({ throwOnError: true })),
      responseData(currentUser({ throwOnError: true })),
    ])

    expect(refreshCalls).toBe(1)
    expect(users.map((user) => user.username)).toEqual(['admin', 'admin'])
  })

  it('adds CSRF to unsafe generated-client requests', async () => {
    document.cookie = 'csrf_token=write-token; path=/'
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const request = new Request(input, init)
      expect(request.headers.get('X-CSRF-Token')).toBe('write-token')
      return new Response(null, { status: 204 })
    })
    vi.stubGlobal('fetch', fetchMock)

    await updateUser({ path: { id: 1 }, body: { status: 1 }, throwOnError: true })

    expect(fetchMock).toHaveBeenCalledOnce()
  })

  it('notifies once when refresh fails', async () => {
    const handler = vi.fn()
    const dispose = onAuthFailure(handler)
    vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const request = new Request(input, init)
      const path = new URL(request.url).pathname
      if (path === '/api/v1/auth/refresh') return new Response(null, { status: 401 })
      return new Response(null, { status: 401 })
    }))

    await Promise.allSettled([
      responseData(currentUser({ throwOnError: true })),
      responseData(currentUser({ throwOnError: true })),
    ])

    expect(handler).toHaveBeenCalledOnce()
    dispose()
  })

  it('notifies when the retried request is still unauthorized', async () => {
    const handler = vi.fn()
    const dispose = onAuthFailure(handler)
    vi.stubGlobal('fetch', vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const request = new Request(input, init)
      const path = new URL(request.url).pathname
      if (path === '/api/v1/auth/refresh') return Response.json({ expiresIn: 900 })
      return new Response(null, { status: 401 })
    }))

    await expect(responseData(currentUser({ throwOnError: true }))).rejects.toBeDefined()

    expect(handler).toHaveBeenCalledOnce()
    dispose()
  })
})
