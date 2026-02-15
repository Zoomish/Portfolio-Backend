import { CACHE_MANAGER } from '@nestjs/cache-manager'
import { Inject, Injectable } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import axios from 'axios'
import { Cache } from 'cache-manager'
import {
    LinkedInContactInfoResponse,
    LinkedInPostResponse,
    LinkedInProfileDetailResponse,
} from './types'
import { CASHED_DATA_TTL } from 'src/const'

@Injectable()
export class LinkedInService {
    private readonly rapidApiKey: string
    private readonly rapidApiHost: string
    private readonly username = 'zoomish'

    constructor(
        @Inject(CACHE_MANAGER) private readonly cacheManager: Cache,
        private readonly configService: ConfigService
    ) {
        this.rapidApiKey = this.configService.get<string>('RAPIDAPI_KEY') || ''
        this.rapidApiHost =
            this.configService.get<string>('RAPIDAPI_HOST') || ''
    }

    private generateCacheKey(
        endpoint: string,
        params?: Record<string, any>
    ): string {
        const paramsStr = params ? `:${JSON.stringify(params)}` : ''
        return `linkedin:${endpoint}:${this.username}${paramsStr}`
    }

    private async getCachedData<T>(cacheKey: string): Promise<T | null> {
        return await this.cacheManager.get<T>(cacheKey)
    }

    private async setCacheData<T>(
        cacheKey: string,
        data: T,
        ttlSeconds: number = CASHED_DATA_TTL
    ): Promise<void> {
        await this.cacheManager.set(cacheKey, data, ttlSeconds)
    }

    private async makeApiRequest<T>(
        endpoint: string,
        params: Record<string, any> = {}
    ): Promise<T> {
        const options = {
            method: 'GET',
            url: `https://${this.rapidApiHost}/${endpoint}`,
            params: { username: this.username, ...params },
            headers: {
                'x-rapidapi-key': this.rapidApiKey,
                'x-rapidapi-host': this.rapidApiHost,
            },
            timeout: 20000,
        }

        const response = await axios.request<T>(options)
        return response.data
    }

    private async fetchWithCache<T>(
        endpoint: string,
        params?: Record<string, any>
    ): Promise<{ data: T; cached: boolean }> {
        const cacheKey = this.generateCacheKey(endpoint, params)

        const cachedData = await this.getCachedData<T>(cacheKey)
        if (cachedData) {
            return { data: cachedData, cached: true }
        }

        const apiData = await this.makeApiRequest<T>(endpoint, params)
        await this.setCacheData(cacheKey, apiData)

        return { data: apiData, cached: false }
    }

    async getProfilePosts(
        pageNumber: string = '1'
    ): Promise<{ data: LinkedInPostResponse; cached: boolean }> {
        return this.fetchWithCache<LinkedInPostResponse>('profile/posts', {
            page_number: pageNumber,
        })
    }

    async getProfileDetails(): Promise<{
        data: LinkedInProfileDetailResponse
        cached: boolean
    }> {
        return this.fetchWithCache<LinkedInProfileDetailResponse>(
            'profile/detail'
        )
    }

    async getProfileContactInfo(): Promise<{
        data: LinkedInContactInfoResponse
        cached: boolean
    }> {
        return this.fetchWithCache<LinkedInContactInfoResponse>(
            'profile/contact'
        )
    }

    async getAllProfileData(): Promise<{
        posts: { data: LinkedInPostResponse; cached: boolean }
        details: { data: LinkedInProfileDetailResponse; cached: boolean }
        contact: { data: LinkedInContactInfoResponse; cached: boolean }
    }> {
        const [posts, details, contact] = await Promise.all([
            this.getProfilePosts('1'),
            this.getProfileDetails(),
            this.getProfileContactInfo(),
        ])

        const combinedKey = this.generateCacheKey('combined')
        const combinedData = { posts, details, contact }
        await this.setCacheData(combinedKey, combinedData)

        return combinedData
    }
}
