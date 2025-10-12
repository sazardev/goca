---
layout: home

hero:
  name: "Goca"
  text: "Go Clean Architecture"
  tagline: "Build production-ready Go applications following Clean Architecture principles. Stop writing boilerplate, start building features."
  image:
    src: /logo.svg
    alt: Goca Logo
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/sazardev/goca

features:
  - icon:
      src: /icons/layers.svg
    title: Clean Architecture by Default
    details: Every line of code follows Uncle Bob's Clean Architecture principles. Proper layer separation, dependency rules, and clear boundaries guaranteed.
  
  - icon:
      src: /icons/zap.svg
    title: Lightning Fast Setup
    details: Generate complete features with all layers in seconds. From entity to handler, repository to use case - everything configured and ready.
  
  - icon:
      src: /icons/target.svg
    title: Best Practices Enforced
    details: Prevents common anti-patterns like fat controllers, god objects, and anemic domain models. Your code stays clean and maintainable.
  
  - icon:
      src: /icons/refresh.svg
    title: Auto Integration
    details: New features are automatically integrated with dependency injection and routing. No manual wiring needed.
  
  - icon:
      src: /icons/package.svg
    title: Multi-Protocol Support
    details: Generate handlers for HTTP REST, gRPC, CLI, Workers, and SOAP. All following the same clean architecture pattern.
  
  - icon:
      src: /icons/flask.svg
    title: Test-Ready
    details: Code generated with clear interfaces and dependency injection makes testing a breeze. TDD-friendly from the start.
  
  - icon:
      src: /icons/database.svg
    title: Repository Pattern
    details: Abstracted data access with interchangeable implementations. Switch from PostgreSQL to MongoDB without touching business logic.
  
  - icon:
      src: /icons/book.svg
    title: Rich Documentation
    details: Comprehensive guides, tutorials, and examples. Learn Clean Architecture while building real applications.
  
  - icon:
      src: /icons/rocket.svg
    title: Production Ready
    details: Used in production systems. Battle-tested patterns and code generation that scales from MVP to enterprise.
---

## Quick Example

```bash
# Initialize a new project
goca init my-api --module github.com/user/my-api

# Generate a complete feature with all layers
goca feature User --fields "name:string,email:string,role:string"

# That's it! You now have:
# → Domain entity with validations
# → Use cases with DTOs
# → Repository with PostgreSQL implementation
# → HTTP handlers with routing
# → Dependency injection configured
```

## Why Clean Architecture?

Clean Architecture ensures your codebase remains:

- **Maintainable**: Changes in one layer don't cascade through the entire system
- **Testable**: Business logic is independent of frameworks and databases
- **Flexible**: Easy to swap implementations without touching core logic
- **Scalable**: Clear boundaries make it easy to add new features

## What Developers Say

> "Goca transformed how we build Go services. What used to take hours now takes minutes, and the code quality is consistently high."
>
> — Production User

> "Finally, a code generator that doesn't just dump code but teaches you proper architecture."
>
> — Go Developer

> "The automatic integration of new features saved us so much time. No more manual wiring!"
>
> — Go Team Lead

## Ready to Build?

<p style="text-align: center; margin: 2rem 0;">
  <a href="/goca/getting-started.html" style="display: inline-block; padding: 0.75rem 2rem; background: #037794ff; color: white; border-radius: 0.75rem; text-decoration: none; font-weight: 600; transition: all 0.3s ease;">Get Started in 5 Minutes</a>
</p>

<style>
.vp-doc a {
  text-decoration: none;
}
</style>
