import { Injectable } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'

@Injectable()
export class HandleService {
    async handleMessage(msg: TelegramBot.Message) {
        const bot: TelegramBot = global.bot
        const chatId = msg.chat.id
        console.log(chatId)
        await bot.sendMessage(chatId, 'hello')
    }
}
