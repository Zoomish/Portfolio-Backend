import { Module } from '@nestjs/common'
import { GithubController } from './github.controller'
import { GithubService } from './github.service'

@Module({
    imports: [],
    providers: [GithubService],
    controllers: [GithubController],
    exports: [],
})
export class GithubModule {}
