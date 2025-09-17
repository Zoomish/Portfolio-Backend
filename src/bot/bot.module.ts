import { Module } from '@nestjs/common'
import { BotService } from './bot.service'
import { HandleService, StartMessageService } from './services'

@Module({
    providers: [BotService, HandleService, StartMessageService],
})
export class BotModule {}
