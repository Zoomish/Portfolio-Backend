import { Module } from '@nestjs/common'
import { BotService } from './bot.service'
import { HandleService } from './services'
import { StartMessageService } from './services/handle/startMessage.service'

@Module({
    providers: [BotService, HandleService, StartMessageService],
})
export class BotModule {}
