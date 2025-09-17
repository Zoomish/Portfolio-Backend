import { Injectable } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'

@Injectable()
export class ContactMeService {
    async contactMe() {
        const bot: TelegramBot = global.bot
        const msg: TelegramBot.Message = global.msg
        await bot.sendMessage(
            msg.chat.id,
            '📩 Хотите связаться со мной напрямую?\n' +
                'Я всегда рад новым знакомствам, интересным предложениям и любым вопросам, связанным с моей работой. 🤝\n' +
                '\n' +
                'Вы можете написать мне прямо в Telegram — просто нажмите на мой профиль:\n' +
                '<b>@ZoomishWork</b> ✉️\n' +
                '\n' +
                'Не стесняйтесь обращаться! Будь то сотрудничество, вопросы по проектам или просто дружеский обмен опытом — я открыт для общения 😊\n' +
                '\n' +
                '<b>📌 Кнопка навигации:</b>\n' +
                '➡️ <b>Связаться со мной</b> — быстрый способ перейти в мой Telegram-профиль и начать диалог. 🚀',
            {
                parse_mode: 'HTML',
                reply_markup: {
                    inline_keyboard: [
                        [
                            {
                                text: 'Связаться со мной',
                                url: 'https://t.me/ZoomishWork',
                            },
                        ],
                    ],
                },
            }
        )
    }
}
