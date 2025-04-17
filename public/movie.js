"use strict";
const reviewForm = document.getElementById('newReviewForm');
const newReviewButton = document.getElementById('newReviewButton');
newReviewButton === null || newReviewButton === void 0 ? void 0 : newReviewButton.addEventListener('click', () => {
    reviewForm === null || reviewForm === void 0 ? void 0 : reviewForm.classList.toggle('h-[340px]');
    reviewForm === null || reviewForm === void 0 ? void 0 : reviewForm.classList.toggle('h-0');
});
const starsHover = (i) => {
    let j = 0;
    for (; j < i + 1; j++) {
        stars[j].firstChild.setAttribute("fill-opacity", "1");
    }
    for (; j < 10; j++) {
        stars[j].firstChild.setAttribute("fill-opacity", "0");
    }
};
window.starsHover = starsHover;
const newBase = (i) => {
    base = i + 1;
    resetStars();
};
window.newBase = newBase;
const stars = document.getElementsByClassName("buttonStar");
const starInput = document.getElementById("starInput");
const resetStars = () => {
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
window.resetStars = resetStars;
reviewForm === null || reviewForm === void 0 ? void 0 : reviewForm.addEventListener('htmx:after-settle', resetStars);
window.isReviewFormVisible = () => {
    return !(reviewForm === null || reviewForm === void 0 ? void 0 : reviewForm.classList.contains('h-0'));
};
