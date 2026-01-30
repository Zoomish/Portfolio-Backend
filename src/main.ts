import { NestFactory } from '@nestjs/core'
import { AppModule } from './app.module'
import { ConsoleLogger } from '@nestjs/common'

async function bootstrap() {
    const app = await NestFactory.create(AppModule, {
        logger: new ConsoleLogger({
            prefix: 'Portfolio',
            compact: 5,
        }),
    })
    app.setGlobalPrefix('api')
    app.enableCors()

    await app.listen(3000, async () => {
        console.log(`Server started on port ${await app.getUrl()}`)
    })
}
bootstrap()
