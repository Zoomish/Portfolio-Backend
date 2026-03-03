import { Injectable, Logger } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'
import { GetAdminCallbackService } from './admin'

@Injectable()
export class HandleGetService {
    constructor(
        private readonly getAdminCallbackService: GetAdminCallbackService
    ) {}
    private readonly logger = new Logger(HandleGetService.name)

    async handleGet(text: string, callbackQuery: TelegramBot.CallbackQuery) {
        const texts = text.split('-')

        switch (texts[0]) {
            case 'admin': {
                return await this.getAdminCallbackService.handleGetAdmin(
                    texts[1],
                    callbackQuery
                )
            }
        }
    }
}
