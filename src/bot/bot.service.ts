import { Injectable, OnModuleInit } from '@nestjs/common'
import TelegramBot, * as telegram from 'node-telegram-bot-api'
import { HandleService } from './services'

@Injectable()
export class BotService implements OnModuleInit {
    constructor(private readonly handleService: HandleService) {}
    async onModuleInit() {
        const telegramToken = process.env.TELEGRAM_TOKEN
        const bot: TelegramBot = new telegram(telegramToken, {
            polling: true,
        })
        await this.initBot(bot)
    }

    async initBot(bot: TelegramBot) {
        global.bot = bot
        await bot.setMyCommands([
            {
                command: '/start',
                description: 'Начать',
            },
            {
                command: '/about',
                description: 'Информация обо мне',
            },
            {
                command: '/contact',
                description: 'Связаться со мной',
            },
        ])

        bot.on('message', async (msg: TelegramBot.Message) => {
            return await this.handleService.handleMessage(msg)
        })
    }
}
