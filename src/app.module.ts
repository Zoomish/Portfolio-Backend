import { Module } from '@nestjs/common'
import { ConfigModule } from '@nestjs/config'
import { AppController } from './app.controller'
import { BotModule } from './bot/bot.module'

@Module({
    imports: [ConfigModule.forRoot({ isGlobal: true }), BotModule],
    controllers: [AppController],
})
export class AppModule {}
