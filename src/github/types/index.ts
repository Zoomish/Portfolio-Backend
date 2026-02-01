export interface GithubRepo {
    id: number
    name: string
    full_name: string
    private: boolean
    html_url: string
    description: string | null
    language: string | null
    stargazers_count: number
    forks_count: number
    readme_content?: string
}

export interface FilteredRepo {
    name: string
    description: string
    language: string
    url: string
    stars: number
    forks: number
}