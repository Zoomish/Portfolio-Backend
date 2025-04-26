import { NestFactory } from '@nestjs/core'
import { AppModule } from './app.module'

async function bootstrap() {
    const app = await NestFactory.create(AppModule)
    app.setGlobalPrefix('api')
    app.enableCors()
    function aaa() {
        setInterval(
            async () => {
                await fetch(
                    'https://activityhelper-mln1.onrender.com/api/repos'
                )
            },
            1000 * 60 * 10
        )
    }

    await app.listen(3000, async () => {
        aaa()
        console.log(`Server started on port ${await app.getUrl()}`)
    })
}
bootstrap()
