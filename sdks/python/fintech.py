import requests
from typing import Dict, Any, Optional

class FintechClient:
    def __init__(self, api_key: str, base_url: str = "http://localhost:8080"):
        self.api_key = api_key
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            "Content-Type": "application/json",
            "X-API-Key": self.api_key
        })

    def _request(self, method: str, path: str, json: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        url = f"{self.base_url}{path}"
        response = self.session.request(method, url, json=json)
        response.raise_for_status()
        return response.json()

    @property
    def ledger(self):
        class LedgerService:
            def __init__(self, client: 'FintechClient'):
                self.client = client

            def record_transaction(self, account_id: str, amount: int, currency: str, description: str, reference_id: str) -> Dict[str, Any]:
                payload = {
                    "accountId": account_id,
                    "amount": amount,
                    "currency": currency,
                    "description": description,
                    "referenceId": reference_id
                }
                return self.client._request("POST", "/v1/ledger/transactions", json=payload)

            def get_account(self, account_id: str) -> Dict[str, Any]:
                return self.client._request("GET", f"/v1/ledger/accounts/{account_id}")

        return LedgerService(self)

    @property
    def auth(self):
        class AuthService:
            def __init__(self, client: 'FintechClient'):
                self.client = client

            def validate_key(self, key_hash: str) -> Dict[str, Any]:
                return self.client._request("POST", "/v1/auth/validate", json={"keyHash": key_hash})

        return AuthService(self)
