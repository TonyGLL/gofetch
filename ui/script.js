document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('search-form');
  const queryInput = document.getElementById('query');
  const resultsList = document.getElementById('results');

  form.addEventListener('submit', async (event) => {
    event.preventDefault();
    const query = queryInput.value.trim();

    if (!query) {
      resultsList.innerHTML = '';
      return;
    }

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/search?q=${encodeURIComponent(query)}`
      );
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const results = await response.json();

      renderResults(results);
    } catch (error) {
      console.error('Fetch error:', error);
      resultsList.innerHTML =
        '<li>Error performing search. See console for details.</li>';
    }
  });

  function renderResults(results) {
    if (!results || results.length === 0) {
      resultsList.innerHTML = '<li>No results found.</li>';
      return;
    }

    resultsList.innerHTML = results
      .map(
        (result) => `
            <li>
                <a href="${result.url}">${result.title}</a>
                <p>Score: ${result.score.toFixed(4)}</p>
            </li>
        `
      )
      .join('');
  }
});
