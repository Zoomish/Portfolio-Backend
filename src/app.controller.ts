import { Controller, Get } from '@nestjs/common';
import { CacheKey, CacheTTL } from '@nestjs/cache-manager';

@Controller('app')
export class AppController {
    @Get()
    @CacheKey('app-hello')
    @CacheTTL(600) 
    getHello(): string {
        return 'Hello World!';
    }
}