---
layout: doc
title: Articles
titleTemplate: Blog | Goca
---

# Articles

In-depth articles covering Clean Architecture patterns, best practices, testing strategies, and advanced techniques for building production-ready Go applications with Goca.

---

<style scoped>
.article-list {
  margin-top: 2rem;
}

.article-item {
  padding: 1.5rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  margin-bottom: 1.5rem;
  transition: all 0.3s ease;
}

.article-item:hover {
  border-color: var(--vp-c-brand);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.article-meta {
  color: var(--vp-c-text-2);
  font-size: 0.9rem;
  margin-top: 0.5rem;
}

.article-tags {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.75rem;
  flex-wrap: wrap;
}

.article-tag {
  background: var(--vp-c-bg-soft);
  color: var(--vp-c-text-2);
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  border: 1px solid var(--vp-c-divider);
}

.coming-soon {
  text-align: center;
  padding: 3rem 1rem;
  color: var(--vp-c-text-2);
}
</style>

<div class="article-list">

<div class="article-item">
  <h3>
    <a href="/goca/blog/articles/mastering-repository-pattern">Mastering the Repository Pattern in Clean Architecture</a>
  </h3>
  <p>A comprehensive guide to the Repository pattern and data access abstraction. Learn what repositories are, how they differ from DAOs, interface design principles, and how Goca generates database-agnostic implementations that maintain Clean Architecture boundaries.</p>
  <div class="article-meta">Architecture • Infrastructure Layer</div>
  <div class="article-tags">
    <span class="article-tag">Repository Pattern</span>
    <span class="article-tag">Data Access</span>
    <span class="article-tag">Infrastructure</span>
    <span class="article-tag">Clean Architecture</span>
  </div>
</div>

<div class="article-item">
  <h3>
    <a href="/goca/blog/articles/mastering-use-cases">Mastering Use Cases in Clean Architecture</a>
  </h3>
  <p>A deep dive into use cases, application services, and DTOs. Learn what use cases are, how they orchestrate business workflows, and how Goca generates production-ready application layer code with comprehensive examples.</p>
  <div class="article-meta">Architecture • Application Layer</div>
  <div class="article-tags">
    <span class="article-tag">Use Cases</span>
    <span class="article-tag">Application Services</span>
    <span class="article-tag">DTOs</span>
    <span class="article-tag">Clean Architecture</span>
  </div>
</div>

<div class="article-item">
  <h3>
    <a href="/goca/blog/articles/understanding-domain-entities">Understanding Domain Entities in Clean Architecture</a>
  </h3>
  <p>A comprehensive guide to domain entities, their role in Clean Architecture, and how Goca generates production-ready entities following DDD principles. Learn the critical distinction between entities and models, best practices, and testing strategies.</p>
  <div class="article-meta">Architecture • Domain-Driven Design</div>
  <div class="article-tags">
    <span class="article-tag">Domain Entities</span>
    <span class="article-tag">Clean Architecture</span>
    <span class="article-tag">DDD</span>
    <span class="article-tag">Best Practices</span>
  </div>
</div>

<div class="article-item">
  <h3>
    <a href="/goca/blog/articles/example-showcase">Advanced Features Showcase</a>
  </h3>
  <p>Demonstration of blog post capabilities including Mermaid diagrams, code blocks, and markdown features. Learn how to leverage VitePress features for technical documentation.</p>
  <div class="article-meta">Example • Tutorial</div>
  <div class="article-tags">
    <span class="article-tag">Mermaid</span>
    <span class="article-tag">Diagrams</span>
    <span class="article-tag">Code Examples</span>
    <span class="article-tag">Clean Architecture</span>
  </div>
</div>

<div class="coming-soon">
  <h2>Coming Soon</h2>
  <p>Articles are in development. Check back soon for in-depth content on:</p>
  <ul style="text-align: left; max-width: 600px; margin: 2rem auto;">
    <li>Building scalable microservices with Clean Architecture</li>
    <li>Advanced testing strategies: Unit, Integration, and E2E</li>
    <li>Database migration patterns and versioning</li>
    <li>Performance optimization techniques</li>
    <li>Domain-Driven Design with Goca</li>
    <li>Implementing event-driven architectures</li>
  </ul>
</div>

</div>

---

## Submit an Article

Have a great article idea or want to share your experience with Goca? We welcome community contributions.

[Open an issue on GitHub](https://github.com/sazardev/goca/issues/new?title=Article%20Proposal:) to propose an article.
