import { Injectable, UnauthorizedException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { ConfidentialClientApplication } from '@azure/msal-node';

@Injectable()
export class MsalService {
  private msalClient: ConfidentialClientApplication;

  constructor(private configService: ConfigService) {
    this.msalClient = new ConfidentialClientApplication({
      auth: {
        clientId: this.configService.get('APP_AZURE_CLIENT_ID'),
        authority: `https://login.microsoftonline.com/${this.configService.get('APP_AZURE_TENANT_ID')}`,
        clientSecret: this.configService.get('AZURE_CLIENT_SECRET'),
      },
    });
  }

  getAuthCodeUrl(): Promise<string> {
    const redirectUri = this.configService.get('APP_AZURE_REDIRECT_URI');
    const tenantId = this.configService.get('APP_AZURE_TENANT_ID') || 'common';

    return this.msalClient.getAuthCodeUrl({
      scopes: [
        'https://graph.microsoft.com/User.Read',
        'openid',
        'profile',
        'email',
      ],
      redirectUri: redirectUri,
      prompt: 'select_account',
    });
  }

  async acquireTokenByCode(code: string): Promise<any> {
    try {
      const redirectUri = this.configService.get('APP_AZURE_REDIRECT_URI');
      const tokenResponse = await this.msalClient.acquireTokenByCode({
        code: code,
        scopes: [
          'https://graph.microsoft.com/User.Read',
          'openid',
          'profile',
          'email',
        ],
        redirectUri: redirectUri,
      });

      if (!tokenResponse) {
        throw new UnauthorizedException('Failed to acquire token');
      }

      return tokenResponse;
    } catch (error) {
      console.error('MSAL Token Error:', error);
      throw new UnauthorizedException('Failed to authenticate with Microsoft');
    }
  }

  async getUserInfo(accessToken: string): Promise<any> {
    try {
      const response = await fetch('https://graph.microsoft.com/v1.0/me', {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to fetch user info');
      }

      return await response.json();
    } catch (error) {
      console.error('Microsoft Graph Error:', error);
      throw new UnauthorizedException('Failed to fetch user information');
    }
  }
}
