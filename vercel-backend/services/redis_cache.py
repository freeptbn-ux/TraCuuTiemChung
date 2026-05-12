import os
from upstash_redis import Redis
from dotenv import load_dotenv

load_dotenv()

class RedisCache:
    def __init__(self):
        url = os.getenv("UPSTASH_REDIS_REST_URL")
        token = os.getenv("UPSTASH_REDIS_REST_TOKEN")
        
        if not url or not token:
            print("WARNING: UPSTASH_REDIS_REST_URL or UPSTASH_REDIS_REST_TOKEN not set.")
            self.client = None
        else:
            self.client = Redis(url=url, token=token)

    def set(self, key: str, value: str, ex: int = 1800):
        """Set value with expiration in seconds (default 30 mins)."""
        if self.client:
            try:
                self.client.set(key, value, ex=ex)
            except Exception as e:
                print(f"Redis set error: {e}")

    def get(self, key: str):
        """Get value from Redis."""
        if self.client:
            try:
                return self.client.get(key)
            except Exception as e:
                print(f"Redis get error: {e}")
        return None

    def delete(self, key: str):
        """Delete key from Redis."""
        if self.client:
            try:
                self.client.delete(key)
            except Exception as e:
                print(f"Redis delete error: {e}")
