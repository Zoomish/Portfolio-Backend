import { CacheModule } from '@nestjs/cache-manager'
import { Global, Module } from '@nestjs/common'
import { ConfigModule, ConfigService } from '@nestjs/config'
import * as redisStore from 'cache-manager-redis-store'

@Global()
@Module({
    imports: [
        CacheModule.registerAsync({
            imports: [ConfigModule],
            useFactory: async (configService: ConfigService) => ({
                store: redisStore,
                url: configService.get<string>('REDIS_URL'),
                ttl: 1000 * 60 * 60 * 24 * 2,
                max: 100,
            }),
            inject: [ConfigService],
        }),
    ],
    exports: [CacheModule],
})
export class RedisModule {}
