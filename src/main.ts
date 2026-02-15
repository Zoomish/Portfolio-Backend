import { ConsoleLogger } from '@nestjs/common'
import { NestFactory } from '@nestjs/core'
import { AppModule } from './app.module'

async function bootstrap() {
    const app = await NestFactory.create(AppModule, {
        logger: new ConsoleLogger({
            prefix: 'Portfolio',
            compact: 5,
        }),
    })
    app.setGlobalPrefix('api')
    app.enableCors()
    const redisClient = require('redis').createClient({
        url: process.env.REDIS_URL,
    })
    redisClient.on('connect', () => console.log('Redis connected'))
    redisClient.on('error', (err) => console.error('Redis error:', err))

    await app.listen(3000, async () => {
        console.log(`Server started on port ${await app.getUrl()}`)
    })
}
bootstrap()
