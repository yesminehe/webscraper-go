import React, { useState } from "react";
import "./App.css";
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Link,
  useNavigate,
  useLocation,
  useParams,
} from "react-router-dom";

function generateSitemapXml(links) {
  const urls = links.map((link) => link.url);
  const xml = `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n${urls
    .map((url) => `  <url>\n    <loc>${url}</loc>\n  </url>`)
    .join("\n")}\n</urlset>`;
  return xml;
}

function downloadSitemap(links) {
  const xml = generateSitemapXml(links);
  const blob = new Blob([xml], { type: "application/xml" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "sitemap.xml";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

function SEOList({ results }) {
  const navigate = useNavigate();
  return (
    <div>
      <button
        style={{
          marginBottom: 16,
          padding: "8px 20px",
          background: "#1976d2",
          color: "#fff",
          border: "none",
          borderRadius: 6,
          cursor: "pointer",
        }}
        onClick={() => downloadSitemap(results.links)}
      >
        Download Sitemap
      </button>
      <h2>Results</h2>
      <p>Total links found: {results.totalLinks}</p>
      <ul style={{ listStyle: "none", padding: 0 }}>
        {results.links.map((link, idx) => (
          <li
            key={idx}
            style={{
              marginBottom: 24,
              border: "1px solid #ccc",
              borderRadius: 8,
              padding: 16,
            }}
          >
            <div style={{ fontWeight: "bold", fontSize: 16 }}>
              <strong>{link.status === 200 ? "✅" : "❌"}</strong>{" "}
              <a href={link.url} target="_blank" rel="noopener noreferrer">
                {link.url}
              </a>
            </div>
            <div>Title: {link.title || "No title"}</div>
            <div>Status: {link.status}</div>
            <button
              style={{ marginTop: 12 }}
              onClick={() => navigate(`/details/${idx}`, { state: { link } })}
            >
              View Details
            </button>
          </li>
        ))}
      </ul>
      <h3>Statistics</h3>
      <ul>
        <li>Visited: {results.stats.visited}</li>
        <li>Success (200): {results.stats.success}</li>
        <li>Errors: {results.stats.errors}</li>
        <li>Not Visited: {results.stats.notVisited}</li>
      </ul>
    </div>
  );
}

function SEODetails() {
  const location = useLocation();
  const { link } = location.state || {};
  const params = useParams();
  const navigate = useNavigate();

  if (!link) {
    return (
      <div>
        No details found for this link. <Link to="/">Back to list</Link>
      </div>
    );
  }
  return (
    <div className="seo-details-card">
      <div className="seo-details-header">
        <span className={`status ${link.status === 200 ? "ok" : "error"}`}>
          {link.status === 200 ? "✅" : "❌"}
        </span>
        <h2>{link.title || "No title"}</h2>
        <a href={link.url} target="_blank" rel="noopener noreferrer">
          {link.url}
        </a>
        <div className="status-code">Status: {link.status}</div>
      </div>

      <div className="seo-section">
        <h3>Meta Information</h3>
        <ul>
          <li className={link.titleEmptyOrShort ? "issue" : "ok"}>
            Title: {link.title || "No title"}{" "}
            {link.titleEmptyOrShort && (
              <span className="warn">(Empty/Short)</span>
            )}
          </li>
          <li
            className={
              link.hasMetaDescription
                ? link.metaDescEmptyOrShort
                  ? "issue"
                  : "ok"
                : "issue"
            }
          >
            Meta Description:{" "}
            {link.hasMetaDescription ? link.metaDescription : "Missing"}
            {link.metaDescEmptyOrShort && (
              <span className="warn">(Empty/Short)</span>
            )}
          </li>
          <li className={link.hasCanonical ? "ok" : "issue"}>
            Canonical Tag: {link.hasCanonical ? "Present" : "Missing"}
          </li>
          <li>
            Robots Meta: {link.hasRobotsMeta ? link.robotsMetaValue : "Missing"}
            {link.hasNoindex && <span className="warn">(Noindex!)</span>}
          </li>
        </ul>
      </div>

      <div className="seo-section">
        <h3>Headings & Content</h3>
        <ul>
          <li
            className={
              link.hasH1 ? (link.multipleH1s ? "issue" : "ok") : "issue"
            }
          >
            H1 Tag: {link.hasH1 ? "Present" : "Missing"}
            {link.multipleH1s && <span className="warn">(Multiple H1s!)</span>}
          </li>
        </ul>
      </div>

      <div className="seo-section">
        <h3>Images & Buttons</h3>
        <ul>
          <li>Total Images: {link.totalImages}</li>
          <li className={link.missingAlts > 0 ? "issue" : "ok"}>
            Images Missing Alts: {link.missingAlts}
          </li>
          <li className={link.buttonsWithoutLabels > 0 ? "issue" : "ok"}>
            Buttons Without Labels: {link.buttonsWithoutLabels}
          </li>
        </ul>
      </div>

      <div className="seo-section">
        <h3>Social & Technical</h3>
        <ul>
          <li className={link.hasFavicon ? "ok" : "issue"}>
            Favicon: {link.hasFavicon ? "Present" : "Missing"}
          </li>
          <li className={link.hasOpenGraph ? "ok" : "issue"}>
            Open Graph Tags: {link.hasOpenGraph ? "Present" : "Missing"}
          </li>
          <li className={link.hasTwitterCard ? "ok" : "issue"}>
            Twitter Card Tags: {link.hasTwitterCard ? "Present" : "Missing"}
          </li>
          <li className={link.hasStructuredData ? "ok" : "issue"}>
            Structured Data: {link.hasStructuredData ? "Present" : "Missing"}
          </li>
          <li className={link.hasViewport ? "ok" : "issue"}>
            Viewport Meta: {link.hasViewport ? "Present" : "Missing"}
          </li>
          <li className={link.hasHtmlLang ? "ok" : "issue"}>
            HTML lang Attribute: {link.hasHtmlLang ? "Present" : "Missing"}
          </li>
        </ul>
      </div>

      <button className="back-btn" onClick={() => navigate("/")}>
        Back to list
      </button>
    </div>
  );
}

function App() {
  const [url, setUrl] = useState("");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setResults(null);
    try {
      const response = await fetch(
        `/api/scrape?url=${encodeURIComponent(url)}`
      );
      if (!response.ok) {
        throw new Error("Failed to fetch results.");
      }
      const data = await response.json();
      setResults(data);
    } catch (err) {
      setError(err.message || "An error occurred.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Router>
      <div className="App">
        <h1>SEO Web Scraper</h1>
        <Routes>
          <Route
            path="/"
            element={
              <>
                <form onSubmit={handleSubmit} style={{ marginBottom: 24 }}>
                  <input
                    type="text"
                    placeholder="Enter website URL"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    style={{ width: 300, padding: 8 }}
                    required
                  />
                  <button
                    type="submit"
                    style={{ marginLeft: 12, padding: "8px 16px" }}
                    disabled={loading}
                  >
                    {loading ? "Scraping..." : "Start Scraping"}
                  </button>
                </form>
                {error && (
                  <div style={{ color: "red", marginBottom: 16 }}>{error}</div>
                )}
                {results && <SEOList results={results} />}
              </>
            }
          />
          <Route path="/details/:index" element={<SEODetails />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
