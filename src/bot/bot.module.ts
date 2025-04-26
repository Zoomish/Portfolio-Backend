import { Module } from '@nestjs/common'
import { BotService } from './bot.service'
import { HandleService } from './services'

@Module({
    providers: [BotService, HandleService],
    exports: [HandleService],
})
export class BotModule {}
