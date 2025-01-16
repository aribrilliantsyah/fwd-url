import { serve } from '@hono/node-server'
import { Hono } from 'hono'

const app = new Hono()

app.get('/', (c) => {
  return c.text('ini adalah home!')
})

app.get('/get', (c) => {
  return c.text('ini adalah get!')
})

app.post('/post', async (c) => {
  let headers = await c.req.header()
  console.log(headers)
  const json = await c.req.json()
  return c.json(json)
})

const port = 3000
console.log(`Server is running on http://localhost:${port}`)

serve({
  fetch: app.fetch,
  port
})
