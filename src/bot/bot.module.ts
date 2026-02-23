import { Module } from '@nestjs/common'
import { UserModule } from 'src/user/user.module'
import { BotService } from './bot.service'
import {
    ContactMeService,
    HandleService,
    StartMessageService,
} from './services'

@Module({
    imports: [UserModule],
    providers: [
        BotService,
        HandleService,
        StartMessageService,
        ContactMeService,
    ],
})
export class BotModule {}
