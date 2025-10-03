import { PassportStrategy } from '@nestjs/passport';
import { Strategy, VerifyCallback } from 'passport-google-oauth20';
import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class GoogleStrategy extends PassportStrategy(Strategy, 'google') {
  constructor(private configService: ConfigService) {
    super({
      clientID: configService.get('GOOGLE_CLIENT_ID'),
      clientSecret: configService.get('GOOGLE_CLIENT_SECRET'),
      callbackURL: configService.get('GOOGLE_REDIRECT_URI'),
      scope: ['email', 'profile'],
    });
  }

  async validate(
    accessToken: string,
    refreshToken: string,
    profile: any,
    done: VerifyCallback,
  ): Promise<any> {
    try {
      const { id, name, emails, photos } = profile;

      // Validate required fields
      if (!id || !emails || emails.length === 0) {
        return done(new Error('Missing required user data from Google'), null);
      }

      const user = {
        provider: 'google',
        provider_user_id: id,
        email: emails[0].value,
        name: name
          ? `${name.givenName || ''} ${name.familyName || ''}`.trim()
          : emails[0].value.split('@')[0],
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
