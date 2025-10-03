/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-call */
import {
  Body,
  Controller,
  Get,
  Post,
  Query,
  Req,
  Res,
  UseGuards,
} from '@nestjs/common';
import type { Request, Response } from 'express';
import { RefreshTokenDto } from '../dto/auth.dto';
import { AuthService } from './auth.service';
import {
  GithubAuthGuard,
  GoogleAuthGuard,
  JwtAuthGuard,
} from './guards/auth.guards';
import { MsalService } from './strategies/msal.service';

@Controller('auth')
export class AuthController {
  constructor(
    private authService: AuthService,
    private msalService: MsalService,
  ) {}

  // Google OAuth
  @Get('google/login')
  @UseGuards(GoogleAuthGuard)
  async googleAuth() {
    // Guard redirects to Google
  }

  @Get('callback/google')
  @UseGuards(GoogleAuthGuard)
  async googleAuthCallback(@Req() req: Request, @Res() res: Response) {
    try {
      // Ensure req.user is of type OAuthCallbackDto
      const user = req.user as import('../dto/auth.dto').OAuthCallbackDto;

      // Validate required fields
      if (!user || !user.provider_user_id || !user.email) {
        const redirectUrl = `http://localhost:7101/auth/login?error=missing_data`;
        return res.redirect(redirectUrl);
      }

      const tokenResponse = await this.authService.handleOAuthCallback(user);

      // Redirect to frontend with tokens and user info
      const redirectUrl = `http://localhost:7101/auth/callback?access_token=${encodeURIComponent(tokenResponse.access_token)}&refresh_token=${encodeURIComponent(tokenResponse.refresh_token)}&user_id=${encodeURIComponent(tokenResponse.user.id)}&email=${encodeURIComponent(tokenResponse.user.email || '')}&name=${encodeURIComponent(tokenResponse.user.name || '')}`;
      return res.redirect(redirectUrl);
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error ? error.message : String(error);
      console.error('Google OAuth callback error:', errorMessage);
      const redirectUrl = `http://localhost:7101/auth/login?error=server_error`;
      return res.redirect(redirectUrl);
    }
  }

  // GitHub OAuth
  @Get('github/login')
  @UseGuards(GithubAuthGuard)
  async githubAuth() {
    // Guard redirects to GitHub
  }

  @Get('callback/github')
  @UseGuards(GithubAuthGuard)
  async githubAuthCallback(
    @Req() req: { user: import('../dto/auth.dto').OAuthCallbackDto },
    @Res() res: Response,
  ) {
    try {
      const user = req.user;

      // Validate required fields
      if (!user || !user.provider_user_id || !user.email) {
        const redirectUrl = `http://localhost:7101/auth/login?error=missing_data`;
        return res.redirect(redirectUrl);
      }

      const tokenResponse = await this.authService.handleOAuthCallback(user);

      // Redirect to frontend with tokens and user info
      const redirectUrl = `http://localhost:7101/auth/callback?access_token=${encodeURIComponent(tokenResponse.access_token)}&refresh_token=${encodeURIComponent(tokenResponse.refresh_token)}&user_id=${encodeURIComponent(tokenResponse.user.id)}&email=${encodeURIComponent(tokenResponse.user.email || '')}&name=${encodeURIComponent(tokenResponse.user.name || '')}`;
      return res.redirect(redirectUrl);
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error ? error.message : String(error);
      console.error('GitHub OAuth callback error:', errorMessage);
      const redirectUrl = `http://localhost:7101/auth/login?error=server_error`;
      return res.redirect(redirectUrl);
    }
  }

  // Microsoft MSAL
  @Get('msal/login')
  async msalAuth(@Res() res: Response) {
    try {
      const authCodeUrl = await this.msalService.getAuthCodeUrl();
      return res.redirect(authCodeUrl);
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : String(error);
      console.error('MSAL auth initiation error:', errorMessage);
      const redirectUrl = `http://localhost:7101/auth/login?error=msal_init_failed`;
      return res.redirect(redirectUrl);
    }
  }

  @Get('msal')
  async msalCallback(
    @Query('code') code: string,
    @Query('error') error: string,
    @Query('error_description') errorDescription: string,
    @Res() res: Response,
  ) {
    try {
      // Check for errors from Microsoft
      if (error) {
        const errorMsg = errorDescription || error;
        console.error('MSAL OAuth error:', errorMsg);
        const redirectUrl = `http://localhost:7101/auth/login?error=${encodeURIComponent(errorMsg)}`;
        return res.redirect(redirectUrl);
      }

      if (!code) {
        const redirectUrl = `http://localhost:7101/auth/login?error=no_code`;
        return res.redirect(redirectUrl);
      }

      type TokenResponse = {
        accessToken: string;
        refreshToken: string;
      };

      type UserInfo = {
        id: string;
        mail?: string;
        userPrincipalName?: string;
        displayName?: string;
      };

      const tokenResponse = (await this.msalService.acquireTokenByCode(
        code,
      )) as TokenResponse;
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
      const userInfo: UserInfo = await this.msalService.getUserInfo(
        tokenResponse.accessToken,
      );

      const oauthData = {
        provider: 'msal',
        provider_user_id: userInfo.id,
        email: userInfo.mail || userInfo.userPrincipalName || '',
        name: userInfo.displayName || '',
        avatar_url: null,
        access_token: tokenResponse.accessToken,
        refresh_token: tokenResponse.refreshToken,
      };

      // Validate required fields
      if (!oauthData.provider_user_id || !oauthData.email) {
        const redirectUrl = `http://localhost:7101/auth/login?error=missing_data`;
        return res.redirect(redirectUrl);
      }

      const authTokens = await this.authService.handleOAuthCallback(oauthData);

      // Redirect to frontend with tokens and user info
      const redirectUrl = `http://localhost:7101/auth/callback?access_token=${encodeURIComponent(authTokens.access_token)}&refresh_token=${encodeURIComponent(authTokens.refresh_token)}&user_id=${encodeURIComponent(authTokens.user.id)}&email=${encodeURIComponent(authTokens.user.email || '')}&name=${encodeURIComponent(authTokens.user.name || '')}`;
      return res.redirect(redirectUrl);
    } catch (error: unknown) {
      const errorMessage =
        error instanceof Error ? error.message : String(error);
      console.error('MSAL OAuth callback error:', errorMessage);
      const redirectUrl = `http://localhost:7101/auth/login?error=server_error`;
      return res.redirect(redirectUrl);
    }
  }

  // Refresh Token
  @Post('refresh')
  async refreshToken(@Body() refreshTokenDto: RefreshTokenDto) {
    return this.authService.refreshTokens(refreshTokenDto.refresh_token);
  }

  // Logout
  @Post('logout')
  @UseGuards(JwtAuthGuard)
  async logout(@Req() req: any) {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const token = req.headers.authorization?.split(' ')[1];
    // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
    await this.authService.logout(req.user.userId, token);
    return { message: 'Logged out successfully' };
  }

  // Get current user (protected route example)
  @Get('me')
  @UseGuards(JwtAuthGuard)
  async getCurrentUser(@Req() req: any) {
    const user = await this.authService.validateUser(req.user.userId);
    return {
      id: user.id,
      email: user.email,
      name: user.name,
      avatar_url: user.avatar_url,
    };
  }

  // Health check
  @Get('health')
  healthCheck() {
    return { status: 'ok', service: 'authentication' };
  }
}
