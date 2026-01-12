---
layout: home
title: Blog
titleTemplate: Blog | Goca

hero:
  name: Goca Blog
  text: Insights, Updates & Architecture
  tagline: Stay updated with the latest releases, tutorials, and architectural insights for building production-ready Go applications
  actions:
    - theme: brand
      text: Latest Release
      link: /blog/releases/v1-17-1
    - theme: alt
      text: All Articles
      link: /blog/articles/

features:
  - title: Release Notes
    details: Comprehensive changelogs and release announcements for each version
    link: /blog/releases/
    linkText: View Releases

  - title: Articles
    details: In-depth articles on Clean Architecture patterns, best practices, and advanced techniques
    link: /blog/articles/
    linkText: Read Articles

  - title: Architecture Insights
    details: Deep dives into design decisions, trade-offs, and architectural patterns used in Goca
    link: /guide/clean-architecture
    linkText: Learn More
---

<style scoped>
:deep(.VPFeature) {
  transition: all 0.3s ease;
}

:deep(.VPFeature:hover) {
  transform: translateY(-4px);
}
</style>

## Latest Release

### [v1.17.1 - Database Driver Configuration Fix](/blog/releases/v1-17-1)
*January 11, 2026*

Critical bug fix resolving database driver configuration issues. SQLite and other non-PostgreSQL databases now properly configured during project initialization. Fixes issue #31.

[Read full release notes](/blog/releases/v1-14-1)

---

## Recent Articles

### [Mastering the Repository Pattern in Clean Architecture](/blog/articles/mastering-repository-pattern)
*October 29, 2025*

A comprehensive guide to the Repository pattern in Clean Architecture. Learn the difference between repositories and DAOs, how to design clean interfaces, implement database-specific code, and how Goca generates production-ready repositories with complete abstraction.

[Read article](/blog/articles/mastering-repository-pattern)

### [Mastering Use Cases in Clean Architecture](/blog/articles/mastering-use-cases)
*October 29, 2025*

A deep dive into use cases and application services. Learn what use cases are, how they differ from controllers, DTOs patterns, and how Goca generates complete application layer code with orchestration logic and best practices.

[Read article](/blog/articles/mastering-use-cases)

### [Understanding Domain Entities in Clean Architecture](/blog/articles/understanding-domain-entities)
*October 29, 2025*

A comprehensive guide to domain entities in Clean Architecture. Learn what entities are, why they're not database models, DDD principles, and how Goca generates production-ready entities with validation and business rules.

[Read article](/blog/articles/understanding-domain-entities)

---

## Coming Soon

Stay tuned for more articles on:
- Building scalable microservices with Clean Architecture
- Advanced testing strategies with Goca
- Database migration patterns and best practices
- Performance optimization in Go applications

---

## Subscribe

Follow the project on [GitHub](https://github.com/sazardev/goca) to stay updated with releases and announcements.
