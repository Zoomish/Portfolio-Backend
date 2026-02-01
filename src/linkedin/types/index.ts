export interface LinkedInPostResponse {
    success: boolean
    message: string
    data: {
        posts: Array<{
            urn: {
                activity_urn: string
                share_urn: string | null
                ugcPost_urn: string | null
            }
            full_urn: string
            posted_at: {
                date: string
                relative: string
                timestamp: number
            }
            text: string
            url: string
            post_type: string
            author: {
                first_name: string
                last_name: string
                headline: string
                username: string
                profile_url: string
                profile_picture: string
            }
            stats: {
                total_reactions: number
                like: number
                support: number
                love: number
                insight: number
                celebrate: number
                funny: number
                comments: number
                reposts: number
            }
        }>
        pagination_token?: string
    }
}

export interface LinkedInProfileDetailResponse {
    success: boolean
    message: string
    data: {
        basic_info: {
            fullname: string
            first_name: string
            last_name: string
            headline: string
            public_identifier: string
            profile_url: string
            profile_picture_url: string
            about: string
            location: {
                country: string
                city: string
                full: string
                country_code: string
            }
            creator_hashtags: string[]
            is_creator: boolean
            is_influencer: boolean
            is_premium: boolean
            open_to_work: boolean
            created_timestamp: number
            show_follower_count: boolean
            background_picture_url: string
            urn: string
            follower_count: number
            connection_count: number
            current_company: string
            current_company_urn: string
            current_company_url: string
            top_skills: string[]
            email: string | null
        }
        experience: Array<{
            title: string
            company: string
            location?: string
            description: string
            duration: string
            start_date: {
                year: number
                month: string
            }
            end_date: {
                year: number
                month: string
            }
            is_current: boolean
            company_linkedin_url?: string
            company_logo_url?: string
            employment_type?: string
            location_type?: string
            skills?: string[]
            company_id?: string
            skills_url?: string
        }>
        education: Array<{
            school: string
            degree: string
            degree_name: string
            field_of_study: string
            duration: string
            school_linkedin_url: string
            activities: string
            start_date: {
                year: number
                month: string
            }
            end_date: {
                year: number
                month: string
            }
        }>
        languages: Array<{
            language: string
            proficiency: string
        }>
    }
}

export interface LinkedInContactInfoResponse {
    success: boolean
    message: string
    data: {
        basic_info: {
            urn: string
            first_name: string
            last_name: string
            full_name: string
            headline: string
            public_identifier: string
            profile_url: string
        }
        email: string | null
        phone_numbers: Array<any>
        websites: Array<{
            url: string
            category: string
            label: string | null
        }>
        birthday: {
            month: number
            day: number
            year: number | null
        }
        address: any | null
        twitter_handles: Array<any>
        wechat: string | null
        instant_messengers: Array<any>
    }
}