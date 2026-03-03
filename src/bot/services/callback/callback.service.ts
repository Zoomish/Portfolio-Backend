import { Injectable, Logger } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'
import { HandleGetService } from './get'

@Injectable()
export class CallbackService {
    constructor(private readonly handleGetService: HandleGetService) {}
    private readonly logger = new Logger(CallbackService.name)

    async handleCallback(callbackQuery: TelegramBot.CallbackQuery) {
        const texts = [
            callbackQuery.data.split('_')[0],
            callbackQuery.data.split('_').slice(1).join('_'),
        ]
        global.msg = callbackQuery.message

        switch (texts[0]) {
            case 'get':
                return await this.handleGetService.handleGet(
                    texts[1],
                    callbackQuery
                )
        }
    }
}
