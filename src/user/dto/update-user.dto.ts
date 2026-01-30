import { IsOptional } from 'class-validator'

export class UpdateUserDto {
    @IsOptional()
    name?: string

    @IsOptional()
    username?: string

    @IsOptional()
    admin?: boolean
}
