import { Module } from '@nestjs/common'
import { RedisModule } from 'src/redis.module'
import { LinkedInController } from './linkedin.controller'
import { LinkedInService } from './linkedin.service'

@Module({
    imports: [RedisModule],
    controllers: [LinkedInController],
    providers: [LinkedInService],
    exports: [],
})
export class LinkedinModule {}
