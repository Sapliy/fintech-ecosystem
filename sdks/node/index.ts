import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';

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
  private client: AxiosInstance;

  constructor(apiKey: string, baseURL: string = 'http://localhost:8080') {
    this.client = axios.create({
      baseURL: baseURL.endsWith('/') ? baseURL.slice(0, -1) : baseURL,
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': apiKey,
      },
    });
  }

  private async request<T>(path: string, options: AxiosRequestConfig = {}): Promise<T> {
    try {
      const response = await this.client.request<T>({
        url: path,
        ...options,
      });
      return response.data;
    } catch (error: any) {
      if (error.response) {
        throw new Error(`Fintech API Error (${error.response.status}): ${JSON.stringify(error.response.data)}`);
      }
      throw error;
    }
  }

  public ledger = {
    recordTransaction: async (req: RecordTransactionRequest): Promise<RecordTransactionResponse> => {
      return this.request<RecordTransactionResponse>('/v1/ledger/transactions', {
        method: 'POST',
        data: req,
      });
    },
    getAccount: async (accountId: string): Promise<GetAccountResponse> => {
      return this.request<GetAccountResponse>(`/v1/ledger/accounts/${accountId}`, {
        method: 'GET',
      });
    },
  };

  public auth = {
    validateKey: async (keyHash: string): Promise<ValidateKeyResponse> => {
      return this.request<ValidateKeyResponse>('/v1/auth/validate', {
        method: 'POST',
        data: { keyHash },
      });
    },
  };
}
