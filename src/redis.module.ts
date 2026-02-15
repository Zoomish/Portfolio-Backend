import { CacheModule } from '@nestjs/cache-manager'
import { Global, Module } from '@nestjs/common'
import { ConfigModule, ConfigService } from '@nestjs/config'
import * as redisStore from 'cache-manager-redis-store'
import { CASHED_DATA_TTL } from './const'

@Global()
@Module({
    imports: [
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
    ],
    exports: [CacheModule],
})
export class RedisModule {}
