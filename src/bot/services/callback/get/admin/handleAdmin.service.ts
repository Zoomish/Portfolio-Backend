import { forwardRef, Inject, Injectable, Logger } from '@nestjs/common'
import TelegramBot from 'node-telegram-bot-api'
import { GetAdminService } from 'src/bot/services/handle'
import { User } from 'src/user/entities/user.entity'
import { UserService } from 'src/user/user.service'

@Injectable()
export class GetAdminCallbackService {
    constructor(
        @Inject(forwardRef(() => GetAdminService))
        private readonly getAdminService: GetAdminService,
        @Inject(forwardRef(() => UserService))
        private readonly userService: UserService
    ) {}
    private readonly logger = new Logger(GetAdminCallbackService.name)

    async handleGetAdmin(
        text: string | undefined,
        callbackQuery: TelegramBot.CallbackQuery
    ) {
        const bot: TelegramBot = global.bot
        const texts = text?.split('=') ?? [text]

        switch (texts[0]) {
            case undefined: {
                bot.answerCallbackQuery(callbackQuery.id, {
                    text: 'Меню администратора',
                }).catch((error) =>
                    this.logger.error('Error answering callback: ' + error)
                )
                return await this.getAdminService.getAdmin(
                    callbackQuery.message.message_id
                )
            }
            case 'users': {
                bot.answerCallbackQuery(callbackQuery.id, {
                    text: 'Получить всех пользователей',
                }).catch((error) =>
                    this.logger.error('Error answering callback: ' + error)
                )
                return await this.handleGetAdminUsers()
            }
        }
    }

    async handleGetAdminUsers() {
        const bot: TelegramBot = global.bot
        const msg: TelegramBot.Message = global.msg
        await bot
            .deleteMessage(msg.chat.id, msg.message_id)
            .catch((error) =>
                this.logger.error('Error deleting message: ' + error)
            )
        const users = await this.userService.findAll()
        for await (const user of users) {
            await this.sendUserInfo(user)
        }
    }

    async sendUserInfo(user: User, msg_id?: number) {
        const bot: TelegramBot = global.bot
        const msg: TelegramBot.Message = global.msg

        const text = `Пользователь <code>${user.tgId}</code>: @` + user.username
        if (msg_id) {
            return await bot
                .editMessageText(text + ' '.repeat(Math.random() * 100), {
                    chat_id: msg.chat.id,
                    message_id: msg_id,
                    parse_mode: 'HTML',
                })
                .catch((error) =>
                    this.logger.error('Error editing message: ' + error)
                )
        }
        return await bot
            .sendMessage(msg.chat.id, text, {
                parse_mode: 'HTML',
            })
            .catch((error) =>
                this.logger.error('Error sending message: ' + error)
            )
    }
}
