import { Injectable } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'

@Injectable()
export class HandleService {
    async handeMessage() {
        const bot: TelegramBot = global.bot
        const chatId = process.env.ADMIN_ID
        console.log(chatId)

        await bot.sendMessage(chatId, 'hello')
    }
}
