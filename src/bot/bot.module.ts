import { Module } from '@nestjs/common'
import { UserModule } from 'src/user/user.module'
import { BotService } from './bot.service'
import {
    CallbackService,
    ContactMeService,
    GetAdminCallbackService,
    GetAdminService,
    HandleGetService,
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
        GetAdminService,
        CallbackService,
        HandleGetService,
        GetAdminCallbackService,
    ],
})
export class BotModule {}
