import { Injectable, Logger } from '@nestjs/common'
@Injectable()
export class GithubService {
    private readonly logger = new Logger(GithubService.name)

    async findOne(tgId: number) {}

    async findAll() {}
}
