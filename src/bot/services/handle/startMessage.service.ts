import { Injectable } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'

@Injectable()
export class StartMessageService {
    async startMessage() {
        const bot: TelegramBot = global.bot
        const msg: TelegramBot.Message = global.msg
        await bot.sendMessage(msg.chat.id, 'hello')
    }
}
