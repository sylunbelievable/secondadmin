import { describe, expect, it } from 'vitest'
import type { Menu } from '#/api'
import { containsPath, routeFor, routeTitle } from './app-shell'

const tree: Menu[] = [{
  id: '1',
  parentId: '0',
  type: 'directory',
  name: '系统管理',
  sort: 1,
  visible: true,
  status: 1,
  children: [{
    id: '2',
    parentId: '1',
    type: 'menu',
    name: '用户管理',
    component: 'users',
    sort: 1,
    visible: true,
    status: 1,
  }],
}]

describe('routeFor', () => {
  it('uses the local route registry and rejects unknown components', () => {
    expect(routeFor({ component: 'users', path: '/wrong' })).toBe('/users')
    expect(routeFor({ component: 'unknown', path: '/unsafe' })).toBeUndefined()
    expect(routeFor({ path: '/custom' })).toBe('/custom')
  })

  it('finds active routes and titles through nested menu trees', () => {
    expect(containsPath(tree[0]!, '/users')).toBe(true)
    expect(routeTitle('/users', tree)).toBe('用户管理')
  })
})
