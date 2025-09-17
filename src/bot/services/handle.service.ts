import { Injectable } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'
import { ContactMeService, StartMessageService } from './handle'

@Injectable()
export class HandleService {
    constructor(
        private readonly startMessageService: StartMessageService,
        private readonly contactMeService: ContactMeService
    ) {}

    async handleMessage(msg: TelegramBot.Message) {
        const text = msg.text

        switch (text) {
            case '/start':
                return await this.startMessageService.startMessage()
            case '/about':
                return await this.startMessageService.startMessage()
            case '/contact':
                return await this.contactMeService.contactMe()

            default:
                break
        }
    }
}
