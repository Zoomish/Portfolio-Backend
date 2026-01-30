import { Injectable, Logger } from '@nestjs/common'
@Injectable()
export class LinkedinService {
    private readonly logger = new Logger(LinkedinService.name)

    async findOne(tgId: number) {}

    async findAll() {}
}
