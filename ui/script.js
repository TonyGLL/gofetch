document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('search-form');
  const queryInput = document.getElementById('query');
  const resultsList = document.getElementById('results');

  form.addEventListener('submit', async (event) => {
    event.preventDefault();
    const query = queryInput.value.trim();
    const start = performance.now();

    if (!query) {
      resultsList.innerHTML = '';
      return;
    }

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/search?q=${encodeURIComponent(
          query
        )}&page=1&limit=10`
      );
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const results = await response.json();
      const end = performance.now();
      const duration = (end - start) / 1000;

      renderResults(results, duration);
    } catch (error) {
      console.error('Fetch error:', error);
      resultsList.innerHTML =
        '<li>Error performing search. See console for details.</li>';
    }
  });

  function renderResults(results, duration) {
    if (!results.data || results.data.length === 0) {
      resultsList.innerHTML = '<li>No results found.</li>';
      return;
    }

    const resultHeader = document.createElement('p');
    resultHeader.textContent = `About ${
      results.total
    } results (${duration.toFixed(4)} seconds)`;
    resultsList.innerHTML = '';
    resultsList.appendChild(resultHeader);

    resultsList.innerHTML += results.data
      .map(
        (result) => `
            <li>
                <a href="${result.url}">${result.title}</a>
            </li>
        `
      )
      .join('');
  }
});
