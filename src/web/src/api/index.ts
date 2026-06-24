import './client'

export * from './generated/sdk.gen'
export type * from './generated/types.gen'
export { errorMessage, onAuthFailure } from './client'

export async function responseData<T>(request: Promise<{ data: T }>) {
  return (await request).data
}
