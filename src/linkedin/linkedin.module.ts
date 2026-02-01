import { Module } from '@nestjs/common'
import { LinkedInController } from './linkedin.controller'
import { LinkedInService } from './linkedin.service'

@Module({
    imports: [],
    controllers: [LinkedInController],
    providers: [LinkedInService],
    exports: [],
})
export class LinkedinModule {}
