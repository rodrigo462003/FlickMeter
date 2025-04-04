"use strict";
const reviewForm = document.getElementById('newReviewForm');
const newReviewButton = document.getElementById('newReviewButton');
newReviewButton === null || newReviewButton === void 0 ? void 0 : newReviewButton.addEventListener('click', () => {
    reviewForm === null || reviewForm === void 0 ? void 0 : reviewForm.classList.toggle('h-[340px]');
    reviewForm === null || reviewForm === void 0 ? void 0 : reviewForm.classList.toggle('h-0');
});
const starsHover = (i) => {
    const stars = document.getElementsByClassName("buttonStar");
    let j = 0;
    for (; j < i + 1; j++) {
        stars[j].firstChild.setAttribute("fill-opacity", "1");
    }
    for (; j < 10; j++) {
        stars[j].firstChild.setAttribute("fill-opacity", "0");
    }
};
const newBase = (i) => {
    base = i + 1;
    resetStars();
};
const resetStars = () => {
    const stars = document.getElementsByClassName("buttonStar");
    const starInput = document.getElementById("starInput");
    let j = 0;
    if (starInput)
        starInput.value = base.toString();
    for (; j < base; j++) {
        stars[j].firstChild.setAttribute("fill-opacity", "1");
    }
    for (; j < 10; j++) {
        stars[j].firstChild.setAttribute("fill-opacity", "0");
    }
};
reviewForm === null || reviewForm === void 0 ? void 0 : reviewForm.addEventListener('htmx:after-settle', resetStars);
window.isReviewFormVisible = () => {
    return !(reviewForm === null || reviewForm === void 0 ? void 0 : reviewForm.classList.contains('h-0'));
};
