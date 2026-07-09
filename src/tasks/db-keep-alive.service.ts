import { Injectable, Logger, OnModuleInit } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { SchedulerRegistry } from '@nestjs/schedule'
import { InjectDataSource } from '@nestjs/typeorm'
import { CronJob } from 'cron'
import { DataSource } from 'typeorm'
import {
    DB_KEEPALIVE_CRON_DEFAULT,
    DB_KEEPALIVE_ENABLED_DEFAULT,
} from '../const'

@Injectable()
export class DbKeepAliveService implements OnModuleInit {
    private readonly logger = new Logger(DbKeepAliveService.name)

    constructor(
        @InjectDataSource()
        private readonly dataSource: DataSource,
        private readonly configService: ConfigService,
        private readonly schedulerRegistry: SchedulerRegistry
    ) {}

    onModuleInit() {
        const cronExpression = this.configService.get(
            'DB_KEEPALIVE_CRON',
            DB_KEEPALIVE_CRON_DEFAULT
        )

        const job = new CronJob(
            cronExpression,
            () => {
                void this.handleCron()
            },
            null,
            true,
            'Europe/Moscow'
        )

        this.schedulerRegistry.addCronJob('db-keep-alive', job)
    }

    async handleCron() {
        await this.runKeepAlive()
    }

    async runKeepAlive(): Promise<void> {
        if (
            this.configService.get(
                'DB_KEEPALIVE_ENABLED',
                DB_KEEPALIVE_ENABLED_DEFAULT
            ) !== 'true'
        ) {
            return
        }

        const queryRunner = this.dataSource.createQueryRunner()
        await queryRunner.connect()

        try {
            await queryRunner.query(
                'SELECT COUNT(*)::int AS count FROM "user"'
            )
            await queryRunner.query(
                'SELECT id FROM "user" WHERE admin = true LIMIT 1'
            )

            await queryRunner.startTransaction()
            await queryRunner.query(
                'INSERT INTO "user" ("tgId", admin, name) VALUES ($1, false, $2)',
                [-1, '__keepalive__']
            )
            await queryRunner.query(
                'UPDATE "user" SET name = $1 WHERE "tgId" = $2',
                ['__keepalive_updated__', -1]
            )
            await queryRunner.query('DELETE FROM "user" WHERE "tgId" = $1', [
                -1,
            ])
            await queryRunner.rollbackTransaction()

            this.logger.log('DB keep-alive completed')
        } catch (error) {
            if (queryRunner.isTransactionActive) {
                await queryRunner.rollbackTransaction()
            }
            this.logger.error('DB keep-alive failed', error)
        } finally {
            await queryRunner.release()
        }
    }
}
