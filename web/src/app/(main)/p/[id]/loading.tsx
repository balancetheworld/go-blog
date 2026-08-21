export default function PostLoading() {
  return (
    <main className="article-detail-main">
      <div className="article-loading-progress" aria-hidden="true">
        <span />
      </div>
      <section
        className="article-detail-shell article-loading"
        role="status"
        aria-label="Loading article"
      >
        <div className="article-loading-back" aria-hidden="true" />
        <div className="article-detail-content" aria-hidden="true">
          <div className="article-loading-title" />
          <div className="article-loading-meta">
            <span />
            <span />
          </div>
          <div className="article-loading-body">
            <span className="article-loading-line" />
            <span className="article-loading-line" />
            <span className="article-loading-line is-short" />
            <span className="article-loading-line" />
            <span className="article-loading-line is-medium" />
          </div>
        </div>
      </section>
    </main>
  )
}
