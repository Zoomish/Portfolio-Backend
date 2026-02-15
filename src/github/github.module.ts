import { Module } from '@nestjs/common'
import { GithubController } from './github.controller'
import { GithubService } from './github.service'
import { RedisModule } from 'src/redis.module'

@Module({
    imports: [RedisModule],
    providers: [GithubService],
    controllers: [GithubController],
    exports: [],
})
export class GithubModule {}
