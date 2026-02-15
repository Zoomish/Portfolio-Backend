import { CACHE_MANAGER } from '@nestjs/cache-manager'
import { Controller, Get, Inject } from '@nestjs/common'
import { Cache } from 'cache-manager'

@Controller('app')
export class AppController {
    constructor(@Inject(CACHE_MANAGER) private readonly cacheManager: Cache) {}

    @Get()
    async getHello(): Promise<string> {
        const cacheKey = 'app_hello'
        const cachedValue = await this.cacheManager.get<string>(cacheKey)

        if (cachedValue) {
            return cachedValue
        }

        const newValue = `Hello World! ${new Date().toISOString()}`
        await this.cacheManager.set(cacheKey, newValue, 900)
        return newValue
    }
}
