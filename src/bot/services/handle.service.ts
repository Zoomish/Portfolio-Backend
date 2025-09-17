import { Injectable } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'
import { StartMessageService } from './handle/startMessage.service'

@Injectable()
export class HandleService {
    constructor(private readonly startMessageService: StartMessageService) {}

    async handleMessage(msg: TelegramBot.Message) {
        const text = msg.text

        switch (text) {
            case '/start':
                return await this.startMessageService.startMessage()
            case '/about':
                return await this.startMessageService.startMessage()

            default:
                break
        }
    }
}
