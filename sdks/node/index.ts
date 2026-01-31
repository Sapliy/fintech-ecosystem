
export interface RecordTransactionRequest {
  accountId: string;
  amount: number; // In cents
  currency: string;
  description: string;
  referenceId: string;
}

export interface RecordTransactionResponse {
  transactionId: string;
  status: string;
}

export interface GetAccountResponse {
  accountId: string;
  balance: number;
  currency: string;
  createdAt: string;
}

export interface ValidateKeyResponse {
  valid: boolean;
  userId: string;
  orgId: string;
  environment: string;
  scopes: string;
}

export class FintechClient {
  private apiKey: string;
  private baseURL: string;

  constructor(apiKey: string, baseURL: string = 'http://localhost:8080') {
    this.apiKey = apiKey;
    this.baseURL = baseURL.endsWith('/') ? baseURL.slice(0, -1) : baseURL;
  }

  private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseURL}${path}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': this.apiKey,
        ...options.headers,
      },
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Fintech API Error (${response.status}): ${errorText || response.statusText}`);
    }

    return response.json() as Promise<T>;
  }

  public ledger = {
    recordTransaction: async (req: RecordTransactionRequest): Promise<RecordTransactionResponse> => {
      return this.request<RecordTransactionResponse>('/v1/ledger/transactions', {
        method: 'POST',
        body: JSON.stringify(req),
      });
    },
    getAccount: async (accountId: string): Promise<GetAccountResponse> => {
      return this.request<GetAccountResponse>(`/v1/ledger/accounts/${accountId}`);
    },
  };

  public auth = {
    validateKey: async (keyHash: string): Promise<ValidateKeyResponse> => {
      return this.request<ValidateKeyResponse>('/v1/auth/validate', {
        method: 'POST',
        body: JSON.stringify({ keyHash }),
      });
    },
  };
}
