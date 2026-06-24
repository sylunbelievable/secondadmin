import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  input: '../server/docs/openapi.json',
  output: 'src/api/generated',
})
