import { Module } from '@nestjs/common'
import { BotService } from './bot.service'
import {
    ContactMeService,
    HandleService,
    StartMessageService,
} from './services'

@Module({
    providers: [
        BotService,
        HandleService,
        StartMessageService,
        ContactMeService,
    ],
})
export class BotModule {}
