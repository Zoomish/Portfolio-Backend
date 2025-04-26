import { Injectable } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'

@Injectable()
export class HandleService {
    async handeMessage() {
        const bot: TelegramBot = global.bot
        const chatId = process.env.ADMIN_ID
        await bot.sendMessage(chatId, 'Hello world')
    }
}
