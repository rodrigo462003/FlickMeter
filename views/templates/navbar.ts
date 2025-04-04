const searchResults = document.getElementById('search-results') as HTMLDivElement | null;
const searchMovies = document.getElementById('search') as HTMLInputElement | null;

function isResultsEmpty() {
    return searchResults?.innerHTML === '';
};

document.addEventListener('click', (event) => {
    if (!searchMovies?.contains(event.target as Node) && !searchResults?.contains(event.target as Node)) {
        searchResults?.replaceChildren()
    }
});





const profDrop = document.getElementById('profileDrop') as HTMLDivElement | null;
const profButton = document.getElementById('profileButton') as HTMLButtonElement | null;

profButton?.addEventListener('click', () => {
    profDrop?.classList.toggle('hidden');
});

document.addEventListener('click', (event) => {
    if (!profDrop?.contains(event.target as Node) && !profButton?.contains(event.target as Node)) {
        profDrop?.classList.add('hidden');
    }
});
