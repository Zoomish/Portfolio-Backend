import { CacheModule } from '@nestjs/cache-manager'
import { Module } from '@nestjs/common'
import { ConfigModule, ConfigService } from '@nestjs/config'
import { ScheduleModule } from '@nestjs/schedule'
import { ThrottlerModule } from '@nestjs/throttler'
import { TypeOrmModule } from '@nestjs/typeorm'
import * as redisStore from 'cache-manager-redis-store'
import * as fs from 'fs'
import { AppController } from './app.controller'
import { BotModule } from './bot/bot.module'
import { CASHED_DATA_TTL } from './const'
import { GithubModule } from './github/github.module'
import { LinkedinModule } from './linkedin/linkedin.module'
import { GetActiveService } from './tasks/getActiveService.service'
import { UserModule } from './user/user.module'

@Module({
    imports: [
        ConfigModule.forRoot({ isGlobal: true }),
        CacheModule.registerAsync({
            imports: [ConfigModule],
            useFactory: async (configService: ConfigService) => ({
                store: redisStore,
                url: configService.get<string>('REDIS_URL'),
                ttl: CASHED_DATA_TTL,
                max: 100,
                isGlobal: true,
            }),
            inject: [ConfigService],
        }),
        ScheduleModule.forRoot(),
        TypeOrmModule.forRootAsync({
            imports: [ConfigModule],
            useFactory: (configService: ConfigService) => ({
                type: 'postgres',
                host: configService.get('POSTGRES_HOST'),
                port: configService.get('POSTGRES_PORT'),
                username: configService.get('POSTGRES_USER'),
                password: configService.get('POSTGRES_PASSWORD'),
                database: configService.get('POSTGRES_DATABASE'),
                entities: [__dirname + '/**/*.entity{.ts,.js}'],
                autoLoadEntities: true,
                synchronize: true,
                migrationsRun: true,
                ssl:
                    configService.get('NODE_ENV') === 'dev'
                        ? false
                        : {
                              rejectUnauthorized: true,
                              ca: fs.readFileSync('./ca.pem').toString(),
                          },
            }),
            inject: [ConfigService],
        }),
        ThrottlerModule.forRoot([
            {
                name: 'short',
                ttl: 1000,
                limit: 10,
            },
            {
                name: 'medium',
                ttl: 10000,
                limit: 50,
            },
            {
                name: 'long',
                ttl: 60000,
                limit: 300,
            },
        ]),
        BotModule,
        UserModule,
        GithubModule,
        LinkedinModule,
    ],
    controllers: [AppController],
    providers: [GetActiveService],
})
export class AppModule {}
