import { PassportStrategy } from '@nestjs/passport';
import { Strategy, Profile } from 'passport-github2';
import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class GithubStrategy extends PassportStrategy(Strategy, 'github') {
  constructor(private configService: ConfigService) {
    super({
      clientID: configService.get('GITHUB_CLIENT_ID'),
      clientSecret: configService.get('GITHUB_CLIENT_SECRET'),
      callbackURL: configService.get('GITHUB_REDIRECT_URI'),
      scope: ['user:email'],
    });
  }

  async validate(
    accessToken: string,
    refreshToken: string,
    profile: Profile,
    done: Function,
  ): Promise<any> {
    try {
      const { id, username, emails, photos } = profile;

      // Validate required fields
      if (!id) {
        return done(new Error('Missing required user data from GitHub'), null);
      }

      // GitHub may not always provide email
      const email =
        emails && emails.length > 0
          ? emails[0].value
          : username
            ? `${username}@github.local`
            : `github_${id}@github.local`;

      const user = {
        provider: 'github',
        provider_user_id: id,
        email: email,
        name: username || profile.displayName || email.split('@')[0],
        avatar_url: photos && photos.length > 0 ? photos[0].value : null,
        access_token: accessToken,
        refresh_token: refreshToken || null,
      };

      done(null, user);
    } catch (error) {
      done(error, null);
    }
  }
}
