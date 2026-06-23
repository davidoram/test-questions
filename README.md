# test-questions

1. Design an In-Memory Rate Limiter: Implement Token Bucket or Leaky Bucket. Code it to be entirely thread-safe for distributed usage.

Token Bucket

Think of a bucket that slowly fills with tokens.

* Tokens are added at a fixed rate (e.g. 10 tokens per second).
* Each request consumes a token.
* If a token is available, the request proceeds.
* If the bucket is empty, the request is rejected, delayed, or throttled.

Example

Bucket size: 100 tokens
Refill rate: 10 tokens/second