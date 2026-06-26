# test-questions

1. Design an In-Memory Rate Limiter: Implement Token Bucket or Leaky Bucket. Code it to be entirely thread-safe for distributed usage.

[Token Bucket](./rate_limit/token-bucket.go)

Think of a bucket that slowly fills with tokens.

* Tokens are added at a fixed rate (e.g. 10 tokens per second).
* Each request consumes a token.
* If a token is available, the request proceeds.
* If the bucket is empty, the request is rejected, delayed, or throttled.

Example

Bucket size: 100 tokens
Refill rate: 10 tokens/second

2. Merge K Sorted Streams. Given multiple huge log files, write an algorithm to merge them sequentially. Explain how to execute this if only a few megabytes of RAM are available. [merge](./merge/merge.go)


3. Build a Thread-Safe Bounded Queue [queue](./queue/queue.go)


4. Ring Hash 
2. SOLID principles

* Single Responsibility Principle (SRP)
    - What it means: A class should have one, and only one, reason to change. It should have exactly one job or purpose within the software.
    - Why it helps: Isolates behaviors so that changes to one part of your system don't trigger unintended side effects in others.

* Open/Closed Principle (OCP)
    - What it means: Software entities (classes, modules, functions) should be open for extension but closed for modification.
    - Why it helps: You can add new features or change behaviors by adding new code, rather than altering existing, tested code.

* Liskov Substitution Principle (LSP) 
    - What it means: Objects of a superclass should be replaceable with objects of its subclasses without breaking the application. Subclasses must behave in a way that the parent class's clients expect.
    - Why it helps: Prevents inheritance misuse and ensures that child classes don't remove functionality or violate the contract of the parent class.
    
* Interface Segregation Principle (ISP) 
    - What it means: No client should be forced to depend on methods it does not use. It is better to have many small, client-specific interfaces than one general-purpose interface.
    - Why it helps: Keeps the system decoupled and prevents "fat" interfaces that force objects to implement behaviors they don't need.
    
* Dependency Inversion Principle (DIP)
    - What it means: High-level modules should not depend on low-level modules. Both should depend on abstractions (e.g., interfaces). Furthermore, abstractions should not depend on details; details should depend on abstractions.
    - Why it helps: Reduces the dependency of high-level business logic on low-level concrete implementations, making components much easier to swap, mock, and test.