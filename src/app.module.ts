import { CacheModule } from '@nestjs/cache-manager'
import { Module } from '@nestjs/common'
import { ConfigModule } from '@nestjs/config'
import { AppController } from './app.controller'
import { BotModule } from './bot/bot.module'
import { GithubModule } from './github/github.module'
import { LinkedinModule } from './linkedin/linkedin.module'
import { RedisModule } from './redis.module'

@Module({
    imports: [
        ConfigModule.forRoot({ isGlobal: true }),
        CacheModule.register({
            ttl: 1000 * 60 * 60 * 24,
            isGlobal: true,
        }),
        BotModule,
        RedisModule,
        GithubModule,
        LinkedinModule,
    ],
    controllers: [AppController],
})
export class AppModule {}
