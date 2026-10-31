import { describe, expect, it } from 'vitest';
import { WebhookAuthMode } from '../../../../../api/v1/schema';
import {
  buildWebhookExamples,
  findActiveAllowedProfile,
  findUnavailableAllowedProfiles,
  updateAllowedProfiles,
} from '../webhookProfileSelection';

describe('webhook profile selection', () => {
  it('adds profiles once in sorted order and removes unchecked profiles', () => {
    expect(updateAllowedProfiles(['prod'], 'staging', true)).toEqual([
      'prod',
      'staging',
    ]);
    expect(updateAllowedProfiles(['staging', 'prod'], 'prod', true)).toEqual([
      'prod',
      'staging',
    ]);
    expect(updateAllowedProfiles(['prod', 'staging'], 'prod', false)).toEqual([
      'staging',
    ]);
  });

  it('finds configured profiles missing from the active profile list', () => {
    expect(
      findUnavailableAllowedProfiles(
        ['prod', 'retired', 'staging'],
        ['prod', 'staging']
      )
    ).toEqual(['retired']);
  });

  it('selects an active allowed profile for generated examples', () => {
    expect(
      findActiveAllowedProfile(['retired', 'staging'], ['prod', 'staging'])
    ).toBe('staging');
    expect(findActiveAllowedProfile(['retired'], ['prod'])).toBe('');
  });

  it('builds profile-bound HMAC request examples', () => {
    const examples = buildWebhookExamples({
      authMode: WebhookAuthMode.token_and_hmac,
      profileName: 'prod',
      webhookUrl: 'https://flow.example/api/v1/webhooks/example',
    });

    expect(examples.curl).toContain('-H "Authorization: Bearer <YOUR_TOKEN>"');
    expect(examples.curl).toContain(
      '-H "X-Flow-Signature: sha256=<SIGNATURE>"'
    );
    expect(examples.curl).toContain('-H "X-Flow-Profile: prod"');
    expect(examples.hmacShell).toContain("profile='prod'");
    expect(examples.hmacShell).toContain(
      'signature_input=$(printf \'x-flow-profile:%s\\n%s\' "$profile" "$body")'
    );
    expect(examples.hmacNode).toContain(
      "const signatureInput = 'x-flow-profile:' + profile + '\\n' + body;"
    );
    expect(examples.hmacNode).toContain(
      "crypto.createHmac('sha256', process.env.FLOW_HMAC_SECRET)"
    );
    expect(examples.hmacNode).toContain(
      "'Authorization': 'Bearer <YOUR_TOKEN>',\n  };"
    );
  });

  it('builds examples for token-only and HMAC-only authentication', () => {
    const tokenOnly = buildWebhookExamples({
      authMode: WebhookAuthMode.token_only,
      profileName: '',
      webhookUrl: 'https://flow.example/token',
    });
    expect(tokenOnly.curl).toContain('Authorization: Bearer <YOUR_TOKEN>');
    expect(tokenOnly.curl).not.toContain('X-Flow-Signature');
    expect(tokenOnly.curl).not.toContain('X-Flow-Profile');

    const hmacOnly = buildWebhookExamples({
      authMode: WebhookAuthMode.hmac_only,
      profileName: '',
      webhookUrl: 'https://flow.example/hmac',
    });
    expect(hmacOnly.curl).not.toContain('Authorization');
    expect(hmacOnly.curl).toContain('X-Flow-Signature');
    expect(hmacOnly.hmacShell).toContain('signature_input="$body"');
    expect(hmacOnly.hmacNode).toContain('const signatureInput = body;');
  });
});
