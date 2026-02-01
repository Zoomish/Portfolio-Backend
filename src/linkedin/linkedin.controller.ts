import { Controller, Get, Query } from '@nestjs/common'
import { LinkedInService } from './linkedin.service'

@Controller('linkedin')
export class LinkedInController {
    constructor(private readonly linkedInService: LinkedInService) {}

    @Get('posts')
    async getPosts(@Query('page') page: string = '1') {
        return this.linkedInService.getProfilePosts(page)
    }

    @Get('profile')
    async getProfileDetails() {
        return this.linkedInService.getProfileDetails()
    }

    @Get('contact')
    async getProfileContact() {
        return this.linkedInService.getProfileContactInfo()
    }

    @Get('all')
    async getAllData() {
        return this.linkedInService.getAllProfileData()
    }
}
