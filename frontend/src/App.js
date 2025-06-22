import React, { useState } from "react";
import "./App.css";

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
    <div className="App">
      <h1>SEO Web Scraper</h1>
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
      {error && <div style={{ color: "red", marginBottom: 16 }}>{error}</div>}
      {results && (
        <div>
          <h2>Results</h2>
          <p>Total links found: {results.totalLinks}</p>
          <ul>
            {results.links.map((link, idx) => (
              <li key={idx} style={{ marginBottom: 8 }}>
                <strong>{link.status === 200 ? "✅" : "❌"}</strong>{" "}
                <a href={link.url} target="_blank" rel="noopener noreferrer">
                  {link.url}
                </a>{" "}
                <br />
                <span>Title: {link.title || "No title"}</span> <br />
                <span>Status: {link.status}</span>
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
      )}
    </div>
  );
}

export default App;
