"use strict";
const searchResults = document.getElementById('search-results');
const searchMovies = document.getElementById('search');
function isResultsEmpty() {
    return (searchResults === null || searchResults === void 0 ? void 0 : searchResults.innerHTML) === '';
}
;
document.addEventListener('click', (event) => {
    if (!(searchMovies === null || searchMovies === void 0 ? void 0 : searchMovies.contains(event.target)) && !(searchResults === null || searchResults === void 0 ? void 0 : searchResults.contains(event.target))) {
        searchResults === null || searchResults === void 0 ? void 0 : searchResults.replaceChildren();
    }
});
const profDrop = document.getElementById('profileDrop');
const profButton = document.getElementById('profileButton');
profButton === null || profButton === void 0 ? void 0 : profButton.addEventListener('click', () => {
    profDrop === null || profDrop === void 0 ? void 0 : profDrop.classList.toggle('hidden');
});
document.addEventListener('click', (event) => {
    if (!(profDrop === null || profDrop === void 0 ? void 0 : profDrop.contains(event.target)) && !(profButton === null || profButton === void 0 ? void 0 : profButton.contains(event.target))) {
        profDrop === null || profDrop === void 0 ? void 0 : profDrop.classList.add('hidden');
    }
});
