import { CACHE_MANAGER } from '@nestjs/cache-manager'
import { Inject, Injectable } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import axios from 'axios'
import { Cache } from 'cache-manager'
import { FilteredRepo, GithubRepo } from './types'

@Injectable()
export class GithubService {
    private readonly githubToken: string
    private readonly username = 'Zoomish'

    constructor(
        @Inject(CACHE_MANAGER) private readonly cacheManager: Cache,
        private readonly configService: ConfigService
    ) {
        this.githubToken = this.configService.get<string>('GITHUB_TOKEN') || ''
    }

    private generateCacheKey(endpoint: string): string {
        return `github:${endpoint}:${this.username}`
    }

    private async getCachedData<T>(cacheKey: string): Promise<T | null> {
        return await this.cacheManager.get<T>(cacheKey)
    }

    private async setCacheData<T>(
        cacheKey: string,
        data: T,
        ttlSeconds: number = 172800
    ): Promise<void> {
        await this.cacheManager.set(cacheKey, data, ttlSeconds)
    }

    private async fetchAllRepos(): Promise<GithubRepo[]> {
        const perPage = 100
        let page = 1
        let allRepos: GithubRepo[] = []
        let hasMore = true

        while (hasMore) {
            const response = await axios.get<GithubRepo[]>(
                `https://api.github.com/user/repos`,
                {
                    headers: {
                        Authorization: `token ${this.githubToken}`,
                        Accept: 'application/vnd.github.v3+json',
                    },
                    params: {
                        per_page: perPage,
                        page: page,
                        sort: 'updated',
                        direction: 'desc',
                    },
                }
            )

            if (response.data.length === 0) {
                hasMore = false
            } else {
                allRepos = allRepos.concat(response.data)
                page++
            }
        }

        return allRepos
    }

    private async fetchRepoReadme(
        repoFullName: string
    ): Promise<string | null> {
        try {
            const response = await axios.get<{ content: string }>(
                `https://api.github.com/repos/${repoFullName}/readme`,
                {
                    headers: {
                        Authorization: `token ${this.githubToken}`,
                        Accept: 'application/vnd.github.v3.raw',
                    },
                }
            )
            return response.data.content
        } catch (error) {
            return null
        }
    }

    private async filterAndProcessRepos(
        repos: GithubRepo[]
    ): Promise<FilteredRepo[]> {
        const filtered: FilteredRepo[] = []
        const badWords = ['portfolio', 'zoomish', 'mesto-ad', 'alko-dubai-ecom', 'bad']

        for (const repo of repos) {
            if (filtered.length >= 6) break

            if (badWords.some((word) => repo.name.toLowerCase().includes(word)))
                continue

            if (!repo.private) {
                filtered.push({
                    name: repo.name,
                    description: repo.description || '',
                    language: repo.language || '',
                    url: repo.html_url,
                    link: repo.homepage,
                    stars: repo.stargazers_count,
                    forks: repo.forks_count,
                })
                continue
            }

            const readme = await this.fetchRepoReadme(repo.full_name)
            if (!readme) continue

            const hasProjectInReadme =
                readme.toLowerCase().includes('project') ||
                readme.toLowerCase().includes('проект')

            if (hasProjectInReadme) {
                filtered.push({
                    name: repo.name,
                    description: repo.description || '',
                    language: repo.language || '',
                    url: repo.html_url,
                    link: repo.homepage,
                    stars: repo.stargazers_count,
                    forks: repo.forks_count,
                })
            }
        }

        return filtered
    }

    async getFilteredRepositories(): Promise<{
        data: FilteredRepo[]
        cached: boolean
    }> {
        const cacheKey = this.generateCacheKey('repos')

        const cachedData = await this.getCachedData<FilteredRepo[]>(cacheKey)
        if (cachedData) {
            return { data: cachedData, cached: true }
        }

        const allRepos = await this.fetchAllRepos()
        const filteredRepos = await this.filterAndProcessRepos(allRepos)

        await this.setCacheData(cacheKey, filteredRepos)

        return { data: filteredRepos, cached: false }
    }
}
